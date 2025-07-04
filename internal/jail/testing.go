package jail

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

func (m *MockFileSystemManager) EnsurePath(path string) error {
	m.EnsurePathCalled = true
	m.EnsurePathPath = path
	return m.EnsurePathError
}

func (m *MockFileSystemManager) Mount(source, target string) error {
	m.MountCalled = true
	m.MountSource = source
	m.MountTarget = target
	return m.MountError
}

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

func (m *MockCommandExecutor) Execute(name string, args ...string) (string, error) {
	m.ExecuteCalled = true
	m.ExecuteName = name
	m.ExecuteArgs = args
	return m.ExecuteOutput, m.ExecuteError
}

// CustomCommandExecutor implements CommandExecutor for testing destroy scenarios
type CustomCommandExecutor struct {
	StopError    error
	DestroyError error
	CallCount    int
}

func (c *CustomCommandExecutor) Execute(name string, args ...string) (string, error) {
	c.CallCount++

	// First call is stop, second call is destroy
	if c.CallCount == 1 {
		return "", c.StopError
	}
	return "", c.DestroyError
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
