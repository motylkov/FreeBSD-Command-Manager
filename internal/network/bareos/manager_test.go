package bareos

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

const inet6Family = "inet6"

func TestBareOSManager_CreateInterface(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		shouldError   bool
	}{
		{
			name:          "successful interface creation",
			interfaceName: "test0",
			shouldError:   false,
		},
		{
			name:          "empty interface name",
			interfaceName: "",
			shouldError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()
			manager := NewManager(mockCmd)

			err := manager.CreateInterface(tc.interfaceName)

			if tc.shouldError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			commands := mockCmd.GetCommands()
			if len(commands) < 2 {
				t.Errorf("expected at least 2 commands, got %d", len(commands))
			}
			expectedCreateCmd := "ifconfig " + tc.interfaceName + " create"
			expectedUpCmd := "ifconfig " + tc.interfaceName + " up"
			if commands[0] != expectedCreateCmd {
				t.Errorf("expected first command %s, got %s", expectedCreateCmd, commands[0])
			}
			if commands[1] != expectedUpCmd {
				t.Errorf("expected second command %s, got %s", expectedUpCmd, commands[1])
			}
		})
	}
}

func TestBareOSManager_CreateBridge(t *testing.T) {
	tests := []struct {
		name        string
		bridgeName  string
		shouldError bool
	}{
		{
			name:        "successful bridge creation",
			bridgeName:  "br0",
			shouldError: false,
		},
		{
			name:        "bridge creation with default name",
			bridgeName:  "",
			shouldError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()

			// Set up mock responses for bridge creation
			if !tc.shouldError {
				mockCmd.SetOutput("ifconfig bridge create", "bridge0")
				mockCmd.SetOutput("ifconfig bridge0 up", "")
			}

			manager := NewManager(mockCmd)

			err := manager.CreateBridge(tc.bridgeName)

			if tc.shouldError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			commands := mockCmd.GetCommands()
			if len(commands) < 2 {
				t.Errorf("expected at least 2 commands, got %d", len(commands))
			}
			expectedCreateCmd := "ifconfig bridge create"
			expectedUpCmd := "ifconfig bridge0 up"
			if commands[0] != expectedCreateCmd {
				t.Errorf("expected first command %s, got %s", expectedCreateCmd, commands[0])
			}
			if commands[1] != expectedUpCmd {
				t.Errorf("expected second command %s, got %s", expectedUpCmd, commands[1])
			}
		})
	}
}

func TestBareOSManager_CreateVLAN(t *testing.T) {
	tests := []struct {
		name        string
		vlanName    string
		parent      string
		vlanID      int
		shouldError bool
	}{
		{
			name:        "successful VLAN creation",
			vlanName:    "vlan100",
			parent:      "em0",
			vlanID:      100,
			shouldError: false,
		},
		{
			name:        "VLAN creation with default name",
			vlanName:    "",
			parent:      "em0",
			vlanID:      100,
			shouldError: false,
		},
		{
			name:        "empty parent interface",
			vlanName:    "vlan100",
			parent:      "",
			vlanID:      100,
			shouldError: true,
		},
		{
			name:        "invalid VLAN ID (too low)",
			vlanName:    "vlan100",
			parent:      "em0",
			vlanID:      0,
			shouldError: true,
		},
		{
			name:        "invalid VLAN ID (too high)",
			vlanName:    "vlan100",
			parent:      "em0",
			vlanID:      4095,
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()

			// Set up mock responses for VLAN creation
			if !tc.shouldError {
				mockCmd.SetOutput("ifconfig vlan create", "vlan0")
				mockCmd.SetOutput("ifconfig vlan0 vlan 100 vlandev em0", "")
				mockCmd.SetOutput("ifconfig vlan0 up", "")
				if tc.vlanName != "" {
					mockCmd.SetOutput("ifconfig vlan0 name "+tc.vlanName, "")
				}
			}

			manager := NewManager(mockCmd)

			err := manager.CreateVLAN(tc.vlanName, tc.parent, tc.vlanID)

			if tc.shouldError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			commands := mockCmd.GetCommands()
			if len(commands) < 2 {
				t.Errorf("expected at least 2 commands, got %d", len(commands))
			}
		})
	}
}

