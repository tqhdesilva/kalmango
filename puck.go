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
		var newYVelocity float64
		newYVelocity = -1 * p.velocity.AtVec(1)
		p.velocity.SetVec(1, newYVelocity)
	case Right, Left:
		var newXVelocity float64
		newXVelocity = -1 * p.velocity.AtVec(0)
		p.velocity.SetVec(0, newXVelocity)
	default:
		p.velocity.ScaleVec(-1, p.velocity)
	}
}

func (p *Puck) GetNoisyPosition() mat.VecDense {
	noise, err := p.positionNoise.Sample()
	if err != nil {
		log.Fatal("error generating position noise")
	}
	var receiver mat.VecDense
	receiver.AddVec(p.position, noise)
	return receiver
}

func (p *Puck) GetNoisyVelocity() mat.VecDense {
	noise, err := p.velocityNoise.Sample()
	if err != nil {
		log.Fatal("error generating position noise")
	}
	var receiver mat.VecDense
	receiver.AddVec(p.velocity, noise)
	return receiver
}
