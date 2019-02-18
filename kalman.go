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

func NewCovMat(n int, data []float64) (*CovMat, error) {
	return &CovMat{
		mat.NewSymDense(n, data),
	}, nil
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
	noise         *CovMat    // Q_k TODO add this to model
	prediction    *mat.Dense // F_k
}

func (k *KalmanFilter) Predict() error {
	//var newMean *mat.VecDense
	newMean := mat.NewVecDense(2, []float64{0.0, 0.0})
	newCovDense := mat.NewDense(2, 2, []float64{0.0, 0.0, 0.0, 0.0})
	covMat, err := NewCovMat(2, []float64{0.0, 0.0, 0.0, 0.0})
	if err != nil {
		return err
	}
	newState := &State{
		covariance: covMat,
	}
	_, cols := k.prediction.Dims()
	if cols != k.State.mean.Len() {
		return errors.New("incorrect dimensions")
	}
	if k.prediction == nil {
		return errors.New("prediction is nil")
	}
	if k.State == nil {
		return errors.New("state is nil")
	}
	if k.State.mean == nil {
		return errors.New("state mean is nil")
	}
	newMean.MulVec(k.prediction, k.State.mean)
	newCovDense.Mul(k.prediction, k.State.covariance)
	newCovDense.Mul(newCovDense, k.prediction.T())

	newState.mean = newMean
	newCovMat, err := NewCovMat(2, []float64{0.0, 0.0, 0.0, 0.0})
	if err != nil {
		return err
	}
	newCovMat.FromDense(newCovDense)
	newState.covariance = newCovMat
	*k.State = *newState
	return nil
}

func (k *KalmanFilter) Update(measurement *mat.VecDense) error {
	// calculate K'
	// (19)
	newKalmanGain := mat.NewDense(2, 2, make([]float64, 4))
	newKalmanGain.Mul(k.stateToSensor, k.State.covariance)
	newKalmanGain.Mul(newKalmanGain, k.stateToSensor.T())
	newKalmanGain.Add(newKalmanGain, k.Sensor.covariance)
	newKalmanGain.Inverse(newKalmanGain)
	newKalmanGain.Mul(k.stateToSensor.T(), newKalmanGain)
	newKalmanGain.Mul(k.State.covariance, newKalmanGain)

	// calculate hat x'_k
	newStateMean := mat.NewVecDense(2, make([]float64, 2))
	newStateMean.MulVec(k.stateToSensor, k.State.mean)
	newStateMean.SubVec(measurement, newStateMean)
	newStateMean.MulVec(newKalmanGain, newStateMean)
	newStateMean.AddVec(k.State.mean, newStateMean)

	newStateCovarianceDense := mat.NewDense(2, 2, make([]float64, 4))
	newStateCovarianceDense.Mul(k.stateToSensor, k.State.covariance)
	newStateCovarianceDense.Mul(newKalmanGain, newStateCovarianceDense)
	newStateCovarianceDense.Sub(k.State.covariance, newStateCovarianceDense)
	newStateCovariance, err := NewCovMat(2, make([]float64, 4))
	if err != nil {
		return err
	}
	newStateCovariance.FromDense(newStateCovarianceDense)

	k.State = &State{newStateCovariance, newStateMean}
	return nil
}
