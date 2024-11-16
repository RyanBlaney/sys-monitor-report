package collectors

import (
	"fmt"
	"sync"
	"sys-monitor-report/internal/utils"
	"time"
)

func PrintSystemLog(config utils.Config) {
	var wg sync.WaitGroup
	metrics := struct {
		CPU      CPUData
		Memory   MemoryData
		DiskData []PartitionData
		IOSpeeds []DiskIOData
		err      error
		sync.Mutex
	}{}

	// Start CPU data collection
	wg.Add(1)
	go func() {
		defer wg.Done()

		cpuData, err := GetCPUData()

		metrics.Lock()
		defer metrics.Unlock()

		if err != nil {
			fmt.Printf("Error collecting CPU data: %v\n", err)
			return
		}
		metrics.CPU = cpuData
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		memoryData, err := GetMemoryData()

		metrics.Lock()
		defer metrics.Unlock()

		if err != nil {
			fmt.Printf("Error collecting memory data: %v\n", err)
			return
		}
		metrics.Memory = memoryData
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		diskData, err := GetPartitionData()

		metrics.Lock()
		defer metrics.Unlock()

		if err != nil {
			fmt.Printf("error collecting disk data: %v\n", err)
			return
		}
		metrics.DiskData = diskData
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ioSpeed, err := GetDiskIOSpeeds(time.Second * 10)

		metrics.Lock()
		defer metrics.Unlock()

		if err != nil {
			fmt.Printf("error collecting disk io speeds: %v\n", err)
		}
		metrics.IOSpeeds = ioSpeed
	}()

	wg.Wait()

	DisplayCPUData(&metrics.CPU)
	DisplayMemoryData(&metrics.Memory)
	DisplayPartitionData(&metrics.DiskData)
	DisplayDiskIOSpeeds(&metrics.IOSpeeds)
	DisplayTopProcesses()
}
