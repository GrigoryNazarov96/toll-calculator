package main

import "github.com/GrigoryNazarov96/toll-calculator/types"

type Aggregator interface {
	AggregateTelemetryData(types.TelemetryData) error
}

type Storer interface {
	Insert(types.TelemetryData) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(s Storer) Aggregator {
	return &InvoiceAggregator{
		store: s,
	}
}

func (a *InvoiceAggregator) AggregateTelemetryData(data types.TelemetryData) error {
	return a.store.Insert(data)
}
