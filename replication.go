package main

import (
	"log"
	"os"
	"strings"

	"github.com/pkg/exec"
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
	output, err := exec.Command("/usr/local/bin/ghe-repl-status").Output()
	if err != nil {
		log.Printf("Error running ghe-repl-status: %v", err)
		return 0
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, subsystem) {
			if strings.HasPrefix(line, "OK:") {
				return 1
			} else {
				return 0
			}
		}
	}

	log.Printf("Subsystem %s not found in ghe-repl-status output", subsystem)
	return 0
}
