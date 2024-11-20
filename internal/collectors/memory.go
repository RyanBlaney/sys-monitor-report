package collectors

import (
	"fmt"
	"sys-monitor-report/internal/report"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryData struct {
	VirtualMemory VirtualMemoryData
	SwapMemory    SwapMemoryData
	Memory        CombinedMemoryData
}

type VirtualMemoryData struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

type SwapMemoryData struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

type CombinedMemoryData struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func GetMemoryData() (MemoryData, error) {
	var data MemoryData

	vmStats, err := mem.VirtualMemory()
	if err != nil {
		return data, fmt.Errorf("error collecting virtual memory data: %v\n", err)
	}
	data.VirtualMemory = VirtualMemoryData{
		Total:       vmStats.Total,
		Used:        vmStats.Used,
		Free:        vmStats.Free,
		UsedPercent: vmStats.UsedPercent,
	}

	swapStats, err := mem.SwapMemory()
	if err != nil {
		return data, fmt.Errorf("error collecting swap memory data: %v", err)
	}
	data.SwapMemory = SwapMemoryData{
		Total:       swapStats.Total,
		Used:        swapStats.Used,
		Free:        swapStats.Free,
		UsedPercent: swapStats.UsedPercent,
	}

	data.Memory = CombinedMemoryData{
		Total:       swapStats.Total + vmStats.Total,
		Used:        swapStats.Used + vmStats.Used,
		Free:        swapStats.Free + vmStats.Free,
		UsedPercent: 0,
	}

	data.Memory.UsedPercent = (float64(data.Memory.Used) / float64(data.Memory.Total)) * 100

	// Update Prometheus Metrics
	report.OverallMemoryUsage.WithLabelValues("used_percent").Set(data.Memory.UsedPercent)
	report.OverallMemoryUsage.WithLabelValues("total_mb").Set(float64(data.Memory.Total) / 1e6)
	report.OverallMemoryUsage.WithLabelValues("used_mb").Set(float64(data.Memory.Used) / 1e6)
	report.OverallMemoryUsage.WithLabelValues("free_mb").Set(float64(data.Memory.Free) / 1e6)

	report.VirtualMemoryUsage.WithLabelValues("used_percent").Set(data.VirtualMemory.UsedPercent)
	report.VirtualMemoryUsage.WithLabelValues("total_mb").
		Set(float64(data.VirtualMemory.Total) / 1e6)
	report.VirtualMemoryUsage.WithLabelValues("used_mb").Set(float64(data.VirtualMemory.Used) / 1e6)
	report.VirtualMemoryUsage.WithLabelValues("free_mb").Set(float64(data.VirtualMemory.Free) / 1e6)

	report.SwapMemoryUsage.WithLabelValues("used_percent").Set(data.SwapMemory.UsedPercent)
	report.SwapMemoryUsage.WithLabelValues("total_mb").Set(float64(data.SwapMemory.Total) / 1e6)
	report.SwapMemoryUsage.WithLabelValues("used_mb").Set(float64(data.SwapMemory.Used) / 1e6)
	report.SwapMemoryUsage.WithLabelValues("free_mb").Set(float64(data.SwapMemory.Free) / 1e6)

	return data, nil
}

func FormatMemoryData(memoryData *MemoryData) string {
	var formattedData string

	// Overall Memory Metrics
	formattedData += "# HELP overall_memory_usage Overall memory usage statistics\n"
	formattedData += "# TYPE overall_memory_usage gauge\n"
	formattedData += fmt.Sprintf(
		"overall_memory_usage{type=\"used_percent\"} %.2f\n", memoryData.Memory.UsedPercent,
	)
	formattedData += fmt.Sprintf(
		"overall_memory_usage{type=\"total_mb\"} %.2f\n",
		float64(memoryData.Memory.Total)/(1024*1024),
	)
	formattedData += fmt.Sprintf(
		"overall_memory_usage{type=\"used_mb\"} %.2f\n",
		float64(memoryData.Memory.Used)/(1024*1024),
	)
	formattedData += fmt.Sprintf(
		"overall_memory_usage{type=\"free_mb\"} %.2f\n",
		float64(memoryData.Memory.Free)/(1024*1024),
	)

	// Virtual Memory Metrics
	formattedData += "# HELP virtual_memory_usage Virtual memory usage statistics\n"
	formattedData += "# TYPE virtual_memory_usage gauge\n"
	formattedData += fmt.Sprintf(
		"virtual_memory_usage{type=\"used_percent\"} %.2f\n", memoryData.VirtualMemory.UsedPercent,
	)
	formattedData += fmt.Sprintf(
		"virtual_memory_usage{type=\"total_mb\"} %.2f\n",
		float64(memoryData.VirtualMemory.Total)/(1024*1024),
	)
	formattedData += fmt.Sprintf(
		"virtual_memory_usage{type=\"used_mb\"} %.2f\n",
		float64(memoryData.VirtualMemory.Used)/(1024*1024),
	)
	formattedData += fmt.Sprintf(
		"virtual_memory_usage{type=\"free_mb\"} %.2f\n",
		float64(memoryData.VirtualMemory.Free)/(1024*1024),
	)

	// Swap Memory Metrics
	formattedData += "# HELP swap_memory_usage Swap memory usage statistics\n"
	formattedData += "# TYPE swap_memory_usage gauge\n"
	formattedData += fmt.Sprintf(
		"swap_memory_usage{type=\"used_percent\"} %.2f\n", memoryData.SwapMemory.UsedPercent,
	)
	formattedData += fmt.Sprintf(
		"swap_memory_usage{type=\"total_mb\"} %.2f\n",
		float64(memoryData.SwapMemory.Total)/(1024*1024),
	)
	formattedData += fmt.Sprintf(
		"swap_memory_usage{type=\"used_mb\"} %.2f\n",
		float64(memoryData.SwapMemory.Used)/(1024*1024),
	)
	formattedData += fmt.Sprintf(
		"swap_memory_usage{type=\"free_mb\"} %.2f\n",
		float64(memoryData.SwapMemory.Free)/(1024*1024),
	)

	return formattedData
}
