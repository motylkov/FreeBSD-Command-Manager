// Package jail provides management for FreeBSD jails.
package jail

import (
	pkgjail "FreeBSD-Command-manager/pkg/jail"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const (
	// DefaultDirectoryPermissions is the default permission mode for creating directories
	DefaultDirectoryPermissions = 0o755
)

// Config represents the configuration for a jail
type Config struct {
	Name  string
	Path  string
	IP    string
	Mount string
}

// Manager defines the interface for jail operations
type Manager interface {
	Create(cfg Config) error
	Start(name string) error
	Stop(name string) error
	Destroy(name string) error
	List() ([]pkgjail.Info, error)
	GetInfo(name string) (*pkgjail.Info, error)
}

// FileSystemManager defines the interface for file system operations
type FileSystemManager interface {
	EnsurePath(path string) error
	Mount(source, target string) error
	Unmount(target string) error
}

// CommandExecutor defines the interface for executing system commands
type CommandExecutor interface {
	Execute(name string, args ...string) (string, error)
}

// FreeBSDJailManager implements Manager for FreeBSD jails
type FreeBSDJailManager struct {
	fsManager FileSystemManager
	cmdExec   CommandExecutor
}

// NewFreeBSDJailManager creates a new FreeBSD jail manager
func NewFreeBSDJailManager(fsManager FileSystemManager, cmdExec CommandExecutor) *FreeBSDJailManager {
	return &FreeBSDJailManager{
		fsManager: fsManager,
		cmdExec:   cmdExec,
	}
}

// Create a new jail with the given configuration
func (j *FreeBSDJailManager) Create(cfg Config) error {
	// Validate parameters
	if cfg.Name == "" || cfg.Path == "" || cfg.IP == "" {
		return errors.New("missing required parameters (name, path, ip)")
	}

	// Ensure jail path exists
	if err := j.fsManager.EnsurePath(cfg.Path); err != nil {
		return fmt.Errorf("failed to create jail path: %v", err)
	}

	// Mount ZFS dataset or image if specified
	if cfg.Mount != "" {
		if err := j.fsManager.Mount(cfg.Mount, cfg.Path); err != nil {
			return fmt.Errorf("failed to mount %s to %s: %v", cfg.Mount, cfg.Path, err)
		}
	}

	// Create jail
	_, err := j.cmdExec.Execute("jail", "-c",
		"name="+cfg.Name,
		"path="+cfg.Path,
		"host.hostname="+cfg.Name,
		"ip4.addr="+cfg.IP,
		"command=/bin/sh")
	if err != nil {
		return fmt.Errorf("failed to create jail: %v", err)
	}

	return nil
}

// Start an existing jail
func (j *FreeBSDJailManager) Start(name string) error {
	if name == "" {
		return errors.New("jail name is required")
	}

	_, err := j.cmdExec.Execute("jail", "-c", "name="+name)
	if err != nil {
		return fmt.Errorf("failed to start jail %s: %v", name, err)
	}

	return nil
}

// Stop a running jail
func (j *FreeBSDJailManager) Stop(name string) error {
	if name == "" {
		return errors.New("jail name is required")
	}

	_, err := j.cmdExec.Execute("jail", "-r", name)
	if err != nil {
		return fmt.Errorf("failed to stop jail %s: %v", name, err)
	}

	return nil
}

// Destroy a jail completely
func (j *FreeBSDJailManager) Destroy(name string) error {
	if name == "" {
		return errors.New("jail name is required")
	}

	// Stop jail if it's running
	if err := j.Stop(name); err != nil {
		// Log warning but continue with destruction
		fmt.Printf("Warning: could not stop jail %s (may already be stopped): %v\n", name, err)
	}

	// Remove jail configuration
	_, err := j.cmdExec.Execute("jail", "-r", name)
	if err != nil {
		return fmt.Errorf("failed to destroy jail %s: %v", name, err)
	}

	return nil
}

// List information about all jails (short)
func (j *FreeBSDJailManager) List() ([]pkgjail.Info, error) {
	output, err := j.cmdExec.Execute("jail", "-l")
	if err != nil {
		return nil, fmt.Errorf("failed to list jails: %v", err)
	}
	jails, err := pkgjail.ParseJailList(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jail list: %v", err)
	}
	return jails, nil
}

// GetInfo returns information about a specific jail.
func (j *FreeBSDJailManager) GetInfo(name string) (*pkgjail.Info, error) {
	if name == "" {
		return nil, errors.New("jail name is required")
	}
	output, err := j.cmdExec.Execute("jail", "-l")
	if err != nil {
		return nil, fmt.Errorf("failed to get jail info for %s: %v", name, err)
	}
	jails, err := pkgjail.ParseJailList(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jail list: %v", err)
	}
	for _, jail := range jails {
		if jail.Name == name {
			return &jail, nil
		}
	}
	return nil, fmt.Errorf("jail %s not found", name)
}

// RealFileSystemManager implements FileSystemManager using real file system operations
type RealFileSystemManager struct{}

// NewRealFileSystemManager creates a new real file system manager
func NewRealFileSystemManager() *RealFileSystemManager {
	return &RealFileSystemManager{}
}

// EnsurePath ensures the given path exists, creating it if necessary
func (r *RealFileSystemManager) EnsurePath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, DefaultDirectoryPermissions); err != nil {
			return fmt.Errorf("failed to create path %s: %v", path, err)
		}
	}
	return nil
}

// Mount a source to a target using nullfs
func (r *RealFileSystemManager) Mount(source, target string) error {
	cmd := exec.Command("mount", "-t", "nullfs", source, target)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to mount %s to %s: %v, output: %s", source, target, err, string(out))
	}
	return nil
}

// Unmount unmounts a target
func (r *RealFileSystemManager) Unmount(target string) error {
	cmd := exec.Command("umount", target)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to unmount %s: %v, output: %s", target, err, string(out))
	}
	return nil
}

// RealCommandExecutor implements CommandExecutor using real command execution
type RealCommandExecutor struct{}

// NewRealCommandExecutor creates a new real command executor
func NewRealCommandExecutor() *RealCommandExecutor {
	return &RealCommandExecutor{}
}

// Execute a command and returns the output
func (r *RealCommandExecutor) Execute(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

// DefaultManager returns a default jail manager with real implementations
func DefaultManager() Manager {
	fsManager := NewRealFileSystemManager()
	cmdExec := NewRealCommandExecutor()
	return NewFreeBSDJailManager(fsManager, cmdExec)
}
