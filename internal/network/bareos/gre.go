package bareos

import (
	"fmt"
)

// CreateGRE creates a GRE tunnel interface
func (n *Manager) CreateGRE(name, remote, local string) error {
	if remote == "" {
		return fmt.Errorf("remote address is required")
	}
	if local == "" {
		return fmt.Errorf("local address is required")
	}

	// Create GRE interface
	greName, err := n.cmdExec.Execute("ifconfig", "gre", "create")
	if err != nil {
		return fmt.Errorf("failed to create GRE interface: %v", err)
	}

	// Configure GRE tunnel
	_, err = n.cmdExec.Execute("ifconfig", greName, "tunnel", local, remote)
	if err != nil {
		return fmt.Errorf("failed to configure GRE tunnel %s: %v", name, err)
	}

	// Bring up the interface
	_, err = n.cmdExec.Execute("ifconfig", greName, "up")
	if err != nil {
		return fmt.Errorf("failed to bring up gre interface %s: %v", name, err)
	}

	if name != "" {
		// Rename to desired name
		_, err = n.cmdExec.Execute("ifconfig", greName, "name", name)
		if err != nil {
			return fmt.Errorf("failed to rename GRE to %s: %v", name, err)
		}
	}

	return nil
}

// DeleteGRE deletes a GRE tunnel interface
func (n *Manager) DeleteGRE(name string) error {
	if name == "" {
		return fmt.Errorf("GRE tunnel name is required")
	}

	return n.DeleteInterface(name)
}
