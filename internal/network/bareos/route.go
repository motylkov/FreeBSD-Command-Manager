package bareos

import (
	"FreeBSD-Command-manager/pkg/netstat"
	"fmt"
	"os/exec"
	"strings"
)

var execCommand = exec.Command

func realListRoutes(family string) (string, error) {
	var cmd *exec.Cmd
	if family == "" {
		cmd = execCommand("netstat", "-rn")
	} else {
		cmd = execCommand("netstat", "-rn", "-f", family)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("route list error: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

var listRoutes = realListRoutes

// AddRoute adds a route for the given network and gateway (IPv4 or IPv6).
func AddRoute(family, network, gw string) error {
	fam := family
	if fam == "" {
		fam = inetFamily
	}
	cmd := execCommand("route", "-n", "add", "-"+fam, network, gw)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("route add error: %v, output: %s", err, string(output))
	}
	return nil
}

// AddRouteWithIface adds a route for the given network, gateway, and optional interface (IPv4 or IPv6).
func AddRouteWithIface(family, network, gw, iface string) error {
	fam := family
	if fam == "" {
		fam = inetFamily
	}
	args := []string{"-n", "add", "-" + fam, network, gw}
	if iface != "" {
		args = append(args, "-ifp", iface)
	}
	cmd := execCommand("route", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("route add error: %v, output: %s", err, string(output))
	}
	return nil
}

// DelRoute deletes a route for the given network (IPv4 or IPv6).
// For default route, prevents deleting the last default route.
func DelRoute(family, network string) error {
	fam := family
	if fam == "" {
		fam = inetFamily
	}
	if network == "default" {
		count, err := countDefaultRoutes(fam)
		if err != nil {
			return err
		}
		if count <= 1 {
			return fmt.Errorf("cannot delete the last default route")
		}
	}
	cmd := execCommand("route", "-n", "delete", "-"+fam, network)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("route delete error: %v, output: %s", err, string(output))
	}
	return nil
}

// ListAllRoutes lists both IPv4 and IPv6 routes as a JSON string using the netstat parser.
func ListAllRoutes(family string) ([]netstat.Route, error) {
	output, err := listRoutes(family)

	if err != nil {
		return nil, fmt.Errorf("route list error: %w, output: %s", err, output)
	}

	routes, err := netstat.ParseNetstat(output)
	if err != nil {
		return nil, fmt.Errorf("parse netstat error: %w", err)
	}
	return routes, nil
}

// countDefaultRoutes returns the number of default routes for the given family.
func countDefaultRoutes(family string) (int, error) {
	out, err := listRoutes(family)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "default" {
			count++
		}
	}
	return count, nil
}
