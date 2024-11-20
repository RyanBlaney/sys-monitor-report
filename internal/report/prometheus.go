package report

import "github.com/prometheus/client_golang/prometheus"

var (
	// CPU Usage
	OverallCPUUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_overall_usage",
		Help: "Current CPU usage percentage",
	})

	PerCoreCPUUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percentage",
			Help: "CPU usage percentage by core",
		},
		[]string{"core"},
	)

	// Memory Usage
	OverallMemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "overall_memory_usage",
			Help: "Overall memory usage statistics",
		},
		[]string{"type"},
	)

	VirtualMemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_memory_usage",
			Help: "Virtual memory usage statistics",
		},
		[]string{"type"},
	)

	SwapMemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "swap_memory_usage",
			Help: "Swap memory usage statistics",
		},
		[]string{"type"},
	)

	// Partition Space Usage
	PartitionSpace = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "partition_space",
			Help: "Partition space usage statistics",
		},
		[]string{"device", "type"},
	)

	PartitionMountpoints = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "partition_mountpoints",
			Help: "List of mountpoints per partition",
		},
		[]string{"device", "mount"},
	)

	// Top Processes
	TopCPUProcesses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_cpu_usage",
			Help: "CPU usage percentage of top processes",
		},
		[]string{"pid", "name"},
	)

	TopMemoryProcesses = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_memory_usage",
			Help: "Memory usage percentage of top processes",
		},
		[]string{"pid", "name"},
	)

	ProcessIOReadCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "process_io_read_count",
			Help: "I/O read operations of top processes",
		},
		[]string{"pid", "name"},
	)

	ProcessIOWriteCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "process_io_write_count",
			Help: "I/O write operations of top processes",
		},
		[]string{"pid", "name"},
	)

	DiskIOReadSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_io_read_speed",
			Help: "Disk I/O read speed in MB/s",
		},
		[]string{"device"},
	)

	DiskIOWriteSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_io_write_speed",
			Help: "Disk I/O write speed in MB/s",
		},
		[]string{"device"},
	)
)

func Init() {
	prometheus.MustRegister(OverallCPUUsage)
	prometheus.MustRegister(PerCoreCPUUsage)
	prometheus.MustRegister(OverallMemoryUsage)
	prometheus.MustRegister(VirtualMemoryUsage)
	prometheus.MustRegister(SwapMemoryUsage)
	prometheus.MustRegister(PartitionSpace)
	prometheus.MustRegister(PartitionMountpoints)
	prometheus.MustRegister(TopCPUProcesses)
	prometheus.MustRegister(TopMemoryProcesses)
	prometheus.MustRegister(ProcessIOReadCount)
	prometheus.MustRegister(ProcessIOWriteCount)
	prometheus.MustRegister(DiskIOReadSpeed)
	prometheus.MustRegister(DiskIOWriteSpeed)
}
