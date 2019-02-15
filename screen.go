package main

import (
	"log"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

type Screen struct {
	*Puck
	ScreenWidth   float64
	ScreenHeight  float64
	positionNoise BivariateGaussian
	velocityNoise BivariateGaussian
}

func (s *Screen) Run() {
	for {
		// check for collisions
		x := s.Puck.position.AtVec(0)
		y := s.Puck.position.AtVec(1)
		switch {
		case x >= s.ScreenWidth:
			s.Puck.EdgeCollide(Right)
		case x <= float64(0):
			s.Puck.EdgeCollide(Left)
		}

		switch {
		case y >= s.ScreenHeight:
			s.Puck.EdgeCollide(Bottom)
		case y <= float64(0):
			s.Puck.EdgeCollide(Top)
		}

		s.Puck.UpdatePosition(timeStep)
	}
}

func (s *Screen) GetNoisyPosition() mat.VecDense {
	noise, err := s.positionNoise.Sample()
	if err != nil {
		log.Fatal("error generating position noise")
	}
	var receiver mat.VecDense
	receiver.AddVec(s.Puck.position, noise)
	return receiver
}

func (s *Screen) GetNoisyVelocity() mat.VecDense {
	noise, err := s.velocityNoise.Sample()
	if err != nil {
		log.Fatal("error generating position noise")
	}
	var receiver mat.VecDense
	receiver.AddVec(s.Puck.velocity, noise)
	return receiver
}

//NewScreen is a factory for Screen objects
func NewScreen(height float64, width float64) *Screen {
	s := &Screen{}
	s.ScreenHeight = height
	s.ScreenWidth = width
	initialXPos := rand.Float64() * width
	initialYPos := rand.Float64() * height
	initialPosition := mat.NewVecDense(2, []float64{initialXPos, initialYPos})
	initialVelocity := mat.NewVecDense(2, []float64{1.0, 1.0})
	s.Puck = &Puck{initialPosition, initialVelocity}
	return s
}
