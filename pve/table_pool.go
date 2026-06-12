package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type poolRow struct {
	PoolID  string
	Comment string
	Members []poolMemberRow
}

type poolMemberRow struct {
	ID     string  `json:"id"`
	Type   string  `json:"type"`
	Node   string  `json:"node"`
	VMID   uint64  `json:"vmid"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
	CPU    float64 `json:"cpu"`
	Mem    uint64  `json:"mem"`
	MaxMem uint64  `json:"max_mem"`
}

func tablePool(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_pool",
		Description: "Proxmox VE resource pools.",
		List: &plugin.ListConfig{
			Hydrate: listPools,
		},
		Columns: []*plugin.Column{
			{Name: "pool_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("PoolID"), Description: "Pool ID."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "Pool description."},
			{Name: "members", Type: proto.ColumnType_JSON, Description: "Pool members (VMs, containers, storage)."},
		},
	}
}

func listPools(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_pool.listPools", "connection_error", err)
		return nil, err
	}

	pools, err := client.Pools(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_pool.listPools", "api_error", err)
		return nil, err
	}

	for _, p := range pools {
		// Fetch full pool details to get members
		pool, err := client.Pool(ctx, p.PoolID)
		if err != nil {
			plugin.Logger(ctx).Error("pve_pool.listPools", "pool_detail_error", err, "pool", p.PoolID)
			// Stream without members
			d.StreamListItem(ctx, poolRow{
				PoolID:  p.PoolID,
				Comment: p.Comment,
			})
			continue
		}

		members := make([]poolMemberRow, 0, len(pool.Members))
		for _, m := range pool.Members {
			members = append(members, poolMemberRow{
				ID:     m.ID,
				Type:   m.Type,
				Node:   m.Node,
				VMID:   m.VMID,
				Name:   m.Name,
				Status: m.Status,
				CPU:    m.CPU,
				Mem:    m.Mem,
				MaxMem: m.MaxMem,
			})
		}

		d.StreamListItem(ctx, poolRow{
			PoolID:  pool.PoolID,
			Comment: pool.Comment,
			Members: members,
		})
	}

	return nil, nil
}