func TestBareOSManager_CreateGRE(t *testing.T) {
	tests := []struct {
		name        string
		greName     string
		remote      string
		local       string
		shouldError bool
	}{
		{
			name:        "successful GRE creation",
			greName:     "gre0",
			remote:      "192.168.1.1",
			local:       "192.168.1.2",
			shouldError: false,
		},
		{
			name:        "GRE creation with default name",
			greName:     "",
			remote:      "192.168.1.1",
			local:       "192.168.1.2",
			shouldError: false,
		},
		{
			name:        "empty remote address",
			greName:     "gre0",
			remote:      "",
			local:       "192.168.1.2",
			shouldError: true,
		},
		{
			name:        "empty local address",
			greName:     "gre0",
			remote:      "192.168.1.1",
			local:       "",
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()

			// Set up mock responses for GRE creation
			if !tc.shouldError {
				mockCmd.SetOutput("ifconfig gre create", "gre0")
				mockCmd.SetOutput("ifconfig gre0 tunnel 192.168.1.2 192.168.1.1", "")
				mockCmd.SetOutput("ifconfig gre0 up", "")
				if tc.greName != "" {
					mockCmd.SetOutput("ifconfig gre0 name "+tc.greName, "")
				}
			}

			manager := NewManager(mockCmd)

			err := manager.CreateGRE(tc.greName, tc.remote, tc.local)

			if tc.shouldError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			commands := mockCmd.GetCommands()
			if len(commands) < 2 {
				t.Errorf("expected at least 2 commands, got %d", len(commands))
			}
		})
	}
}

func TestBareOSManager_List(t *testing.T) {
	mockCmd := NewMockCommandExecutor()

	// Set up mock ifconfig output
	mockIfconfigOutput := `em0: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
	options=48505bb<RXCSUM,TXCSUM,VLAN_MTU,VLAN_HWTAGGING,JUMBO_MTU,VLAN_HWCSUM,TSO4,LRO,VLAN_HWFILTER,VLAN_HWTSO,HWSTATS,MEXTPG>
	ether 08:00:27:20:af:31
	inet 192.168.88.241 netmask 0xffffff00 broadcast 192.168.88.255
	inet6 fe80::a00:27ff:fe20:af31%em0 prefixlen 64 scopeid 0x1
	media: Ethernet autoselect (1000baseT <full-duplex>)
	status: active
	nd6 options=23<PERFORMNUD,ACCEPT_RTADV,AUTO_LINKLOCAL>
lo0: flags=1008049<UP,LOOPBACK,RUNNING,MULTICAST,LOWER_UP> metric 0 mtu 16384
	options=680003<RXCSUM,TXCSUM,LINKSTATE,RXCSUM_IPV6,TXCSUM_IPV6>
	inet 127.0.0.1 netmask 0xff000000
	inet6 ::1 prefixlen 128
	inet6 fe80::1%lo0 prefixlen 64 scopeid 0x2
	groups: lo
	nd6 options=21<PERFORMNUD,AUTO_LINKLOCAL>`

	mockCmd.SetOutput("ifconfig", mockIfconfigOutput)

	manager := NewManager(mockCmd)

	interfaces, err := manager.List()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(interfaces) == 0 {
		t.Error("expected at least one interface")
	}

	expectedCmd := "ifconfig"
	if len(mockCmd.GetCommands()) == 0 || mockCmd.GetCommands()[0] != expectedCmd {
		t.Errorf("expected command %s, got %v", expectedCmd, mockCmd.GetCommands())
	}
}

