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

func (c *Connection) ReadMessage() (messageType int, p []byte, err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.connection.ReadMessage()
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
			log.Print(err)
		}
		c := make(chan time.Time)
		bc := make(chan Edge)
		go s.Run(td, c, bc)
		initialMeasurement := s.GetNoisyState()
		kf, err := NewKalmanFilter(initialMeasurement, td)
		if err != nil {
			log.Print(err)
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		readconn := &Connection{connection: conn}
		writeconn := &Connection{connection: conn}
		if err != nil {
			log.Print(err)
		}
		go func() {
			defer conn.Close()
			for {
				mt, msg, err := readconn.ReadMessage()
				if err != nil {
					log.Print(err)
					return
				}
				if (mt == websocket.TextMessage) &&
					(string(msg) == "update") {
					measure := s.GetNoisyState()
					um := NewUpdateMessage(measure)
					if err := SendMessage(um, writeconn); err != nil {
						log.Print(err)
						return
					}
					if err := kf.Update(measure); err != nil {
						log.Print(err)
						return
					}
				}

			}
		}()
		t := <-c
		Bk := mat.NewDense(4, 4, []float64{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, -2, 0,
			0, 0, 0, -2,
		})
		for {
			time.Sleep(time.Duration(td/.001) * time.Millisecond)
			xvel := kf.State.mean.AtVec(2)
			yvel := kf.State.mean.AtVec(3)
			if err != nil {
				log.Print(err)
				return
			}
			uk := mat.NewVecDense(4, make([]float64, 4))
			select {
			case b := <-bc:
				switch b := b; {
				case (b == Top) || (b == Bottom):
					uk.SetVec(3, yvel)
				case (b == Left) || (b == Right):
					uk.SetVec(2, xvel)
				}
			case t = <-c:
			}

			if err = kf.Predict(Bk, uk); err != nil {
				log.Print(err)
				return
			}
			msg := NewStateMessage(kf, s, t)
			if err := SendMessage(msg, writeconn); err != nil {
				log.Print(err)
				return
			}
		}
	}
}
