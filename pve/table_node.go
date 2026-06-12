package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableNode(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_node",
		Description: "Proxmox VE cluster nodes with status and resource usage.",
		List: &plugin.ListConfig{
			Hydrate: listNodes,
		},
		Columns: []*plugin.Column{
			{Name: "node", Type: proto.ColumnType_STRING, Description: "Node name."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Node status (online/offline)."},
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Node ID (e.g. node/pve1)."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Resource type."},
			{Name: "cpu", Type: proto.ColumnType_DOUBLE, Transform: transform.FromField("CPU"), Description: "Current CPU usage (fraction 0.0 - 1.0)."},
			{Name: "max_cpu", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxCPU"), Description: "Number of CPU cores."},
			{Name: "mem", Type: proto.ColumnType_INT, Description: "Current memory usage in bytes."},
			{Name: "max_mem", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxMem"), Description: "Total memory in bytes."},
			{Name: "disk", Type: proto.ColumnType_INT, Description: "Current disk usage in bytes."},
			{Name: "max_disk", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxDisk"), Description: "Total disk size in bytes."},
			{Name: "uptime", Type: proto.ColumnType_INT, Description: "Node uptime in seconds."},
			{Name: "ssl_fingerprint", Type: proto.ColumnType_STRING, Transform: transform.FromField("SSLFingerprint"), Description: "SSL certificate fingerprint."},
			{Name: "level", Type: proto.ColumnType_STRING, Description: "Support level."},
		},
	}
}

func listNodes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_node.listNodes", "connection_error", err)
		return nil, err
	}

	nodes, err := client.Nodes(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_node.listNodes", "api_error", err)
		return nil, err
	}

	for _, node := range nodes {
		d.StreamListItem(ctx, node)
	}

	return nil, nil
}
