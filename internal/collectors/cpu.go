package collectors

import (
	"fmt"
	"sys-monitor-report/internal/report"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

// CPUData holds CPU metrics
type CPUData struct {
	TotalUsage float64
	PerCore    []float64
	NumCores   int32
}

func FormatCPUData(cpuData *CPUData) string {
	var formattedData string

	// Overall Usage
	formattedData = formattedData +
		"# HELP cpu_overall_usage CPU overall usage percentage"
	formattedData = formattedData +
		"# TYPE cpu_overall_usage gauge"
	formattedData = formattedData + fmt.Sprintf(
		"cpu_overall_usage %.2f%%\n\n", cpuData.TotalUsage,
	)

	// Per-Core Usage
	formattedData = formattedData +
		"# HELP cpu_usage_percentage CPU usage percentage by core"
	formattedData = formattedData +
		"# TYPE cpu_usage_percentage gauge"
	for i, usage := range cpuData.PerCore {
		formattedData = formattedData + fmt.Sprintf(
			"cpu_usage_percentage{core=\"%d\"} %.2f%%\n", i+1, usage,
		)
	}

	formattedData += "\n"

	return formattedData
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

		report.OverallCPUUsage.Set(data.TotalUsage)
	}

	// Per-Core Usage
	perCoreUsage, err := cpu.Percent(time.Second, true)
	if err != nil {
		return data, fmt.Errorf("error collecting per-core CPU usage: %v\n", err)
	}
	data.PerCore = perCoreUsage

	for i, usage := range perCoreUsage {
		report.PerCoreCPUUsage.WithLabelValues(fmt.Sprintf("core_%d", i+1)).Set(usage)
	}

	// Number of Cores
	numCores, err := cpu.Counts(true)
	if err != nil {
		return data, fmt.Errorf("error collecting number of CPU cores: %v\n", err)
	}
	data.NumCores = int32(numCores)

	return data, nil
}
