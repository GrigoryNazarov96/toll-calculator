package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GrigoryNazarov96/toll-calculator/types"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env file")
	}
	port := os.Getenv("ws_port")
	dr, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", dr.handleWS)
	http.ListenAndServe(port, nil)
}

type DataReceiver struct {
	conn *websocket.Conn
	prod DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p     DataProducer
		err   error
		topic = "obudata"
	)
	p, err = NewKafkaProducer(topic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReceiver{
		prod: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("OBU connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error: ", err)
			continue
		}
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka produce error: ", err)
		}
	}
}
