package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-pve",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		DefaultTransform: transform.FromGo(),
		SchemaMode:       plugin.SchemaModeDynamic,
		TableMapFunc:     pluginTableMap,
	}
	return p
}

func pluginTableMap(ctx context.Context, d *plugin.TableMapData) (map[string]*plugin.Table, error) {
	return map[string]*plugin.Table{
		"pve_node":           tableNode(ctx),
		"pve_cluster_status": tableClusterStatus(ctx),
		"pve_vm":             tableVM(ctx),
		"pve_container":      tableContainer(ctx),
		"pve_vm_snapshot":    tableVMSnapshot(ctx),
		"pve_storage":        tableStorage(ctx),
		"pve_network":        tableNetwork(ctx),
		"pve_user":           tableUser(ctx),
		"pve_role":           tableRole(ctx),
		"pve_acl":            tableACL(ctx),
		"pve_pool":           tablePool(ctx),
		"pve_task":           tableTask(ctx),
		"pve_backup":         tableBackup(ctx),
		"pve_firewall_rule":  tableFirewallRule(ctx),
		"pve_ha_resource":    tableHAResource(ctx),
		"pve_replication":    tableReplication(ctx),
	}, nil
}
