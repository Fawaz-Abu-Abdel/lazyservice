package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lazyservice/lazyservice/internal/app"
	"github.com/lazyservice/lazyservice/internal/detectors"
	"github.com/lazyservice/lazyservice/internal/metrics"
	"github.com/lazyservice/lazyservice/internal/ui"
)

func main() {
	// Parse command line flags
	minimal := flag.Bool("minimal", false, "Use minimal dashboard (values-only updates)")
	flag.Parse()
	
	fmt.Println("🚀 Starting LazyService...")
	
	// Create application
	application := app.NewApp()

	detectorCount := 0

	// Register Docker detector and metrics collector
	dockerDetector, err := detectors.NewDockerDetector()
	if err != nil {
		fmt.Printf("⚠️  Docker detector not available: %v\n", err)
	} else {
		application.RegisterDetector(dockerDetector)
		dockerMetrics, err := metrics.NewDockerMetricsCollector()
		if err == nil {
			application.RegisterMetricsCollector(app.ServiceTypeDocker, dockerMetrics)
		}
		fmt.Println("✓ Docker detector registered")
		detectorCount++
	}

	// Register Kubernetes detector and metrics collector
	k8sDetector, err := detectors.NewKubernetesDetector()
	if err != nil {
		// K8s is optional, don't show warning on Windows
	} else {
		application.RegisterDetector(k8sDetector)
		k8sMetrics, err := metrics.NewK8sMetricsCollector()
		if err == nil {
			application.RegisterMetricsCollector(app.ServiceTypeKubernetes, k8sMetrics)
		}
		fmt.Println("✓ Kubernetes detector registered")
		detectorCount++
	}

	// Register Systemd detector
	systemdDetector := detectors.NewSystemdDetector()
	application.RegisterDetector(systemdDetector)
	fmt.Println("✓ Systemd detector registered")
	detectorCount++

	// Register Process detector and metrics collector
	processDetector := detectors.NewProcessDetector(nil)
	application.RegisterDetector(processDetector)
	processMetrics := metrics.NewProcessMetricsCollector()
	application.RegisterMetricsCollector(app.ServiceTypeProcess, processMetrics)
	fmt.Println("✓ Process detector registered")
	detectorCount++

	fmt.Printf("\n📊 %d detectors active. Starting dashboard...\n", detectorCount)
	time.Sleep(2 * time.Second)
	
	// Start application
	application.Start()

	// Create dashboard based on mode
	var dashboard interface {
		Run() error
		Stop()
	}
	
	if *minimal {
		fmt.Println("📺 Starting MINIMAL dashboard (values-only updates)...")
		dashboard = ui.NewMinimalDashboard(application)
	} else {
		fmt.Println("📺 Starting PRO dashboard...")
		dashboard = ui.NewProDashboard(application)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		dashboard.Stop()
		application.Stop()
		os.Exit(0)
	}()

	// Run dashboard
	if err := dashboard.Run(); err != nil {
		log.Fatalf("Error running dashboard: %v", err)
	}

	// Clean shutdown
	application.Stop()
}
