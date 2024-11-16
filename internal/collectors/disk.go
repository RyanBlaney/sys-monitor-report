package collectors

import (
	"fmt"
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
	}

	// Convert map to slice
	var partitionData []PartitionData
	for _, data := range partitionMap {
		partitionData = append(partitionData, *data)
	}

	return partitionData, nil
}

// DisplayPartitionData displays all of the relevent metrics
func DisplayPartitionData(partitions *[]PartitionData) {
	fmt.Println("=== Partition Data ===")
	for _, part := range *partitions {
		fmt.Printf("Device: %s\n", part.Device)
		fmt.Printf("  Filesystem: %s\n", part.Filesystem)
		fmt.Printf("  Total Space: %.2f GB\n", float64(part.Total)/1e9)
		fmt.Printf("  Used Space: %.2f GB\n", float64(part.Used)/1e9)
		fmt.Printf("  Free Space: %.2f GB\n", float64(part.Free)/1e9)
		fmt.Println("  Mountpoints:")
		for _, mount := range part.Mountpoints {
			fmt.Printf("    - %s\n", mount)
		}
	}
	fmt.Println("")
}

type DiskIOData struct {
	Device     string
	ReadSpeed  float64
	WriteSpeed float64
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
		}
	}

	return ioData, nil
}
