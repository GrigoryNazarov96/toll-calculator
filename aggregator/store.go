package main

import "github.com/GrigoryNazarov96/toll-calculator/types"

type InMemoryStore struct {
	data map[int]float64
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[int]float64),
	}
}

func (s *InMemoryStore) Insert(d types.TelemetryData) error {
	s.data[d.OBUID] += d.Distance
	return nil
}
