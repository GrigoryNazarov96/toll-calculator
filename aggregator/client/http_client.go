package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GrigoryNazarov96/toll-calculator/types"
)

type HTTPClient struct {
	Endpoint string
	Client
}

func NewHttpClient(e string) *HTTPClient {
	return &HTTPClient{
		Endpoint: e,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, data *types.TelemetryDataRequest) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with status code %d", res.StatusCode)
	}
	defer res.Body.Close()
	return nil
}

func (c *HTTPClient) GetInvoice(ctx context.Context, obuId int) (*types.Invoice, error) {
	invReq := &types.GetInvoiceRequest{
		ObuID: int32(obuId),
	}
	b, err := json.Marshal(invReq)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/%s?id=%d", c.Endpoint, "invoice", obuId)
	req, err := http.NewRequest("GET", endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responded with status code %d", res.StatusCode)
	}
	var inv *types.Invoice
	if err := json.NewDecoder(res.Body).Decode(&inv); err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return inv, nil
}
