package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	blacklistStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "blacklist_status",
			Help: "Blacklist status for IP addresses (1 for blacklisted, 0 for not blacklisted)",
		},
		[]string{"ip_address", "blacklist_name", "status"},
	)
)

type BlacklistEntry struct {
	Name string `json:"Name"`
}

type BlacklistResponse struct {
	Passed   []BlacklistEntry `json:"Passed"`
	Failed   []BlacklistEntry `json:"Failed"`
	Warnings []BlacklistEntry `json:"Warnings"`
	Timeouts []BlacklistEntry `json:"Timeouts"`
}

func getMetrics() {
	godotenv.Load()
	token := os.Getenv("API_TOKEN")
	ipAddresses := os.Getenv("IP_ADDRESSES")

	IPs := strings.Split(ipAddresses, ",")

	for _, ipAddress := range IPs {
		url := "https://api.mxtoolbox.com/api/v1/lookup/blacklist/" + ipAddress

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("Authorization", token)

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
			continue
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Error while reading the response bytes: ", err)
		}

		var blacklistData BlacklistResponse
		if err := json.Unmarshal(body, &blacklistData); err != nil {
			log.Println("Error parsing JSON response:", err)
			continue
		}

		for _, entry := range blacklistData.Passed {
			blacklistStatus.WithLabelValues(ipAddress, entry.Name, "Passed").Set(0)
		}
		for _, entry := range blacklistData.Failed {
			blacklistStatus.WithLabelValues(ipAddress, entry.Name, "Failed").Set(1)
		}
		for _, entry := range blacklistData.Warnings {
			blacklistStatus.WithLabelValues(ipAddress, entry.Name, "Warnings").Set(0)
		}
		for _, entry := range blacklistData.Timeouts {
			blacklistStatus.WithLabelValues(ipAddress, entry.Name, "Timeouts").Set(0)
		}
	}
}
func main() {
	prometheus.MustRegister(blacklistStatus)

	go func() {
		for {
			getMetrics()
			time.Sleep(24 * time.Hour)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
