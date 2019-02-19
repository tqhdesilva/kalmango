package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	const timedelta float64 = .1
	rand.Seed(time.Now().UTC().UnixNano())

	handler := mkHandler(timedelta)
	http.HandleFunc("/websocket", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
