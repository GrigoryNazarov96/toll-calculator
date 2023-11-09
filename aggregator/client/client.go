package client

import (
	"context"

	"github.com/GrigoryNazarov96/toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.TelemetryDataRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}
