package collectors

import (
	"fmt"
	"sys-monitor-report/internal/report"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

// PartitionData holds summary information about a physical partition
type PartitionData struct {
	Device      string   // Device name (e.g., /dev/sda1)
	Mountpoints []string // List of mount points sharing the device
	Filesystem  string   // File system type (e.g., ext4, btrfs)
	Total       uint64   // Total space in bytes
	Used        uint64   // Used space in bytes
	Free        uint64   // Free space in bytes
}

// GetPartitionData retrieves grouped disk usage statistics
func GetPartitionData() ([]PartitionData, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("error collecting disk partitions: %v", err)
	}

	partitionMap := make(map[string]*PartitionData)

	// Iterate over each partition
	for _, part := range partitions {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			fmt.Printf("Error retrieving usage for %s: %v\n", part.Mountpoint, err)
			continue
		}

		// Group by Device
		if _, exists := partitionMap[part.Device]; !exists {
			partitionMap[part.Device] = &PartitionData{
				Device:      part.Device,
				Mountpoints: []string{part.Mountpoint},
				Filesystem:  part.Fstype,
				Total:       usage.Total,
				Used:        usage.Used,
				Free:        usage.Free,
			}
		} else {
			// Append additional mount points for the same device
			partitionMap[part.Device].Mountpoints = append(partitionMap[part.Device].Mountpoints, part.Mountpoint)
		}

		// Update Prometheus Metrics
		report.PartitionSpace.WithLabelValues(part.Device, "total_gb").
			Set(float64(usage.Total) / 1e9)
		report.PartitionSpace.WithLabelValues(part.Device, "used_gb").
			Set(float64(usage.Used) / 1e9)
		report.PartitionSpace.WithLabelValues(part.Device, "free_gb").
			Set(float64(usage.Free) / 1e9)
		report.PartitionSpace.WithLabelValues(part.Device, "used_percent").
			Set(usage.UsedPercent)

		report.PartitionMountpoints.WithLabelValues(part.Device, part.Mountpoint).Set(1)

	}

	// Convert map to slice
	var partitionData []PartitionData
	for _, data := range partitionMap {
		partitionData = append(partitionData, *data)
	}

	return partitionData, nil
}

type DiskIOData struct {
	Device     string
	ReadSpeed  float64
	WriteSpeed float64
}

func FormatDiskIOSpeeds(diskData *[]DiskIOData) string {
	var formattedData string

	formattedData += "# Help disk_io_read_speed Disk I/O read speed in MB/s\n"
	formattedData += "# Type disk_io_read_speed gauge\n"
	formattedData += "# Help disk_io_write_speed Disk I/O write speed in MB/s\n"
	formattedData += "# Type disk_io_write_speed gauge\n"

	for _, data := range *diskData {
		// diskIOReadSpeed.With(prometheus.Labels{"device": data.Device}).Set(data.ReadSpeed)
		// diskIOWriteSpeed.With(prometheus.Labels{"device": data.Device}).Set(data.WriteSpeed)

		formattedData += fmt.Sprintf(
			"disk_io_read_speed{device=\"%s\"} %.2f\n",
			data.Device, data.ReadSpeed,
		)
		formattedData += fmt.Sprintf(
			"disk_io_write_speed{device=\"%s\"} %.2f\n",
			data.Device, data.WriteSpeed,
		)
	}

	return formattedData
}

func DisplayDiskIOSpeeds(diskData *[]DiskIOData) {
	fmt.Println("=== Disk I/O Speeds ===")
	for _, data := range *diskData {
		fmt.Printf("Device: %s\n", data.Device)
		fmt.Printf("  Read Speed: %.2f MB/s\n", data.ReadSpeed)
		fmt.Printf("  Write Speed: %.2f MB/s\n", data.WriteSpeed)
	}
}

func GetDiskIOSpeeds(interval time.Duration) ([]DiskIOData, error) {
	initialStats, err := disk.IOCounters()
	if err != nil {
		return nil, fmt.Errorf("error collecting initial disk I/O stats: %v\n", err)
	}

	time.Sleep(interval)

	finalStats, err := disk.IOCounters()
	if err != nil {
		return nil, fmt.Errorf("error collecting final disk I/O stats: %v\n", err)
	}

	// Calculate Speeds
	var ioData []DiskIOData
	for device, initial := range initialStats {
		if final, exists := finalStats[device]; exists {
			readSpeed := float64(final.ReadBytes-initial.ReadBytes) / interval.Seconds()
			writeSpeed := float64(final.WriteBytes-initial.WriteBytes) / interval.Seconds()

			ioData = append(ioData, DiskIOData{
				Device:     device,
				ReadSpeed:  readSpeed,
				WriteSpeed: writeSpeed,
			})

			// Update Prometheus Data
			report.DiskIOReadSpeed.WithLabelValues(device).Set(readSpeed / 1e6)
			report.DiskIOWriteSpeed.WithLabelValues(device).Set(writeSpeed / 1e6)
		}
	}

	return ioData, nil
}

// FormatPartitionData formats and displays partition data in Prometheus-compatible format
func FormatPartitionData(partitions *[]PartitionData) string {
	var formattedData string

	// Add HELP and TYPE directives
	formattedData += "# HELP partition_space Partition space usage statistics\n"
	formattedData += "# TYPE partition_space gauge\n"

	// Loop through partitions to format their data
	for _, part := range *partitions {
		deviceLabel := fmt.Sprintf(`device="%s"`, part.Device)

		// Add total space metric
		formattedData += fmt.Sprintf(
			"partition_space{%s,type=\"total_gb\"} %.2f\n",
			deviceLabel, float64(part.Total)/1e9,
		)

		// Add used space metric
		formattedData += fmt.Sprintf(
			"partition_space{%s,type=\"used_gb\"} %.2f\n",
			deviceLabel, float64(part.Used)/1e9,
		)

		// Add free space metric
		formattedData += fmt.Sprintf(
			"partition_space{%s,type=\"free_gb\"} %.2f\n",
			deviceLabel, float64(part.Free)/1e9,
		)

		// Add used percentage metric
		usedPercent := (float64(part.Used) / float64(part.Total)) * 100
		formattedData += fmt.Sprintf(
			"partition_space{%s,type=\"used_percent\"} %.2f\n",
			deviceLabel, usedPercent,
		)
	}

	// Add mountpoint data
	formattedData += "# HELP partition_mountpoints List of mountpoints per partition\n"
	formattedData += "# TYPE partition_mountpoints gauge\n"

	for _, part := range *partitions {
		for _, mount := range part.Mountpoints {
			formattedData += fmt.Sprintf(
				"partition_mountpoints{device=\"%s\",mount=\"%s\"} 1\n",
				part.Device, mount,
			)
		}
	}

	formattedData += "\n"

	return formattedData
}
