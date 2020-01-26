package main

import (
    "image/color"
	"github.com/faiface/pixel"
)

type Fountain struct {
    Pos *Point
    Parts int
    Particles []*Particle
    Theta float64
    Phi float64
    V float64
    Decay float64
    Step float64
}

func NewFountain(pos *Point, parts int, theta, phi, v float64, delay float64, col color.RGBA) *Fountain {
    decay := float64(parts) * delay
    particles := make([]*Particle, parts)
    for i:=0;i<parts;i++ {
        particles[i] = NewParticle(pos, theta, phi, v, decay, float64(i)*delay, col)
    }

    return &Fountain{pos, parts, particles, theta, phi, v, decay, 0}
}

func (self *Fountain) Draw(win pixel.Target, ob *Observer, dt float64) {
    for _,p := range self.Particles {
        p.Draw(win, ob, dt)
    }

}


