package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sys-monitor-report/internal/collectors"
	"sys-monitor-report/internal/utils"
	"syscall"
	"time"
)

func main() {
	config, err := utils.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logInterval := time.Duration(config.LogInterval) * time.Second

	// Graceful shutdown on quit
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Starting system monitor...")
	ticker := time.NewTicker(logInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				collectors.PrintSystemLog(config)
			case <-stop:
				fmt.Println("Shutting down system monitor...")
				return
			}
		}
	}()

	<-stop
	fmt.Println("System monitor terminated.")
}
