package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type replicationRow struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	Guest    int    `json:"guest"`
	Schedule string `json:"schedule"`
	Rate     int    `json:"rate"`
	Disable  int    `json:"disable"`
	Comment  string `json:"comment"`
	Digest   string `json:"digest"`
}

func tableReplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_replication",
		Description: "Proxmox VE storage replication jobs.",
		List: &plugin.ListConfig{
			Hydrate: listReplications,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Replication job ID."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Replication type (local)."},
			{Name: "source", Type: proto.ColumnType_STRING, Description: "Source node."},
			{Name: "target", Type: proto.ColumnType_STRING, Description: "Target node."},
			{Name: "guest", Type: proto.ColumnType_INT, Description: "Guest VMID being replicated."},
			{Name: "schedule", Type: proto.ColumnType_STRING, Description: "Replication schedule (e.g. */15 = every 15 min)."},
			{Name: "rate", Type: proto.ColumnType_INT, Description: "Rate limit in MiB/s."},
			{Name: "disable", Type: proto.ColumnType_INT, Description: "Whether replication is disabled (1/0)."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "Job comment."},
			{Name: "digest", Type: proto.ColumnType_STRING, Description: "Config digest."},
		},
	}
}

func listReplications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_replication.listReplications", "connection_error", err)
		return nil, err
	}

	var replications []replicationRow
	err = client.Get(ctx, "/cluster/replication", &replications)
	if err != nil {
		plugin.Logger(ctx).Error("pve_replication.listReplications", "api_error", err)
		return nil, err
	}

	for _, r := range replications {
		d.StreamListItem(ctx, r)
	}

	return nil, nil
}
