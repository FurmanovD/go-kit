package prommetrics

import (
	"sync"
)

// metricsInstance handles an info to be able to use metrics singleton globally.
var (
	mu       sync.Mutex
	instance Metrics
)

// Instance creates a metrics singletone
func Instance() Metrics {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = NewMetrics()
	}
	return instance
}

// InstanceReset stops and deletes a Metrics singletone
func InstanceReset() Metrics {
	mu.Lock()
	defer mu.Unlock()

	if instance != nil {
		instance.StopHTTP()
	}
	instance = NewMetrics()
	return instance
}
