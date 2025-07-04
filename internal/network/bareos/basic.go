package bareos

import (
	"fmt"
)

// CreateInterface creates a generic network interface
func (n *BareOSManager) CreateInterface(name string) error {
	if name == "" {
		return fmt.Errorf("interface name is required")
	}

	_, err := n.cmdExec.Execute("ifconfig", name, "create")
	if err != nil {
		return fmt.Errorf("failed to create interface %s: %v", name, err)
	}

	// Bring up the interface
	_, err = n.cmdExec.Execute("ifconfig", name, "up")
	if err != nil {
		return fmt.Errorf("failed to bring up interface %s: %v", name, err)
	}

	return nil
}

// DeleteInterface deletes a network interface
func (n *BareOSManager) DeleteInterface(name string) error {
	if name == "" {
		return fmt.Errorf("interface name is required")
	}

	_, err := n.cmdExec.Execute("ifconfig", name, "destroy")
	if err != nil {
		return fmt.Errorf("failed to delete interface %s: %v", name, err)
	}

	return nil
}
