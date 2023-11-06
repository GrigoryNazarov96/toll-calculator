package main

import (
	"log"

	"github.com/GrigoryNazarov96/toll-calculator/aggregator/client"
)

const (
	topic    = "obudata"
	endpoint = "http://localhost:4020/aggregate"
)

func main() {
	s := NewCalculatorService()
	s = NewLogMiddleware(s)
	cl := client.NewClient(endpoint)
	c, err := NewKafkaConsumer(topic, s, cl)
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
}
