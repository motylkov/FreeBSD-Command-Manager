package bareos

import (
	"fmt"
)

// CreateBridge creates a bridge interface
func (n *Manager) CreateBridge(name string) error {
	// Create bridge interface
	bridgeName, err := n.cmdExec.Execute("ifconfig", "bridge", "create")
	if err != nil {
		return fmt.Errorf("failed to create bridge interface: %v", err)
	}

	// Bring up the bridge
	_, err = n.cmdExec.Execute("ifconfig", bridgeName, "up")
	if err != nil {
		return fmt.Errorf("failed to bring up bridge %s: %v", bridgeName, err)
	}

	if name != "" {
		// Rename to desired name
		_, err = n.cmdExec.Execute("ifconfig", bridgeName, "name", name)
		if err != nil {
			return fmt.Errorf("failed to rename bridge to %s: %v", name, err)
		}
	}

	return nil
}

// DeleteBridge deletes a bridge interface
func (n *Manager) DeleteBridge(name string) error {
	if name == "" {
		return fmt.Errorf("bridge name is required")
	}

	return n.DeleteInterface(name)
}
