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
    Contains(x, y float64) bool
    Rotate()
    Snap()
    Copy() FlatShape
}

type Atom struct {
    Cfg byte
    *Rectangle
}

func (self *Atom) Draw(win *pixelgl.Window) {
    cfg := self.Cfg
    self.Rectangle.Draw(win)
    col := color.RGBA{0x0, 0xff, 0xff, 0xff}
    if cfg & 0x1 != 0 {
        //c1 := NewCircle(self.X + self.Width - self.Width/10, self.Y + self.Height/2, self.Width/10, 0, color.RGBA{0xff, 0xff, 0xff, 0xff})
        c1 := NewCircle(self.X + self.Width, self.Y + self.Height/2, self.Width/10, 0, col)
        c1.Draw(win)
    }
    if cfg & 0x2 != 0 {
        //c2 := NewCircle(self.X + self.Width/2, self.Y + self.Height - self.Width/10, self.Width/10, 0, color.RGBA{0xff, 0xff, 0xff, 0xff})
        c2 := NewCircle(self.X + self.Width/2, self.Y + self.Height, self.Width/10, 0, col)
        c2.Draw(win)
    }
    if cfg & 0x4 != 0 {
        //c3 := NewCircle(self.X + self.Width/10, self.Y + self.Height/2, self.Width/10, 0, color.RGBA{0xff, 0xff, 0xff, 0xff})
        c3 := NewCircle(self.X, self.Y + self.Height/2, self.Width/10, 0, col)
        c3.Draw(win)
    }
    if cfg & 0x8 != 0 {
        //c4 := NewCircle(self.X + self.Width/2, self.Y + self.Width/10, self.Width/10, 0, color.RGBA{0xff, 0xff, 0xff, 0xff})
        c4 := NewCircle(self.X + self.Width/2, self.Y, self.Width/10, 0, col)
        c4.Draw(win)
    }
}

func (self *Atom) Rotate() {
    self.Cfg <<= 1
    c := (self.Cfg >> 4) & 0x1
    self.Cfg |= c
}

type Carbon struct {
    *Atom
}

func (self *Carbon) Draw(win *pixelgl.Window) {
    self.Atom.Draw(win)
}

func (self *Carbon) Copy() FlatShape {
    return NewCarbon(self.Width, self.X, self.Y, self.Thickness, self.Color)
}

func NewCarbon(s, x, y, t float64, col color.RGBA) *Carbon {
    return &Carbon{&Atom{0xff, &Rectangle{s, s, x, y, t, false, col}}}
}

type Hydrogen struct {
    *Atom
}

func (self *Hydrogen) Draw(win *pixelgl.Window) {
    self.Atom.Draw(win)
}

func (self *Hydrogen) Copy() FlatShape {
    return NewHydrogen(self.Width, self.X, self.Y, self.Thickness, self.Color)
}

func NewHydrogen(s, x, y, t float64, col color.RGBA) *Hydrogen {
    return &Hydrogen{&Atom{0xf1, &Rectangle{s, s, x, y, t, false, col}}}
}

type Nitrogen struct {
    *Atom
}

func (self *Nitrogen) Draw(win *pixelgl.Window) {
    self.Atom.Draw(win)
}

func (self *Nitrogen) Copy() FlatShape {
    return NewNitrogen(self.Width, self.X, self.Y, self.Thickness, self.Color)
}

func NewNitrogen(s, x, y, t float64, col color.RGBA) *Nitrogen {
    return &Nitrogen{&Atom{0xf7, &Rectangle{s, s, x, y, t, false, col}}}
}

type OxygenA struct {
    *Atom
}

func (self *OxygenA) Draw(win *pixelgl.Window) {
    self.Atom.Draw(win)
}

func (self *OxygenA) Copy() FlatShape {
    return NewOxygenA(self.Width, self.X, self.Y, self.Thickness, self.Color)
}

func NewOxygenA(s, x, y, t float64, col color.RGBA) *OxygenA {
    return &OxygenA{&Atom{0xf5, &Rectangle{s, s, x, y, t, false, col}}}
}

type OxygenB struct {
    *Atom
}

func (self *OxygenB) Draw(win *pixelgl.Window) {
    self.Atom.Draw(win)
}

func (self *OxygenB) Copy() FlatShape {
    return NewOxygenB(self.Width, self.X, self.Y, self.Thickness, self.Color)
}

func NewOxygenB(s, x, y, t float64, col color.RGBA) *OxygenB {
    return &OxygenB{&Atom{0xf3, &Rectangle{s, s, x, y, t, false, col}}}
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

func (self *Rectangle) Rotate() {
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

func (self *Rectangle) Contains(x, y float64) bool {
    return (x > self.X) && (x < self.X+self.Width) && (y > self.Y) && (y < self.Y + self.Height)
}

func (self *Rectangle) Center() (float64, float64) {
    return self.X + self.Width/2, self.Y + self.Height/2
}

func (self *Rectangle) Copy() FlatShape {
    return NewRectangle(self.Width, self.Height, self.X, self.Y, self.Thickness, self.Color)
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