func TestBareOSManager_GetInfo(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		shouldError   bool
	}{
		{
			name:          "successful info retrieval",
			interfaceName: "em0",
			shouldError:   false,
		},
		{
			name:          "empty interface name",
			interfaceName: "",
			shouldError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()

			if !tc.shouldError {
				// Set up mock ifconfig output for specific interface
				mockIfconfigOutput := `em0: flags=1008843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST,LOWER_UP> metric 0 mtu 1500
options=48505bb<RXCSUM,TXCSUM,VLAN_MTU,VLAN_HWTAGGING,JUMBO_MTU,VLAN_HWCSUM,TSO4,LRO,VLAN_HWFILTER,VLAN_HWTSO,HWSTATS,MEXTPG>
ether 08:00:27:20:af:31
inet 192.168.88.241 netmask 0xffffff00 broadcast 192.168.88.255
inet6 fe80::a00:27ff:fe20:af31%em0 prefixlen 64 scopeid 0x1
media: Ethernet autoselect (1000baseT <full-duplex>)
status: active
nd6 options=23<PERFORMNUD,ACCEPT_RTADV,AUTO_LINKLOCAL>`

				mockCmd.SetOutput("ifconfig em0", mockIfconfigOutput)
			}

			manager := NewManager(mockCmd)

			info, err := manager.GetInfo(tc.interfaceName)

			if tc.shouldError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if info == nil {
				t.Error("expected info but got nil")
				return
			}
			if info.Name != tc.interfaceName {
				t.Errorf("expected name %s, got %s", tc.interfaceName, info.Name)
			}
			expectedCmd := "ifconfig " + tc.interfaceName
			if len(mockCmd.GetCommands()) == 0 || mockCmd.GetCommands()[0] != expectedCmd {
				t.Errorf("expected command %s, got %v", expectedCmd, mockCmd.GetCommands())
			}
		})
	}
}

func TestNewTestManager(t *testing.T) {
	manager := NewTestManager()

	// Test that it's a BareOSManager with mock implementations
	_, ok := manager.(*Manager)
	if !ok {
		t.Error("NewTestManager did not return a BareOSManager")
	}

	// Test that we can call methods without errors
	err := manager.CreateInterface("test-bridge")
	if err != nil {
		t.Errorf("NewTestManager CreateInterface failed: %v", err)
	}
}

func TestMockCommandExecutor(t *testing.T) {
	mock := NewMockCommandExecutor()

	// Test setting and getting output
	mock.SetOutput("test command", "test output")
	output, err := mock.Execute("test", "command")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if output != "test output" {
		t.Errorf("expected output 'test output', got '%s'", output)
	}

	// Test setting and getting error
	testError := errors.New("test error")
	mock.SetError("error command", testError)
	_, err = mock.Execute("error", "command")
	if err != testError {
		t.Errorf("expected error %v, got %v", testError, err)
	}

	// Test command tracking
	commands := mock.GetCommands()
	if len(commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(commands))
	}

	// Test clearing commands
	mock.ClearCommands()
	commands = mock.GetCommands()
	if len(commands) != 0 {
		t.Errorf("expected 0 commands after clear, got %d", len(commands))
	}
}

