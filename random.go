package main

import (
	"errors"
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

type GaussianRandomVariable struct {
	mean     float64
	variance float64
}

func (grvs *GaussianRandomVariable) Sample() float64 {
	// use box-muller transform from uniform to standard normal
	u1 := rand.Float64()
	u2 := rand.Float64()
	x := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	//scale and shift standard normal
	x = x * math.Sqrt(grvs.variance)
	x = x + grvs.mean

	return x
}

type BivariateGaussian struct {
	covariance *mat.SymDense
}

func (bg *BivariateGaussian) Sample() (mat.Vector, error) {
	// see https://stackoverflow.com/questions/6142576/sample-from-multivariate-normal-gaussian-distribution-in-c
	// see https://en.wikipedia.org/wiki/Multivariate_normal_distribution#Drawing_values_from_the_distribution
	// TODO we don't need to recalculate the eigenvalues eigenvectors each time
	rows, columns := bg.covariance.Dims()
	if rows != columns {
		return nil, errors.New("covariance matrix is not square")
	}
	iidSamples := make([]float64, rows)
	standardNormal := GaussianRandomVariable{0.0, 1.0}
	for i := 0; i < rows; i++ {
		iidSamples[i] = standardNormal.Sample()
	}
	eigenSym := mat.EigenSym{}
	if !eigenSym.Factorize(bg.covariance, true) {
		return nil, errors.New("factorization of covariance matrix failed")
	}
	// do we need to add orthonormality assertion, or is that guaranteed?
	// I think it is guaranteed
	q := mat.Dense{}
	eigenValues := eigenSym.Values(nil)
	q.EigenvectorsSym(&eigenSym)
	for i := 0; i < rows; i++ {
		iidSamples[i] = iidSamples[i] * math.Sqrt(eigenValues[i])
	}
	iidSamplesVector := mat.NewVecDense(len(iidSamples), iidSamples)
	var r mat.Dense
	r.Mul(&q, iidSamplesVector)
	rVector := r.ColView(0)
	return rVector, nil
}
