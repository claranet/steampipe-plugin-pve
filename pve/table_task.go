package pve

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type taskRow struct {
	UPID         string
	ID           string
	Type         string
	User         string
	Status       string
	Node         string
	PID          uint64
	PStart       uint64
	ExitStatus   string
	IsCompleted  bool
	IsRunning    bool
	IsFailed     bool
	IsSuccessful bool
	StartTime    string
	EndTime      string
}

func tableTask(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_task",
		Description: "Proxmox VE cluster tasks (recent operations).",
		List: &plugin.ListConfig{
			Hydrate: listTasks,
		},
		Columns: []*plugin.Column{
			{Name: "upid", Type: proto.ColumnType_STRING, Transform: transform.FromField("UPID"), Description: "Unique process ID."},
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Task ID."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Task type (qmstart, vzdump, etc)."},
			{Name: "user", Type: proto.ColumnType_STRING, Description: "User who initiated the task."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Task status."},
			{Name: "node", Type: proto.ColumnType_STRING, Description: "Node the task ran on."},
			{Name: "pid", Type: proto.ColumnType_INT, Transform: transform.FromField("PID"), Description: "Process ID."},
			{Name: "pstart", Type: proto.ColumnType_INT, Transform: transform.FromField("PStart"), Description: "Process start counter."},
			{Name: "exit_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("ExitStatus"), Description: "Task exit status."},
			{Name: "is_completed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsCompleted"), Description: "Whether task has completed."},
			{Name: "is_running", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsRunning"), Description: "Whether task is currently running."},
			{Name: "is_failed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsFailed"), Description: "Whether task has failed."},
			{Name: "is_successful", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsSuccessful"), Description: "Whether task completed successfully."},
			{Name: "start_time", Type: proto.ColumnType_STRING, Transform: transform.FromField("StartTime"), Description: "Task start time."},
			{Name: "end_time", Type: proto.ColumnType_STRING, Transform: transform.FromField("EndTime"), Description: "Task end time."},
		},
	}
}

func listTasks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_task.listTasks", "connection_error", err)
		return nil, err
	}

	cluster, err := client.Cluster(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_task.listTasks", "cluster_error", err)
		return nil, err
	}

	tasks, err := cluster.Tasks(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_task.listTasks", "api_error", err)
		return nil, err
	}

	for _, t := range tasks {
		startTime := ""
		if !t.StartTime.IsZero() {
			startTime = t.StartTime.Format("2006-01-02T15:04:05Z07:00")
		}
		endTime := ""
		if !t.EndTime.IsZero() {
			endTime = t.EndTime.Format("2006-01-02T15:04:05Z07:00")
		}

		d.StreamListItem(ctx, taskRow{
			UPID:         string(t.UPID),
			ID:           t.ID,
			Type:         t.Type,
			User:         t.User,
			Status:       t.Status,
			Node:         t.Node,
			PID:          t.PID,
			PStart:       t.PStart,
			ExitStatus:   t.ExitStatus,
			IsCompleted:  t.IsCompleted,
			IsRunning:    t.IsRunning,
			IsFailed:     t.IsFailed,
			IsSuccessful: t.IsSuccessful,
			StartTime:    startTime,
			EndTime:      endTime,
		})
	}

	return nil, nil
}
