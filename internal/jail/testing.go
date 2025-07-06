package jail

import (
	"fmt"
	"strings"
)

const (
	jailCommand          = "jail"
	jailDestroyedMessage = "jail destroyed"
)

// MockFileSystemManager implements FileSystemManager for testing
type MockFileSystemManager struct {
	EnsurePathCalled bool
	EnsurePathPath   string
	EnsurePathError  error

	MountCalled bool
	MountSource string
	MountTarget string
	MountError  error

	UnmountCalled bool
	UnmountTarget string
	UnmountError  error
}

// EnsurePath ensures the given path exists (mock implementation).
func (m *MockFileSystemManager) EnsurePath(path string) error {
	m.EnsurePathCalled = true
	m.EnsurePathPath = path
	return m.EnsurePathError
}

// Mount mounts a source to a target (mock implementation).
func (m *MockFileSystemManager) Mount(source, target string) error {
	m.MountCalled = true
	m.MountSource = source
	m.MountTarget = target
	return m.MountError
}

// Unmount unmounts a target (mock implementation).
func (m *MockFileSystemManager) Unmount(target string) error {
	m.UnmountCalled = true
	m.UnmountTarget = target
	return m.UnmountError
}

// MockCommandExecutor implements CommandExecutor for testing
type MockCommandExecutor struct {
	ExecuteCalled bool
	ExecuteName   string
	ExecuteArgs   []string
	ExecuteOutput string
	ExecuteError  error
}

// Execute executes a command (mock implementation).
func (m *MockCommandExecutor) Execute(name string, args ...string) (string, error) {
	m.ExecuteCalled = true
	m.ExecuteName = name
	m.ExecuteArgs = args
	return m.ExecuteOutput, m.ExecuteError
}

// CustomCommandExecutor implements CommandExecutor for testing destroy scenarios
type CustomCommandExecutor struct {
	StopError        error
	DestroyError     error
	CallCount        int
	ExecutedCommands []string // Track what commands were executed
	IsDestroyMode    bool     // Track if we're in destroy mode
}

// Execute executes a command (mock implementation for destroy scenarios).
func (c *CustomCommandExecutor) Execute(name string, args ...string) (string, error) {
	c.CallCount++

	// Build the full command string for tracking
	command := name
	if len(args) > 0 {
		command += " " + strings.Join(args, " ")
	}
	c.ExecutedCommands = append(c.ExecutedCommands, command)

	// Simulate different jail commands based on the actual command being executed
	if name == jailCommand && len(args) > 0 {
		switch args[0] {
		case "stop":
			return "jail stopped", c.StopError
		case "destroy":
			return jailDestroyedMessage, c.DestroyError
		case "start":
			return "jail started", nil
		case "create":
			return "jail created", nil
		case "list":
			return "jail1\njail2\njail3", nil
		case "info":
			return "jail information", nil
		case "-r":
			if c.IsDestroyMode {
				return jailDestroyedMessage, c.DestroyError
			}
			// First -r call is stop, second is destroy
			if c.CallCount == 1 {
				return "jail stopped", c.StopError
			}
			return jailDestroyedMessage, c.DestroyError
		}
	}

	// Default response for unknown commands
	return fmt.Sprintf("executed: %s", command), nil
}

// NewMockManager creates a new jail manager with mock implementations for testing
func NewMockManager(fsManager FileSystemManager, cmdExec CommandExecutor) Manager {
	return NewFreeBSDJailManager(fsManager, cmdExec)
}

// NewTestManager creates a new jail manager with default mock implementations
func NewTestManager() Manager {
	mockFS := &MockFileSystemManager{}
	mockCmd := &MockCommandExecutor{}
	return NewFreeBSDJailManager(mockFS, mockCmd)
}
