package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/GrigoryNazarov96/toll-calculator/types"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	interval = time.Second * 5
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
		ids[i] = int(rand.Intn(2799999999999999999))
	}
	return ids
}

func sendOBUData(conn *websocket.Conn, data types.OBUData) error {
	return conn.WriteJSON(data)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load the env file")
	}
	wsEndpoint := os.Getenv("ws_endpoint")
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
