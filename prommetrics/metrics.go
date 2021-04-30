package prommetrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	DefaultShutdownTimeout = 5 * time.Second

	// HTTP server defaults.
	DefaultReadTimeout       = 5 * time.Second
	DefaultReadHeaderTimeout = 5 * time.Second
	DefaultWriteTimeout      = 5 * time.Second
)

var (
	// Errors
	ErrorServerAlreadyRunning error = errors.New("Prometheus HTTP server is already running")
)

// Metrics contains all set of methods to manage metrics collector instance behavior
type Metrics interface {
	StartHTTP(port uint, namespace string) error
	StopHTTP()

	PushCustom(url, pushJob string) error
	PushCollected() error

	gaugeMethods
	counterMethods
	histogramMethods
	summaryMethods
}

// metricsHandler implements Metrics interface.
type metricsHandler struct {
	sync.Mutex

	metricsNamespace string
	registry         *prometheus.Registry

	// prometheus publisher HTTP server parameters
	port   uint
	server *http.Server

	// pusher parameters
	pushUrl          string
	pushJob          string
	collectedMetrics map[string]prometheus.Collector
}

// NewMetrics creates a new metrics handler and returns its interface.
func NewMetrics() Metrics {
	return &metricsHandler{
		registry:         prometheus.NewRegistry(),
		collectedMetrics: make(map[string]prometheus.Collector),
	}
}

// StartHTTP opens and listens HTTP route to let external Prometheus to collect metrics.
func (m *metricsHandler) StartHTTP(port uint, namespace string) error {
	m.Lock()
	defer m.Unlock()

	if port == 0 {
		log.Debug("Prometheus server.Start(...) called with port == 0")
		return errors.New("bad parameter: port")
	}
	if namespace == "" {
		log.Debug("Prometheus server.Start(...) called with an empty namespace")
		return errors.New("bad parameter: namespace")
	}

	if m.server != nil {
		if port == m.port && namespace == m.metricsNamespace {
			log.Warningf(
				"Prometheus server.Start(%d, '%s') called again while server is running",
				port,
				namespace,
			)
			return nil
		}
		log.Debugf(
			"Error: Prometheus server.Start(%d, '%s') called again while server is running with (%d, '%s')",
			port,
			namespace,
			m.port,
			m.metricsNamespace,
		)
		return ErrorServerAlreadyRunning
	}

	m.port = port
	m.metricsNamespace = namespace

	listen := fmt.Sprintf(":%d", port)
	m.server = &http.Server{
		Addr:              listen,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		WriteTimeout:      DefaultWriteTimeout,

		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
				}
			}),
	}

	go func(srv *http.Server) {
		for {
			log.Infof("prometheus service start to listen on port %s", listen)
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Errorf("prometheus service failed: %s", err)
			}
			time.Sleep(time.Second) //let service stop
		}
	}(m.server)

	return nil
}

// StopHTTP stops running metrics provider HTTP server.
func (m *metricsHandler) StopHTTP() {
	m.Lock()
	defer m.Unlock()

	if m.server == nil {
		log.Debug("metrics HTTP server is not running")
		return
	}

	log.Debug("metrics HTTP server is stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer func() {
		m.server = nil
		m.port = 0
		m.metricsNamespace = ""
		cancel()
	}()

	if err := m.server.Shutdown(ctx); err != nil {
		log.Errorf("stopping prometheus HTTP server error: %s", err)
	} else {
		log.Debug("prometheus HTTP server has been stopped")
	}
}

// addMetricToPush adds a metrics item to be pushed later
func (m *metricsHandler) addMetricItem(subsystem, name string, metric prometheus.Collector) {
	m.Lock()
	defer m.Unlock()

	m.registry.MustRegister(metric)
	m.collectedMetrics[subsystem+"."+name] = metric
}

// PushCustom pushes all metrics collected to a PushGatewayClient
func (m *metricsHandler) PushCustom(url, pushJob string) error {

	// save for next pushes
	m.pushUrl = url
	m.pushJob = pushJob

	pusher := push.New(url, pushJob)
	m.Lock()
	defer m.Unlock()

	var err error
	for name, met := range m.collectedMetrics {
		if err = pusher.Collector(met).Push(); err != nil {
			return fmt.Errorf("error pushing metric '%s' to a pushgateway '%s': %w", name, url, err)
		}
	}

	return nil
}

// PushCollected pushes all colelcted metrics to the same pushgateway that was used last time
func (m *metricsHandler) PushCollected() error {
	return m.PushCustom(m.pushUrl, m.pushJob)
}
