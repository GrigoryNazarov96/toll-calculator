package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/GrigoryNazarov96/toll-calculator/aggregator/client"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load env file")
	}
	var (
		e  = os.Getenv("agg_http_url")
		p  = os.Getenv("gateway_port")
		c  = client.NewHttpClient(e)
		ih = newInvoiceHandler(c)
	)

	http.HandleFunc("/invoice", makeApiFunc(ih.handleGetInvoice))
	http.ListenAndServe(p, nil)
}

// handlers
func (i *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	values, ok := r.URL.Query()["id"]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU id"})
		return fmt.Errorf("missing OBU id")
	}
	obuId, err := strconv.Atoi(values[0])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect OBU id"})
		return fmt.Errorf("incorrect OBU id")
	}
	inv, err := i.client.GetInvoice(context.TODO(), obuId)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

// helpers
func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"url":  r.RequestURI,
			}).Info("http request sent through the gateway")
		}(start)
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
