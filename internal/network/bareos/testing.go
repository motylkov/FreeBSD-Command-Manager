package bareos

// MockCommandExecutor implements CommandExecutor for testing
type MockCommandExecutor struct {
	commands []string
	outputs  map[string]string
	errors   map[string]error
}

// NewMockCommandExecutor creates a new mock command executor
func NewMockCommandExecutor() *MockCommandExecutor {
	return &MockCommandExecutor{
		commands: []string{},
		outputs:  make(map[string]string),
		errors:   make(map[string]error),
	}
}

// Execute simulates command execution for testing
func (m *MockCommandExecutor) Execute(name string, args ...string) (string, error) {
	// Build command string for tracking
	cmdStr := name
	for _, arg := range args {
		cmdStr += " " + arg
	}
	m.commands = append(m.commands, cmdStr)

	// Check if we have a predefined error for this command
	if err, exists := m.errors[cmdStr]; exists {
		return "", err
	}

	// Return predefined output or empty string
	if output, exists := m.outputs[cmdStr]; exists {
		return output, nil
	}

	return "", nil
}

// SetOutput sets the output for a specific command
func (m *MockCommandExecutor) SetOutput(command, output string) {
	m.outputs[command] = output
}

// SetError sets an error for a specific command
func (m *MockCommandExecutor) SetError(command string, err error) {
	m.errors[command] = err
}

// GetCommands returns all executed commands
func (m *MockCommandExecutor) GetCommands() []string {
	return m.commands
}

// ClearCommands clears the command history
func (m *MockCommandExecutor) ClearCommands() {
	m.commands = []string{}
}

// NewTestManager creates a new BareOS manager for testing
func NewTestManager() ManagerInterface {
	mockCmd := NewMockCommandExecutor()
	return NewManager(mockCmd)
}

// NewTestManagerWithExecutor creates a new BareOS manager with a specific command executor
func NewTestManagerWithExecutor(cmdExec CommandExecutor) ManagerInterface {
	return NewManager(cmdExec)
}
