package main

import (
    "image/color"
    "math"
	"github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
)

type Sphere struct {
    Center *Point
    Radius float64
    Color color.RGBA
}

func NewSphere(c *Point, r float64, col color.RGBA) *Sphere {
    return &Sphere{c, r, col}
}

func (self *Sphere) Draw(win pixel.Target, ob *Observer) {
    //if !self.InView(ob) {
    //    return
    //}
    imd := imdraw.New(nil)

    relative_C := Snap(self.Center, ob.Pos, ob.Theta, ob.Phi)
    orth := &Point{relative_C[1], -relative_C[0], 0.0}
    m := self.Radius/Magnitude(orth)
    orth = &Point{orth[0]*m, orth[1]*m}
    edge := Add(orth, relative_C)

    rh1,th1,ph1 := RecToSphere(relative_C)
    if th1 > math.Pi/2 || th1 < -math.Pi/2 {
        return
    }
    x1,y1 := ob.Project(th1,ph1)

    _,th2,ph2 := RecToSphere(edge)
    x2,y2 := ob.Project(th2,ph2)

    r := math.Sqrt(math.Pow(x1-x2,2) + math.Pow(y1-y2,2))

    imd.Color = self.Color

    /*
    if d < r {
        imd.Color = colornames.Lightblue
        self.Active = true
    } else {
        imd.Color = colornames.Navy
        self.Active = false
    }
    */

    imd.Push(pixel.V(x1, y1))
    _ = rh1
    imd.Circle(r, 0)

    imd.Draw(win)

}
