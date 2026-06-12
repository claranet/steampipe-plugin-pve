package pve

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type snapshotRow struct {
	Node        string
	VMID        uint64
	VMName      string
	Name        string
	Description string
	Snaptime    time.Time
	Parent      string
	VMState     int
}

func tableVMSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_vm_snapshot",
		Description: "Proxmox VE VM and container snapshots.",
		List: &plugin.ListConfig{
			Hydrate: listVMSnapshots,
		},
		Columns: []*plugin.Column{
			{Name: "node", Type: proto.ColumnType_STRING, Description: "Node name."},
			{Name: "vmid", Type: proto.ColumnType_INT, Transform: transform.FromField("VMID"), Description: "VM/Container ID."},
			{Name: "vm_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("VMName"), Description: "VM/Container name."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Snapshot name."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Snapshot description."},
			{Name: "snaptime", Type: proto.ColumnType_TIMESTAMP, Description: "Snapshot creation time."},
			{Name: "parent", Type: proto.ColumnType_STRING, Description: "Parent snapshot name."},
			{Name: "vm_state", Type: proto.ColumnType_INT, Transform: transform.FromField("VMState"), Description: "Whether VM state is included (1/0)."},
		},
	}
}

func listVMSnapshots(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_vm_snapshot.listVMSnapshots", "connection_error", err)
		return nil, err
	}

	nodes, err := client.Nodes(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_vm_snapshot.listVMSnapshots", "nodes_error", err)
		return nil, err
	}

	for _, ns := range nodes {
		node, err := client.Node(ctx, ns.Node)
		if err != nil {
			plugin.Logger(ctx).Error("pve_vm_snapshot.listVMSnapshots", "node_error", err, "node", ns.Node)
			continue
		}

		// VM snapshots via raw API (go-proxmox doesn't expose a typed Snapshots() on VirtualMachine)
		vms, err := node.VirtualMachines(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("pve_vm_snapshot.listVMSnapshots", "vm_error", err, "node", ns.Node)
			continue
		}
		for _, vm := range vms {
			vmid := uint64(vm.VMID)
			var snapshots []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Snaptime    int64  `json:"snaptime"`
				Parent      string `json:"parent"`
				VMState     int    `json:"vmstate"`
			}
			err := client.Get(ctx, fmt.Sprintf("/nodes/%s/qemu/%d/snapshot", ns.Node, vmid), &snapshots)
			if err != nil {
				plugin.Logger(ctx).Error("pve_vm_snapshot.listVMSnapshots", "snapshot_error", err, "vmid", vmid)
				continue
			}
			for _, snap := range snapshots {
				if snap.Name == "current" {
					continue
				}
				d.StreamListItem(ctx, snapshotRow{
					Node:        ns.Node,
					VMID:        vmid,
					VMName:      vm.Name,
					Name:        snap.Name,
					Description: snap.Description,
					Snaptime:    time.Unix(snap.Snaptime, 0),
					Parent:      snap.Parent,
					VMState:     snap.VMState,
				})
			}
		}

		// Container snapshots
		containers, err := node.Containers(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("pve_vm_snapshot.listVMSnapshots", "container_error", err, "node", ns.Node)
			continue
		}
		for _, ct := range containers {
			snaps, err := ct.Snapshots(ctx)
			if err != nil {
				plugin.Logger(ctx).Error("pve_vm_snapshot.listVMSnapshots", "ct_snapshot_error", err, "vmid", uint64(ct.VMID))
				continue
			}
			for _, snap := range snaps {
				if strings.Compare(snap.Name, "current") == 1 || len(snap.Name) == 0 {
					continue
				}
				d.StreamListItem(ctx, snapshotRow{
					Node:        ns.Node,
					VMID:        uint64(ct.VMID),
					VMName:      ct.Name,
					Name:        snap.Name,
					Description: snap.Description,
					Snaptime:    time.Unix(snap.SnapshotCreationTime, 0),
					Parent:      snap.Parent,
				})
			}
		}
	}

	return nil, nil
}
