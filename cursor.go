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
    Radius float64
    Thickness float64
    Color color.RGBA
}

func NewCursor(x, y, r, t float64, col color.RGBA) *Cursor {
    return &Cursor{x, y, r, t, col}
}

func (self *Cursor) Draw(win *pixelgl.Window) {
    imd := imdraw.New(nil)
    imd.Color = self.Color
    imd.Push(pixel.V(self.X, self.Y))
    imd.Circle(self.Radius, self.Thickness)
    imd.Draw(win)
}


