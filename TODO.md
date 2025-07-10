# FreeBSD Management Tool Implementation Plan

- Create CLI commands for FreeBSD management


## Queue 1: Jail Implementation

- Implement basic jail configuration
    - Create CLI commands for jail management (create/start/stop/destroy) (done)
    - Jail templates

- Add Terraform output generator (done for testing)


## Queue 2: Basic OS Networking

- Implement native network
    - VXLAN (done)
    - VLAN (done)
    - GRE (done)
    - interface (done)
    - bridge (done)
    - ip (done)
    - other

- Add linter (done)


## Queue 3: ZFS Management 

- Implement ZFS management operations
    - snapshot management (Create/delete)


## Queue 4: Bhyve VM Implementation

- Create VM commands


## Queue 5: Advanced Networking
- Advanced Network
    - VALE
    - VPP
    - OVS
