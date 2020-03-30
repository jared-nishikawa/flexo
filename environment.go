package main

import (
    "github.com/faiface/pixel/pixelgl"
    "github.com/faiface/pixel"
)

/*
type MetaShape interface {
    Draw(win pixel.Target)
}
*/

/*
type Shape interface {
    Draw(win pixel.Target, ob *Observer, dt float64)
    Dist(ob *Observer) float64
}

type StaticShape interface {
    //Draw(win *pixelgl.Window, ob *Observer)
    Draw(win pixel.Target, ob *Observer, dt float64)
    Dist(ob *Observer) float64
}
*/

//type DynamicShape interface {
//    Draw(win pixel.Target, ob *Observer, dt float64)
//    Dist(ob *Observer) float64
//}

type Environment struct {
    Observer *Observer
    Window *pixelgl.Window
    Batch *pixel.Batch
    Cursor *Cursor
    World *World
    //Dynamic []DynamicShape
    //Movable []FlatShape
    //Immovable []FlatShape
    //Templates []FlatShape
    Dt float64
}

func NewEnvironment(ob *Observer, win *pixelgl.Window, bat *pixel.Batch, cursor *Cursor, world *World, dt float64) *Environment {
    return &Environment{ob, win, bat, cursor, world, dt}
}



// start at the observer, draw a straight line forward, find the first thing in the world we hit
// return the square encountered right before we hit the object
func (self *Environment) Cast() *Point {
    ob := self.Observer
    stepSize := 0.25
    wpos := ob.Pos
    step := SphereToRec(stepSize, ob.Theta, ob.Phi)
    for {
        wpos = Add(wpos, step)
        if !self.World.Exists(wpos) {
            break
        }
        if obj := self.World.Get(wpos); obj != nil {
            wpos = Subtract(wpos, step)
            return self.World.SnapPoint(wpos)
        }
    }
    return nil
}
