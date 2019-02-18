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

func (cm *CovMat) FromDense(d *mat.Dense) error {
	// TODO implement this
	t := d.T()
	if !mat.Equal(d, t) {
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
	covariance *CovMat       // P_k
	mean       *mat.VecDense // x_hat_k
}

type KalmanFilter struct {
	*Sensor
	*State
	stateToSensor *mat.Dense // H_k
	noise         *mat.Dense // Q_k
	prediction    *mat.Dense // F_k
}

func (k *KalmanFilter) Predict() {
	var newMean *mat.VecDense
	var newCovDense *mat.Dense
	var newState *State
	newMean.MulVec(k.prediction, k.State.mean)
	newCovDense.Mul(k.prediction, k.State.covariance)
	newCovDense.Mul(newCovDense, k.prediction.T())

	newState.mean = newMean
	var newCovMat *CovMat
	newCovMat.FromDense(newCovDense)
	newState.covariance = newCovMat
	*k.State = *newState
}

func (k *KalmanFilter) Update(measurement *mat.VecDense) {
	// calculate K'
	// (19)
	var newKalmanGain *mat.Dense
	newKalmanGain.Mul(k.stateToSensor, k.State.covariance)
	newKalmanGain.Mul(newKalmanGain, k.stateToSensor.T())
	newKalmanGain.Add(newKalmanGain, k.Sensor.covariance)
	newKalmanGain.Inverse(newKalmanGain)
	newKalmanGain.Mul(k.stateToSensor.T(), newKalmanGain)
	newKalmanGain.Mul(k.State.covariance, newKalmanGain)

	// calculate hat x'_k
	var newStateMean *mat.VecDense
	newStateMean.MulVec(k.stateToSensor, k.State.mean)
	newStateMean.SubVec(measurement, newStateMean)
	newStateMean.MulVec(newKalmanGain, newStateMean)
	newStateMean.AddVec(k.State.mean, newStateMean)

	var newStateCovarianceDense *mat.Dense
	newStateCovarianceDense.Mul(k.stateToSensor, k.State.covariance)
	newStateCovarianceDense.Mul(newKalmanGain, newStateCovarianceDense)
	newStateCovarianceDense.Sub(k.State.covariance, newStateCovarianceDense)
	var newStateCovariance *CovMat
	newStateCovariance.FromDense(newStateCovarianceDense)

	k.State = &State{newStateCovariance, newStateMean}
}
