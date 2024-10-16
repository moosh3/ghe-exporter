package main

import (
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	version = os.Getenv("VERSION")

	replicationInfoMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ghe_repl_status_info",
		Help: "GHE Replication Status Exporter Info",
		ConstLabels: prometheus.Labels{
			"version": version,
		},
	})

	replicationCreatedMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ghe_repl_status_created",
		Help: "GHE Replication Status Exporter Creation Timestamp",
	})

	replicationStatusMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ghe_repl_status",
		Help: "GHE Replication Status",
	}, []string{"subsystem"})
)

func registerReplicationMetrics(reg *prometheus.Registry) {
	reg.MustRegister(replicationInfoMetric)
	reg.MustRegister(replicationCreatedMetric)
	reg.MustRegister(replicationStatusMetric)
}

func exportReplicationMetrics() {
	subsystems := []string{"mysql", "mssql", "elk", "redis", "git", "pages", "alambic", "githooks", "consul"}
	for _, subsystem := range subsystems {
		value := getSubsystemStatus(subsystem)
		replicationStatusMetric.WithLabelValues(subsystem).Set(value)
	}
}

func getSubsystemStatus(subsystem string) float64 {
	// TODO: Implement actual status checking logic for each subsystem
	// For now, we'll return a random value between 0 and 1
	return float64(time.Now().UnixNano() % 2)
}
