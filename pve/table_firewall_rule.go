package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type firewallRuleRow struct {
	Pos      int
	Type     string
	Action   string
	Comment  string
	Dest     string
	Dport    string
	Enable   int
	IcmpType string
	Iface    string
	Log      string
	Macro    string
	Proto    string
	Source   string
	Sport    string
}

func tableFirewallRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_firewall_rule",
		Description: "Proxmox VE cluster-level firewall rules.",
		List: &plugin.ListConfig{
			Hydrate: listFirewallRules,
		},
		Columns: []*plugin.Column{
			{Name: "pos", Type: proto.ColumnType_INT, Description: "Rule position."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Rule type (in, out, group)."},
			{Name: "action", Type: proto.ColumnType_STRING, Description: "Rule action (ACCEPT, DROP, REJECT)."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "Rule comment."},
			{Name: "dest", Type: proto.ColumnType_STRING, Description: "Destination address/alias."},
			{Name: "dport", Type: proto.ColumnType_STRING, Description: "Destination port(s)."},
			{Name: "enable", Type: proto.ColumnType_INT, Description: "Whether rule is enabled (1/0)."},
			{Name: "icmp_type", Type: proto.ColumnType_STRING, Description: "ICMP type filter."},
			{Name: "iface", Type: proto.ColumnType_STRING, Description: "Network interface filter."},
			{Name: "log", Type: proto.ColumnType_STRING, Description: "Log level (emerg, alert, crit, err, warning, notice, info, debug, nolog)."},
			{Name: "macro", Type: proto.ColumnType_STRING, Description: "Macro name (e.g. SSH, HTTP, DNS)."},
			{Name: "proto", Type: proto.ColumnType_STRING, Description: "IP protocol (tcp, udp, icmp, etc)."},
			{Name: "source", Type: proto.ColumnType_STRING, Description: "Source address/alias."},
			{Name: "sport", Type: proto.ColumnType_STRING, Description: "Source port(s)."},
		},
	}
}

func listFirewallRules(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_firewall_rule.listFirewallRules", "connection_error", err)
		return nil, err
	}

	var rules []firewallRuleRow
	err = client.Get(ctx, "/cluster/firewall/rules", &rules)
	if err != nil {
		plugin.Logger(ctx).Error("pve_firewall_rule.listFirewallRules", "api_error", err)
		return nil, err
	}

	for _, rule := range rules {
		d.StreamListItem(ctx, rule)
	}

	return nil, nil
}
