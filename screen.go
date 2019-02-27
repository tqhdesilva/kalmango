package main

import (
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

type Screen struct {
	*Puck
	ScreenWidth  float64
	ScreenHeight float64
}

// const tickDuration time.Duration = 100 * time.Millisecond
// const timeStep float64 = .1

func (s *Screen) Run(timeStep float64, c chan time.Time, b chan Edge) {
	var tickDuration time.Duration = time.Duration(timeStep*1000) * time.Millisecond
	tick := time.Tick(tickDuration)
	for {
		select {
		case timeStamp := <-tick:
			func() {
				x := s.Puck.position.AtVec(0)
				y := s.Puck.position.AtVec(1)
				switch {
				case x >= s.ScreenWidth:
					s.Puck.EdgeCollide(Right)
					defer func() { b <- Right }()
				case x <= float64(0):
					s.Puck.EdgeCollide(Left)
					defer func() { b <- Left }()
				}

				switch {
				case y >= s.ScreenHeight:
					s.Puck.EdgeCollide(Bottom)
					defer func() { b <- Bottom }()
				case y <= float64(0):
					s.Puck.EdgeCollide(Top)
					defer func() { b <- Top }()
				}

				s.Puck.UpdatePosition(timeStep)
				c <- timeStamp
			}()
		}
	}
}

//NewScreen is a factory for Screen objects
func NewScreen(height float64, width float64) (*Screen, error) {
	s := &Screen{}
	s.ScreenHeight = height
	s.ScreenWidth = width
	initialXPos := rand.Float64() * width
	initialYPos := rand.Float64() * height
	initialPosition := mat.NewVecDense(2, []float64{initialXPos, initialYPos})
	initialVelocity := mat.NewVecDense(2, []float64{1.0, 1.0})
	posNoise, err := NewMultivariateGaussian(mat.NewSymDense(2, []float64{1.0, 0.0, 0.0, 1.0}))
	if err != nil {
		return nil, err
	}
	velNoise, err := NewMultivariateGaussian(mat.NewSymDense(2, []float64{1.0, 0.0, 0.0, 1.0}))
	if err != nil {
		return nil, err
	}
	s.Puck = &Puck{
		initialPosition,
		initialVelocity,
		posNoise,
		velNoise,
	}
	return s, nil
}
