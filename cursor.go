package main

import (
    "image/color"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
    "github.com/faiface/pixel/imdraw"
)

type Cursor struct {
    X float64
    Y float64
    MaxX float64
    MaxY float64
    Multiplier float64
    Radius float64
    Thickness float64
    Color color.RGBA
    Active byte
    Inactive byte
}

func NewCursor(x, y, maxX, maxY, m, r, t float64, col color.RGBA, active, inactive byte) *Cursor {
    return &Cursor{x, y, maxX, maxY, m, r, t, col, active, inactive}
}

func (self *Cursor) SetActive() {
    self.Color.A = self.Active
}

func (self *Cursor) SetInactive() {
    self.Color.A = self.Inactive
}

func (self *Cursor) Move(dx, dy, dt float64) {
    self.X += dx*dt*self.Multiplier
    self.Y += dy*dt*self.Multiplier
    if self.X > self.MaxX {
        self.X = self.MaxX
    }
    if self.X < 0 {
        self.X = 0
    }
    if self.Y > self.MaxY {
        self.Y = self.MaxY
    }
    if self.Y < 0 {
        self.Y = 0
    }
}

func (self *Cursor) Draw(win *pixelgl.Window) {
    imd := imdraw.New(nil)
    imd.Color = self.Color
    imd.Push(pixel.V(self.X, self.Y))
    imd.Circle(self.Radius, self.Thickness)
    imd.Draw(win)
}


