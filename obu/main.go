package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/GrigoryNazarov96/toll-calculator/types"
	"github.com/gorilla/websocket"
)

const (
	interval   = time.Second * 5
	wsEndpoint = "ws://localhost:30000/ws"
)

func genCoordinates() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func genLoc() (float64, float64) {
	return genCoordinates(), genCoordinates()
}

func genOBUID(n int) []int {
	ids := make([]int, n)
	for i := range ids {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func sendOBUData(conn *websocket.Conn, data types.OBUData) error {
	return conn.WriteJSON(data)
}

func main() {
	OBUIDs := genOBUID(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for _, id := range OBUIDs {
			lat, long := genLoc()
			data := types.OBUData{
				OBUID: id,
				Lat:   lat,
				Long:  long,
			}
			if err := sendOBUData(conn, data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(interval)
	}
}
