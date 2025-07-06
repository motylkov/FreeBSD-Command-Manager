# FreeBSD Command Manager

A comprehensive Go CLI tool for managing FreeBSD infrastructure resources.

[![CI](https://github.com/motylkov/FreeBSD-Command-Manager/actions/workflows/ci.yml/badge.svg)](https://github.com/motylkov/FreeBSD-Command-Manager/actions/workflows/ci.yml)

## Features

- **Jail Management**
- **Network Management**

## Examples

### Jail Management

#### Complete Jail Workflow

```bash
# 1. Create a jail
./freebsd-cmd-manager jail create \
  --name web-server \
  --path /jails/web-server \
  --ip 192.168.1.10 \
  --mount /zfs/jails/web-server

# 2. Start the jail
./freebsd-cmd-manager jail start --name web-server

# 3. Check jail status
./freebsd-cmd-manager jail info --name web-server

# 4. List all jails
./freebsd-cmd-manager jail list

# 5. Stop the jail
./freebsd-cmd-manager jail stop --name web-server

# 6. Destroy the jail when no longer needed
./freebsd-cmd-manager jail destroy --name web-server
```

### Network Management

#### Basic Network Interface Operations

```bash
# List all network interfaces
./freebsd-cmd-manager network list

# Get detailed information about a specific interface
./freebsd-cmd-manager network info --name em0

# Create a generic network interface
./freebsd-cmd-manager network iface --name test0

# Delete a network interface
./freebsd-cmd-manager network delete-iface --name test0
```

#### Bridge Interface Management

```bash
# Create a bridge interface for connecting multiple networks
./freebsd-cmd-manager network bridge --name br0

# Create a bridge with specific configuration
./freebsd-cmd-manager network bridge --name br-lan

# Create a bridge with default name (bridge0, bridge1, etc.)
./freebsd-cmd-manager network bridge

# Delete a bridge interface
./freebsd-cmd-manager network delete-bridge --name br0
```

#### VLAN Configuration

```bash
# Create a VLAN interface on em0 with VLAN ID 100
./freebsd-cmd-manager network vlan \
  --name vlan100 \
  --parent em0 \
  --vlan-id 100

# Create a VLAN with default name (vlan0, vlan1, etc.)
./freebsd-cmd-manager network vlan \
  --parent em0 \
  --vlan-id 100

# Create multiple VLANs for network segmentation
./freebsd-cmd-manager network vlan --name vlan10 --parent em0 --vlan-id 10
./freebsd-cmd-manager network vlan --name vlan20 --parent em0 --vlan-id 20
./freebsd-cmd-manager network vlan --name vlan30 --parent em0 --vlan-id 30

# Delete a VLAN interface
./freebsd-cmd-manager network delete-vlan --name vlan100
```

#### GRE Tunnel Configuration

```bash
# Create a GRE tunnel for site-to-site VPN
./freebsd-cmd-manager network gre \
  --name gre0 \
  --local 192.168.1.1 \
  --remote 203.0.113.1

# Create a GRE tunnel with default name (gre0, gre1, etc.)
./freebsd-cmd-manager network gre \
  --local 192.168.1.1 \
  --remote 203.0.113.1

# Create multiple GRE tunnels
./freebsd-cmd-manager network gre --name gre-site1 --local 192.168.1.1 --remote 10.0.1.1
./freebsd-cmd-manager network gre --name gre-site2 --local 192.168.1.1 --remote 10.0.2.1

# Delete a GRE tunnel
./freebsd-cmd-manager network delete-gre --name gre0
```

#### VXLAN Tunnel Configuration

```bash
# Create a VXLAN tunnel for network virtualization
./freebsd-cmd-manager network vxlan \
  --name vxlan0 \
  --local 192.168.1.1 \
  --remote 192.168.1.2 \
  --vni 1000

# Create a VXLAN with default name (vxlan0, vxlan1, etc.)
./freebsd-cmd-manager network vxlan \
  --local 192.168.1.1 \
  --remote 192.168.1.2 \
  --vni 1000

# Create a VXLAN with multicast group and device specification
./freebsd-cmd-manager network vxlan \
  --name vxlan-multicast \
  --local 192.168.1.1 \
  --remote 192.168.1.2 \
  --vni 1000 \
  --group 224.0.0.1 \
  --dev em0

# Create multiple VXLAN tunnels for different network segments
./freebsd-cmd-manager network vxlan --name vxlan-tenant1 --local 192.168.1.1 --remote 192.168.1.2 --vni 1001
./freebsd-cmd-manager network vxlan --name vxlan-tenant2 --local 192.168.1.1 --remote 192.168.1.3 --vni 1002

# Delete a VXLAN tunnel
./freebsd-cmd-manager network delete-vxlan --name vxlan0
```

#### Complete Network Setup Example

```bash
# 1. Create a bridge for LAN connectivity
./freebsd-cmd-manager network bridge --name br-lan

# 2. Create VLANs for different network segments
./freebsd-cmd-manager network vlan --name vlan10 --parent em0 --vlan-id 10  # Management
./freebsd-cmd-manager network vlan --name vlan20 --parent em0 --vlan-id 20  # Production
./freebsd-cmd-manager network vlan --name vlan30 --parent em0 --vlan-id 30  # Development

# 3. Create GRE tunnels for remote site connectivity
./freebsd-cmd-manager network gre --name gre-datacenter --local 192.168.1.1 --remote 203.0.113.1
./freebsd-cmd-manager network gre --name gre-backup --local 192.168.1.1 --remote 198.51.100.1

# 4. Create VXLAN tunnels for network virtualization
./freebsd-cmd-manager network vxlan --name vxlan-tenant1 --local 192.168.1.1 --remote 192.168.1.2 --vni 1001
./freebsd-cmd-manager network vxlan --name vxlan-tenant2 --local 192.168.1.1 --remote 192.168.1.3 --vni 1002

# 5. List all interfaces to verify configuration
./freebsd-cmd-manager network list

# 6. Get detailed info about specific interfaces
./freebsd-cmd-manager network info --name br-lan
./freebsd-cmd-manager network info --name vlan10
./freebsd-cmd-manager network info --name gre-datacenter
./freebsd-cmd-manager network info --name vxlan-tenant1
```

#### Network Cleanup Example

```bash
# Remove all VXLAN tunnels
./freebsd-cmd-manager network delete-vxlan --name vxlan-tenant1
./freebsd-cmd-manager network delete-vxlan --name vxlan-tenant2

# Remove all GRE tunnels
./freebsd-cmd-manager network delete-gre --name gre-datacenter
./freebsd-cmd-manager network delete-gre --name gre-backup

# Remove all VLANs
./freebsd-cmd-manager network delete-vlan --name vlan10
./freebsd-cmd-manager network delete-vlan --name vlan20
./freebsd-cmd-manager network delete-vlan --name vlan30

# Remove bridge
./freebsd-cmd-manager network delete-bridge --name br-lan

# Verify cleanup
./freebsd-cmd-manager network list
```

## Installation

## Configuration

## Development

### Prerequisites

