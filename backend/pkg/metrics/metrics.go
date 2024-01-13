package metrics

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	server *http.Server
}

func NewMetrics(address string) *Metrics {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return &Metrics{
		server: server,
	}
}

func (m *Metrics) Run() error {
	if err := m.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (m *Metrics) Stop(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}
