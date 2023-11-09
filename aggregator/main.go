package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/GrigoryNazarov96/toll-calculator/types"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load env file")
	}
	var (
		httpPort = os.Getenv("agg_http_port")
		grpcPort = os.Getenv("agg_grpc_port")
		store    = NewInMemoryStore()
		a        = NewDistanceAggregator(store)
	)
	a = NewMetricsMiddleware(a)
	a = NewLogMiddleware(a)
	go makeGRPCtransport(grpcPort, a)
	makeHTTPTransport(httpPort, a)
}

func makeGRPCtransport(port string, a Aggregator) error {
	fmt.Println("GRPC transport is ready to handle requests on port ", port)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer ln.Close()
	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewGRPCServer(a))
	return server.Serve(ln)
}

func makeHTTPTransport(port string, a Aggregator) {
	var (
		aggmh            = NewHTTPMetricHandler("aggregator")
		invmh            = NewHTTPMetricHandler("invoice")
		aggregateHandler = makeHTTPHandlerFunc(aggmh.instrument(handleAggregate(a)))
		invoiceHandler   = makeHTTPHandlerFunc(invmh.instrument(handleGetInvoice(a)))
	)
	http.HandleFunc("/aggregate", aggregateHandler)
	http.HandleFunc("/invoice", invoiceHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(port, nil)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
