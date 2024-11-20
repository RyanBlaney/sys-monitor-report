package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sys-monitor-report/internal/collectors"
	"sys-monitor-report/internal/report"
	"sys-monitor-report/internal/utils"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	report.Init()

	fmt.Println("Starting system monitor...")

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("Prometheus metrics available at http://localhost:8080/metrics")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	ticker := time.NewTicker(logInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				collectors.CollectSystemMetrics(config)
			case <-stop:
				fmt.Println("Shutting down system monitor...")
				return
			}
		}
	}()

	<-stop
	fmt.Println("\nSystem monitor terminated.")
}
