package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GrigoryNazarov96/toll-calculator/types"
)

func main() {
	var (
		port  = ":4020"
		store = NewInMemoryStore()
		a     = NewInvoiceAggregator(store)
	)
	a = NewLogMiddleware(a)
	makeHTTPTransport(port, a)
}

func makeHTTPTransport(port string, a Aggregator) {
	fmt.Println("HTTP transport is ready to handle requests on port ", port)
	http.HandleFunc("aggregate", handleAggregate(a))
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

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
