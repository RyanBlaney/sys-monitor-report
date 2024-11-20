package collectors

import (
	"fmt"
	"sync"
	"sys-monitor-report/internal/utils"
	"time"
)

// CollectSystemMetrics gathers system metrics and updates Prometheus metrics
func CollectSystemMetrics(config utils.Config) {
	var wg sync.WaitGroup

	// Dynamic sampling for CPU and Memory
	go DynamicSampling(
		"cpu",
		float64(config.Thresholds.CPU),
		time.Second*time.Duration(config.LogInterval),
		time.Second*time.Duration(config.LogIntervalHighFreq),
		30*time.Second,
	)

	go DynamicSampling(
		"memory",
		float64(config.Thresholds.Memory),
		time.Second*time.Duration(config.LogInterval),
		time.Second*time.Duration(config.LogIntervalHighFreq),
		30*time.Second,
	)

	// Collect Partition Data
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := GetPartitionData()
		if err != nil {
			fmt.Printf("Error collecting partition data: %v\n", err)
			return
		}
		// UpdatePartitionMetrics(partitionData) // Updates Prometheus metrics
	}()

	// Collect Disk I/O Speeds
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := GetDiskIOSpeeds(time.Second * 10)
		if err != nil {
			fmt.Printf("Error collecting disk I/O speeds: %v\n", err)
			return
		}
		// UpdateDiskIOMetrics(diskIOData) // Updates Prometheus metrics
	}()

	// Collect Top Processes (CPU and Memory)
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := GetTopProcesses("cpu", 10)
		if err != nil {
			fmt.Printf("Error collecting top CPU processes: %v\n", err)
			return
		}
		// UpdateTopProcessesMetrics(topCPUProcesses, "cpu")

		_, err = GetTopProcesses("memory", 10)
		if err != nil {
			fmt.Printf("Error collecting top memory processes: %v\n", err)
			return
		}
		// UpdateTopProcessesMetrics(topMemoryProcesses, "memory")
	}()

	// Wait for all tasks to complete
	wg.Wait()
	fmt.Println("System metrics collection completed.")
}

/* func PrintSystemLog(config utils.Config) {
	var wg sync.WaitGroup
	var results []string
	var mu sync.Mutex // Mutex to protect results slice

	// Start dynamic sampling for CPU and memory in separate goroutines
	go DynamicSampling(
		"cpu",
		float64(config.Thresholds.CPU),
		time.Second*time.Duration(config.LogInterval),
		time.Second*time.Duration(config.LogIntervalHighFreq),
		30*time.Second,
	)

	go DynamicSampling(
		"memory",
		float64(config.Thresholds.Memory),
		time.Second*time.Duration(config.LogInterval),
		time.Second*time.Duration(config.LogIntervalHighFreq),
		30*time.Second,
	)

	// Add tasks for disk data, disk IO speeds, and top processes
	wg.Add(2)

	go func() {
		defer wg.Done()
		diskData, err := GetPartitionData()
		output := ""
		if err != nil {
			output = fmt.Sprintf("Error collecting disk data: %v\n", err)
		} else {
			output = FormatPartitionData(&diskData)
		}
		mu.Lock()
		results = append(results, output)
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		ioSpeed, err := GetDiskIOSpeeds(time.Second * 10)
		output := ""
		if err != nil {
			output = fmt.Sprintf("error collecting disk IO speeds: %v\n", err)
		} else {
			output = formatDiskIOSpeeds(ioSpeed)
		}
		mu.Lock()
		results = append(results, output)
		mu.Unlock()
	}()
	go func() {
		defer wg.Done()

		output := ""
		topMemory, err := GetTopProcesses("memory", 10)
		if err != nil {
			output = fmt.Sprintf("error collecting top processes: %v\n", err)
			return
		}

		topCPU, err := GetTopProcesses("cpu", 10)
		if err != nil {
			output = fmt.Sprintf("error collecting top processes: %v\n", err)
			return
		}

		output += FormatTopProcesses(topCPU, "cpu")
		output += FormatTopProcesses(topMemory, "memory")

		mu.Lock()
		results = append(results, output)
		mu.Unlock()
	}()

	// Wait for all tasks to finish
	wg.Wait()

	// Print results in order
	for _, result := range results {
		fmt.Print(result)
	}
} */

// DynamicSampling monitors metrics dynamically
func DynamicSampling(
	metric string,
	threshold float64,
	normalInterval, highFreqInterval,
	highFreqDuration time.Duration,
) {
	currentInterval := normalInterval
	highFreqTimer := time.NewTimer(0)
	highFreqActive := false

	for {
		select {
		case <-time.After(currentInterval):
			var spikeDetected bool
			switch metric {
			case "cpu":
				cpuData, err := GetCPUData()
				if err != nil {
					fmt.Printf("Error collecting CPU data: %v\n", err)
					continue
				}

				// fmt.Printf("=== CPU Metrics ===\n")
				// DisplayCPUData(&cpuData)
				if cpuData.TotalUsage > threshold {
					fmt.Printf("CPU Spike Detected: %.2f%%\n", cpuData.TotalUsage)
					spikeDetected = true
				}
			case "memory":
				memoryData, err := GetMemoryData()
				if err != nil {
					fmt.Printf("Error collecting memory data: %v\n", err)
					continue
				}

				// fmt.Printf("=== Memory Metrics ===\n")
				// FormatMemoryData(&memoryData)
				if memoryData.Memory.UsedPercent > threshold {
					fmt.Printf("Memory Spike Detected: %.2f%%\n", memoryData.Memory.UsedPercent)
					spikeDetected = true
				}
			}

			// Adjust sampling rate if a spike is detected
			if spikeDetected && !highFreqActive {
				fmt.Println("Switching to high frequency sampling...")
				currentInterval = highFreqInterval
				highFreqTimer.Reset(highFreqDuration)
				highFreqActive = true
			}
		case <-highFreqTimer.C:
			// Revert to normal sampling after high-frequency duration
			if highFreqActive {
				fmt.Println("Reverting to normal sampling...")
				currentInterval = normalInterval
				highFreqActive = false
			}
		}
	}
}
