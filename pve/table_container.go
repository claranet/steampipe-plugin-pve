package pve

import (
	"context"

	"github.com/luthermonson/go-proxmox"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type containerRow struct {
	VMID    uint64
	Name    string
	Node    string
	Status  string
	Tags    string
	CPUs    int
	MaxDisk uint64
	MaxMem  uint64
	MaxSwap uint64
	Uptime  uint64
	// Config fields
	Hostname     string
	Description  string
	Arch         string
	OSType       string
	Cores        int
	Memory       int
	Swap         int
	RootFS       string
	Nameserver   string
	SearchDomain string
	OnBoot       bool
	Protection   bool
	Unprivileged bool
	Features     string
	Timezone     string
	// Merged config
	Networks    map[string]string
	MountPoints map[string]string
}

func tableContainer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_container",
		Description: "Proxmox VE LXC containers with configuration and status.",
		List: &plugin.ListConfig{
			Hydrate: listContainers,
		},
		Columns: []*plugin.Column{
			{Name: "vmid", Type: proto.ColumnType_INT, Transform: transform.FromField("VMID"), Description: "Container ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Container name."},
			{Name: "node", Type: proto.ColumnType_STRING, Description: "Node running the container."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Container status (running, stopped)."},
			{Name: "tags", Type: proto.ColumnType_STRING, Description: "Semicolon-separated tags."},
			{Name: "cpus", Type: proto.ColumnType_INT, Description: "Allocated CPUs."},
			{Name: "max_disk", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxDisk"), Description: "Max root disk size in bytes."},
			{Name: "max_mem", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxMem"), Description: "Max memory in bytes."},
			{Name: "max_swap", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxSwap"), Description: "Max swap in bytes."},
			{Name: "uptime", Type: proto.ColumnType_INT, Description: "Uptime in seconds."},
			// Config
			{Name: "hostname", Type: proto.ColumnType_STRING, Description: "Container hostname."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Container description."},
			{Name: "arch", Type: proto.ColumnType_STRING, Description: "Architecture (amd64, arm64, etc)."},
			{Name: "os_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("OSType"), Description: "OS type."},
			{Name: "cores", Type: proto.ColumnType_INT, Description: "CPU cores limit."},
			{Name: "memory", Type: proto.ColumnType_INT, Description: "Memory in MiB."},
			{Name: "swap", Type: proto.ColumnType_INT, Description: "Swap in MiB."},
			{Name: "rootfs", Type: proto.ColumnType_STRING, Transform: transform.FromField("RootFS"), Description: "Root filesystem config."},
			{Name: "nameserver", Type: proto.ColumnType_STRING, Description: "DNS nameserver."},
			{Name: "search_domain", Type: proto.ColumnType_STRING, Transform: transform.FromField("SearchDomain"), Description: "DNS search domain."},
			{Name: "on_boot", Type: proto.ColumnType_BOOL, Transform: transform.FromField("OnBoot"), Description: "Start on boot."},
			{Name: "protection", Type: proto.ColumnType_BOOL, Description: "Protection from removal."},
			{Name: "unprivileged", Type: proto.ColumnType_BOOL, Description: "Unprivileged container flag."},
			{Name: "features", Type: proto.ColumnType_STRING, Description: "Container features."},
			{Name: "timezone", Type: proto.ColumnType_STRING, Description: "Container timezone."},
			{Name: "networks", Type: proto.ColumnType_JSON, Description: "Network interfaces (net0, net1, ...)."},
			{Name: "mount_points", Type: proto.ColumnType_JSON, Transform: transform.FromField("MountPoints"), Description: "Mount points (mp0, mp1, ...)."},
		},
	}
}

func listContainers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_container.listContainers", "connection_error", err)
		return nil, err
	}

	nodes, err := client.Nodes(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_container.listContainers", "nodes_error", err)
		return nil, err
	}

	for _, ns := range nodes {
		node, err := client.Node(ctx, ns.Node)
		if err != nil {
			plugin.Logger(ctx).Error("pve_container.listContainers", "node_error", err, "node", ns.Node)
			continue
		}

		containers, err := node.Containers(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("pve_container.listContainers", "container_error", err, "node", ns.Node)
			continue
		}

		for _, ct := range containers {
			row := containerRowFromCT(ct)
			d.StreamListItem(ctx, row)
		}
	}

	return nil, nil
}

func containerRowFromCT(ct *proxmox.Container) containerRow {
	row := containerRow{
		VMID:    uint64(ct.VMID),
		Name:    ct.Name,
		Node:    ct.Node,
		Status:  ct.Status,
		Tags:    ct.Tags,
		CPUs:    ct.CPUs,
		MaxDisk: ct.MaxDisk,
		MaxMem:  ct.MaxMem,
		MaxSwap: ct.MaxSwap,
		Uptime:  ct.Uptime,
	}

	if ct.ContainerConfig != nil {
		cfg := ct.ContainerConfig
		row.Hostname = cfg.Hostname
		row.Description = cfg.Description
		row.Arch = cfg.Arch
		row.OSType = cfg.OSType
		row.Cores = cfg.Cores
		row.Memory = cfg.Memory
		row.Swap = cfg.Swap
		row.RootFS = cfg.RootFS
		row.Nameserver = cfg.Nameserver
		row.SearchDomain = cfg.SearchDomain
		row.OnBoot = bool(cfg.OnBoot)
		row.Protection = bool(cfg.Protection)
		row.Unprivileged = bool(cfg.Unprivileged)
		row.Features = cfg.Features
		row.Timezone = cfg.Timezone
		row.Networks = cfg.MergeNets()
		row.MountPoints = cfg.MergeMps()
	}

	return row
}
