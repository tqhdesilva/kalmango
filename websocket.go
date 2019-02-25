package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"gonum.org/v1/gonum/mat"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func MakeHandler(td float64) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup s and kalman filter
		rand.Seed(time.Now().UTC().UnixNano())
		s := NewScreen(10, 10)
		c := make(chan time.Time)
		bc := make(chan Edge)
		go s.Run(td, c, bc)
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
		<-c
		for {
			time.Sleep(time.Duration(td/.001) * time.Millisecond)
			m := s.GetNoisyState()
			xvel := kf.State.mean.AtVec(2)
			yvel := kf.State.mean.AtVec(3)
			err = conn.WriteMessage(
				websocket.TextMessage,
				[]byte(fmt.Sprintf("Measurement: %+v", m)),
			)
			if err != nil {
				log.Fatal(err)
			}
			Bk := mat.NewDense(4, 4, []float64{
				0, 0, 0, 0,
				0, 0, 0, 0,
				-2 * xvel, 0, 0, 0,
				0, -2 * yvel, 0, 0,
			})
			uk := mat.NewVecDense(4, make([]float64, 4))
			select {
			case b := <-bc:
				err = conn.WriteMessage(
					websocket.TextMessage,
					[]byte(fmt.Sprintf("Bounced %+v", b)),
				)
				switch b := b; {
				case (b == Top) || (b == Bottom):
					uk.SetVec(1, 1)
				case (b == Left) || (b == Right):
					uk.SetVec(0, 1)
				}
			case <-c:
			}
			err = kf.Predict(Bk, uk)
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
