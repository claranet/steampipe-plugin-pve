package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type haResourceRow struct {
	SID        string `json:"sid"`
	Type       string `json:"type"`
	State      string `json:"state"`
	Status     string `json:"status"`
	Group      string `json:"group"`
	MaxRestart int    `json:"max_restart"`
	MaxRelocate int   `json:"max_relocate"`
	Comment    string `json:"comment"`
	Digest     string `json:"digest"`
}

func tableHAResource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_ha_resource",
		Description: "Proxmox VE High Availability managed resources.",
		List: &plugin.ListConfig{
			Hydrate: listHAResources,
		},
		Columns: []*plugin.Column{
			{Name: "sid", Type: proto.ColumnType_STRING, Transform: transform.FromField("SID"), Description: "HA resource SID (e.g. vm:100, ct:101)."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Resource type (vm, ct)."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "Requested HA state (started, stopped, disabled, ignored)."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Current HA status."},
			{Name: "group", Type: proto.ColumnType_STRING, Description: "HA group name."},
			{Name: "max_restart", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxRestart"), Description: "Max restart attempts."},
			{Name: "max_relocate", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxRelocate"), Description: "Max relocate attempts."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "Resource comment."},
			{Name: "digest", Type: proto.ColumnType_STRING, Description: "Config digest."},
		},
	}
}

func listHAResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_ha_resource.listHAResources", "connection_error", err)
		return nil, err
	}

	var resources []haResourceRow
	err = client.Get(ctx, "/cluster/ha/resources", &resources)
	if err != nil {
		plugin.Logger(ctx).Error("pve_ha_resource.listHAResources", "api_error", err)
		return nil, err
	}

	for _, r := range resources {
		d.StreamListItem(ctx, r)
	}

	return nil, nil
}
