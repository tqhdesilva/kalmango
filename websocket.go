package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func mkHandler(t chan time.Time,
	s *Screen,
	kf *KalmanFilter) func(http.ResponseWriter,
	*http.Request,
	float64) {

	return func(w http.ResponseWriter, r *http.Request, timedelta float64) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read:", err)
				return
			}
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Printf("write:", err)
			}
		}
	}
}
