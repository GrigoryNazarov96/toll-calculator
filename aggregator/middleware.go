package main

import (
	"time"

	"github.com/GrigoryNazarov96/toll-calculator/types"
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
