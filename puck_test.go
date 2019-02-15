package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestUpdatePosition(t *testing.T) {
	testPuck := Puck{
		mat.NewVecDense(2, []float64{0.0, 0.0}),
		mat.NewVecDense(2, []float64{1.0, 2.0})}
	testPuck.UpdatePosition(2.0)
	expectedPosition := mat.NewVecDense(2, []float64{2.0, 4.0})
	if !mat.Equal(
		testPuck.position, expectedPosition) {
		t.Errorf(
			"Position is wrong, got: %+v, expected: %+v",
			testPuck.position,
			expectedPosition,
		)
	}
}

func TestEdgeCollide(t *testing.T) {
	testPuck := Puck{
		mat.NewVecDense(2, []float64{0.0, 0.0}),
		mat.NewVecDense(2, []float64{1.0, 1.0})}
	testPuck.EdgeCollide(Top)
	expectedVelocity := mat.NewVecDense(2, []float64{1.0, -1.0})
	if !mat.Equal(testPuck.velocity, expectedVelocity) {
		t.Errorf(
			"Velocity was incorrect, got: %+v, expected: %+v",
			testPuck.velocity,
			expectedVelocity,
		)
	}
}
