// https://www.bzarg.com/p/how-a-kalman-filter-works-in-pictures/
// TODO add control vector and matrix

package main

import (
	"errors"

	"gonum.org/v1/gonum/mat"
)

type CovMat struct {
	*mat.SymDense
}

func (cm *CovMat) FromDense(d mat.Dense) error {
	// TODO implement this
	t := d.T()
	if !mat.Equal(&d, t) {
		return errors.New("can't convert non-symmetric Dense matrix to SymDense")
	}
	n, _ := d.Dims()
	rowColumnView := make([]float64, n*n)
	for i := 0; i < n; i++ {
		row := d.RawRowView(i)
		for j, k := range row {
			rowColumnView[i*n+j] = k
		}
	}
	cm.SymDense = mat.NewSymDense(n, rowColumnView)
	return nil
}

type Sensor struct {
	covariance *CovMat // R_k
}

type State struct {
	covariance *CovMat
	mean       *mat.VecDense // x_hat_k
}

type KalmanFilter struct {
	*Sensor
	*State
	stateToSensor *mat.Dense // H_k
	noise         *mat.Dense // Q_k
	prediction    *mat.Dense // F_k
}

func (k *KalmanFilter) Predict() *State {
	var newMean *mat.VecDense
	var newCovDense *mat.Dense
	var newState *State
	newMean.MulVec(k.prediction, k.State.mean)
	newCovDense.Mul(k.prediction, k.State.covariance)
	newCovDense.Mul(newCovDense, k.prediction.T())

	newState.mean = newMean
	var newCovMat *CovMat
	newCovMat.FromDense(*newCovDense)
	newState.covariance = newCovMat
	return newState
}

func (k *KalmanFilter) Update() {
	// updates state
	k.Sensor.noise.Sample()
}
