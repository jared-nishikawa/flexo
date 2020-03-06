package main

import (
    "image/color"
	"github.com/faiface/pixel"
)

const WALLHEIGHT = 10.0

type Wall struct {
    Area    *Polygon
    Border  *WallBorder
}

func NewWall(x1, y1, x2, y2, h float64) *Wall {
    if h == 0 {
        h = WALLHEIGHT
    }
    p1 := &Point{x1, y1, 0.0}
    p2 := &Point{x1, y1, h}
    p3 := &Point{x2, y2, h}
    p4 := &Point{x2, y2, 0.0}
    points := []*Point{p1, p2, p3, p4}
    col := color.RGBA{0xff, 0xff, 0xff, 0xff}
    poly := NewPolygon(points, col)
    border := NewWallBorder(p1, p2, p3, p4)
    return &Wall{
        Area: poly,
        Border: border}
}

func (self *Wall) Draw(win pixel.Target, ob *Observer, dt float64) {
    self.Area.Draw(win, ob, dt)
    self.Border.Draw(win, ob, dt)
}

func (self *Wall) Dist(ob *Observer) float64 {
    return self.Area.Dist(ob)
}

type WallBorder struct {
    Lines   []*Line
}

func NewWallBorder(p1, p2, p3, p4 *Point) *WallBorder {
    col := color.RGBA{0x0, 0x0, 0x0, 0xff}
    l1 := NewLine(p1, p2, col)
    l2 := NewLine(p2, p3, col)
    l3 := NewLine(p3, p4, col)
    l4 := NewLine(p4, p1, col)
    return &WallBorder{Lines: []*Line{l1, l2, l3, l4}}
}

func (self *WallBorder) Draw(win pixel.Target, ob *Observer, dt float64) {
    for _,line := range self.Lines {
        line.Draw(win, ob, dt)
    }
}
