package main

import (
	"log"
	"net/http"
)

func main() {
	const td float64 = .1

	h := MakeHandler(td)
	http.HandleFunc("/websocket", h)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
