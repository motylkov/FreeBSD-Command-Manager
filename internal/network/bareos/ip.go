package bareos

import (
	"fmt"
	"os/exec"
)

const defaultFamily = "inet"

// AddIP adds an IP address to an interface (IPv4 or IPv6).
func AddIP(iface, ip string, mask int, family string) error {
	if iface == "" || ip == "" || mask == 0 {
		return fmt.Errorf("iface, ip, and mask are required")
	}
	fam := family
	if fam == "" {
		fam = defaultFamily
	}
	addr := fmt.Sprintf("%s/%d", ip, mask)
	execCmd := exec.Command("ifconfig", iface, fam, addr, "add") //nolint:gosec
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ifconfig error: %v, output: %s", err, string(output))
	}
	return nil
}

// AliasIP adds an alias IP address to an interface (IPv4 or IPv6).
func AliasIP(iface, ip string, mask int, family string) error {
	if iface == "" || ip == "" || mask == 0 {
		return fmt.Errorf("iface, ip, and mask are required")
	}
	fam := family
	if fam == "" {
		fam = defaultFamily
	}
	addr := fmt.Sprintf("%s/%d", ip, mask)
	execCmd := exec.Command("ifconfig", iface, fam, addr, "alias") //nolint:gosec
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ifconfig error: %v, output: %s", err, string(output))
	}
	return nil
}

// DeleteIP deletes an IP address from an interface (IPv4 or IPv6).
func DeleteIP(iface, ip string, mask int, family string) error {
	if iface == "" || ip == "" || mask == 0 {
		return fmt.Errorf("iface, ip, and mask are required")
	}
	fam := family
	if fam == "" {
		fam = defaultFamily
	}
	addr := fmt.Sprintf("%s/%d", ip, mask)
	execCmd := exec.Command("ifconfig", iface, fam, addr, "delete") //nolint:gosec
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ifconfig error: %v, output: %s", err, string(output))
	}
	return nil
}