func TestBareOSManager_CreateVXLAN(t *testing.T) {
	tests := []struct {
		name        string
		vxlanName   string
		local       string
		remote      string
		group       string
		dev         string
		vxlanID     int
		shouldError bool
	}{
		{
			name:        "successful VXLAN creation",
			vxlanName:   "vxlan0",
			local:       "192.168.1.1",
			remote:      "192.168.1.2",
			group:       "",
			dev:         "",
			vxlanID:     1000,
			shouldError: false,
		},
		{
			name:        "VXLAN creation with default name",
			vxlanName:   "",
			local:       "192.168.1.1",
			remote:      "192.168.1.2",
			group:       "",
			dev:         "",
			vxlanID:     1000,
			shouldError: false,
		},
		{
			name:        "VXLAN creation with group and dev",
			vxlanName:   "vxlan0",
			local:       "192.168.1.1",
			remote:      "192.168.1.2",
			group:       "224.0.0.1",
			dev:         "em0",
			vxlanID:     1000,
			shouldError: false,
		},
		{
			name:        "empty local address",
			vxlanName:   "vxlan0",
			local:       "",
			remote:      "192.168.1.2",
			group:       "",
			dev:         "",
			vxlanID:     1000,
			shouldError: true,
		},
		{
			name:        "empty remote address",
			vxlanName:   "vxlan0",
			local:       "192.168.1.1",
			remote:      "",
			group:       "",
			dev:         "",
			vxlanID:     1000,
			shouldError: true,
		},
		{
			name:        "invalid VXLAN ID (too low)",
			vxlanName:   "vxlan0",
			local:       "192.168.1.1",
			remote:      "192.168.1.2",
			group:       "",
			dev:         "",
			vxlanID:     0,
			shouldError: true,
		},
		{
			name:        "invalid VXLAN ID (too high)",
			vxlanName:   "vxlan0",
			local:       "192.168.1.1",
			remote:      "192.168.1.2",
			group:       "",
			dev:         "",
			vxlanID:     16777216,
			shouldError: true,
		},
		{
			name:        "invalid local IP address",
			vxlanName:   "vxlan0",
			local:       "invalid-ip",
			remote:      "192.168.1.2",
			group:       "",
			dev:         "",
			vxlanID:     1000,
			shouldError: true,
		},
		{
			name:        "invalid remote IP address",
			vxlanName:   "vxlan0",
			local:       "192.168.1.1",
			remote:      "invalid-ip",
			group:       "",
			dev:         "",
			vxlanID:     1000,
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()

			// Set up mock responses for VXLAN creation
			if !tc.shouldError {
				mockCmd.SetOutput("ifconfig vxlan create", "vxlan0")
				// Build expected command based on parameters
				expectedCmd := "ifconfig vxlan0 vxlan vni 1000 remote 192.168.1.2 local 192.168.1.1"
				if tc.group != "" {
					expectedCmd += " group " + tc.group
				}
				if tc.dev != "" {
					expectedCmd += " dev " + tc.dev
				}
				mockCmd.SetOutput(expectedCmd, "")
				mockCmd.SetOutput("ifconfig vxlan0 up", "")
				if tc.vxlanName != "" {
					mockCmd.SetOutput("ifconfig vxlan0 name "+tc.vxlanName, "")
				}
			}

			manager := NewManager(mockCmd)

			err := manager.CreateVXLAN(tc.vxlanName, tc.local, tc.remote, tc.group, tc.dev, tc.vxlanID)

			if tc.shouldError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			commands := mockCmd.GetCommands()
			if len(commands) < 3 {
				t.Errorf("expected at least 3 commands, got %d", len(commands))
			}
		})
	}
}

// Test helpers for Delete* and bridge interface management
func runDeleteTest(t *testing.T, mockCmd *MockCommandExecutor, deleteFunc func(string) error, ifName, expectedCmd string, shouldError bool) {
	t.Helper()
	err := deleteFunc(ifName)
	if shouldError {
		if err == nil {
			t.Error("expected error but got none")
		}
		return
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(mockCmd.GetCommands()) == 0 || mockCmd.GetCommands()[0] != expectedCmd {
		t.Errorf("expected command %s, got %v", expectedCmd, mockCmd.GetCommands())
	}
}

func runBridgeMemberTest(t *testing.T, mockCmd *MockCommandExecutor, bridgeFunc func(string, string) error, bridgeName, ifaceName, expectedCmd string, shouldError bool) {
	t.Helper()
	err := bridgeFunc(bridgeName, ifaceName)
	if shouldError {
		if err == nil {
			t.Error("expected error but got none")
		}
		return
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(mockCmd.GetCommands()) == 0 || mockCmd.GetCommands()[0] != expectedCmd {
		t.Errorf("expected command %s, got %v", expectedCmd, mockCmd.GetCommands())
	}
}

// Update Delete* tests to use new helpers
func TestBareOSManager_DeleteInterface(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		shouldError   bool
	}{
		{"successful interface deletion", "test0", false},
		{"empty interface name", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()
			manager := NewManager(mockCmd)
			runDeleteTest(t, mockCmd, manager.DeleteInterface, tc.interfaceName, "ifconfig "+tc.interfaceName+" destroy", tc.shouldError)
		})
	}
}

