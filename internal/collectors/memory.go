package collectors

import (
	"fmt"

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

	return data, nil
}

// DisplayMemoryData prints the relevant memory metrics in a human-readable format
func DisplayMemoryData(memoryData *MemoryData) {
	fmt.Println("=== Memory Metrics ===")

	// Virtual Memory (RAM)
	fmt.Println("Virtual Memory (RAM):")
	fmt.Printf("  Total: %.2f GB\n", float64(memoryData.VirtualMemory.Total)/1e9)
	fmt.Printf("  Used: %.2f GB\n", float64(memoryData.VirtualMemory.Used)/1e9)
	fmt.Printf("  Free: %.2f GB\n", float64(memoryData.VirtualMemory.Free)/1e9)
	fmt.Printf("  Usage: %.2f%%\n", memoryData.VirtualMemory.UsedPercent)

	// Swap Memory
	fmt.Println("\nSwap Memory:")
	fmt.Printf("  Total: %.2f GB\n", float64(memoryData.SwapMemory.Total)/1e9)
	fmt.Printf("  Used: %.2f GB\n", float64(memoryData.SwapMemory.Used)/1e9)
	fmt.Printf("  Free: %.2f GB\n", float64(memoryData.SwapMemory.Free)/1e9)
	fmt.Printf("  Usage: %.2f%%\n", memoryData.SwapMemory.UsedPercent)
	fmt.Println("")
}
