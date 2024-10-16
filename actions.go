package main

import (
	"log"
	"os"
	"strings"
	"time"

	exec "github.com/pkg/exec"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Define the info metric
	actionsInfoMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ghe_actions_check_info",
			Help: "GHE Actions Check Exporter Info",
		},
		[]string{"version"},
	)

	// Define the created metric
	actionsCreatedMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ghe_actions_check_created",
			Help: "GHE Actions Check Exporter Creation Timestamp",
		},
	)

	// Define the main check metric
	actionsCheckMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ghe_actions_check",
			Help: "GHE Actions Subsystems Status",
		},
		[]string{"subsystem"},
	)
)

func registerActionsMetrics(reg *prometheus.Registry) {
	reg.MustRegister(actionsInfoMetric)
	reg.MustRegister(actionsCreatedMetric)
	reg.MustRegister(actionsCheckMetric)
}

func exportActionsMetrics() {
	// Set info metric
	actionsInfoMetric.With(prometheus.Labels{"version": os.Getenv("VERSION")}).Set(1)

	// Set created metric
	actionsCreatedMetric.Set(float64(time.Now().Unix()))

	// Run ghe-actions-check command and parse output
	output, err := exec.Command("/usr/local/bin/ghe-actions-check").Output()
	if err != nil {
		log.Printf("Error running ghe-actions-check: %v", err)
		return
	}

	// Parse the output and set check metrics
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) < 4 || parts[len(parts)-1] != "healthy!" {
			log.Printf("Unexpected line format: %s", line)
			continue
		}

		subsystem := strings.ToLower(strings.Join(parts[:len(parts)-3], "_"))
		actionsCheckMetric.With(prometheus.Labels{"subsystem": subsystem}).Set(1)
	}
}
