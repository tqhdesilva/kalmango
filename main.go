// TODO add control vector, control matrix correction for edge collisions
// maybe have a channel just for those events?
// TODO add covariance matrix to websocket stream
package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	const td float64 = .1
	rand.Seed(time.Now().UTC().UnixNano())

	h := MakeHandler(td)
	http.HandleFunc("/websocket", h)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
