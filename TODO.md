# FreeBSD Management Tool Implementation Plan

- Create CLI commands for FreeBSD management

## Queue 1: Jail Implementation
- Create CLI commands for jail management (create/start/stop/destroy/list/info)
- Implement basic jail configuration (name, path, IP)
- Add Terraform output generator for jails (done)
- Develop jail status monitoring
- Write tests for jail operations


## Next
- Implement native network
    - VXLAN
    - VLAN
    - GRE
    - interface
    - bridge
    - other
- Advanced Network
    - VALE
    - VPP
    - OVS

- Implement ZFS management operations
    - snapshot management (Create/delete)

- Create VM commands