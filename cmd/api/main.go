package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	config := LoadConfig()
	h := MakeHandler(config.TimeDelta)
	http.HandleFunc("/websocket", h)
	http.Handle("/", http.FileServer(http.Dir(config.Assets)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
