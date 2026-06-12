package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type clusterStatusRow struct {
	ID         string
	Type       string
	Name       string
	Node       string
	Status     string
	CGroupMode uint64
	Content    string
	CPU        float64
	Disk       uint64
	DiskRead   uint64
	DiskWrite  uint64
	HAState    string
	Level      string
	MaxCPU     uint64
	MaxDisk    uint64
	MaxMem     uint64
	Mem        uint64
	NetIn      uint64
	NetOut     uint64
	Pool       string
	Storage    string
	Tags       string
	Template   uint64
	Uptime     uint64
	VMID       uint64
}

func tableClusterStatus(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_cluster_status",
		Description: "Proxmox VE cluster resources overview (VMs, containers, nodes, storage).",
		List: &plugin.ListConfig{
			Hydrate: listClusterResources,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Resource ID."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Resource type (qemu, lxc, node, storage)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Resource name."},
			{Name: "node", Type: proto.ColumnType_STRING, Description: "Node the resource belongs to."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Resource status."},
			{Name: "cpu", Type: proto.ColumnType_DOUBLE, Transform: transform.FromField("CPU"), Description: "Current CPU usage."},
			{Name: "max_cpu", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxCPU"), Description: "Max CPU cores."},
			{Name: "mem", Type: proto.ColumnType_INT, Description: "Current memory usage in bytes."},
			{Name: "max_mem", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxMem"), Description: "Max memory in bytes."},
			{Name: "disk", Type: proto.ColumnType_INT, Description: "Current disk usage in bytes."},
			{Name: "max_disk", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxDisk"), Description: "Max disk size in bytes."},
			{Name: "disk_read", Type: proto.ColumnType_INT, Transform: transform.FromField("DiskRead"), Description: "Disk read bytes."},
			{Name: "disk_write", Type: proto.ColumnType_INT, Transform: transform.FromField("DiskWrite"), Description: "Disk write bytes."},
			{Name: "net_in", Type: proto.ColumnType_INT, Transform: transform.FromField("NetIn"), Description: "Network bytes received."},
			{Name: "net_out", Type: proto.ColumnType_INT, Transform: transform.FromField("NetOut"), Description: "Network bytes sent."},
			{Name: "uptime", Type: proto.ColumnType_INT, Description: "Uptime in seconds."},
			{Name: "ha_state", Type: proto.ColumnType_STRING, Transform: transform.FromField("HAState"), Description: "HA manager state."},
			{Name: "pool", Type: proto.ColumnType_STRING, Description: "Pool membership."},
			{Name: "storage", Type: proto.ColumnType_STRING, Description: "Storage name (for storage resources)."},
			{Name: "tags", Type: proto.ColumnType_STRING, Description: "Resource tags."},
			{Name: "template", Type: proto.ColumnType_INT, Description: "Whether this is a template (1) or not (0)."},
			{Name: "vmid", Type: proto.ColumnType_INT, Transform: transform.FromField("VMID"), Description: "VM/container ID."},
			{Name: "level", Type: proto.ColumnType_STRING, Description: "Support level."},
			{Name: "content", Type: proto.ColumnType_STRING, Description: "Storage content types."},
		},
	}
}

func listClusterResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_cluster_status.listClusterResources", "connection_error", err)
		return nil, err
	}

	cluster, err := client.Cluster(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_cluster_status.listClusterResources", "cluster_error", err)
		return nil, err
	}

	resources, err := cluster.Resources(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_cluster_status.listClusterResources", "api_error", err)
		return nil, err
	}

	for _, r := range resources {
		row := clusterStatusRow{
			ID:         r.ID,
			Type:       r.Type,
			Name:       r.Name,
			Node:       r.Node,
			Status:     r.Status,
			CGroupMode: r.CGroupMode,
			Content:    r.Content,
			CPU:        r.CPU,
			Disk:       r.Disk,
			DiskRead:   r.DiskRead,
			DiskWrite:  r.DiskWrite,
			HAState:    r.HAstate,
			Level:      r.Level,
			MaxCPU:     r.MaxCPU,
			MaxDisk:    r.MaxDisk,
			MaxMem:     r.MaxMem,
			Mem:        r.Mem,
			NetIn:      r.NetIn,
			NetOut:     r.NetOut,
			Pool:       r.Pool,
			Storage:    r.Storage,
			Tags:       r.Tags,
			Template:   r.Template,
			Uptime:     r.Uptime,
			VMID:       r.VMID,
		}
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}
