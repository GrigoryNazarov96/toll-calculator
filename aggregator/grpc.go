package main

import (
	"context"

	"github.com/GrigoryNazarov96/toll-calculator/types"
)

type GRPCServer struct {
	types.UnimplementedAggregatorServer
	a Aggregator
}

func NewGRPCServer(a Aggregator) *GRPCServer {
	return &GRPCServer{
		a: a,
	}
}

func (s *GRPCServer) Aggregate(ctx context.Context, req *types.TelemetryDataRequest) (*types.None, error) {
	data := types.TelemetryData{
		OBUID:    int(req.ObuID),
		Distance: req.Distance,
		Unix:     req.Unix,
	}
	return &types.None{}, s.a.AggregateTelemetryData(data)
}
