package main

import "gonum.org/v1/gonum/mat"

type Edge int

const (
	Top Edge = iota
	Right
	Bottom
	Left
)

const timeStep float64 = .001

type Puck struct {
	position *mat.VecDense
	velocity *mat.VecDense
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
