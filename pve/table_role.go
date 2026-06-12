package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type roleRow struct {
	RoleID  string
	Privs   string
	Special bool
}

func tableRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_role",
		Description: "Proxmox VE access control roles.",
		List: &plugin.ListConfig{
			Hydrate: listRoles,
		},
		Columns: []*plugin.Column{
			{Name: "role_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("RoleID"), Description: "Role ID."},
			{Name: "privs", Type: proto.ColumnType_STRING, Description: "Comma-separated list of privileges."},
			{Name: "special", Type: proto.ColumnType_BOOL, Description: "Whether this is a built-in role."},
		},
	}
}

func listRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_role.listRoles", "connection_error", err)
		return nil, err
	}

	roles, err := client.Roles(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_role.listRoles", "api_error", err)
		return nil, err
	}

	for _, r := range roles {
		d.StreamListItem(ctx, roleRow{
			RoleID:  r.RoleID,
			Privs:   r.Privs,
			Special: bool(r.Special),
		})
	}

	return nil, nil
}
