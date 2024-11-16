package collectors

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

// ProcessData holds information about a process
type ProcessData struct {
	PID        int32
	Name       string
	CPUUsage   float64 // CPU usage percentage
	MemUsage   float64 // Memory usage percentage
	ReadCount  uint64  // Number of read operations
	WriteCount uint64  // Number of write operations
}

// GetTopProcesses retrieves the top N processes for the specified metric
func GetTopProcesses(metric string, topN int) ([]ProcessData, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("error fetching processes: %v", err)
	}

	var processData []ProcessData
	for _, proc := range processes {
		var usage float64
		switch metric {
		case "cpu":
			usage, err = proc.CPUPercent()
		case "memory":
			var memUsage float32
			memUsage, err = proc.MemoryPercent()
			usage = float64(memUsage)
		default:
			return nil, fmt.Errorf("unsupported metric: %s", metric)
		}
		if err != nil || usage == 0 {
			continue
		}

		name, err := proc.Name()
		if err != nil {
			continue
		}

		ioCounter, err := proc.IOCounters()
		if err != nil {
			continue
		}

		processData = append(processData, ProcessData{
			PID:        proc.Pid,
			Name:       name,
			CPUUsage:   usage,
			MemUsage:   usage,
			ReadCount:  ioCounter.ReadCount,
			WriteCount: ioCounter.WriteCount,
		})
	}

	// Sort processes by usage in descending order
	sort.Slice(processData, func(i, j int) bool {
		switch metric {
		case "cpu":
			return processData[i].CPUUsage > processData[j].CPUUsage
		case "memory":
			return processData[i].MemUsage > processData[j].MemUsage
		}
		return false
	})

	// Return the top N processes
	if len(processData) > topN {
		processData = processData[:topN]
	}
	return processData, nil
}

// DisplayTopProcesses calls sub-functions to display top processes by category
func DisplayTopProcesses() {
	fmt.Println("\n=== Top Processes ===")

	// Top CPU Processes
	fmt.Println("Top CPU Processes:")
	topCPUProcesses, err := GetTopProcesses("cpu", 10)
	if err != nil {
		fmt.Printf("Error collecting top CPU processes: %v\n", err)
	} else {
		displayTopProcessesByMetric(topCPUProcesses, "CPU")
	}

	// Top Virtual Memory Processes
	fmt.Println("\nTop Memory Processes:")
	topVirtMemProcesses, err := GetTopProcesses("memory", 10)
	if err != nil {
		fmt.Printf("Error collecting top memory processes: %v\n", err)
	} else {
		displayTopProcessesByMetric(topVirtMemProcesses, "Memory")
	}
}

// displayTopProcessesByMetric prints the top N processes for a given metric
func displayTopProcessesByMetric(processes []ProcessData, metric string) {
	fmt.Printf(
		"  %-13s | %-26s | %-15s | %-16s | %-17s | %-10s\n",
		"PID", "Name", "Metric", "Usage", "I/O Reads", "I/O Writes",
	)
	fmt.Println(strings.Repeat("-", 120)) // Separator line

	for _, proc := range processes {
		fmt.Printf(
			"  PID: %-8d | Name: %-20s | %-15s | Usage: %-8.2f%% | Reads: %-10d | Writes: %-10d\n",
			proc.PID,
			proc.Name,
			metric,
			proc.CPUUsage,
			proc.ReadCount,
			proc.WriteCount,
		)
	}
}
