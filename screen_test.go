package main

import (
	"testing"
	"time"

	"gonum.org/v1/gonum/mat"
)

func TestNewScreen(t *testing.T) {
	screen := NewScreen(10.0, 10.0)
	if screen.Puck.position.AtVec(0) > 10.0 ||
		screen.Puck.position.AtVec(0) < 0 ||
		screen.Puck.position.AtVec(1) < 0 ||
		screen.Puck.position.AtVec(1) > 10 {
		t.Error("Puck out of bounds")
	}
}

func TestRun(t *testing.T) {
	screen := NewScreen(10.0, 10.0)
	startPosition := screen.Puck.position
	expectedPosition := mat.VecDense{}
	expectedPosition.AddVec(mat.NewVecDense(2, []float64{.1, .1}), startPosition)
	go screen.Run()
	time.Sleep(110 * time.Millisecond)
	//sometimes this breaks, race condition
	if !mat.Equal(screen.Puck.position, &expectedPosition) {
		t.Errorf("Position was incorrect, got: %+v, expected: %+v", screen.Puck.position, &expectedPosition)
	}
}
