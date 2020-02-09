package main

import (
    "image/color"
    "math"
	"github.com/faiface/pixel"
)

type Bouncing struct {
    OrigPos *Point
    Pos *Point
    Parts int
    Theta float64
    Phi float64
    V float64
    Decay int
    Step float64
    Color color.RGBA
}

func NewBouncing(pos *Point, parts int, theta, phi, v float64, decay int, col color.RGBA) *Bouncing {
    return &Bouncing{pos, pos, parts, theta, phi, v, decay, 0, col}
}

func (self *Bouncing) Draw(win pixel.Target, ob *Observer, dt float64) {
    z := ob.Gravity*math.Pow(self.Step, 2) + self.V*self.Step + self.OrigPos[2]
    if z < 0 {
        z = 0
        self.Step = 0
    }
    newPos := &Point{self.OrigPos[0], self.OrigPos[1], z}
    self.Pos = newPos
    c := NewSphere(newPos, 0.5, self.Color)
    self.Step += dt
    c.Draw(win, ob)
}

func (self *Bouncing) Dist(ob *Observer) float64 {
    d := Distance(ob.Pos, self.Pos)
    return d
}
