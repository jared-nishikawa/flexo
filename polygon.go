package main

import (
    "image/color"
    "github.com/faiface/pixel/imdraw"
	//"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
)

type Polygon struct {
    Points []*Point
    Color color.RGBA
}

func NewPolygon(points []*Point, col color.RGBA) *Polygon {
    return &Polygon{points, col}
}

func (self *Polygon) Draw(win pixel.Target, ob *Observer) {
    for _,p := range self.Points {
        if !(ob.PointInView(p)) {
            return
        }
    }

    imd := imdraw.New(nil)
    imd.Color = self.Color

    sumDists := 0.0
    numPoints := 0
    for _,p := range self.Points {
        relative_p := Snap(p, ob.Pos, ob.Theta, ob.Phi)
        rh, th, ph := RecToSphere(relative_p)
        x, y := ob.Project(th, ph)
        sumDists += rh
        numPoints += 1
        imd.Push(pixel.V(x,y))
    }

    imd.Polygon(0)

    imd.Draw(win)
}


