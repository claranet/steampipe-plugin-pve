package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type aclRow struct {
	Path      string
	RoleID    string
	Type      string
	UGID      string
	Propagate bool
}

func tableACL(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_acl",
		Description: "Proxmox VE access control lists.",
		List: &plugin.ListConfig{
			Hydrate: listACLs,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "ACL path (e.g. /, /vms/100, /storage/local)."},
			{Name: "role_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("RoleID"), Description: "Role ID assigned."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "ACL subject type (user, group, token)."},
			{Name: "ugid", Type: proto.ColumnType_STRING, Transform: transform.FromField("UGID"), Description: "User/Group/Token ID."},
			{Name: "propagate", Type: proto.ColumnType_BOOL, Description: "Whether ACL propagates to child paths."},
		},
	}
}

func listACLs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_acl.listACLs", "connection_error", err)
		return nil, err
	}

	acls, err := client.ACL(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_acl.listACLs", "api_error", err)
		return nil, err
	}

	for _, a := range acls {
		d.StreamListItem(ctx, aclRow{
			Path:      a.Path,
			RoleID:    a.RoleID,
			Type:      a.Type,
			UGID:      a.UGID,
			Propagate: bool(a.Propagate),
		})
	}

	return nil, nil
}
