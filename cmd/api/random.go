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
	covariance     *mat.SymDense
	eigenmatrix    *mat.EigenSym
	eigenvalues    []float64
	eigenvectors   *mat.Dense
	standardNormal *GaussianRandomVariable
}

func NewMultivariateGaussian(covariance *mat.SymDense) (*MultivariateGaussian, error) {
	eigenmatrix := &mat.EigenSym{}
	if !eigenmatrix.Factorize(covariance, true) {
		return nil, errors.New("factorization of covariance matrix failed")
	}
	eigenvalues := eigenmatrix.Values(nil)
	eigenvectors := &mat.Dense{}
	eigenmatrix.VectorsTo(eigenvectors)
	standardNormal := &GaussianRandomVariable{0.0, 1.0}
	return &MultivariateGaussian{
		covariance,
		eigenmatrix,
		eigenvalues,
		eigenvectors,
		standardNormal,
	}, nil
}

func (bg *MultivariateGaussian) Sample() (mat.Vector, error) {
	// see https://stackoverflow.com/questions/6142576/sample-from-multivariate-normal-gaussian-distribution-in-c
	// see https://en.wikipedia.org/wiki/Multivariate_normal_distribution#Drawing_values_from_the_distribution
	rows, columns := bg.covariance.Dims()
	if rows != columns {
		return nil, errors.New("covariance matrix is not square")
	}
	iid := make([]float64, rows)
	for i := 0; i < rows; i++ {
		iid[i] = bg.standardNormal.Sample()
	}
	for i := 0; i < rows; i++ {
		iid[i] = iid[i] * math.Sqrt(bg.eigenvalues[i])
	}
	iidVector := mat.NewVecDense(len(iid), iid)
	var r mat.Dense
	r.Mul(bg.eigenvectors, iidVector)
	rVector := r.ColView(0)
	return rVector, nil
}
