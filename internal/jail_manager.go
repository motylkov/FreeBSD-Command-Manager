package internal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type JailConfig struct {
	Name  string
	Path  string
	IP    string
	Mount string
}

func CreateJail(cfg JailConfig) error {
	// Validate parameters
	if cfg.Name == "" || cfg.Path == "" || cfg.IP == "" {
		return errors.New("missing required parameters (name, path, ip)")
	}

	// Ensure jail path exists
	if _, err := os.Stat(cfg.Path); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.Path, 0755); err != nil {
			return fmt.Errorf("failed to create jail path: %v", err)
		}
	}

	// Mount ZFS dataset or image
	if cfg.Mount != "" {
		mountCmd := exec.Command("mount", "-t", "nullfs", cfg.Mount, cfg.Path)
		if out, err := mountCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to mount %s to %s: %v, output: %s", cfg.Mount, cfg.Path, err, string(out))
		}
	}

	// Create jail
	jailCmd := exec.Command("jail", "-c",
		"name="+cfg.Name,
		"path="+cfg.Path,
		"host.hostname="+cfg.Name,
		"ip4.addr="+cfg.IP,
		"command=/bin/sh")
	if out, err := jailCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create jail: %v, output: %s", err, string(out))
	}

	return nil
}

func StartJail(name string) error {
	// Start jail
	startCmd := exec.Command("jail", "-c", "name="+name)
	if out, err := startCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start jail %s: %v, output: %s", name, err, string(out))
	}

	return nil
}

func StopJail(name string) error {
	// Stop jail
	stopCmd := exec.Command("jail", "-r", name)
	if out, err := stopCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop jail %s: %v, output: %s", name, err, string(out))
	}

	return nil
}

func DestroyJail(name string) error {
	// Stop jail if it's running
	if err := StopJail(name); err != nil {
		// already stopped
		fmt.Printf("Warning: could not stop jail %s (may already be stopped): %v\n", name, err)
	}

	// Remove jail configuration
	removeCmd := exec.Command("jail", "-r", name)
	if out, err := removeCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to destroy jail %s: %v, output: %s", name, err, string(out))
	}

	return nil
}
