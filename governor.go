package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	governorUploadedBytesIncrease = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ghe_governor_uploaded_bytes_increase10m",
			Help: "Uploaded bytes increase over 10 minutes",
		},
		[]string{"job", "environment", "role", "activities", "key_field", "sort_field"},
	)

	governorReceivedBytesIncrease = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ghe_governor_received_bytes_increase10m",
			Help: "Received bytes increase over 10 minutes",
		},
		[]string{"job", "environment", "role", "activities", "key_field", "sort_field"},
	)

	governorSamplesIncrease = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ghe_governor_samples_increase10m",
			Help: "Samples increase over 10 minutes",
		},
		[]string{"job", "environment", "role", "activities", "key_field", "sort_field"},
	)

	governorRuntimeSecondsIncrease = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ghe_governor_runtime_seconds_increase10m",
			Help: "Runtime seconds increase over 10 minutes",
		},
		[]string{"job", "environment", "role", "activities", "key_field", "sort_field"},
	)
)

func registerGovernorMetrics(reg *prometheus.Registry) {
	reg.MustRegister(governorUploadedBytesIncrease)
	reg.MustRegister(governorReceivedBytesIncrease)
	reg.MustRegister(governorSamplesIncrease)
	reg.MustRegister(governorRuntimeSecondsIncrease)
}

func exportGovernorMetrics(job, environment, role, activities string, uploadedBytes, receivedBytes, samples, runtimeSeconds float64) {
	labels := prometheus.Labels{
		"job":         job,
		"environment": environment,
		"role":        role,
		"activities":  activities,
		"key_field":   "repo",
		"sort_field":  "net", // Default for uploaded and received bytes
	}

	governorUploadedBytesIncrease.With(labels).Set(uploadedBytes)
	governorReceivedBytesIncrease.With(labels).Set(receivedBytes)

	labels["sort_field"] = "count"
	governorSamplesIncrease.With(labels).Set(samples)

	labels["sort_field"] = "rt"
	governorRuntimeSecondsIncrease.With(labels).Set(runtimeSeconds)
}
