package bareos

import (
	"fmt"
)

// CreateVLAN creates a VLAN interface
func (n *Manager) CreateVLAN(name, parent string, vlanID int) error {
	if parent == "" {
		return fmt.Errorf("parent interface is required")
	}
	if vlanID < 1 || vlanID > 4094 {
		return fmt.Errorf("VLAN ID must be between 1 and 4094")
	}

	// Create VLAN interface
	vlanName, err := n.cmdExec.Execute("ifconfig", "vlan", "create")
	if err != nil {
		return fmt.Errorf("failed to create VLAN interface: %v", err)
	}

	// Configure VLAN
	_, err = n.cmdExec.Execute("ifconfig", vlanName, "vlan", fmt.Sprintf("%d", vlanID), "vlandev", parent)
	if err != nil {
		return fmt.Errorf("failed to configure VLAN %s: %v", name, err)
	}

	// Bring up the interface
	_, err = n.cmdExec.Execute("ifconfig", vlanName, "up")
	if err != nil {
		return fmt.Errorf("failed to bring up VLAN %s: %v", name, err)
	}

	if name != "" {
		// Rename to desired name
		_, err = n.cmdExec.Execute("ifconfig", vlanName, "name", name)
		if err != nil {
			return fmt.Errorf("failed to rename VLAN to %s: %v", name, err)
		}
	}

	return nil
}

// DeleteVLAN deletes a VLAN interface
func (n *Manager) DeleteVLAN(name string) error {
	if name == "" {
		return fmt.Errorf("VLAN name is required")
	}

	return n.DeleteInterface(name)
}
