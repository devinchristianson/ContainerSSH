package backend

import (
	"sync"

	"go.containerssh.io/containerssh/config"
	"go.containerssh.io/containerssh/http"
	internalConfig "go.containerssh.io/containerssh/internal/config"
	"go.containerssh.io/containerssh/internal/metrics"
	"go.containerssh.io/containerssh/internal/sshserver"
	"go.containerssh.io/containerssh/log"
)

// New creates a new backend handler.
//
//goland:noinspection GoUnusedExportedFunction
func New(
	config config.AppConfig,
	logger log.Logger,
	metricsCollector metrics.Collector,
	defaultAuthResponse sshserver.AuthResponse,
) (sshserver.Handler, error) {
	loader, err := internalConfig.NewHTTPLoader(
		config.ConfigServer,
		logger,
		metricsCollector,
	)
	if err != nil {
		return nil, err
	}

	var cleanupClient http.Client
	if config.CleanupServer.URL != "" {
		cleanupClient, err = http.NewClient(
			config.CleanupServer,
			logger,
		)
		if err != nil {
			return nil, err
		}
	}

	backendRequestsCounter := metricsCollector.MustCreateCounter(
		MetricNameBackendRequests,
		MetricUnitBackendRequests,
		MetricHelpBackendRequests,
	)
	backendErrorCounter := metricsCollector.MustCreateCounter(
		MetricNameBackendError,
		MetricUnitBackendError,
		MetricHelpBackendError,
	)

	return &handler{
		config:                 config,
		configLoader:           loader,
		cleanupClient:          cleanupClient,
		authResponse:           defaultAuthResponse,
		metricsCollector:       metricsCollector,
		logger:                 logger,
		backendRequestsCounter: backendRequestsCounter,
		backendErrorCounter:    backendErrorCounter,
		lock:                   &sync.Mutex{},
	}, nil
}
