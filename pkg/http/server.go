package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mwantia/metrics-merger/pkg/common"
)

func CreateAndServe(cfg *common.ServerConfig) error {
	cache := common.NewCache()
	go StartScraping(cfg, cache)

	http.HandleFunc("/", HandleMetrics(cache))

	log.Printf("Starting server on '%s'", cfg.Address)
	return http.ListenAndServe(cfg.Address, nil)
}

func StartScraping(cfg *common.ServerConfig, cache *common.MetricsCache) {
	interval, err := time.ParseDuration(cfg.ScrapeInterval)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to parse scrape interval '%s'", cfg.ScrapeInterval))
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		for _, endpoint := range cfg.Endpoints {
			go func(endpoint common.EndpointConfig) {
				metrics, err := FetchEndpointBody(endpoint, cfg.MergeLabel)
				if err != nil {
					log.Printf("Error fetching metrics from %s: %v", endpoint.Address, err)
					return
				}
				cache.SetMetrics(endpoint.Name, metrics)
			}(endpoint)
		}
		<-ticker.C
	}
}
