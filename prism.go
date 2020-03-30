package main

import (
    "image/color"
	//"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
)

type SolidPrism struct {
    Polygons []*Polygon
}

// P will be the exact center
func NewSolidPrism(r, s, t float64, P *Point, col *color.RGBA) *SolidPrism {
    x := P[0] - r/2
    y := P[1] - s/2
    z := P[2] - t/2
    points := []*Point{
        &Point{x,y,z},
        &Point{x+r,y,z},
        &Point{x+r,y+s,z},
        &Point{x,y+s,z},
        &Point{x,y,z+t},
        &Point{x+r,y,z+t},
        &Point{x+r,y+s,z+t},
        &Point{x,y+s,z+t}}
    faces := [][]*Point{
        []*Point{
            points[0], points[1], points[2], points[3]},
        []*Point{
            points[0], points[1], points[5], points[4]},
        []*Point{
            points[1], points[2], points[6], points[5]},
        []*Point{
            points[2], points[3], points[7], points[6]},
        []*Point{
            points[3], points[0], points[4], points[7]},
        []*Point{
            points[4], points[5], points[6], points[7]},
        }
    polygons := []*Polygon{}
    for _,face := range faces {
        poly := NewPolygon(face, col)
        polygons = append(polygons, poly)
    }
    return &SolidPrism{polygons}
}

func NewSolidCube(s float64, P *Point, col *color.RGBA) *SolidPrism {
    return NewSolidPrism(s, s, s, P, col)
}
func (self *SolidPrism) Draw(win pixel.Target, ob *Observer, dt float64) {
    for _,poly := range self.Polygons {
        poly.Draw(win, ob, dt)
    }
}

type Prism struct {
    Lines []*Line
}

// P will be the exact center
func NewPrism(r, s, t float64, P *Point, col *color.RGBA) *Prism {
    x := P[0] - r/2
    y := P[1] - s/2
    z := P[2] - t/2
    bottom := []*Point{
        &Point{x,y,z},
        &Point{x+r,y,z},
        &Point{x+r,y+s,z},
        &Point{x,y+s,z}}

    top := []*Point{
        &Point{x,y,z+t},
        &Point{x+r,y,z+t},
        &Point{x+r,y+s,z+t},
        &Point{x,y+s,z+t}}

    lines := []*Line{}
    for i:=0;i<4;i++ {
        l1 := NewLine(bottom[i], top[i], col)
        l2 := NewLine(top[i], top[(i+1)%4], col)
        l3 := NewLine(bottom[i], bottom[(i+1)%4], col)
        lines = append(lines, l1)
        lines = append(lines, l2)
        lines = append(lines, l3)
    }

    return &Prism{lines}
}

func NewCube(s float64, P *Point, col *color.RGBA) *Prism {
    return NewPrism(s, s, s, P, col)
}

//func (self *Prism) Draw(win *pixelgl.Window, ob *Observer) {
func (self *Prism) Draw(win pixel.Target, ob *Observer, dt float64) {
    for _,line := range self.Lines {
        line.Draw(win, ob, dt)
    }
}

type BorderedPrism struct {
    *Prism
    *SolidPrism
}

// P will be the exact center
func NewBorderedPrism(r, s, t float64, P *Point, faceCol, borderCol *color.RGBA) *BorderedPrism {
    pr := NewPrism(r, s, t, P, borderCol)
    sr := NewSolidPrism(r, s, t, P, faceCol)
    return &BorderedPrism{pr, sr}
}

func (self *BorderedPrism) Draw(win pixel.Target, ob *Observer, dt float64) {
    self.SolidPrism.Draw(win, ob, dt)
    self.Prism.Draw(win, ob, dt)
}

func NewBorderedCube(s float64, P *Point, faceCol, borderCol *color.RGBA) *BorderedPrism {
    return NewBorderedPrism(s, s, s, P, faceCol, borderCol)
}

