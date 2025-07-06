package ifconfig

import (
	"fmt"
	"math/bits"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	Ethernet   string = "ethernet"
	Wireless   string = "wireless"
	Loopback   string = "loopback"
	Bridge     string = "bridge"
	VLAN       string = "vlan"
	VXLAN      string = "vxlan"
	Tunnel     string = "tunnel"
	PPP        string = "ppp"
	LAGG       string = "lagg"
	GIF        string = "gif"
	GRE        string = "gre"
	Tap        string = "tap"
	Stf        string = "stf"
	Enc        string = "enc"
	Unknown    string = "unknown"
	StatusUp          = "up"
	StatusDown        = "down"

	// Array indices for parsing
	MinFieldsForIP      = 2
	IPAddressIndex      = 1
	NetmaskFieldIndex   = 2
	NetmaskValueIndex   = 3
	PrefixLenFieldIndex = 2
	PrefixLenValueIndex = 3

	// Bit masks for netmask processing
	ByteMask = 0xFF
)

// Info represents information about a network interface
type Info struct {
	Name   string
	Type   string
	Status string
	IPv4   []string
	IPv6   []string
	MAC    string
}

// Parses FreeBSD ifconfig output into a slice of Info structs.
func ParseIfconfig(output string) []Info {
	var result []Info
	lines := strings.Split(output, "\n")

	var currentInfo *Info
	var flags string
	var media string
	var groups string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if isNewInterfaceLine(line) {
			if currentInfo != nil {
				currentInfo.Type = determineInterfaceType(flags, media, groups)
				currentInfo.Status = determineInterfaceStatus(flags)
				result = append(result, *currentInfo)
			}
			currentInfo = &Info{Name: extractInterfaceName(line)}
			flags, media, groups = extractFlagsMediaGroups(line)
			continue
		}

		if currentInfo == nil {
			continue
		}

		flags, media, groups = updateFlagsMediaGroups(line, flags, media, groups)
		parseMAC(line, currentInfo)
		parseIPv4(line, currentInfo)
		parseIPv6(line, currentInfo)
	}

	if currentInfo != nil {
		currentInfo.Type = determineInterfaceType(flags, media, groups)
		currentInfo.Status = determineInterfaceStatus(flags)
		result = append(result, *currentInfo)
	}

	return result
}

func isNewInterfaceLine(line string) bool {
	return regexp.MustCompile(`^([a-z0-9]+):\s+flags=`).MatchString(line)
}

func extractInterfaceName(line string) string {
	matches := regexp.MustCompile(`^([a-z0-9]+):\s+flags=`).FindStringSubmatch(line)
	if matches != nil {
		return matches[1]
	}
	return ""
}

func extractFlagsMediaGroups(line string) (flags, media, groups string) {
	if strings.Contains(line, "flags=") {
		flags = line
	}
	return
}

func updateFlagsMediaGroups(line, flags, media, groups string) (newFlags, newMedia, newGroups string) {
	if flags == "" && strings.Contains(line, "flags=") {
		flags = line
	}
	if strings.Contains(line, "media:") {
		media = line
	}
	if strings.Contains(line, "groups:") {
		groups = line
	}
	return flags, media, groups
}

func parseMAC(line string, currentInfo *Info) {
	if strings.Contains(line, "ether ") {
		parts := strings.Fields(line)
		if len(parts) >= MinFieldsForIP {
			if mac, err := net.ParseMAC(parts[IPAddressIndex]); err == nil {
				currentInfo.MAC = mac.String()
			}
		}
	}
}

func parseIPv4(line string, currentInfo *Info) {
	if strings.Contains(line, "inet ") && !strings.Contains(line, "inet6") {
		parts := strings.Fields(line)
		if len(parts) >= MinFieldsForIP {
			ip := net.ParseIP(parts[IPAddressIndex])
			if ip != nil && ip.To4() != nil {
				cidr := parts[IPAddressIndex]
				if len(parts) >= 3 && strings.Contains(parts[NetmaskFieldIndex], "netmask") {
					netmaskStr := parts[NetmaskValueIndex]
					if netmaskStr != "" && strings.HasPrefix(netmaskStr, "0x") {
						if cidrBits := hexNetmaskToCIDR(netmaskStr); cidrBits > 0 {
							cidr = parts[IPAddressIndex] + "/" + fmt.Sprintf("%d", cidrBits)
						}
					}
				}
				currentInfo.IPv4 = append(currentInfo.IPv4, cidr)
			}
		}
	}
}

