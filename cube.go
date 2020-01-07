package main

import (
    "image/color"
	"github.com/faiface/pixel/pixelgl"
)
type Cube struct {
    Lines []*Line
}

func NewCube(s float64, P *Point, col color.RGBA) *Cube {
    x := P[0]
    y := P[1]
    z := P[2]
    bottom := []*Point{
        &Point{x,y,z},
        &Point{x+s,y,z},
        &Point{x+s,y+s,z},
        &Point{x,y+s,z}}

    top := []*Point{
        &Point{x,y,z+s},
        &Point{x+s,y,z+s},
        &Point{x+s,y+s,z+s},
        &Point{x,y+s,z+s}}

    lines := []*Line{}
    for i:=0;i<4;i++ {
        v := NewLine(bottom[i], top[i], col)
        t := NewLine(top[i], top[(i+1)%4], col)
        b := NewLine(bottom[i], bottom[(i+1)%4], col)
        lines = append(lines, v)
        lines = append(lines, t)
        lines = append(lines, b)
    }

    return &Cube{lines}
}

func (self *Cube) Draw(win *pixelgl.Window, ob *Observer) {
    for _,line := range self.Lines {
        line.Draw(win, ob)
    }
}

