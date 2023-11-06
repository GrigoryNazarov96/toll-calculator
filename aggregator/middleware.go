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
