package main

import (
	"log"
	"os"

	"github.com/GrigoryNazarov96/toll-calculator/aggregator/client"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env file")
	}
	grpc_endpoint := os.Getenv("grpc_endpoint")
	topic := os.Getenv("kafka_topic")
	s := NewCalculatorService()
	s = NewLogMiddleware(s)

	grpc_cl, err := client.NewGRPCClient(grpc_endpoint)
	if err != nil {
		log.Fatal(err)
	}

	c, err := NewKafkaConsumer(topic, s, grpc_cl)
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
}
