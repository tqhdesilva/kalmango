package main

import (
	"log"
	"net/http"
)

func main() {
	config := LoadConfig()
	h := MakeHandler(config.TimeDelta)
	http.HandleFunc("/websocket", h)
	http.Handle("/", http.FileServer(http.Dir(config.Assets)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
