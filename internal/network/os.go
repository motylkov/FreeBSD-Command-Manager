package network

import (
	"fmt"
	"os/exec"
)

type BareOSManager struct{}

func NewBareOSManager() *BareOSManager {
	return &BareOSManager{}
}

// iface
func (m *BareOSManager) CreateInterface(name string) error {
	cmd := exec.Command("ifconfig", name, "create")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to create interface %s: %v, output: %s", name, err, string(out))
	}
	return nil
}

func (m *BareOSManager) DeleteInterface(name string) error {
	return nil
}

// bridge
func (m *BareOSManager) CreateBridge(name string) error {
	return nil
}

func (m *BareOSManager) DeleteBridge(name string) error {
	return nil
}

// vlan
func (m *BareOSManager) CreateVLAN(name, parent string, vlanID int) error {
	return nil
}

func (m *BareOSManager) DeleteVLAN(name string) error {
	return nil
}

// vxlan

// gre
func (m *BareOSManager) CreateGRE(name string, remote, local string) error {
	return nil
}

func (m *BareOSManager) DeleteGRE(name string) error {
	return nil
}