func TestBareOSManager_DeleteVLAN(t *testing.T) {
	tests := []struct {
		name        string
		vlanName    string
		shouldError bool
	}{
		{"successful VLAN deletion", "vlan100", false},
		{"empty VLAN name", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()
			manager := NewManager(mockCmd)
			runDeleteTest(t, mockCmd, manager.DeleteVLAN, tc.vlanName, "ifconfig "+tc.vlanName+" destroy", tc.shouldError)
		})
	}
}

func TestBareOSManager_DeleteGRE(t *testing.T) {
	tests := []struct {
		name        string
		greName     string
		shouldError bool
	}{
		{"successful GRE deletion", "gre0", false},
		{"empty GRE name", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()
			manager := NewManager(mockCmd)
			runDeleteTest(t, mockCmd, manager.DeleteGRE, tc.greName, "ifconfig "+tc.greName+" destroy", tc.shouldError)
		})
	}
}

func TestBareOSManager_DeleteVXLAN(t *testing.T) {
	tests := []struct {
		name        string
		vxlanName   string
		shouldError bool
	}{
		{"successful VXLAN deletion", "vxlan0", false},
		{"empty VXLAN name", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()
			manager := NewManager(mockCmd)
			runDeleteTest(t, mockCmd, manager.DeleteVXLAN, tc.vxlanName, "ifconfig "+tc.vxlanName+" destroy", tc.shouldError)
		})
	}
}

// Update Add/RemoveInterfaceToBridge tests to use new helpers
func TestBareOSManager_AddInterfaceToBridge(t *testing.T) {
	tests := []struct {
		name          string
		bridgeName    string
		interfaceName string
		shouldError   bool
	}{
		{"successful interface addition", "br0", "em0", false},
		{"empty bridge name", "", "em0", true},
		{"empty interface name", "br0", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()
			manager := NewManager(mockCmd)
			runBridgeMemberTest(t, mockCmd, manager.AddInterfaceToBridge, tc.bridgeName, tc.interfaceName, "ifconfig "+tc.bridgeName+" addm "+tc.interfaceName, tc.shouldError)
		})
	}
}

func TestBareOSManager_RemoveInterfaceFromBridge(t *testing.T) {
	tests := []struct {
		name          string
		bridgeName    string
		interfaceName string
		shouldError   bool
	}{
		{"successful interface removal", "br0", "em0", false},
		{"empty bridge name", "", "em0", true},
		{"empty interface name", "br0", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCmd := NewMockCommandExecutor()
			manager := NewManager(mockCmd)
			runBridgeMemberTest(t, mockCmd, manager.RemoveInterfaceFromBridge, tc.bridgeName, tc.interfaceName, "ifconfig "+tc.bridgeName+" deletem "+tc.interfaceName, tc.shouldError)
		})
	}
}

func TestAddIP_InvalidInput(t *testing.T) {
	t.Run("empty iface", func(t *testing.T) {
		err := AddIP("", "192.168.1.10", 24, defaultFamily)
		if err == nil {
			t.Error("expected error for empty iface")
		}
	})
	t.Run("empty ip", func(t *testing.T) {
		err := AddIP("em0", "", 24, defaultFamily)
		if err == nil {
			t.Error("expected error for empty ip")
		}
	})
	t.Run("zero mask", func(t *testing.T) {
		err := AddIP("em0", "192.168.1.10", 0, defaultFamily)
		if err == nil {
			t.Error("expected error for zero mask")
		}
	})
}

func TestAliasIP_InvalidInput(t *testing.T) {
	t.Run("empty iface", func(t *testing.T) {
		err := AliasIP("", "192.168.1.20", 24, defaultFamily)
		if err == nil {
			t.Error("expected error for empty iface")
		}
	})
	t.Run("empty ip", func(t *testing.T) {
		err := AliasIP("em0", "", 24, defaultFamily)
		if err == nil {
			t.Error("expected error for empty ip")
		}
	})
	t.Run("zero mask", func(t *testing.T) {
		err := AliasIP("em0", "192.168.1.20", 0, defaultFamily)
		if err == nil {
			t.Error("expected error for zero mask")
		}
	})
}

