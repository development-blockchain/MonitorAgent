package systeminfo

import (
	"bytes"
	"fmt"
	"github.com/develope/MonitorAgent/metrics"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"os"
)

var (
	subnamespace   = "systeminfo"
	memTotal       = metrics.NewRegisteredGauge(subnamespace, "mem_total", []string{"hostname"})
	memAvailable   = metrics.NewRegisteredGauge(subnamespace, "mem_available", []string{"hostname"})
	memUsed        = metrics.NewRegisteredGauge(subnamespace, "mem_used", []string{"hostname"})
	memUsedPercent = metrics.NewRegisteredGauge(subnamespace, "mem_used_percent", []string{"hostname"})

	hostUptime     = metrics.NewRegisteredGauge(subnamespace, "host_uptime", []string{"hostname"})
	hostProcessors = metrics.NewRegisteredGauge(subnamespace, "host_processor_count", []string{"hostname"})
	hostCpuCount   = metrics.NewRegisteredGauge(subnamespace, "host_cpu_count", []string{"hostname"})
	hostCpuLoad    = metrics.NewRegisteredGauge(subnamespace, "host_cpu_load", []string{"hostname"})

	diskUsageTotal       = metrics.NewRegisteredGauge(subnamespace, "disk_usage_total", []string{"hostname", "path"})
	diskUsageFree        = metrics.NewRegisteredGauge(subnamespace, "disk_usage_free", []string{"hostname", "path"})
	diskUsageUsed        = metrics.NewRegisteredGauge(subnamespace, "disk_usage_used", []string{"hostname", "path"})
	diskUsageUsedPercent = metrics.NewRegisteredGauge(subnamespace, "disk_usage_used_percent", []string{"hostname", "path"})
)

type CPUInformation struct {
	CPUCount int     `json:"cpus"`
	CPULoad  float64 `json:"cpuLoad"`
}

type SystemInformation struct {
	Memory  mem.VirtualMemoryStat `json:"memory"`
	Host    host.InfoStat         `json:"host"`
	Load    load.AvgStat          `json:"load"`
	CPU     CPUInformation        `json:"cpu"`
	Storage []disk.UsageStat      `json:"storage"`
	Network []net.IOCountersStat  `json:"network"`
}

func Metrics() {
	hostname, _ := os.Hostname()
	info := SystemInfo()
	memTotal.WithLabelValues(hostname).Set(float64(info.Memory.Total))
	memAvailable.WithLabelValues(hostname).Set(float64(info.Memory.Available))
	memUsed.WithLabelValues(hostname).Set(float64(info.Memory.Used))
	memUsedPercent.WithLabelValues(hostname).Set(info.Memory.UsedPercent)

	hostUptime.WithLabelValues(hostname).Set(float64(info.Host.Uptime))
	hostProcessors.WithLabelValues(hostname).Set(float64(info.Host.Procs))
	hostCpuCount.WithLabelValues(hostname).Set(float64(info.CPU.CPUCount))
	hostCpuLoad.WithLabelValues(hostname).Set(info.CPU.CPULoad)

	for _, diskUsage := range info.Storage {
		diskUsageTotal.WithLabelValues(hostname, diskUsage.Path).Set(float64(diskUsage.Total))
		diskUsageFree.WithLabelValues(hostname, diskUsage.Path).Set(float64(diskUsage.Free))
		diskUsageUsed.WithLabelValues(hostname, diskUsage.Path).Set(float64(diskUsage.Used))
		diskUsageUsedPercent.WithLabelValues(hostname, diskUsage.Path).Set(float64(diskUsage.UsedPercent))
	}
}

