package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type backupJobRow struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	Enabled        int    `json:"enabled"`
	Schedule       string `json:"schedule"`
	Storage        string `json:"storage"`
	Node           string `json:"node"`
	Pool           string `json:"pool"`
	VMID           string `json:"vmid"`
	Compress       string `json:"compress"`
	Mode           string `json:"mode"`
	MailTo         string `json:"mailto"`
	Mailnotification string `json:"mailnotification"`
	PruneBackups   string `json:"prune-backups"`
	NotesTemplate  string `json:"notes-template"`
	All            int    `json:"all"`
	Exclude        string `json:"exclude"`
	Comment        string `json:"comment"`
}

func tableBackup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_backup",
		Description: "Proxmox VE scheduled backup jobs.",
		List: &plugin.ListConfig{
			Hydrate: listBackupJobs,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Backup job ID."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Backup type."},
			{Name: "enabled", Type: proto.ColumnType_INT, Description: "Whether backup job is enabled (1/0)."},
			{Name: "schedule", Type: proto.ColumnType_STRING, Description: "Schedule (e.g. daily, weekly, or cron expression)."},
			{Name: "storage", Type: proto.ColumnType_STRING, Description: "Target storage for backups."},
			{Name: "node", Type: proto.ColumnType_STRING, Description: "Node filter (empty = all nodes)."},
			{Name: "pool", Type: proto.ColumnType_STRING, Description: "Pool filter."},
			{Name: "vmid", Type: proto.ColumnType_STRING, Transform: transform.FromField("VMID"), Description: "VM/CT IDs to back up (comma-separated)."},
			{Name: "compress", Type: proto.ColumnType_STRING, Description: "Compression type (0, gzip, lzo, zstd)."},
			{Name: "mode", Type: proto.ColumnType_STRING, Description: "Backup mode (snapshot, suspend, stop)."},
			{Name: "mail_to", Type: proto.ColumnType_STRING, Transform: transform.FromField("MailTo"), Description: "Mail recipients."},
			{Name: "prune_backups", Type: proto.ColumnType_STRING, Transform: transform.FromField("PruneBackups"), Description: "Prune/retention settings."},
			{Name: "notes_template", Type: proto.ColumnType_STRING, Transform: transform.FromField("NotesTemplate"), Description: "Notes template."},
			{Name: "all", Type: proto.ColumnType_INT, Description: "Backup all VMs (1/0)."},
			{Name: "exclude", Type: proto.ColumnType_STRING, Description: "Excluded VM IDs."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "Job comment."},
		},
	}
}

func listBackupJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_backup.listBackupJobs", "connection_error", err)
		return nil, err
	}

	var jobs []backupJobRow
	err = client.Get(ctx, "/cluster/backup", &jobs)
	if err != nil {
		plugin.Logger(ctx).Error("pve_backup.listBackupJobs", "api_error", err)
		return nil, err
	}

	for _, job := range jobs {
		d.StreamListItem(ctx, job)
	}

	return nil, nil
}
