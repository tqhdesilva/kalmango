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
	return rand.NormFloat64()*math.Sqrt(grvs.variance) + grvs.mean
}

type MultivariateGaussian struct {
	covariance *mat.SymDense
}

func (bg *MultivariateGaussian) Sample() (mat.Vector, error) {
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
