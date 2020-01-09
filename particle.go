package main

import (
    "image/color"
    "math"
    "math/rand"
	"github.com/faiface/pixel/pixelgl"
)

type Particle struct {
    Pos *Point
    Theta float64
    Phi float64
    V float64
    Decay float64
    Step float64
    Color color.RGBA
}

func NewParticle(pos *Point, theta, phi, v, decay, step float64, col color.RGBA) *Particle {
    return &Particle{pos, rand.Float64()*math.Pi*2, phi, v, decay, step, col}
}

func (self *Particle) Draw(win *pixelgl.Window, ob *Observer, dt float64) {
    dx := self.Step*math.Cos(self.Theta)*math.Sqrt(self.V)
    dy := self.Step*math.Sin(self.Theta)*math.Sqrt(self.V)
    dz := 0.001*ob.Gravity*math.Pow(self.Step, 2) + self.V*self.Step + self.Pos[2]
    x := self.Pos[0]
    y := self.Pos[1]
    z := self.Pos[2]

    if self.Step > self.Decay {
        dx = 0
        dy = 0
        dz = 0
        self.Theta = rand.Float64()*math.Pi*2
        self.Step = 0
    }
    newPos := &Point{x+dx, y+dy, z+dz}
    c := NewSphere(newPos, 0.5, self.Color)
    self.Step += dt
    c.Draw(win, ob)
}


