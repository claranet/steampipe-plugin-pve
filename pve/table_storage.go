package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type storageRow struct {
	Node         string
	Storage      string
	Type         string
	Content      string
	Enabled      int
	Active       int
	Shared       int
	Total        uint64
	Used         uint64
	Avail        uint64
	UsedFraction float64
}

func tableStorage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_storage",
		Description: "Proxmox VE storage configurations per node.",
		List: &plugin.ListConfig{
			Hydrate: listStorages,
		},
		Columns: []*plugin.Column{
			{Name: "node", Type: proto.ColumnType_STRING, Description: "Node name."},
			{Name: "storage", Type: proto.ColumnType_STRING, Description: "Storage name."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Storage type (dir, lvm, zfspool, nfs, cifs, etc)."},
			{Name: "content", Type: proto.ColumnType_STRING, Description: "Allowed content types (images, rootdir, iso, backup, etc)."},
			{Name: "enabled", Type: proto.ColumnType_INT, Description: "Whether storage is enabled (1/0)."},
			{Name: "active", Type: proto.ColumnType_INT, Description: "Whether storage is active (1/0)."},
			{Name: "shared", Type: proto.ColumnType_INT, Description: "Whether storage is shared across nodes (1/0)."},
			{Name: "total", Type: proto.ColumnType_INT, Description: "Total storage capacity in bytes."},
			{Name: "used", Type: proto.ColumnType_INT, Description: "Used storage in bytes."},
			{Name: "avail", Type: proto.ColumnType_INT, Description: "Available storage in bytes."},
			{Name: "used_fraction", Type: proto.ColumnType_DOUBLE, Transform: transform.FromField("UsedFraction"), Description: "Used storage as a fraction (0.0 - 1.0)."},
		},
	}
}

func listStorages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_storage.listStorages", "connection_error", err)
		return nil, err
	}

	nodes, err := client.Nodes(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_storage.listStorages", "nodes_error", err)
		return nil, err
	}

	for _, ns := range nodes {
		node, err := client.Node(ctx, ns.Node)
		if err != nil {
			plugin.Logger(ctx).Error("pve_storage.listStorages", "node_error", err, "node", ns.Node)
			continue
		}

		storages, err := node.Storages(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("pve_storage.listStorages", "storage_error", err, "node", ns.Node)
			continue
		}

		for _, s := range storages {
			d.StreamListItem(ctx, storageRow{
				Node:         ns.Node,
				Storage:      s.Storage,
				Type:         s.Type,
				Content:      s.Content,
				Enabled:      s.Enabled,
				Active:       s.Active,
				Shared:       s.Shared,
				Total:        s.Total,
				Used:         s.Used,
				Avail:        s.Avail,
				UsedFraction: s.UsedFraction,
			})
		}
	}

	return nil, nil
}
