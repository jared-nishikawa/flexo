package main

import (
    "image/color"
	//"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
)

type Prism struct {
    Lines []*Line
}

func NewPrism(r, s, t float64, P *Point, col color.RGBA) *Prism {
    x := P[0]
    y := P[1]
    z := P[2]
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

func NewCube(s float64, P *Point, col color.RGBA) *Prism {
    return NewPrism(s, s, s, P, col)
}


//func (self *Prism) Draw(win *pixelgl.Window, ob *Observer) {
func (self *Prism) Draw(win pixel.Target, ob *Observer) {
    for _,line := range self.Lines {
        line.Draw(win, ob)
    }
}

