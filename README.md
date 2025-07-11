# FreeBSD Command Manager

A comprehensive Go CLI tool for managing FreeBSD infrastructure resources.

[![CI](https://github.com/motylkov/FreeBSD-Command-Manager/actions/workflows/ci.yml/badge.svg)](https://github.com/motylkov/FreeBSD-Command-Manager/actions/workflows/ci.yml) 

## Features

- [Jail Management](#jail-management)
- Network Management
  - OS level
    - [Interfaces Management](#network-management)
    - [IP Address Management](#ip-address-management) 
    - [Route Management](#route-management)


## Examples

### Jail Management

#### Complete Jail Workflow

```bash
# 1. Create a jail
./fcom jail create \
  --name web-server \
  --path /jails/web-server \
  --ip 192.168.1.10 \
  --mount /zfs/jails/web-server

# 2. Start the jail
./fcom jail start --name web-server

# 3. Check jail status
./fcom jail info --name web-server

# 4. List all jails
./fcom jail list

# 5. Stop the jail
./fcom jail stop --name web-server

# 6. Destroy the jail when no longer needed
./fcom jail destroy --name web-server
```

### Network Management

#### Basic Network Interface Operations

```bash
# List all network interfaces
./fcom network list

# Get detailed information about a specific interface
./fcom network info --name em0

# Create a generic network interface
./fcom network iface --name test0

# Delete a network interface
./fcom network delete-iface --name test0
```

#### Bridge Interface Management

```bash
# Create a bridge interface for connecting multiple networks
./fcom network bridge --name br0

# Create a bridge with specific configuration
./fcom network bridge --name br-lan

# Create a bridge with default name (bridge0, bridge1, etc.)
./fcom network bridge

# Delete a bridge interface
./fcom network delete-bridge --name br0
```

#### VLAN Configuration

```bash
# Create a VLAN interface on em0 with VLAN ID 100
./fcom network vlan \
  --name vlan100 \
  --parent em0 \
  --vlan-id 100

# Create a VLAN with default name (vlan0, vlan1, etc.)
./fcom network vlan \
  --parent em0 \
  --vlan-id 100

# Create multiple VLANs for network segmentation
./fcom network vlan --name vlan10 --parent em0 --vlan-id 10
./fcom network vlan --name vlan20 --parent em0 --vlan-id 20
./fcom network vlan --name vlan30 --parent em0 --vlan-id 30

# Delete a VLAN interface
./fcom network delete-vlan --name vlan100
```

#### GRE Tunnel Configuration

```bash
# Create a GRE tunnel for site-to-site VPN
./fcom network gre \
  --name gre0 \
  --local 192.168.1.1 \
  --remote 203.0.113.1

# Create a GRE tunnel with default name (gre0, gre1, etc.)
./fcom network gre \
  --local 192.168.1.1 \
  --remote 203.0.113.1

# Create multiple GRE tunnels
./fcom network gre --name gre-site1 --local 192.168.1.1 --remote 10.0.1.1
./fcom network gre --name gre-site2 --local 192.168.1.1 --remote 10.0.2.1

# Delete a GRE tunnel
./fcom network delete-gre --name gre0
```

#### VXLAN Tunnel Configuration

```bash
# Create a VXLAN tunnel for network virtualization
./fcom network vxlan \
  --name vxlan0 \
  --local 192.168.1.1 \
  --remote 192.168.1.2 \
  --vni 1000

# Create a VXLAN with default name (vxlan0, vxlan1, etc.)
./fcom network vxlan \
  --local 192.168.1.1 \
  --remote 192.168.1.2 \
  --vni 1000

# Create a VXLAN with multicast group and device specification
./fcom network vxlan \
  --name vxlan-multicast \
  --local 192.168.1.1 \
  --remote 192.168.1.2 \
  --vni 1000 \
  --group 224.0.0.1 \
  --dev em0

# Create multiple VXLAN tunnels for different network segments
./fcom network vxlan --name vxlan-tenant1 --local 192.168.1.1 --remote 192.168.1.2 --vni 1001
./fcom network vxlan --name vxlan-tenant2 --local 192.168.1.1 --remote 192.168.1.3 --vni 1002

# Delete a VXLAN tunnel
./fcom network delete-vxlan --name vxlan0
```

#### Complete Network Setup Example

```bash
# 1. Create a bridge for LAN connectivity
./fcom network bridge --name br-lan

# 2. Create VLANs for different network segments
./fcom network vlan --name vlan10 --parent em0 --vlan-id 10  # Management
./fcom network vlan --name vlan20 --parent em0 --vlan-id 20  # Production
./fcom network vlan --name vlan30 --parent em0 --vlan-id 30  # Development

# 3. Create GRE tunnels for remote site connectivity
./fcom network gre --name gre-datacenter --local 192.168.1.1 --remote 203.0.113.1
./fcom network gre --name gre-backup --local 192.168.1.1 --remote 198.51.100.1

# 4. Create VXLAN tunnels for network virtualization
./fcom network vxlan --name vxlan-tenant1 --local 192.168.1.1 --remote 192.168.1.2 --vni 1001
./fcom network vxlan --name vxlan-tenant2 --local 192.168.1.1 --remote 192.168.1.3 --vni 1002

# 5. List all interfaces to verify configuration
./fcom network list

# 6. Get detailed info about specific interfaces
./fcom network info --name br-lan
./fcom network info --name vlan10
./fcom network info --name gre-datacenter
./fcom network info --name vxlan-tenant1
```

#### Network Cleanup Example

```bash
# Remove all VXLAN tunnels
./fcom network delete-vxlan --name vxlan-tenant1
./fcom network delete-vxlan --name vxlan-tenant2

# Remove all GRE tunnels
./fcom network delete-gre --name gre-datacenter
./fcom network delete-gre --name gre-backup

# Remove all VLANs
./fcom network delete-vlan --name vlan10
./fcom network delete-vlan --name vlan20
./fcom network delete-vlan --name vlan30

# Remove bridge
./fcom network delete-bridge --name br-lan

# Verify cleanup
./fcom network list
```

### IP Address Management

```bash
# Add IPv4 address to interface
./fcom ip add --iface em0 --ip 192.168.1.10 --mask 24 --family inet

# Add IPv6 address to interface
./fcom ip add --iface em0 --ip 2001:db8::1 --mask 64 --family inet6

# Add alias IPv4 address
./fcom ip alias --iface em0 --ip 192.168.1.20 --mask 24 --family inet

# Delete IPv6 address from interface
./fcom ip delete --iface em0 --ip 2001:db8::1 --mask 64 --family inet6
```

### Route Management

Easily manage IPv4 and IPv6 routes, including adding, deleting, and listing routes. The CLI ensures you cannot accidentally remove the last default route, protecting system connectivity.

#### Add a Route

```bash
# Add an IPv4 network route
./fcom network route add --family inet --net 10.0.0.0/24 --gw 10.0.0.1

# Add an IPv6 network route
./fcom network route add --family inet6 --net 2001:db8::/64 --gw 2001:db8::1

# Add an IPv4 network route via a specific interface
./fcom network route add --family inet --net 10.0.0.0/24 --gw 10.0.0.1 --iface em0

# Add an IPv6 network route via a specific interface
./fcom network route add --family inet6 --net 2001:db8::/64 --gw 2001:db8::1 --iface em1

# Add a default IPv4 route
./fcom network route add --family inet --net default --gw 192.168.1.1

# Add a host route (single IP)
./fcom network route add --family inet --net 192.168.1.50/32 --gw 10.0.0.1
```

#### Delete a Route

```bash
# Delete an IPv4 network route
./fcom network route del --family inet --net 10.0.0.0/24

# Delete an IPv6 network route
./fcom network route del --family inet6 --net 2001:db8::/64

# Delete a default route (if more than one exists)
./fcom network route del --family inet --net default
```

#### List Routes

```bash
# List all IPv4 routes
./fcom network route list --family inet

# List all IPv6 routes
./fcom network route list --family inet6

# List all routes (both families)
./fcom network route list
```

#### Example: Safe Default Route Handling

```bash
# Attempting to delete the last default route will fail:
./fcom network route del --family inet --net default
# Output: cannot delete the last default route
```

#### Troubleshooting
- Ensure you specify the correct `--family` (either `inet` for IPv4 or `inet6` for IPv6).
- The `--net` flag accepts both network prefixes (e.g., `10.0.0.0/24`) and `default`.
- You cannot delete the last default route for a family; this is a safety feature.

### Version Information

```bash
# Show version, commit, and build date
./fcom version
# Example output:
# 0.02.abcdef.2025-07-06
```

## Contributing

Contributions are welcome! To contribute:

1. **Fork** the repository and create your branch from `main`.
2. **Write clear, well-tested code** and ensure it passes all linters and tests:
   ```sh
   golangci-lint run
   go test ./...
   ```
3. **Document** any new features or changes in the README if needed.
4. **Open a Pull Request** with a clear description of your changes and why they are needed.

### Development Environment

- Use Go 1.22 or newer (see `.golangci.yml` and CI workflow for the current version).
- Use `golangci-lint` v2.x for linting.
- Run `go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...` before submitting.

### Code Style

- Follow Go best practices and idioms.
- Keep functions short and focused.
- Add comments for exported functions and types.
- Use descriptive commit messages.

If you have any questions, feel free to open an issue or discussion!

![FreeBSD](https://img.shields.io/badge/-FreeBSD-%23870000?style=for-the-badge&logo=freebsd&logoColor=white) ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)