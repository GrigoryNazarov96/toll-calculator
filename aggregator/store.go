package main

import (
	"fmt"

	"github.com/GrigoryNazarov96/toll-calculator/types"
)

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

func (s *InMemoryStore) Get(id int) (float64, error) {
	d, ok := s.data[id]
	if !ok {
		return 0.0, fmt.Errorf("couldn't find distance for OBU id %d", id)
	}
	return d, nil
}
