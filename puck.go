package main

import (
	"log"

	"gonum.org/v1/gonum/mat"
)

type Edge int

const (
	Top Edge = iota
	Right
	Bottom
	Left
)

type Puck struct {
	position      *mat.VecDense
	velocity      *mat.VecDense
	positionNoise *MultivariateGaussian
	velocityNoise *MultivariateGaussian
}

func (p *Puck) UpdatePosition(timestep float64) {
	p.position.AddScaledVec(p.position, timestep, p.velocity)
}

func (p *Puck) EdgeCollide(e Edge) {
	switch e {
	case Top, Bottom:
		vy := -1 * p.velocity.AtVec(1)
		p.velocity.SetVec(1, vy)
	case Right, Left:
		vx := -1 * p.velocity.AtVec(0)
		p.velocity.SetVec(0, vx)
	default:
		p.velocity.ScaleVec(-1, p.velocity)
	}
}

func (p *Puck) GetNoisyPosition() mat.VecDense {
	n := p.position.Len()
	noise, err := p.positionNoise.Sample()
	if err != nil {
		log.Fatal("error generating position noise")
	}
	pos := mat.NewVecDense(n, make([]float64, n))
	pos.AddVec(p.position, noise)
	return *pos
}

func (p *Puck) GetNoisyVelocity() mat.VecDense {
	n := p.position.Len()
	noise, err := p.velocityNoise.Sample()
	if err != nil {
		log.Fatal("error generating position noise")
	}
	pos := mat.NewVecDense(n, make([]float64, n))
	pos.AddVec(p.velocity, noise)
	return *pos
}

func (p *Puck) GetNoisyState() *mat.VecDense {
	pos := p.GetNoisyPosition()
	vel := p.GetNoisyVelocity()
	n := pos.Len()
	m := vel.Len()
	data := make([]float64, n+m)
	for i := 0; i < n; i++ {
		data[i] = pos.AtVec(i)
	}
	for i := 0; i < m; i++ {
		data[i+n] = vel.AtVec(i)
	}
	return mat.NewVecDense(n+m, data)
}
