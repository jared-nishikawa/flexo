package main

import (
    "image/color"
    "math"
	"github.com/faiface/pixel/pixelgl"
)

type Bouncing struct {
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
    return &Bouncing{pos, parts, theta, phi, v, decay, 0, col}
}

func (self *Bouncing) Draw(win *pixelgl.Window, ob *Observer, dt float64) {
    z := ob.Gravity*math.Pow(self.Step, 2) + self.V*self.Step + self.Pos[2]
    if z < 0 {
        z = 0
        self.Step = 0
    }
    newPos := &Point{self.Pos[0], self.Pos[1], z}
    c := NewCircle(newPos, 0.5, self.Color)
    self.Step += dt
    c.Draw(win, ob)
}
