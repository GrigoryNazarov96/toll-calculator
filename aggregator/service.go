package main

import "github.com/GrigoryNazarov96/toll-calculator/types"

const basePrice = 3.15

type Aggregator interface {
	AggregateTelemetryData(types.TelemetryData) error
	GetInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.TelemetryData) error
	Get(int) (float64, error)
}

type DistanceAggregator struct {
	store Storer
}

func NewDistanceAggregator(s Storer) Aggregator {
	return &DistanceAggregator{
		store: s,
	}
}

func (a *DistanceAggregator) AggregateTelemetryData(data types.TelemetryData) error {
	return a.store.Insert(data)
}

func (a *DistanceAggregator) GetInvoice(id int) (*types.Invoice, error) {
	d, err := a.store.Get(id)
	if err != nil {
		return nil, err
	}
	i := &types.Invoice{
		OBUID:         id,
		TotalDistance: d,
		Fee:           basePrice * d,
	}
	return i, nil
}
