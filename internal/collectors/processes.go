package collectors

import (
	"fmt"
	"sort"
	"strings"
	"sys-monitor-report/internal/report"

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

	UpdatePrometheusTopProcesses(&processData, metric)

	return processData, nil
}

// FormatTopProcesses formats top processes data in Prometheus-compatible format
func FormatTopProcesses(processes []ProcessData, metric string) string {
	var formattedData string

	// HELP and TYPE for CPU and Memory Metrics
	formattedData += fmt.Sprintf(
		"# HELP process_%s_usage Process %s usage statistics\n",
		metric,
		metric,
	)
	formattedData += fmt.Sprintf("# TYPE process_%s_usage gauge\n", metric)

	// HELP and TYPE for I/O Reads and Writes
	formattedData += "# HELP process_io_read_count Process I/O read operations\n"
	formattedData += "# TYPE process_io_read_count counter\n"
	formattedData += "# HELP process_io_write_count Process I/O write operations\n"
	formattedData += "# TYPE process_io_write_count counter\n"

	// Loop through processes to add metrics
	for _, proc := range processes {
		// Add CPU/Memory usage
		formattedData += fmt.Sprintf(
			"process_%s_usage{pid=\"%d\",name=\"%s\"} %.2f\n",
			metric, proc.PID, escapeQuotes(proc.Name), proc.CPUUsage,
		)

		// Add I/O Read Count
		formattedData += fmt.Sprintf(
			"process_io_read_count{pid=\"%d\",name=\"%s\"} %d\n",
			proc.PID, escapeQuotes(proc.Name), proc.ReadCount,
		)

		// Add I/O Write Count
		formattedData += fmt.Sprintf(
			"process_io_write_count{pid=\"%d\",name=\"%s\"} %d\n",
			proc.PID, escapeQuotes(proc.Name), proc.WriteCount,
		)
	}

	return formattedData
}

func UpdatePrometheusTopProcesses(processes *[]ProcessData, metric string) {
	for _, proc := range *processes {
		name := escapeQuotes(proc.Name)

		switch metric {
		case "cpu":
			report.TopCPUProcesses.WithLabelValues(fmt.Sprintf("%d", proc.PID), name).
				Set(proc.CPUUsage)
		case "memory":
			report.TopMemoryProcesses.WithLabelValues(fmt.Sprintf("%d", proc.PID), name).
				Set(proc.MemUsage)
		}

		report.ProcessIOReadCount.WithLabelValues(fmt.Sprintf("%d", proc.PID), name).
			Add(float64(proc.ReadCount))
		report.ProcessIOWriteCount.WithLabelValues(fmt.Sprintf("%d", proc.PID), name).
			Add(float64(proc.WriteCount))
	}
}

// escapeQuotes escapes double quotes in process names for Prometheus labels
func escapeQuotes(input string) string {
	return strings.ReplaceAll(input, `"`, `\"`)
}
