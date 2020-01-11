package main

import (
    "github.com/faiface/pixel/pixelgl"
)

type MetaShape interface {
    Draw(win *pixelgl.Window)
}

type StaticShape interface {
    Draw(win *pixelgl.Window, ob *Observer)
}

type DynamicShape interface {
    Draw(win *pixelgl.Window, ob *Observer, dt float64)
}

type Environment struct {
    Observer *Observer
    Window *pixelgl.Window
    Cursor *Cursor
    Static []StaticShape
    Dynamic []DynamicShape
    Movable []FlatShape
    Immovable []FlatShape
    Templates []FlatShape
    Dt float64
}

func NewEnvironment(ob *Observer, win *pixelgl.Window, cursor *Cursor, static []StaticShape, dynamic []DynamicShape, movable, immovable, templates []FlatShape, dt float64) *Environment {
    return &Environment{ob, win, cursor, static, dynamic, movable, immovable, templates, dt}
}