// SystemInfo - export Data structure form of the SystemInfo
func SystemInfo() SystemInformation {
	info := SystemInformation{}

	vmem, err := mem.VirtualMemory()
	if err == nil {
		info.Memory = *vmem
	}

	hostdata, cerr := host.Info()
	if cerr == nil {
		info.Host = *hostdata
	}

	loaddata, lerr := load.Avg()
	if lerr == nil {
		info.Load = *loaddata
	}

	info.CPU.CPUCount, _ = cpu.Counts(false)
	cpuLoad, cperr := cpu.Percent(0, false)
	if cperr == nil {
		info.CPU.CPULoad = cpuLoad[0]
	}

	// OK look up physical storage/partition info
	disks, derr := disk.Partitions(false)
	if derr == nil {
		// OK we have something - make space
		info.Storage = make([]disk.UsageStat, len(disks))
		// Range through the partitions and retrieve the data
		for idx, partition := range disks {
			diskInfo, err := disk.Usage(partition.Mountpoint)
			if err == nil {
				info.Storage[idx] = *diskInfo
			}

		}
	}

	// OK lookup network stats
	info.Network, _ = net.IOCounters(false)

	return info
}

// Prometheus - exports a prometheus compatible string
func Prometheus(metrics SystemInformation) string {
	var buffer bytes.Buffer

	buffer.WriteString("# HELP systeminfo_memory_bytes How much memory.\n")
	buffer.WriteString("# TYPE systeminfo_memory_bytes gauge\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_memory_bytes %v\n", metrics.Memory.Total))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_memory_available_bytes How much memory available.\n")
	buffer.WriteString("# TYPE systeminfo_memory_available_bytes gauge\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_memory_available_bytes %v\n", metrics.Memory.Available))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_memory_used_bytes How much memory available.\n")
	buffer.WriteString("# TYPE systeminfo_memory_used_bytes gauge\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_memory_used_bytes %v\n", metrics.Memory.Used))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_memory_used_percentage How much memory used as percentage.\n")
	buffer.WriteString("# TYPE systeminfo_memory_used_percentage gauge\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_memory_used_percentage %v\n", metrics.Memory.UsedPercent))
	buffer.WriteString("\n")
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_host_uptime How long been on\n")
	buffer.WriteString("# TYPE systeminfo_host_uptime counter\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_host_uptime %v\n", metrics.Host.Uptime))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_host_procs How Many Processes\n")
	buffer.WriteString("# TYPE systeminfo_host_procs guage\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_procs %v\n", metrics.Host.Procs))
	buffer.WriteString("\n")
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_load_1 Load Last Minute\n")
	buffer.WriteString("# TYPE systeminfo_load_1 guage\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_load_1 %v\n", metrics.Load.Load1))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_load_5 Load Last 5 Minute\n")
	buffer.WriteString("# TYPE systeminfo_load_5 guage\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_load_5 %v\n", metrics.Load.Load5))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_load_15 Load Last 15 Minute\n")
	buffer.WriteString("# TYPE systeminfo_load_15 guage\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_load_15 %v\n", metrics.Load.Load15))
	buffer.WriteString("\n")
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_cpu_cores Number or Cores\n")
	buffer.WriteString("# TYPE systeminfo_cpu_cores counter\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_cpu_cores %v\n", metrics.CPU.CPUCount))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_cpu_load load\n")
	buffer.WriteString("# TYPE systeminfo_cpu_load guage\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_cpu_load %v\n", metrics.CPU.CPULoad))
	buffer.WriteString("\n")
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_network_bytes_sent Network bytes sent\n")
	buffer.WriteString("# TYPE systeminfo_network_bytes_sent counter\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_network_bytes_sent %v\n", metrics.Network[0].BytesSent))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_network_bytes_received Network bytes received\n")
	buffer.WriteString("# TYPE systeminfo_network_bytes_received counter\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_network_bytes_received %v\n", metrics.Network[0].BytesRecv))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_network_packets_sent Network packets sent\n")
	buffer.WriteString("# TYPE systeminfo_network_packets_sent counter\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_network_packets_sent %v\n", metrics.Network[0].PacketsSent))
	buffer.WriteString("\n")
	buffer.WriteString("# HELP systeminfo_network_packets_received Network packets received\n")
	buffer.WriteString("# TYPE systeminfo_network_packets_received counter\n")
	buffer.WriteString(fmt.Sprintf("systeminfo_network_packets_received %v\n", metrics.Network[0].PacketsRecv))
	buffer.WriteString("\n")

	return buffer.String()
}