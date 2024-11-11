package common

import (
	"strings"
	"sync"
)

type MetricsCache struct {
	Metrics map[string]string
	Mutex   sync.RWMutex
}

func NewCache() *MetricsCache {
	return &MetricsCache{
		Metrics: make(map[string]string),
	}
}

func (c *MetricsCache) SetMetrics(name string, metrics string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Metrics[name] = metrics
}

func (c *MetricsCache) GetAllMetrics() string {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var metrics []string
	for _, m := range c.Metrics {
		metrics = append(metrics, m)
	}

	return strings.Join(metrics, "\n")
}
