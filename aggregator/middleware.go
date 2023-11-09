package main

import (
	"time"

	"github.com/GrigoryNazarov96/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

type MetricsMiddleware struct {
	aggr_counter     prometheus.Counter
	inv_counter      prometheus.Counter
	aggr_err_counter prometheus.Counter
	inv_err_counter  prometheus.Counter
	aggr_latency     prometheus.Histogram
	inv_latency      prometheus.Histogram
	next             Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	rca := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_req_counter",
	})
	rci := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "invoice_req_counter",
	})
	eca := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
	})
	eci := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "invoice_error_counter",
	})
	la := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_req_latency",

		Buckets: []float64{0.1, 0.5, 1},
	})
	li := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "invoice_req_latency",

		Buckets: []float64{0.1, 0.5, 1},
	})
	return &MetricsMiddleware{
		aggr_counter:     rca,
		inv_counter:      rci,
		aggr_err_counter: eca,
		inv_err_counter:  eci,
		aggr_latency:     la,
		inv_latency:      li,
		next:             next,
	}
}

func (m *LogMiddleware) AggregateTelemetryData(data types.TelemetryData) (err error) {
	start := time.Now()
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"took": time.Since(start),
		}).Info("aggregate telemetry data")
	}(start)
	err = m.next.AggregateTelemetryData(data)
	return
}

func (m *LogMiddleware) GetInvoice(id int) (i *types.Invoice, err error) {
	start := time.Now()
	defer func(start time.Time) {
		var (
			dist float64
			fee  float64
		)
		if i != nil {
			dist = i.TotalDistance
			fee = i.Fee
		}
		logrus.WithFields(logrus.Fields{
			"err":           err,
			"took":          time.Since(start),
			"OBUID":         id,
			"totalDistance": dist,
			"fee":           fee,
		}).Info("get invoice function")
	}(start)
	i, err = m.next.GetInvoice(id)
	return
}

func (m *MetricsMiddleware) AggregateTelemetryData(data types.TelemetryData) (err error) {
	start := time.Now()
	defer func(start time.Time) {
		m.aggr_latency.Observe(float64(time.Since(start).Seconds()))
		m.aggr_counter.Inc()
		if err != nil {
			m.aggr_err_counter.Inc()
		}
	}(start)
	err = m.next.AggregateTelemetryData(data)
	return
}

func (m *MetricsMiddleware) GetInvoice(id int) (i *types.Invoice, err error) {
	start := time.Now()
	defer func(start time.Time) {
		m.inv_latency.Observe(float64(time.Since(start).Seconds()))
		m.inv_counter.Inc()
		if err != nil {
			m.inv_err_counter.Inc()
		}
	}(start)
	i, err = m.next.GetInvoice(id)
	return
}
