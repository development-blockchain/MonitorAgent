package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	Namespace = "monitor"
)

func NewRegisteredCounterVec(subname string, name string, labels []string) *prometheus.CounterVec {
	counterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: subname,
			Name:      name,
		},
		labels,
	)
	prometheus.MustRegister(counterVec)
	return counterVec
}

func NewRegisteredGauge(subname string, name string, labels []string) *prometheus.GaugeVec {
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: subname,
			Name:      name,
		},
		labels,
	)
	prometheus.MustRegister(gaugeVec)
	return gaugeVec
}

func NewRegisteredHistogram(subname string, name string, labels []string) *prometheus.HistogramVec {
	histogramVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: subname,
			Name:      name,
		},
		labels,
	)
	prometheus.MustRegister(histogramVec)
	return histogramVec
}

func NewRegisteredSummary(subname string, name string, labels []string) *prometheus.SummaryVec {
	summaryVec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: Namespace,
			Subsystem: subname,
			Name:      name,
		},
		labels,
	)
	prometheus.MustRegister(summaryVec)
	return summaryVec
}

func StartMetrics(port int) {
	addr := fmt.Sprintf(":%d", port)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