func TestDeleteIP_InvalidInput(t *testing.T) {
	t.Run("empty iface", func(t *testing.T) {
		err := DeleteIP("", "192.168.1.10", 24, defaultFamily)
		if err == nil {
			t.Error("expected error for empty iface")
		}
	})
	t.Run("empty ip", func(t *testing.T) {
		err := DeleteIP("em0", "", 24, defaultFamily)
		if err == nil {
			t.Error("expected error for empty ip")
		}
	})
	t.Run("zero mask", func(t *testing.T) {
		err := DeleteIP("em0", "192.168.1.10", 0, defaultFamily)
		if err == nil {
			t.Error("expected error for zero mask")
		}
	})
}

func TestAddRoute(t *testing.T) {
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	var calledArgs []string
	execCommand = func(name string, args ...string) *exec.Cmd {
		calledArgs = append([]string{name}, args...)
		return exec.Command("echo") // always succeeds
	}

	err := AddRoute(defaultFamily, "10.0.0.0/24", "10.0.0.1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(calledArgs) == 0 || calledArgs[0] != "route" {
		t.Errorf("expected route command, got %v", calledArgs)
	}
}

func TestAddRoute_Error(t *testing.T) {
	oldExecCommand := execCommand
	defer func() { execCommand = oldExecCommand }()

	execCommand = func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("false") // always fails
	}

	err := AddRoute(defaultFamily, "10.0.0.0/24", "10.0.0.1")
	if err == nil || !strings.Contains(err.Error(), "route add error") {
		t.Errorf("expected route add error, got %v", err)
	}
}

func TestDelRoute_LastDefault(t *testing.T) {
	oldListRoutes := listRoutes
	oldExecCommand := execCommand
	defer func() { listRoutes = oldListRoutes; execCommand = oldExecCommand }()

	listRoutes = func(family string) (string, error) {
		if family == defaultFamily {
			return "default 10.0.0.1\n", nil // only one default
		}
		if family == inet6Family {
			return "default fe80::1\n2001:db8::/64 fe80::1\n", nil
		}
		return "", nil
	}

	err := DelRoute(defaultFamily, "default")
	if err == nil || !strings.Contains(err.Error(), "cannot delete the last default route") {
		t.Errorf("expected last default route error, got %v", err)
	}
}

func TestDelRoute_Success(t *testing.T) {
	oldListRoutes := listRoutes
	oldExecCommand := execCommand
	defer func() { listRoutes = oldListRoutes; execCommand = oldExecCommand }()

	listRoutes = func(family string) (string, error) {
		if family == defaultFamily {
			return "default 10.0.0.1\ndefault 10.0.0.2\n", nil // two defaults
		}
		if family == inet6Family {
			return "default fe80::1\ndefault fe80::2\n", nil // two defaults
		}
		return "", nil
	}
	execCommand = func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("echo", "") // always succeeds
	}

	err := DelRoute(defaultFamily, "default")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestListAllRoutes(t *testing.T) {
	oldListRoutes := listRoutes
	defer func() { listRoutes = oldListRoutes }()

	listRoutes = func(family string) (string, error) {
		if family == defaultFamily {
			return "default 10.0.0.1\n10.0.0.0/24 10.0.0.1\n", nil
		}
		if family == inet6Family {
			return "default fe80::1\n2001:db8::/64 fe80::1\n", nil
		}
		return "", nil
	}

	out, err := ListAllRoutes("")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !strings.Contains(out, "IPv4 routes:") || !strings.Contains(out, "IPv6 routes:") {
		t.Errorf("expected both IPv4 and IPv6 routes, got %s", out)
	}
}

func TestListAllRoutes_Error(t *testing.T) {
	oldListRoutes := listRoutes
	defer func() { listRoutes = oldListRoutes }()

	listRoutes = func(_ string) (string, error) {
		return "", errors.New("fail")
	}

	_, err := ListAllRoutes(defaultFamily)
	if err == nil || !strings.Contains(err.Error(), "fail") {
		t.Errorf("expected error, got %v", err)
	}
}
