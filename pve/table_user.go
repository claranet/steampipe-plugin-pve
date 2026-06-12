package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type userRow struct {
	UserID    string
	Comment   string
	Email     string
	Enable    bool
	Expire    int
	Firstname string
	Lastname  string
	Groups    []string
	Keys      string
	RealmType string
}

func tableUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_user",
		Description: "Proxmox VE users.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
		},
		Columns: []*plugin.Column{
			{Name: "user_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("UserID"), Description: "User ID (e.g. root@pam)."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "User comment."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "User email address."},
			{Name: "enable", Type: proto.ColumnType_BOOL, Description: "Whether user is enabled."},
			{Name: "expire", Type: proto.ColumnType_INT, Description: "Account expiration date (Unix epoch, 0 = never)."},
			{Name: "firstname", Type: proto.ColumnType_STRING, Description: "First name."},
			{Name: "lastname", Type: proto.ColumnType_STRING, Description: "Last name."},
			{Name: "groups", Type: proto.ColumnType_JSON, Description: "Groups the user belongs to."},
			{Name: "keys", Type: proto.ColumnType_STRING, Description: "User keys."},
			{Name: "realm_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("RealmType"), Description: "Realm type (pam, pve, ldap, ad)."},
		},
	}
}

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_user.listUsers", "connection_error", err)
		return nil, err
	}

	users, err := client.Users(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_user.listUsers", "api_error", err)
		return nil, err
	}

	for _, u := range users {
		d.StreamListItem(ctx, userRow{
			UserID:    u.UserID,
			Comment:   u.Comment,
			Email:     u.Email,
			Enable:    bool(u.Enable),
			Expire:    u.Expire,
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
			Groups:    u.Groups,
			Keys:      u.Keys,
			RealmType: u.RealmType,
		})
	}

	return nil, nil
}
