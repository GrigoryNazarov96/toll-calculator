package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/GrigoryNazarov96/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Code int
	Err  error
}

func (e APIError) Error() string {
	return e.Err.Error()
}

type HTTPMetricHandler struct {
	req_counter prometheus.Counter
	err_counter prometheus.Counter
	req_latency prometheus.Histogram
}

func makeHTTPHandlerFunc(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writeJSON(w, apiErr.Code, map[string]string{"error": apiErr.Error()})
			}
		}
	}
}

func NewHTTPMetricHandler(namespace string) *HTTPMetricHandler {
	rc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", namespace, "req_counter"),
	})
	ec := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", namespace, "err_counter"),
	})
	rl := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", namespace, "req_latency"),
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &HTTPMetricHandler{
		req_counter: rc,
		err_counter: ec,
		req_latency: rl,
	}
}

func (h *HTTPMetricHandler) instrument(next APIFunc) APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		start := time.Now()
		defer func(start time.Time) {
			h.req_latency.Observe(time.Since(start).Seconds())
			h.req_counter.Inc()
			if err != nil {
				h.err_counter.Inc()
			}
		}(start)
		err = next(w, r)
		return err
	}
}

func handleAggregate(a Aggregator) APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "POST" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("method not supported (%s)", r.Method),
			}
		}
		var td types.TelemetryData
		if err := json.NewDecoder(r.Body).Decode(&td); err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("failed to decode the response body: %s", err),
			}
		}
		if err := a.AggregateTelemetryData(td); err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return writeJSON(w, http.StatusOK, map[string]string{"msg": "ok"})
	}
}

func handleGetInvoice(a Aggregator) APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid HTTP method %s", r.Method),
			}
		}
		values, ok := r.URL.Query()["id"]
		if !ok {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("missing OBU id"),
			}
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid OBU id %s", values[0]),
			}
		}
		invoice, err := a.GetInvoice(obuID)
		if err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return writeJSON(w, http.StatusOK, invoice)
	}
}
