package main

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestPredict(t *testing.T) {
	const timeDelta float64 = .1
	sensor := &Sensor{
		&CovMat{mat.NewSymDense(2, []float64{1.0, 0.0, 0.0, 1.0})},
	}
	state := &State{
		&CovMat{mat.NewSymDense(2, []float64{1.0, 0.0, 0.0, 1.0})},
		mat.NewVecDense(2, []float64{5.0, 1.0}),
	}
	stateToSensor := mat.NewDense(2, 2, []float64{1.0, 0.0, 0.0, 1.0})
	noise := &CovMat{mat.NewSymDense(2, []float64{0.0, 0.0, 0.0, 0.0})}
	prediction := mat.NewDense(2, 2, []float64{1.0, timeDelta, 0.0, 1.0})
	kf := &KalmanFilter{
		sensor,
		state,
		stateToSensor,
		noise,
		prediction,
	}
	if kf.State == nil {
		t.Error("state is nil")
	}
	if kf.prediction == nil {
		t.Error("prediction is nil")
	}
	oldState := *kf.State
	kf.Predict()
	if mat.Equal(oldState.mean, kf.State.mean) {
		t.Error("mean state is still the same")
	}
}

func TestUpdate(t *testing.T) {
	const timeDelta float64 = .1
	sensor := &Sensor{
		&CovMat{mat.NewSymDense(2, []float64{1.0, 0.0, 0.0, 1.0})},
	}
	state := &State{
		&CovMat{mat.NewSymDense(2, []float64{1.0, 0.0, 0.0, 1.0})},
		mat.NewVecDense(2, []float64{5.0, 1.0}),
	}
	stateToSensor := mat.NewDense(2, 2, []float64{1.0, 0.0, 0.0, 1.0})
	noise := &CovMat{mat.NewSymDense(2, []float64{0.0, 0.0, 0.0, 0.0})}
	prediction := mat.NewDense(2, 2, []float64{1.0, timeDelta, 0.0, 1.0})
	kf := &KalmanFilter{
		sensor,
		state,
		stateToSensor,
		noise,
		prediction,
	}
	measurement := mat.NewVecDense(2, []float64{4.0, 1.3})
	oldState := *kf.State
	kf.Update(measurement)

	if mat.Equal(oldState.mean, kf.State.mean) {
		t.Error("mean state is still the same")
	}
}
