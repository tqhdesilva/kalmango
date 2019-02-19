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

func NewKalmanFilter(initialMeasurement *mat.VecDense, timedelta float64) (*KalmanFilter, error) {
	sensorCovMat, err := NewCovMat(4, []float64{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 1.0,
	})
	if err != nil {
		return nil, err
	}
	sensor := &Sensor{
		sensorCovMat,
	}

	stateToSensor := mat.NewDense(4, 4, []float64{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	})
	noise, err := NewCovMat(4, []float64{
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
	})
	prediction := mat.NewDense(4, 4, []float64{
		1.0, 0.0, timedelta, 0.0,
		0.0, 1.0, 0.0, timedelta,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	})
	// initially ought to be the same as sensorCovMat
	stateCovMat, err := NewCovMat(4, []float64{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 1.0,
	})
	if err != nil {
		return nil, err
	}

	initialState := &State{
		stateCovMat,
		initialMeasurement,
	}

	kf := &KalmanFilter{
		sensor,
		initialState,
		stateToSensor,
		noise,
		prediction,
	}
	return kf, nil
}

func mkHandler(timedelta float64) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup screen and kalman filter
		rand.Seed(time.Now().UTC().UnixNano())
		screen := NewScreen(10, 10)
		c := make(chan time.Time)
		go screen.Run(timedelta, c)
		initialMeasurement := screen.GetNoisyState()
		kf, err := NewKalmanFilter(initialMeasurement, timedelta)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		go func() {
			for {
				mt, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("read error: %s", err)
				}
				if (mt == websocket.TextMessage) &&
					(string(message) == "update") {
					kf.Update(screen.GetNoisyState())
					fmt.Println("predictions updated")
				}

			}
		}()
		for {
			// mt, message, err := conn.ReadMessage()
			// if err != nil {
			// 	log.Printf("read:", err)
			// 	return
			// }
			time.Sleep(time.Duration(timedelta/.001) * time.Millisecond)
			<-c
			err = conn.WriteMessage(
				websocket.TextMessage,
				[]byte(fmt.Sprintf("Measurement: %+v", screen.GetNoisyPosition())),
			)
			err = kf.Predict()
			if err != nil {
				log.Fatal(err)
			}
			err = kf.Update(screen.GetNoisyState())
			err = conn.WriteMessage(
				websocket.TextMessage,
				[]byte(fmt.Sprintf("Prediction: %+v", kf.State.mean)),
			)
			if err != nil {
				log.Printf("write: %s", err)
			}
		}
	}
}
