package main

import (
	"testing"
	"time"

	"gonum.org/v1/gonum/mat"
)

func TestNewScreen(t *testing.T) {
	posCov, err := NewCovMat(2, []float64{
		1.0, 0.0,
		0.0, 1.0,
	})
	velCov, err := NewCovMat(2, []float64{
		0.2, 0.0,
		0.0, 0.2,
	})
	if err != nil {
		t.Error(err)
	}
	screen, err := NewScreen(10.0, 10.0, posCov, velCov)
	if err != nil {
		t.Error(err)
	}
	if screen.Puck.position.AtVec(0) > 10.0 ||
		screen.Puck.position.AtVec(0) < 0 ||
		screen.Puck.position.AtVec(1) < 0 ||
		screen.Puck.position.AtVec(1) > 10 {
		t.Error("Puck out of bounds")
	}
}

func TestRun(t *testing.T) {
	// TODO sometimes this breaks, race condition
	posCov, err := NewCovMat(2, []float64{
		1.0, 0.0,
		0.0, 1.0,
	})
	velCov, err := NewCovMat(2, []float64{
		0.2, 0.0,
		0.0, 0.2,
	})
	screen, err := NewScreen(10.0, 10.0, posCov, velCov)
	if err != nil {
		t.Error(err)
	}
	screen.Puck.position.SetVec(0, 5.0)
	screen.Puck.position.SetVec(1, 5.0)
	startPosition := screen.Puck.position
	expectedPosition := mat.VecDense{}
	expectedPosition.AddVec(mat.NewVecDense(2, []float64{.1, .1}), startPosition)
	c := make(chan time.Time)
	bc := make(chan Edge)
	go screen.Run(.1, c, bc)
	<-c
	select {
	case <-bc:
		t.Error("Bounced in the middle")
	default:
	}
	if !mat.Equal(screen.Puck.position, &expectedPosition) {
		t.Errorf("Position was incorrect, got: %+v, expected: %+v", screen.Puck.position, &expectedPosition)
	}
}

func TestBounce(t *testing.T) {
	posCov, err := NewCovMat(2, []float64{
		1.0, 0.0,
		0.0, 1.0,
	})
	velCov, err := NewCovMat(2, []float64{
		0.2, 0.0,
		0.0, 0.2,
	})
	screen, err := NewScreen(10.0, 10.0, posCov, velCov)
	if err != nil {
		t.Error(err)
	}
	screen.Puck.position.SetVec(0, 10.5)
	screen.Puck.position.SetVec(1, 5.0)
	screen.Puck.velocity.SetVec(0, .1)
	screen.Puck.velocity.SetVec(1, .1)
	c := make(chan time.Time)
	bc := make(chan Edge)
	go screen.Run(.1, c, bc)
	<-c
	select {
	case e := <-bc:
		if e != Right {
			t.Errorf("bounced on the wrong side: %+v", e)
		}
	case <-c:
		t.Error("didn't send bounce in channel")
	}

}
