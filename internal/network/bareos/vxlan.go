package bareos

import (
	"fmt"
	"net"
	"strconv"
)

// CreateVXLAN creates a VXLAN interface
func (n *BareOSManager) CreateVXLAN(name, local, remote, group, dev string, vxlanID int) error {
	if local == "" {
		return fmt.Errorf("local address is required")
	}
	if remote == "" {
		return fmt.Errorf("remote address is required")
	}
	if vxlanID < 1 || vxlanID > 16777215 {
		return fmt.Errorf("VXLAN ID must be between 1 and 16777215")
	}

	// Validate IP addresses
	if net.ParseIP(local) == nil {
		return fmt.Errorf("invalid local IP address: %s", local)
	}
	if net.ParseIP(remote) == nil {
		return fmt.Errorf("invalid remote IP address: %s", remote)
	}

	// Create VXLAN interface
	vxlanName, err := n.cmdExec.Execute("ifconfig", "vxlan", "create")
	if err != nil {
		return fmt.Errorf("failed to create VXLAN interface: %v", err)
	}

	// Build VXLAN configuration command
	args := []string{vxlanName, "vxlan", "vni", strconv.Itoa(vxlanID), "remote", remote, "local", local}

	// Add optional parameters
	if group != "" {
		args = append(args, "group", group)
	}
	if dev != "" {
		args = append(args, "dev", dev)
	}

	// Configure VXLAN tunnel
	_, err = n.cmdExec.Execute("ifconfig", args...)
	if err != nil {
		return fmt.Errorf("failed to configure VXLAN tunnel %s: %v", name, err)
	}

	// Bring up the interface
	_, err = n.cmdExec.Execute("ifconfig", vxlanName, "up")
	if err != nil {
		return fmt.Errorf("failed to bring up VXLAN interface %s: %v", name, err)
	}

	if name != "" {
		// Rename to desired name
		_, err = n.cmdExec.Execute("ifconfig", vxlanName, "name", name)
		if err != nil {
			return fmt.Errorf("failed to rename VXLAN to %s: %v", name, err)
		}
	}

	return nil
}

// DeleteVXLAN deletes a VXLAN interface
func (n *BareOSManager) DeleteVXLAN(name string) error {
	if name == "" {
		return fmt.Errorf("VXLAN name is required")
	}

	return n.DeleteInterface(name)
}
