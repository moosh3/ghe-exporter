package main

import (
	"log"
	"os"
	"strconv"
	"time"

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

	// Set check metrics for each subsystem
	subsystems := []string{
		"mssql", "mps_database", "mps_nomad_job", "mps_at_health_api",
		"token_database", "token_nomad_job", "token_at_health_api",
		"actions_database", "actions_nomad_job", "actions_at_health_api",
		"artifactcache_database", "artifactcache_nomad_job", "artifactcache_at_health_api",
	}

	for _, subsystem := range subsystems {
		value, err := strconv.ParseFloat(os.Getenv(subsystem), 64)
		if err != nil {
			log.Printf("Error parsing value for subsystem %s: %v", subsystem, err)
			continue
		}
		actionsCheckMetric.With(prometheus.Labels{"subsystem": subsystem}).Set(value)
	}
}

// The main function is removed as it's now handled in the main.go file
