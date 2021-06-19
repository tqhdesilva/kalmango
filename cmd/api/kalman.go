// https://www.bzarg.com/p/how-a-kalman-filter-works-in-pictures/
package main

import (
	"errors"
	"sync"

	"gonum.org/v1/gonum/mat"
)

type CovMat struct {
	mat.SymDense
}

func NewCovMat(n int, data []float64) (*CovMat, error) {
	return &CovMat{
		*mat.NewSymDense(n, data),
	}, nil
}

func (cm *CovMat) FromDense(d *mat.Dense) error {
	t := d.T()
	if !mat.EqualApprox(d, t, .0001) {
		return errors.New("can't convert non-symmetric Dense matrix to SymDense")
	}
	n, _ := d.Dims()
	rcv := make([]float64, n*n)
	for i := 0; i < n; i++ {
		row := d.RawRowView(i)
		for j, k := range row {
			rcv[i*n+j] = k
		}
	}
	cm.SymDense = *mat.NewSymDense(n, rcv)
	return nil
}

type Sensor struct {
	covariance CovMat // R_k
}

type State struct {
	covariance CovMat       // P_k
	mean       mat.VecDense // x_hat_k
	m          sync.Mutex
}

type KalmanFilter struct {
	Sensor
	State         *State
	stateToSensor mat.Dense // H_k
	noise         CovMat    // Q_k TODO add this to model
	prediction    mat.Dense // F_k
}

func (k *KalmanFilter) Predict(Bk *mat.Dense, uk *mat.VecDense) error {
	k.State.m.Lock()
	defer k.State.m.Unlock() // fatal error: sync: unlock of unlocked mutex
	n := k.State.mean.Len()
	x := mat.NewVecDense(n, make([]float64, n))
	pd := mat.NewDense(n, n, make([]float64, n*n))
	_, cols := k.prediction.Dims()
	if cols != k.State.mean.Len() {
		return errors.New("incorrect dimensions")
	}
	if &k.prediction == nil {
		return errors.New("prediction is nil")
	}
	if &k.State == nil {
		return errors.New("state is nil")
	}
	if &k.State.mean == nil {
		return errors.New("state mean is nil")
	}
	control := mat.NewVecDense(n, make([]float64, n))
	control.MulVec(Bk, uk)
	x.MulVec(&k.prediction, &k.State.mean)
	x.AddVec(x, control)

	pd.Mul(&k.prediction, &k.State.covariance)
	pd.Mul(pd, k.prediction.T())

	p, err := NewCovMat(n, make([]float64, n*n))
	if err != nil {
		return err
	}
	err = p.FromDense(pd)
	if err != nil {
		return err
	}
	k.State.covariance = *p
	k.State.mean = *x
	return nil
}

func (k *KalmanFilter) Update(measurement *mat.VecDense) error {
	// calculate K'
	k.State.m.Lock()
	defer k.State.m.Unlock()
	n := k.State.mean.Len()
	kg := mat.NewDense(n, n, make([]float64, n*n))
	kg.Mul(&k.stateToSensor, &k.State.covariance)
	kg.Mul(kg, k.stateToSensor.T())
	kg.Add(kg, &k.Sensor.covariance)
	kg.Inverse(kg)
	kg.Mul(k.stateToSensor.T(), kg)
	kg.Mul(&k.State.covariance, kg) // P_k is 0

	// calculate hat x'_k
	x := mat.NewVecDense(n, make([]float64, n))
	x.MulVec(&k.stateToSensor, &k.State.mean)
	x.SubVec(measurement, x)
	x.MulVec(kg, x)
	x.AddVec(&k.State.mean, x)

	// calculate p'_k
	pd := mat.NewDense(n, n, make([]float64, n*n))
	pd.Mul(&k.stateToSensor, &k.State.covariance)
	pd.Mul(kg, pd)
	pd.Sub(&k.State.covariance, pd)
	p, err := NewCovMat(n, make([]float64, n*n))
	if err != nil {
		return err
	}
	err = p.FromDense(pd)
	if err != nil {
		return err
	}

	k.State = &State{
		covariance: *p,
		mean:       *x}
	return nil
}

func matchShape(size int, m mat.Matrix) bool {
	r, c := m.Dims()
	return (r == c || c == 1) && r == size
}

func NewKalmanFilter(
	z *mat.VecDense,
	H *mat.Dense,
	Q *CovMat,
	R *CovMat,
	F *mat.Dense,
	td float64,
) (*KalmanFilter, error) {
	r, _ := z.Dims()
	if !matchShape(r, H) {
		return nil, errors.New("Invalid shape of H.")
	}
	if !matchShape(r, Q) {
		return nil, errors.New("Invalid shape of Q.")
	}
	if !matchShape(r, R) {
		return nil, errors.New("Invalid shape of R.")
	}
	if !matchShape(r, F) {
		return nil, errors.New("Invalid shape of F.")
	}

	sensor := Sensor{*R}

	// TODO initially ought to be the same as r
	p, err := NewCovMat(4, []float64{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 1.0, 0.0, 1.0,
	})
	if err != nil {
		return nil, err
	}

	s := &State{
		covariance: *p,
		mean:       *z,
	}

	kf := &KalmanFilter{
		Sensor:        sensor,
		State:         s,
		stateToSensor: *H,
		noise:         *Q,
		prediction:    *F,
	}
	return kf, nil
}
