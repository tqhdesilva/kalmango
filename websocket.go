package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func MakeHandler(td float64) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup s and kalman filter
		rand.Seed(time.Now().UTC().UnixNano())
		s := NewScreen(10, 10)
		c := make(chan time.Time)
		go s.Run(td, c)
		initialMeasurement := s.GetNoisyState()
		kf, err := NewKalmanFilter(initialMeasurement, td)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		go func() {
			for {
				mt, msg, err := conn.ReadMessage()
				if err != nil {
					log.Printf("read error: %s", err)
				}
				if (mt == websocket.TextMessage) &&
					(string(msg) == "update") {
					err = kf.Update(s.GetNoisyState())
					if err != nil {
						log.Fatal(err)
					}
				}

			}
		}()
		for {
			time.Sleep(time.Duration(td/.001) * time.Millisecond)
			<-c
			err = conn.WriteMessage(
				websocket.TextMessage,
				[]byte(fmt.Sprintf("Measurement: %+v", s.GetNoisyState())),
			)
			err = kf.Predict()
			if err != nil {
				log.Fatal(err)
			}
			err = kf.Update(s.GetNoisyState())
			if err != nil {
				log.Fatal(err)
			}
			err = conn.WriteMessage(
				websocket.TextMessage,
				[]byte(fmt.Sprintf("Mean: %+v", kf.State.mean)),
			)
			conn.WriteMessage(
				websocket.TextMessage,
				[]byte(fmt.Sprintf("Covariance: %+v", kf.State.covariance.SymDense)),
			)
			if err != nil {
				log.Printf("write: %s", err)
			}
		}
	}
}

// TODO handle socket.close()
