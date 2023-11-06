package main

import (
	"time"

	"github.com/GrigoryNazarov96/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Calculator
}

func NewLogMiddleware(next Calculator) Calculator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data types.OBUData) (dist float64) {
	start := time.Now()
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"dist": dist,
		}).Info()
	}(start)
	dist = m.next.CalculateDistance(data)
	return
}
