package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GrigoryNazarov96/toll-calculator/types"
)

type Client struct {
	Endpoint string
}

func NewClient(e string) *Client {
	return &Client{
		Endpoint: e,
	}
}

func (c *Client) AggregateInvoice(data types.TelemetryData) error {
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
	return nil
}
