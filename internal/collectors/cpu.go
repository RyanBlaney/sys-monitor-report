package collectors

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

// CPUData holds CPU metrics
type CPUData struct {
	TotalUsage float64
	PerCore    []float64
	NumCores   int32
}

// DisplayCPUData prints the relevant CPU metrics in a human-readable format
func DisplayCPUData(cpuData *CPUData) {
	fmt.Println("=== CPU Metrics ===")

	// Total CPU Usage
	fmt.Printf("Total CPU Usage: %.2f%%\n", cpuData.TotalUsage)

	// Per-Core CPU Usage
	fmt.Println("\nPer-Core CPU Usage:")
	for i, usage := range cpuData.PerCore {
		fmt.Printf("  Core %d: %.2f%%\n", i+1, usage)
	}
	fmt.Println("")
}

// GetCPUData collects overall, per-core, and top process CPU usage
func GetCPUData() (CPUData, error) {
	var data CPUData

	// Total CPU usage
	totalUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return data, fmt.Errorf("error collecting total CPU usage: %v\n", err)
	}
	if len(totalUsage) > 0 {
		data.TotalUsage = totalUsage[0]
	}

	// Per-Core Usage
	perCoreUsage, err := cpu.Percent(time.Second, true)
	if err != nil {
		return data, fmt.Errorf("error collecting per-core CPU usage: %v\n", err)
	}
	data.PerCore = perCoreUsage

	// Number of Cores
	numCores, err := cpu.Counts(true)
	if err != nil {
		return data, fmt.Errorf("error collecting number of CPU cores: %v\n", err)
	}
	data.NumCores = int32(numCores)

	return data, nil
}
