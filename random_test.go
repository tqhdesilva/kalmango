package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestSample(t *testing.T) {
	tables := []struct {
		length     int
		covariance *mat.SymDense
	}{
		{2, mat.NewSymDense(2, []float64{1, 0, 0, 1})},
		{3, mat.NewSymDense(3, []float64{1, 0, .4, .1, 1, 0, 0, 0, 1})},
	}
	for _, table := range tables {
		distribution := MultivariateGaussian{table.covariance}
		result, err := distribution.Sample()
		if err != nil {
			t.Error("error in sampling")
		}
		if result.Len() != table.length {
			t.Errorf("length of result vector is wrong length, got: %d, expected: %d", result.Len(), table.length)
		}
	}
}
