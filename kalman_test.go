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
	Bk := mat.NewDense(2, 2, make([]float64, 4))
	uk := mat.NewVecDense(2, make([]float64, 2))
	kf.Predict(Bk, uk)
	if expected := mat.NewVecDense(2, []float64{5.1, 1.0}); !mat.Equal(expected, state.mean) {
		t.Errorf("expected: %+v, got: %+v", expected, state.mean)
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

func TestFromDense(t *testing.T) {
	denseMat := mat.NewDense(4, 4, []float64{
		0.20784313725490197, 0, 0.0392156862745098, 0,
		0, 0.20784313725490197, 0, 0.0392156862745098,
		0.03921568627450981, 0, 0.19607843137254902, 0,
		0, 0.03921568627450981, 0, 0.19607843137254902,
	})
	covMat, _ := NewCovMat(4, make([]float64, 16))
	err := covMat.FromDense(denseMat)
	if err != nil {
		t.Error(err)
	}
	var sum float64
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sum = sum + covMat.At(i, j)
		}
	}
	if sum == 0.0 {
		t.Error("Got zero matrix")
	}
}
