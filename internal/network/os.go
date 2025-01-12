package network

type BareOSManager struct{}

func NewBareOSManager() *BareOSManager {
	return &BareOSManager{}
}

// iface
func (m *BareOSManager) CreateInterface(ifType, name string) error {
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
