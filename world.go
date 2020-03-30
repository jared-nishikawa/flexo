package main

import (
	"github.com/faiface/pixel"
)

type Object interface {
    Draw(pixel.Target, *Observer, float64)
}

type World struct {
    Grid    [][][]Object
    Snap    float64
}

func (self *World) SnapPoint(p *Point) *Point {
    c := self.Snap
    x,_ := Align(p[0], c)
    y,_ := Align(p[1], c)
    z,_ := Align(p[2], c)
    return &Point{x,y,z}
}
func (self *World) Convert(p *Point) (int, int, int) {
    l := len(self.Grid)
    c := self.Snap
    _,i := Align(p[0], c)
    _,j := Align(p[1], c)
    _,k := Align(p[2], c)
    x := i + int(l/2)
    y := j + int(l/2)
    z := k + int(l/2)
    if x >= l {
        x = -1
    }
    if y >= l {
        y = -1
    }
    if z >= l {
        z = -1
    }
    return x,y,z
}

func (self *World) Exists(p *Point) bool {
    x,y,z := self.Convert(p)
    return x >= 0 && y >= 0 && z >= 0
}

func (self *World) Set(p *Point, o Object) {
    x,y,z := self.Convert(p)
    self.Grid[x][y][z] = o
}

func (self World) Get(p *Point) Object {
    x,y,z := self.Convert(p)
    return self.Grid[x][y][z]
}

func (self *World) Ray(P *Point, rho, theta, phi float64) *Point {
    Q := SphereToRec(rho, theta, phi)
    T := Add(P, Q)
    return self.SnapPoint(T)
}
