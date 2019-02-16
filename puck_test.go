package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestUpdatePosition(t *testing.T) {
	tables := []struct {
		startPosition *mat.VecDense
		velocity      *mat.VecDense
		endPosition   *mat.VecDense
		timestep      float64
	}{
		{
			mat.NewVecDense(2, []float64{0.0, 0.0}),
			mat.NewVecDense(2, []float64{1.0, 2.0}),
			mat.NewVecDense(2, []float64{2.0, 4.0}),
			2.0,
		},
		{
			mat.NewVecDense(2, []float64{-1.0, 3.0}),
			mat.NewVecDense(2, []float64{1.0, 0.0}),
			mat.NewVecDense(2, []float64{0.0, 3.0}),
			1.0,
		},
	}
	for _, table := range tables {
		testPuck := Puck{
			table.startPosition,
			table.velocity,
			&BivariateGaussian{},
			&BivariateGaussian{},
		}
		testPuck.UpdatePosition(table.timestep)
		if !mat.Equal(testPuck.position, table.endPosition) {
			t.Errorf(
				"Position is wrong, got: %+v, expected: %+v",
				testPuck.position,
				table.endPosition,
			)
		}
	}
}

func TestEdgeCollide(t *testing.T) {
	testPuck := Puck{
		mat.NewVecDense(2, []float64{0.0, 0.0}),
		mat.NewVecDense(2, []float64{1.0, 1.0}),
		&BivariateGaussian{},
		&BivariateGaussian{},
	}
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

func TestGetNoisyPosition(t *testing.T) {

}

func TestGetNoisyVelocity(t *testing.T) {

}
