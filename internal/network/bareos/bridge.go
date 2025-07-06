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

// AddInterfaceToBridge adds an interface to a bridge
func (n *Manager) AddInterfaceToBridge(bridgeName, interfaceName string) error {
	if bridgeName == "" {
		return fmt.Errorf("bridge name is required")
	}
	if interfaceName == "" {
		return fmt.Errorf("interface name is required")
	}

	// Add interface to bridge using ifconfig
	_, err := n.cmdExec.Execute("ifconfig", bridgeName, "addm", interfaceName)
	if err != nil {
		return fmt.Errorf("failed to add interface %s to bridge %s: %v", interfaceName, bridgeName, err)
	}

	return nil
}

// RemoveInterfaceFromBridge removes an interface from a bridge
func (n *Manager) RemoveInterfaceFromBridge(bridgeName, interfaceName string) error {
	if bridgeName == "" {
		return fmt.Errorf("bridge name is required")
	}
	if interfaceName == "" {
		return fmt.Errorf("interface name is required")
	}

	// Remove interface from bridge using ifconfig
	_, err := n.cmdExec.Execute("ifconfig", bridgeName, "deletem", interfaceName)
	if err != nil {
		return fmt.Errorf("failed to remove interface %s from bridge %s: %v", interfaceName, bridgeName, err)
	}

	return nil
}
