package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	const timeDelta float64 = .1
	rand.Seed(time.Now().UTC().UnixNano())
	screen := NewScreen(10, 10)
	c := make(chan time.Time)
	go screen.Run(timeDelta, c)

	kf := &KalmanFilter{}
	handler := mkHandler(c, screen, kf)
	http.HandleFunc("/websocket", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
