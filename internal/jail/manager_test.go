package jail

import (
	"errors"
	"testing"
)

func TestFreeBSDJailManager_Create(t *testing.T) {
	tests := []struct {
		name        string
		cfg         Config
		fsError     error
		cmdError    error
		expectError bool
	}{
		{
			name: "successful creation",
			cfg: Config{
				Name: "test-jail",
				Path: "/jails/test-jail",
				IP:   "192.168.1.100",
			},
			expectError: false,
		},
		{
			name: "missing name",
			cfg: Config{
				Path: "/jails/test-jail",
				IP:   "192.168.1.100",
			},
			expectError: true,
		},
		{
			name: "filesystem error",
			cfg: Config{
				Name: "test-jail",
				Path: "/jails/test-jail",
				IP:   "192.168.1.100",
			},
			fsError:     errors.New("filesystem error"),
			expectError: true,
		},
		{
			name: "command execution error",
			cfg: Config{
				Name: "test-jail",
				Path: "/jails/test-jail",
				IP:   "192.168.1.100",
			},
			cmdError:    errors.New("command error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockFileSystemManager{EnsurePathError: tt.fsError}
			mockCmd := &MockCommandExecutor{ExecuteError: tt.cmdError}

			manager := NewFreeBSDJailManager(mockFS, mockCmd)

			err := manager.Create(tt.cfg)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				// Verify filesystem manager was called
				if !mockFS.EnsurePathCalled {
					t.Error("EnsurePath was not called")
				}
				if mockFS.EnsurePathPath != tt.cfg.Path {
					t.Errorf("expected path %s, got %s", tt.cfg.Path, mockFS.EnsurePathPath)
				}

				// Verify command executor was called
				if !mockCmd.ExecuteCalled {
					t.Error("Execute was not called")
				}
				if mockCmd.ExecuteName != "jail" {
					t.Errorf("expected command 'jail', got %s", mockCmd.ExecuteName)
				}
			}
		})
	}
}

func TestFreeBSDJailManager_Start(t *testing.T) {
	tests := []struct {
		name        string
		jailName    string
		cmdError    error
		expectError bool
	}{
		{
			name:        "successful start",
			jailName:    "test-jail",
			expectError: false,
		},
		{
			name:        "empty name",
			jailName:    "",
			expectError: true,
		},
		{
			name:        "command error",
			jailName:    "test-jail",
			cmdError:    errors.New("command error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockFileSystemManager{}
			mockCmd := &MockCommandExecutor{ExecuteError: tt.cmdError}

			manager := NewFreeBSDJailManager(mockFS, mockCmd)

			err := manager.Start(tt.jailName)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if !mockCmd.ExecuteCalled {
					t.Error("Execute was not called")
				}
				if mockCmd.ExecuteName != "jail" {
					t.Errorf("expected command 'jail', got %s", mockCmd.ExecuteName)
				}
			}
		})
	}
}

func TestFreeBSDJailManager_Stop(t *testing.T) {
	tests := []struct {
		name        string
		jailName    string
		cmdError    error
		expectError bool
	}{
		{
			name:        "successful stop",
			jailName:    "test-jail",
			expectError: false,
		},
		{
			name:        "empty name",
			jailName:    "",
			expectError: true,
		},
		{
			name:        "command error",
			jailName:    "test-jail",
			cmdError:    errors.New("command error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockFileSystemManager{}
			mockCmd := &MockCommandExecutor{ExecuteError: tt.cmdError}

			manager := NewFreeBSDJailManager(mockFS, mockCmd)

			err := manager.Stop(tt.jailName)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if !mockCmd.ExecuteCalled {
					t.Error("Execute was not called")
				}
				if mockCmd.ExecuteName != "jail" {
					t.Errorf("expected command 'jail', got %s", mockCmd.ExecuteName)
				}
			}
		})
	}
}

func TestFreeBSDJailManager_Destroy(t *testing.T) {
	tests := []struct {
		name         string
		jailName     string
		stopError    error
		destroyError error
		expectError  bool
	}{
		{
			name:        "successful destroy",
			jailName:    "test-jail",
			expectError: false,
		},
		{
			name:        "empty name",
			jailName:    "",
			expectError: true,
		},
		{
			name:        "stop error (should continue)",
			jailName:    "test-jail",
			stopError:   errors.New("stop error"),
			expectError: false, // Should continue with destroy
		},
		{
			name:         "destroy error",
			jailName:     "test-jail",
			destroyError: errors.New("destroy error"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockFileSystemManager{}
			customCmd := &CustomCommandExecutor{
				StopError:    tt.stopError,
				DestroyError: tt.destroyError,
			}

			manager := NewFreeBSDJailManager(mockFS, customCmd)

			err := manager.Destroy(tt.jailName)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestNewTestManager(t *testing.T) {
	// Test that NewTestManager creates a working manager
	manager := NewTestManager()
	if manager == nil {
		t.Error("NewTestManager returned nil")
	}

	// Test that it's a FreeBSDJailManager with mock implementations
	_, ok := manager.(*FreeBSDJailManager)
	if !ok {
		t.Error("NewTestManager did not return a FreeBSDJailManager")
	}

	// Test that the manager can be used (with mock implementations)
	cfg := Config{
		Name: "test-jail",
		Path: "/jails/test-jail",
		IP:   "192.168.1.100",
	}

	err := manager.Create(cfg)
	if err != nil {
		t.Errorf("NewTestManager Create failed: %v", err)
	}

	err = manager.Start("test-jail")
	if err != nil {
		t.Errorf("NewTestManager Start failed: %v", err)
	}

	err = manager.Stop("test-jail")
	if err != nil {
		t.Errorf("NewTestManager Stop failed: %v", err)
	}

	err = manager.Destroy("test-jail")
	if err != nil {
		t.Errorf("NewTestManager Destroy failed: %v", err)
	}
}
