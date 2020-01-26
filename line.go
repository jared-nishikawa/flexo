package main

import (
    "image/color"
    "github.com/faiface/pixel/imdraw"
	//"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
)

type Line struct {
    P *Point
    Q *Point
    Color color.RGBA
}

func NewLine(P,Q *Point, col color.RGBA) *Line {
    return &Line{P, Q, col}
}

func (self *Line) Draw(win pixel.Target, ob *Observer) {
    if !(ob.PointInView(self.P) && ob.PointInView(self.Q)) {
        return
    }
    imd := imdraw.New(nil)
    imd.Color = self.Color

    relative_P := Snap(self.P, ob.Pos, ob.Theta, ob.Phi)
    relative_Q := Snap(self.Q, ob.Pos, ob.Theta, ob.Phi)

    rh1,th1,ph1 := RecToSphere(relative_P)
    x1,y1 := ob.Project(th1,ph1)

    rh2,th2,ph2 := RecToSphere(relative_Q)
    x2,y2 := ob.Project(th2,ph2)

    avg := (rh1+rh2) / 2
    imd.Push(pixel.V(x1,y1), pixel.V(x2,y2))
    foo := ob.Width / 9.6
    imd.Line(foo/avg)

    imd.Draw(win)
}


