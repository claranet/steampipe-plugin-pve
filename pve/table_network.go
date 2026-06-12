package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type networkRow struct {
	Node      string
	Iface     string
	Type      string
	Active    int
	Autostart int
	CIDR      string
	CIDR6     string
	Gateway   string
	Gateway6  string
	Netmask   string
	Netmask6  string
	MTU       string
	Address   string
	Address6  string
	Method    string
	Method6   string
}

func tableNetwork(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_network",
		Description: "Proxmox VE network interfaces per node.",
		List: &plugin.ListConfig{
			Hydrate: listNetworks,
		},
		Columns: []*plugin.Column{
			{Name: "node", Type: proto.ColumnType_STRING, Description: "Node name."},
			{Name: "iface", Type: proto.ColumnType_STRING, Description: "Interface name."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Interface type (bridge, bond, eth, vlan, OVS*)."},
			{Name: "active", Type: proto.ColumnType_INT, Description: "Whether interface is active (1/0)."},
			{Name: "autostart", Type: proto.ColumnType_INT, Description: "Whether interface starts on boot (1/0)."},
			{Name: "cidr", Type: proto.ColumnType_STRING, Transform: transform.FromField("CIDR"), Description: "IPv4 CIDR address."},
			{Name: "cidr6", Type: proto.ColumnType_STRING, Transform: transform.FromField("CIDR6"), Description: "IPv6 CIDR address."},
			{Name: "gateway", Type: proto.ColumnType_STRING, Description: "IPv4 default gateway."},
			{Name: "gateway6", Type: proto.ColumnType_STRING, Description: "IPv6 default gateway."},
			{Name: "netmask", Type: proto.ColumnType_STRING, Description: "IPv4 subnet mask."},
			{Name: "netmask6", Type: proto.ColumnType_STRING, Description: "IPv6 prefix length."},
			{Name: "mtu", Type: proto.ColumnType_STRING, Transform: transform.FromField("MTU"), Description: "Maximum transmission unit."},
			{Name: "address", Type: proto.ColumnType_STRING, Description: "IPv4 address."},
			{Name: "address6", Type: proto.ColumnType_STRING, Description: "IPv6 address."},
			{Name: "method", Type: proto.ColumnType_STRING, Description: "IPv4 address method (static, dhcp, manual)."},
			{Name: "method6", Type: proto.ColumnType_STRING, Description: "IPv6 address method."},
		},
	}
}

func listNetworks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_network.listNetworks", "connection_error", err)
		return nil, err
	}

	nodes, err := client.Nodes(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_network.listNetworks", "nodes_error", err)
		return nil, err
	}

	for _, ns := range nodes {
		node, err := client.Node(ctx, ns.Node)
		if err != nil {
			plugin.Logger(ctx).Error("pve_network.listNetworks", "node_error", err, "node", ns.Node)
			continue
		}

		networks, err := node.Networks(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("pve_network.listNetworks", "network_error", err, "node", ns.Node)
			continue
		}

		for _, n := range networks {
			d.StreamListItem(ctx, networkRow{
				Node:      ns.Node,
				Iface:     n.Iface,
				Type:      n.Type,
				Active:    n.Active,
				Autostart: n.Autostart,
				CIDR:      n.CIDR,
				CIDR6:     n.CIDR6,
				Gateway:   n.Gateway,
				Gateway6:  n.Gateway6,
				Netmask:   n.Netmask,
				Netmask6:  n.Netmask6,
				MTU:       n.MTU,
				Address:   n.Address,
				Address6:  n.Address6,
				Method:    n.Method,
				Method6:   n.Method6,
			})
		}
	}

	return nil, nil
}
