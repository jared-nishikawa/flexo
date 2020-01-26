package main

import (
    "github.com/faiface/pixel/pixelgl"
    "github.com/faiface/pixel"
)

type MetaShape interface {
    Draw(win pixel.Target)
}

type StaticShape interface {
    //Draw(win *pixelgl.Window, ob *Observer)
    Draw(win pixel.Target, ob *Observer)
}

type DynamicShape interface {
    Draw(win pixel.Target, ob *Observer, dt float64)
}

type Environment struct {
    Observer *Observer
    Window *pixelgl.Window
    Batch *pixel.Batch
    Cursor *Cursor
    Static []StaticShape
    Dynamic []DynamicShape
    Movable []FlatShape
    Immovable []FlatShape
    Templates []FlatShape
    Dt float64
}

func NewEnvironment(ob *Observer, win *pixelgl.Window, bat *pixel.Batch, cursor *Cursor, static []StaticShape, dynamic []DynamicShape, movable, immovable, templates []FlatShape, dt float64) *Environment {
    return &Environment{ob, win, bat, cursor, static, dynamic, movable, immovable, templates, dt}
}

