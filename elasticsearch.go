package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

// Define structs to hold Elasticsearch cluster data
type ClusterHealth struct {
	ClusterName             string `json:"cluster_name"`
	Status                  string `json:"status"`
	TimedOut                bool   `json:"timed_out"`
	NumberOfNodes           int    `json:"number_of_nodes"`
	NumberOfDataNodes       int    `json:"number_of_data_nodes"`
	ActivePrimaryShards     int    `json:"active_primary_shards"`
	ActiveShards            int    `json:"active_shards"`
	RelocatingShards        int    `json:"relocating_shards"`
	InitializingShards      int    `json:"initializing_shards"`
	UnassignedShards        int    `json:"unassigned_shards"`
	DelayedUnassignedShards int    `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks    int    `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch   int    `json:"number_of_in_flight_fetch"`
}

// Define Prometheus metrics
var (
	elasticsearchInfoMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ghe_elk_info",
			Help: "Information about the Elasticsearch cluster",
		},
		[]string{"version"},
	)

	// ... Define other metrics here ...

	elasticsearchStatusMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ghe_elk_status",
			Help: "Current R.A.G. status in labels",
		},
		[]string{"status"},
	)

	// ... Define other metrics here ...
)

func registerElasticsearchMetrics(reg *prometheus.Registry) {
	reg.MustRegister(elasticsearchInfoMetric)
	reg.MustRegister(elasticsearchStatusMetric)
	// Register other Elasticsearch metrics here
}

func exportElasticsearchMetrics() {
	// Fetch Elasticsearch cluster health
	health, err := fetchClusterHealth()
	if err != nil {
		log.Printf("Error fetching cluster health: %v", err)
		return
	}

	// Update metrics
	elasticsearchInfoMetric.With(prometheus.Labels{"version": "1.0.0"}).Set(1)
	// ... Update other metrics based on the health data ...

	elasticsearchStatusMetric.With(prometheus.Labels{"status": "green"}).Set(boolToFloat64(health.Status == "green"))
	elasticsearchStatusMetric.With(prometheus.Labels{"status": "yellow"}).Set(boolToFloat64(health.Status == "yellow"))
	elasticsearchStatusMetric.With(prometheus.Labels{"status": "red"}).Set(boolToFloat64(health.Status == "red"))

	// ... Update other metrics ...
}

func fetchClusterHealth() (*ClusterHealth, error) {
	// Fetch cluster health from Elasticsearch API
	resp, err := http.Get("http://localhost:9200/_cluster/health")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode JSON response into ClusterHealth struct
	var health ClusterHealth
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, err
	}

	return &health, nil
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
