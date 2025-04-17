package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	blacklistStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "blacklist_status",
			Help: "Blacklist status for IP addresses (1 for blacklisted, 0 for not blacklisted)",
		},
		[]string{"ip_address", "blacklist_name", "status"},
	)
	token       string
	ipAddresses string
	host        string
	port        string
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
	ipAddressesIter := strings.SplitSeq(ipAddresses, ",")

	for ipAddress := range ipAddressesIter {
		url := "https://api.mxtoolbox.com/api/v1/lookup/blacklist/" + ipAddress

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("Authorization", token)

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
			return
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Error while reading the response bytes: ", err)
			return
		}

		var blacklistData BlacklistResponse
		if err := json.Unmarshal(body, &blacklistData); err != nil {
			log.Println("Error parsing JSON response:", err)
			return
		}

		for _, entry := range blacklistData.Passed {
			blacklistStatus.WithLabelValues(ipAddress, entry.Name, "Passed").Set(0)
		}
		for _, entry := range blacklistData.Failed {
			blacklistStatus.WithLabelValues(ipAddress, entry.Name, "Failed").Set(1)
		}
		for _, entry := range blacklistData.Warnings {
			blacklistStatus.WithLabelValues(ipAddress, entry.Name, "Warning").Set(0)
		}
		for _, entry := range blacklistData.Timeouts {
			blacklistStatus.WithLabelValues(ipAddress, entry.Name, "Timeout").Set(0)
		}
	}
}
func main() {
	token = os.Getenv("API_TOKEN")
	ipAddresses = os.Getenv("IP_ADDRESSES")
	host = os.Getenv("HOST")
	port = os.Getenv("PORT")
	if token == "" || ipAddresses == "" {
		log.Fatal("Please provide API_TOKEN and IP_ADDRESSES as environment variables")
	}
	if host == "" {
		log.Println("No host provided, using default (bind to all interfaces)")
	}
	if port == "" {
		log.Println("No port provided, using default (2112)")
		port = "2112"
	}

	prometheus.MustRegister(blacklistStatus)

	go func() {
		for {
			getMetrics()
			time.Sleep(24 * time.Hour)
		}
	}()

	listenAddr := net.JoinHostPort(host, port)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(listenAddr, nil)
}
