package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Create a new registry
	reg := prometheus.NewRegistry()

	// Register metrics from other files
	registerGovernorMetrics(reg)
	registerActionsMetrics(reg)
	registerElasticsearchMetrics(reg)
	registerReplicationMetrics(reg)

	// Create an HTTP handler for the metrics
	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	// Set up the HTTP server
	http.Handle("/metrics", handler)

	// Start background goroutines for updating metrics
	go updateElasticsearchMetrics()
	go updateReplicationMetrics()

	// Start the server
	fmt.Println("Starting Prometheus exporter on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updateElasticsearchMetrics() {
	for {
		exportElasticsearchMetrics() // This function is defined in elasticsearch.go
		time.Sleep(15 * time.Second)
	}
}

func updateReplicationMetrics() {
	for {
		exportReplicationMetrics() // This function is defined in replication.go
		time.Sleep(15 * time.Second)
	}
}

func updateGovernorMetrics() {
	for {
		exportGovernorMetrics() // This function is defined in governor.go
		time.Sleep(15 * time.Second)
	}
}

func updateActionsMetrics() {
	for {
		exportActionsMetrics() // This function is defined in actions.go
		time.Sleep(15 * time.Second)
	}
}