func parseIPv6(line string, currentInfo *Info) {
	if strings.Contains(line, "inet6 ") {
		parts := strings.Fields(line)
		if len(parts) >= MinFieldsForIP {
			ipStr := parts[IPAddressIndex]
			if idx := strings.Index(ipStr, "%"); idx != -1 {
				ipStr = ipStr[:idx]
			}
			ip := net.ParseIP(ipStr)
			if ip != nil && ip.To16() != nil && ip.To4() == nil {
				cidr := ipStr
				if len(parts) >= 3 && strings.Contains(parts[PrefixLenFieldIndex], "prefixlen") {
					prefixLen := parts[PrefixLenValueIndex]
					if prefixLen != "" {
						cidr = ipStr + "/" + prefixLen
					}
				}
				currentInfo.IPv6 = append(currentInfo.IPv6, cidr)
			}
		}
	}
}

func determineInterfaceType(flags, media, groups string) string {
	if isLoopback(flags) {
		return Loopback
	}
	if isBridge(media) {
		return Bridge
	}
	if isVLAN(media, groups) {
		return VLAN
	}
	if isVXLAN(media) {
		return VXLAN
	}
	if isLAGG(media) {
		return LAGG
	}
	if isPPP(flags) {
		return PPP
	}
	if isWireless(media) {
		return Wireless
	}
	if isEthernet(media) {
		return Ethernet
	}
	if isTunnel(flags) {
		return Tunnel
	}
	if isGIF(flags) {
		return GIF
	}
	if isGRE(media) {
		return GRE
	}
	if isTap(flags) {
		return Tap
	}
	if isStf(flags) {
		return Stf
	}
	if isEnc(flags) {
		return Enc
	}
	return Unknown
}

func isLoopback(flags string) bool {
	return strings.Contains(flags, "LOOPBACK")
}

func isBridge(media string) bool {
	return strings.Contains(media, "bridge members:")
}

func isVLAN(media, groups string) bool {
	return strings.Contains(media, "vlan:") || strings.Contains(groups, "vlan")
}

func isVXLAN(media string) bool {
	return strings.Contains(media, "vxlan:")
}

func isLAGG(media string) bool {
	return strings.Contains(media, "laggproto:")
}

func isPPP(flags string) bool {
	return strings.Contains(flags, "POINTOPOINT") && strings.Contains(flags, "MULTICAST")
}

func isWireless(media string) bool {
	return strings.Contains(media, "IEEE 802.11")
}

func isEthernet(media string) bool {
	return strings.Contains(media, "Ethernet") || strings.Contains(media, "autoselect")
}

func isTunnel(flags string) bool {
	return strings.Contains(flags, "TUNNEL")
}

func isGIF(flags string) bool {
	return strings.Contains(flags, "SIMPLEX") && strings.Contains(flags, "MULTICAST")
}

func isGRE(media string) bool {
	return strings.Contains(media, "gre")
}

func isTap(flags string) bool {
	return strings.Contains(flags, "TAP")
}

func isStf(flags string) bool {
	return strings.Contains(flags, "stf")
}

func isEnc(flags string) bool {
	return strings.Contains(flags, "ENCAP")
}

func determineInterfaceStatus(flags string) string {
	if strings.Contains(flags, "RUNNING") {
		return StatusUp
	}
	return StatusDown
}

func hexNetmaskToCIDR(netmaskStr string) int {
	netmask, err := strconv.ParseUint(netmaskStr, 0, 64)
	if err != nil {
		return 0
	}

	cidrBits := 0
	for netmask > 0 {
		cidrBits += bits.OnesCount64(netmask & ByteMask)
		netmask >>= 8
	}

	return cidrBits
}
