package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GrigoryNazarov96/toll-calculator/types"
)

func main() {
	var (
		port  = ":4020"
		store = NewInMemoryStore()
		a     = NewDistanceAggregator(store)
	)
	a = NewLogMiddleware(a)
	makeHTTPTransport(port, a)
}

func makeHTTPTransport(port string, a Aggregator) {
	fmt.Println("HTTP transport is ready to handle requests on port ", port)
	http.HandleFunc("/aggregate", handleAggregate(a))
	http.HandleFunc("/invoice", handleGetInvoice(a))
	http.ListenAndServe(port, nil)
}

func handleAggregate(a Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var td types.TelemetryData
		if err := json.NewDecoder(r.Body).Decode(&td); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := a.AggregateTelemetryData(td); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func handleGetInvoice(a Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["id"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU id"})
			return
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect OBU id"})
			return
		}
		invoice, err := a.GetInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
