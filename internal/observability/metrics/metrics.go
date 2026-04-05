package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	GetRatesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_get_rates_total",
			Help: "Total number of GetRates calls",
		},
		[]string{"method", "status"},
	)

	GetRatesDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_get_rates_duration_seconds",
			Help:    "Duration of GetRates calls in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	RatesAsk = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "grinex_rates_ask",
		Help: "Last ask price in minor units",
	})

	RatesBid = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "grinex_rates_bid",
		Help: "Last bid price in minor units",
	})
)

func Register() {
	prometheus.MustRegister(GetRatesTotal, GetRatesDuration, RatesAsk, RatesBid)
}

func RunServer(ctx context.Context, port int) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	return srv.ListenAndServe()
}
