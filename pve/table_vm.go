package pve

import (
	"context"

	"github.com/luthermonson/go-proxmox"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type vmRow struct {
	// Core fields
	VMID      uint64
	Name      string
	Node      string
	Status    string
	Lock      string
	Tags      string
	Template  bool
	PID       uint64
	QMPStatus string
	// Resource usage
	CPU       float64
	CPUs      int
	Mem       uint64
	MaxMem    uint64
	Disk      uint64
	MaxDisk   uint64
	DiskRead  uint64
	DiskWrite uint64
	NetIn     uint64
	NetOut    uint64
	Uptime    uint64
	// Config (hydrated)
	Description string
	OSType      string
	Machine     string
	Bios        string
	Boot        string
	Cores       int
	Sockets     int
	CPUType     string
	Memory      int
	Balloon     int
	// Note: VirtualMachineConfig.Memory is proxmox.StringOrInt - we cast to int
	SCSIHW      string
	VGA         string
	OnBoot      int
	Agent       string
	Hotplug     string
	Numa        int
	Protection  int
	// Network and disks as merged maps
	Networks map[string]string
	Disks    map[string]string
}

func tableVM(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pve_vm",
		Description: "Proxmox VE QEMU virtual machines with configuration and status.",
		List: &plugin.ListConfig{
			Hydrate: listVMs,
		},
		Columns: []*plugin.Column{
			// Identity
			{Name: "vmid", Type: proto.ColumnType_INT, Transform: transform.FromField("VMID"), Description: "VM ID."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "VM name."},
			{Name: "node", Type: proto.ColumnType_STRING, Description: "Node running the VM."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "VM status (running, stopped, etc)."},
			{Name: "lock", Type: proto.ColumnType_STRING, Description: "Lock type if VM is locked."},
			{Name: "tags", Type: proto.ColumnType_STRING, Description: "Semicolon-separated tags."},
			{Name: "template", Type: proto.ColumnType_BOOL, Description: "True if this is a template."},
			{Name: "pid", Type: proto.ColumnType_INT, Transform: transform.FromField("PID"), Description: "QEMU process ID."},
			{Name: "qmp_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("QMPStatus"), Description: "QMP status."},
			// Resource usage
			{Name: "cpu", Type: proto.ColumnType_DOUBLE, Transform: transform.FromField("CPU"), Description: "Current CPU usage (fraction)."},
			{Name: "cpus", Type: proto.ColumnType_INT, Description: "Number of allocated CPUs."},
			{Name: "mem", Type: proto.ColumnType_INT, Description: "Current memory usage in bytes."},
			{Name: "max_mem", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxMem"), Description: "Max memory in bytes."},
			{Name: "disk", Type: proto.ColumnType_INT, Description: "Current disk usage in bytes."},
			{Name: "max_disk", Type: proto.ColumnType_INT, Transform: transform.FromField("MaxDisk"), Description: "Max disk size in bytes."},
			{Name: "disk_read", Type: proto.ColumnType_INT, Transform: transform.FromField("DiskRead"), Description: "Disk read bytes."},
			{Name: "disk_write", Type: proto.ColumnType_INT, Transform: transform.FromField("DiskWrite"), Description: "Disk write bytes."},
			{Name: "net_in", Type: proto.ColumnType_INT, Transform: transform.FromField("NetIn"), Description: "Network bytes received."},
			{Name: "net_out", Type: proto.ColumnType_INT, Transform: transform.FromField("NetOut"), Description: "Network bytes sent."},
			{Name: "uptime", Type: proto.ColumnType_INT, Description: "Uptime in seconds."},
			// Config
			{Name: "description", Type: proto.ColumnType_STRING, Description: "VM description."},
			{Name: "os_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("OSType"), Description: "OS type (l26, win10, etc)."},
			{Name: "machine", Type: proto.ColumnType_STRING, Description: "Machine type."},
			{Name: "bios", Type: proto.ColumnType_STRING, Description: "BIOS type (seabios, ovmf)."},
			{Name: "boot", Type: proto.ColumnType_STRING, Description: "Boot order configuration."},
			{Name: "cores", Type: proto.ColumnType_INT, Description: "Number of CPU cores."},
			{Name: "sockets", Type: proto.ColumnType_INT, Description: "Number of CPU sockets."},
			{Name: "cpu_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("CPUType"), Description: "CPU type emulation."},
			{Name: "memory", Type: proto.ColumnType_INT, Description: "Memory in MiB."},
			{Name: "balloon", Type: proto.ColumnType_INT, Description: "Balloon memory target in MiB."},
			{Name: "scsihw", Type: proto.ColumnType_STRING, Transform: transform.FromField("SCSIHW"), Description: "SCSI controller type."},
			{Name: "vga", Type: proto.ColumnType_STRING, Transform: transform.FromField("VGA"), Description: "VGA type."},
			{Name: "on_boot", Type: proto.ColumnType_INT, Transform: transform.FromField("OnBoot"), Description: "Start on boot (1/0)."},
			{Name: "agent", Type: proto.ColumnType_STRING, Description: "QEMU guest agent config."},
			{Name: "hotplug", Type: proto.ColumnType_STRING, Description: "Hotplug devices."},
			{Name: "numa", Type: proto.ColumnType_INT, Description: "NUMA enabled (1/0)."},
			{Name: "protection", Type: proto.ColumnType_INT, Description: "Protection from removal (1/0)."},
			// Merged config
			{Name: "networks", Type: proto.ColumnType_JSON, Description: "Network interfaces (net0, net1, ...)."},
			{Name: "disks", Type: proto.ColumnType_JSON, Description: "Disk devices (scsi0, virtio0, ide0, ...)."},
		},
	}
}

func listVMs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("pve_vm.listVMs", "connection_error", err)
		return nil, err
	}

	nodes, err := client.Nodes(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("pve_vm.listVMs", "nodes_error", err)
		return nil, err
	}

	for _, ns := range nodes {
		node, err := client.Node(ctx, ns.Node)
		if err != nil {
			plugin.Logger(ctx).Error("pve_vm.listVMs", "node_error", err, "node", ns.Node)
			continue
		}

		vms, err := node.VirtualMachines(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("pve_vm.listVMs", "vm_error", err, "node", ns.Node)
			continue
		}

		for _, vm := range vms {
			row := vmRowFromVM(vm)
			d.StreamListItem(ctx, row)
		}
	}

	return nil, nil
}

func vmRowFromVM(vm *proxmox.VirtualMachine) vmRow {
	row := vmRow{
		VMID:      uint64(vm.VMID),
		Name:      vm.Name,
		Node:      vm.Node,
		Status:    vm.Status,
		Lock:      vm.Lock,
		Tags:      vm.Tags,
		Template:  bool(vm.Template),
		PID:       uint64(vm.PID),
		QMPStatus: vm.QMPStatus,
		CPU:       vm.CPU,
		CPUs:      vm.CPUs,
		Mem:       vm.Mem,
		MaxMem:    vm.MaxMem,
		Disk:      vm.Disk,
		MaxDisk:   vm.MaxDisk,
		DiskRead:  vm.DiskRead,
		DiskWrite: vm.DiskWrite,
		NetIn:     vm.NetIn,
		NetOut:    vm.Netout,
		Uptime:    vm.Uptime,
	}

	if vm.VirtualMachineConfig != nil {
		cfg := vm.VirtualMachineConfig
		row.Description = cfg.Description
		row.OSType = cfg.OSType
		row.Machine = cfg.Machine
		row.Bios = cfg.Bios
		row.Boot = cfg.Boot
		row.Cores = cfg.Cores
		row.Sockets = cfg.Sockets
		row.CPUType = cfg.CPU
		row.Memory = int(cfg.Memory)
		row.Balloon = cfg.Balloon
		row.SCSIHW = cfg.SCSIHW
		row.VGA = cfg.VGA
		row.OnBoot = cfg.OnBoot
		row.Agent = cfg.Agent
		row.Hotplug = cfg.Hotplug
		row.Numa = cfg.Numa
		row.Protection = cfg.Protection

		row.Networks = cfg.MergeNets()

		disks := make(map[string]string)
		for k, v := range cfg.MergeIDEs() {
			disks[k] = v
		}
		for k, v := range cfg.MergeSCSIs() {
			disks[k] = v
		}
		for k, v := range cfg.MergeSATAs() {
			disks[k] = v
		}
		for k, v := range cfg.MergeVirtIOs() {
			disks[k] = v
		}
		row.Disks = disks
	}

	return row
}
