package main

import (
    "image/color"
	"github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type FlatShape interface {
    GetX() float64
    GetY() float64
    GetDragging() bool
    SetDragging(b bool)
    SetLoc(x, y float64)
    Center() (float64, float64)
    Draw(win *pixelgl.Window)
    Snap()
}

type Rectangle struct {
    Width float64
    Height float64
    X float64
    Y float64
    Thickness float64
    Dragging bool
    Color color.RGBA
}

func NewRectangle(w, h, x, y, t float64, col color.RGBA) *Rectangle {
    return &Rectangle{w, h, x, y, t, false, col}
}

func NewSquare(s, x, y, t float64, col color.RGBA) *Rectangle {
    return NewRectangle(s, s, x, y, t, col)
}

func (self *Rectangle) GetX() float64 {
    return self.X
}

func (self *Rectangle) GetY() float64 {
    return self.Y
}

func (self *Rectangle) GetDragging() bool {
    return self.Dragging
}

func (self *Rectangle) SetDragging(b bool) {
    self.Dragging = b
}

func (self *Rectangle) Center() (float64, float64) {
    return self.X + self.Width/2, self.Y + self.Height/2
}

func (self *Rectangle) Draw(win *pixelgl.Window) {
    w := self.Width
    h := self.Height
    x := self.X
    y := self.Y
    t := self.Thickness
    imd := imdraw.New(nil)
    imd.Color = self.Color

    imd.Push(pixel.V(x,y), pixel.V(x+w,y+h))
    imd.Rectangle(t)

    imd.Draw(win)
}

func (self *Rectangle) SetLoc(x, y float64) {
    self.X = x
    self.Y = y
}

func (self *Rectangle) Snap() {
    self.X = float64(int(self.X/50 + 0.5)*50)
    self.Y = float64(int(self.Y/50 + 0.5)*50)
}


type Circle struct {
    X float64
    Y float64
    Radius float64
    Thickness float64
    Dragging bool
    Color color.RGBA
}

func NewCircle(x, y, r, t float64, col color.RGBA) *Circle {
    return &Circle{x, y, r, t, false, col}
}

func (self *Circle) GetX() float64 {
    return self.X - self.Radius
}

func (self *Circle) GetY() float64 {
    return self.Y - self.Radius
}

func (self *Circle) GetDragging() bool {
    return self.Dragging
}

func (self *Circle) SetDragging(b bool) {
    self.Dragging = b
}

func (self *Circle) Center() (float64, float64) {
    return self.X, self.Y
}

func (self *Circle) SetLoc(x, y float64) {
    x += self.Radius
    y += self.Radius
    self.X = x
    self.Y = y
}


func (self *Circle) Draw(win *pixelgl.Window) {
    imd := imdraw.New(nil)
    imd.Color = self.Color
    imd.Push(pixel.V(self.X,self.Y))
    imd.Circle(self.Radius, self.Thickness)
    imd.Draw(win)
}

func (self *Circle) Snap() {
}

