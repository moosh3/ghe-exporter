package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	licenseAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_license_available",
		Help: "Number of available GitHub Enterprise licenses",
	})

	licenseUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_license_used",
		Help: "Number of used GitHub Enterprise licenses",
	})

	licenseExpiration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "github_license_expiration",
		Help: "Days until expiration of the GitHub Enterprise license",
	})
)

func registerLicenseMetrics(reg *prometheus.Registry) {
	reg.MustRegister(licenseAvailable)
	reg.MustRegister(licenseUsed)
	reg.MustRegister(licenseExpiration)
}

type LicenseInfo struct {
	DaysUntilExpiration string `json:"days_until_expiration"`
	SeatsUsed           string `json:"seats_used"`
	SeatsAvailable      string `json:"seats_available"`
}

func fetchLicenseInfo() (*LicenseInfo, error) {
	githubHost := os.Getenv("GITHUB_HOST")
	if githubHost == "" {
		return nil, fmt.Errorf("GITHUB_HOST environment variable is not set")
	}

	url := fmt.Sprintf("https://%s/api/v3/enterprise/settings/license", githubHost)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching license info: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var licenseInfo LicenseInfo
	err = json.Unmarshal(body, &licenseInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return &licenseInfo, nil
}

func exportLicenseMetrics() {
	licenseInfo, err := fetchLicenseInfo()
	if err != nil {
		fmt.Printf("Error fetching license info: %v\n", err)
		return
	}

	seatsAvailable, err := strconv.ParseFloat(licenseInfo.SeatsAvailable, 64)
	if err == nil {
		licenseAvailable.Set(seatsAvailable)
	}

	seatsUsed, err := strconv.ParseFloat(licenseInfo.SeatsUsed, 64)
	if err == nil {
		licenseUsed.Set(seatsUsed)
	}

	daysUntilExpiration, err := strconv.ParseFloat(licenseInfo.DaysUntilExpiration, 64)
	if err == nil {
		licenseExpiration.Set(daysUntilExpiration)
	}
}
