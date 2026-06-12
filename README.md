# Steampipe Plugin for Proxmox VE

Query your [Proxmox VE](https://www.proxmox.com/) infrastructure using SQL with [Steampipe](https://steampipe.io).

- **Author:** Martin Weber (martin.weber@claranet.com)
- **Version:** 0.1.0
- **License:** MIT

## Quick Start

1. Build and install the plugin:

   ```sh
   make install
   ```

2. Configure a connection in `~/.steampipe/config/pve.spc`:

   ```hcl
   connection "pve" {
     plugin           = "local/pve"
     api_url          = "https://your-pve-host:8006/api2/json"
     api_token_id     = "user@realm!tokenname"
     api_token_secret = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
     tls_insecure     = true
   }
   ```

3. Start querying:

   ```sh
   steampipe query
   ```

   ```sql
   select name, node, status, cpus, mem, max_mem
   from pve_vm
   where status = 'running';
   ```

## Configuration Reference

| Field | Type | Required | Description |
|---|---|---|---|
| `api_url` | string | Yes | Full API base URL (e.g. `https://host:8006/api2/json`) |
| `api_token_id` | string | No | API token ID (`user@realm!tokenname`) |
| `api_token_secret` | string | No | API token secret (UUID) |
| `user` | string | No | Username for password auth (`user@realm`) |
| `password` | string | No | Password for password auth |
| `tls_insecure` | bool | No | Skip TLS certificate verification (default `false`) |

Either `api_token_id` + `api_token_secret` or `user` + `password` must be provided.

## Authentication

### API Token (recommended)

Create an API token in the web UI under Datacenter → Permissions → API Tokens, then configure:

```hcl
connection "pve" {
  plugin           = "local/pve"
  api_url          = "https://pve.example.com:8006/api2/json"
  api_token_id     = "monitor@pve!steampipe"
  api_token_secret = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

### Username / Password

```hcl
connection "pve" {
  plugin   = "local/pve"
  api_url  = "https://pve.example.com:8006/api2/json"
  user     = "root@pam"
  password = "your-password"
}
```

## Multi-Cluster (Aggregator)

Query multiple PVE clusters as a single connection using Steampipe's aggregator:

```hcl
connection "pve_dc1" {
  plugin           = "local/pve"
  api_url          = "https://cluster-a:8006/api2/json"
  api_token_id     = "monitor@pve!steampipe"
  api_token_secret = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

connection "pve_dc2" {
  plugin           = "local/pve"
  api_url          = "https://cluster-b:8006/api2/json"
  api_token_id     = "monitor@pve!steampipe"
  api_token_secret = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

connection "pve_all" {
  plugin      = "local/pve"
  type        = "aggregator"
  connections = ["pve_dc1", "pve_dc2"]
}
```

Then query across both clusters:

```sql
select _ctx ->> 'connection_name' as cluster, name, node, status
from pve_all.pve_vm
where status = 'running';
```

## Tables

16 tables covering the full Proxmox VE API surface:

| Table | Description |
|---|---|
| `pve_node` | Cluster nodes with status and resource usage |
| `pve_cluster_status` | Cluster-wide resource overview |
| `pve_vm` | QEMU virtual machines with full configuration |
| `pve_container` | LXC containers with configuration |
| `pve_vm_snapshot` | VM and container snapshots |
| `pve_storage` | Storage configurations per node |
| `pve_network` | Network interfaces per node |
| `pve_user` | Users |
| `pve_role` | Access control roles and privileges |
| `pve_acl` | Access control list entries |
| `pve_pool` | Resource pools with members |
| `pve_task` | Cluster tasks (recent operations) |
| `pve_backup` | Scheduled backup jobs |
| `pve_firewall_rule` | Cluster-level firewall rules |
| `pve_ha_resource` | HA managed resources |
| `pve_replication` | Storage replication jobs |

See [docs/tables.md](docs/tables.md) for column descriptions and example queries.

## Example Queries

**Running VMs with resource usage:**

```sql
select name, node, status, cpus,
  mem / 1048576 as mem_mb,
  max_mem / 1048576 as max_mem_mb
from pve_vm
where status = 'running'
order by node, name;
```

**Storage utilization over 80%:**

```sql
select node, storage, type,
  round(used_fraction * 100, 1) as used_pct,
  avail / 1073741824 as avail_gb
from pve_storage
where used_fraction > 0.8
order by used_fraction desc;
```

**Recent failed tasks:**

```sql
select type, user, node, exit_status, start_time
from pve_task
where is_failed = true
order by start_time desc
limit 10;
```

**Multi-cluster VM count:**

```sql
select _ctx ->> 'connection_name' as cluster,
  count(*) as total,
  sum(case when status = 'running' then 1 else 0 end) as running
from pve_all.pve_vm
group by 1;
```
## Installation

```bash
mkdir -p ~/.steampipe/plugins/local/pve
curl -L -o ~/.steampipe/plugins/local/pve/steampipe-plugin-pve.plugin \
  https://github.com/claranet/steampipe-plugin-pve/releases/download/v0.1.0/steampipe-plugin-pve_v0.1.0.darwin-arm64.plugin
chmod +x ~/.steampipe/plugins/local/pve/steampipe-plugin-pve.plugin
```

## Building

Requirements: Go 1.23+, Steampipe

```sh
make build     # compile plugin binary
make install   # build and install to ~/.steampipe/plugins/local/pve/
make clean     # remove compiled binary
make test      # run tests
```
