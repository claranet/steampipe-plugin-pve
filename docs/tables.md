# Tables

## Proxmox VE Tables

### Cluster & Nodes

| Table | Description |
|-------|-------------|
| `pve_node` | Cluster nodes with status and resource usage |
| `pve_cluster_status` | Cluster-wide resource overview (VMs, containers, nodes, storage) |

### Virtual Machines & Containers

| Table | Description |
|-------|-------------|
| `pve_vm` | QEMU virtual machines with config (CPU, memory, disks, network) |
| `pve_container` | LXC containers with config (hostname, memory, mount points) |
| `pve_vm_snapshot` | VM and container snapshots |

### Storage & Network

| Table | Description |
|-------|-------------|
| `pve_storage` | Storage configurations per node |
| `pve_network` | Network interfaces per node |

### Access Control

| Table | Description |
|-------|-------------|
| `pve_user` | Users |
| `pve_role` | Access control roles and their privileges |
| `pve_acl` | Access control list entries |
| `pve_pool` | Resource pools with their members |

### Operations

| Table | Description |
|-------|-------------|
| `pve_task` | Cluster tasks (recent operations) |
| `pve_backup` | Scheduled backup jobs |
| `pve_firewall_rule` | Cluster-level firewall rules |
| `pve_ha_resource` | HA managed resources |
| `pve_replication` | Storage replication jobs |

---

## Example Queries

### List all running VMs across all nodes

```sql
select vmid, name, node, cpus, max_mem / 1048576 as memory_mb, status
from pve_vm
where status = 'running'
order by node, name;
```

### Find VMs with more than 4 CPU cores

```sql
select vmid, name, node, cores, sockets, cores * sockets as total_cores
from pve_vm
where cores > 4;
```

### Storage usage above 80%

```sql
select node, storage, type,
  round(used_fraction * 100, 1) as used_pct,
  avail / 1073741824 as avail_gb
from pve_storage
where used_fraction > 0.8
order by used_fraction desc;
```

### List all containers and their network config

```sql
select vmid, name, node, hostname, status, networks
from pve_container
order by node, vmid;
```

### Show cluster resource overview

```sql
select type, count(*) as count,
  sum(case when status = 'running' then 1 else 0 end) as running
from pve_cluster_status
where type in ('qemu', 'lxc')
group by type;
```

### Find users with no expiry

```sql
select user_id, email, enable, groups
from pve_user
where expire = 0;
```

### List recent failed tasks

```sql
select upid, type, user, node, exit_status, start_time
from pve_task
where is_failed = true
order by start_time desc
limit 10;
```

### Check HA resource status

```sql
select sid, state, status, "group", max_restart
from pve_ha_resource;
```

### Backup jobs targeting specific storage

```sql
select id, schedule, storage, vmid, compress, mode
from pve_backup
where enabled = 1;
```

### Multi-cluster VM count (using aggregator connection)

```sql
select _ctx ->> 'connection_name' as cluster,
  count(*) as vm_count,
  sum(case when status = 'running' then 1 else 0 end) as running
from pve_all.pve_vm
group by 1;
```
