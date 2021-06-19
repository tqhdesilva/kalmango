package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	filepath := os.Getenv("API_CONFIG_PATH")
	var config *Config
	var err error
	if filepath == "" {
		config = DefaultConfig()
	} else {
		config, err = LoadConfig(filepath)
	}
	if err != nil {
		log.Fatal(err)
	}
	h := MakeHandler(config.TimeDelta)
	http.HandleFunc("/websocket", h)
	http.Handle("/", http.FileServer(http.Dir(config.Assets)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
