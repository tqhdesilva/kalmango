package main

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"gonum.org/v1/gonum/mat"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Message interface{}

type StateMessage struct {
	EstimatedPosition   []float64   `json:"estimated_position"`
	EstimatedCovariance [][]float64 `json:"estimated_covariance"`
	ActualPosition      []float64   `json:"actual_position"`
	Time                time.Time   `json:"timestamp"`
}

type UpdateMessage struct {
	NoisyPosition []float64 `json:"noisy_position"`
	Time          time.Time `json:"timestamp"`
}

type Connection struct {
	connection *websocket.Conn
	mux        sync.Mutex
}

func NewUpdateMessage(m *mat.VecDense) UpdateMessage {
	return UpdateMessage{
		NoisyPosition: vecToSlice(m),
		Time:          time.Now(),
	}
}

func vecToSlice(m mat.Vector) []float64 {
	n := m.Len()
	s := make([]float64, n)
	for i := 0; i < n; i++ {
		s[i] = m.AtVec(i)
	}
	return s
}

func denseToSlice(d mat.Matrix) [][]float64 {
	r, c := d.Dims()
	m := make([][]float64, r)
	for i := 0; i < r; i++ {
		m[i] = make([]float64, c)
		for j := 0; j < c; j++ {
			m[i][j] = d.At(i, j)
		}
	}
	return m
}

func NewStateMessage(kf *KalmanFilter, s *Screen, t time.Time) StateMessage {
	return StateMessage{
		EstimatedPosition:   vecToSlice(kf.State.mean.SliceVec(0, 2)),
		EstimatedCovariance: denseToSlice(kf.State.covariance.SliceSym(0, 2)),
		ActualPosition:      vecToSlice(s.Puck.position),
		Time:                t,
	}
}

func SendMessage(m Message, conn *Connection) error {
	conn.mux.Lock()
	defer conn.mux.Unlock()
	if err := conn.connection.WriteJSON(m); err != nil {
		return err
	}
	return nil
}

func MakeHandler(td float64) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UTC().UnixNano())
		s, err := NewScreen(10, 10)
		if err != nil {
			log.Fatal(err)
		}
		c := make(chan time.Time)
		bc := make(chan Edge)
		go s.Run(td, c, bc)
		initialMeasurement := s.GetNoisyState()
		kf, err := NewKalmanFilter(initialMeasurement, td)
		if err != nil {
			log.Fatal(err)
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		connl := &Connection{connection: conn}
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for {
				mt, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
					return
				}
				if (mt == websocket.TextMessage) &&
					(string(msg) == "update") {
					measure := s.GetNoisyState()
					um := NewUpdateMessage(measure)
					if err := SendMessage(um, connl); err != nil {
						log.Fatal(err)
					}
					if err := kf.Update(measure); err != nil {
						log.Fatal(err)
					}
				}

			}
		}()
		t := <-c
		for {
			time.Sleep(time.Duration(td/.001) * time.Millisecond)
			xvel := kf.State.mean.AtVec(2)
			yvel := kf.State.mean.AtVec(3)
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
				switch b := b; {
				case (b == Top) || (b == Bottom):
					uk.SetVec(1, 1)
				case (b == Left) || (b == Right):
					uk.SetVec(0, 1)
				}
			case t = <-c:
			}

			if err = kf.Predict(Bk, uk); err != nil {
				log.Fatal(err)
			}
			msg := NewStateMessage(kf, s, t)
			if err := SendMessage(msg, connl); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// TODO handle socket.close()
