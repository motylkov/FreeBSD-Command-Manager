package bareos

import (
	"FreeBSD-Command-manager/pkg/ifconfig"
	"fmt"
	"os/exec"
)

// Config represents network interface configuration
type Config struct {
	Name   string
	Type   string // "iface", "bridge", "vlan", "gre"
	Parent string // For VLAN interfaces
	VLANID int    // For VLAN interfaces
	Local  string // For GRE tunnels
	Remote string // For GRE tunnels
}

// ManagerInterface defines the interface for network operations
type ManagerInterface interface {
	CreateInterface(name string) error
	DeleteInterface(name string) error
	CreateBridge(name string) error
	DeleteBridge(name string) error
	AddInterfaceToBridge(bridgeName, interfaceName string) error
	RemoveInterfaceFromBridge(bridgeName, interfaceName string) error
	CreateVLAN(name, parent string, vlanID int) error
	DeleteVLAN(name string) error
	CreateGRE(name, remote, local string) error
	DeleteGRE(name string) error
	CreateVXLAN(name, local, remote, group, dev string, vxlanID int) error
	DeleteVXLAN(name string) error
	List() ([]ifconfig.Info, error)
	GetInfo(name string) (*ifconfig.Info, error)
}

// CommandExecutor defines the interface for executing system commands
type CommandExecutor interface {
	Execute(name string, args ...string) (string, error)
}

// Manager implements Manager for BareOS network operations
type Manager struct {
	cmdExec CommandExecutor
}

// NewManager creates a new BareOS network manager
func NewManager(cmdExec CommandExecutor) *Manager {
	return &Manager{
		cmdExec: cmdExec,
	}
}

// List returns information about all network interfaces
func (n *Manager) List() ([]ifconfig.Info, error) {
	output, err := n.cmdExec.Execute("ifconfig")
	if err != nil {
		return nil, fmt.Errorf("failed to list interfaces: %v", err)
	}
	return ifconfig.ParseIfconfig(output), nil
}

// GetInfo returns information about a specific network interface
func (n *Manager) GetInfo(name string) (*ifconfig.Info, error) {
	if name == "" {
		return nil, fmt.Errorf("interface name is required")
	}
	output, err := n.cmdExec.Execute("ifconfig", name)
	if err != nil {
		return nil, fmt.Errorf("failed to get info for interface %s: %v", name, err)
	}
	infos := ifconfig.ParseIfconfig(output)
	if len(infos) == 0 {
		return nil, fmt.Errorf("interface %s not found", name)
	}
	return &infos[0], nil
}

// RealCommandExecutor implements CommandExecutor for real system commands
type RealCommandExecutor struct{}

// NewRealCommandExecutor creates a new real command executor
func NewRealCommandExecutor() *RealCommandExecutor {
	return &RealCommandExecutor{}
}

// Execute runs a system command and returns the output
func (r *RealCommandExecutor) Execute(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// DefaultManager returns the default network manager instance
func DefaultManager() ManagerInterface {
	cmdExec := NewRealCommandExecutor()
	return NewManager(cmdExec)
}
