package bareos

import (
	"fmt"
	"os/exec"
	"strings"
)

// AddRoute adds a route for the given network and gateway (IPv4 or IPv6).
func AddRoute(family, network, gw string) error {
	fam := family
	if fam == "" {
		fam = "inet"
	}
	cmd := exec.Command("route", "-n", "add", "-"+fam, network, gw)
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
		fam = "inet"
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
	cmd := exec.Command("route", "-n", "delete", "-"+fam, network)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("route delete error: %v, output: %s", err, string(output))
	}
	return nil
}

// ListRoutes lists all routes for the given family (IPv4 or IPv6).
func ListRoutes(family string) (string, error) {
	fam := family
	if fam == "" {
		fam = "inet"
	}
	cmd := exec.Command("netstat", "-rn", "-f", fam)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("route list error: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

// countDefaultRoutes returns the number of default routes for the given family.
func countDefaultRoutes(family string) (int, error) {
	out, err := ListRoutes(family)
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
