package client

import (
	"context"

	"github.com/GrigoryNazarov96/toll-calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	endpoint string
	client   types.AggregatorClient
}

func NewGRPCClient(e string) (*GRPCClient, error) {
	conn, err := grpc.Dial(e, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		endpoint: e,
		client:   c,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, data *types.TelemetryDataRequest) error {
	_, err := c.client.Aggregate(ctx, data)
	return err
}

func (c *GRPCClient) GetInvoice(ctx context.Context, obuID int) (*types.Invoice, error) {
	return &types.Invoice{}, nil
}
