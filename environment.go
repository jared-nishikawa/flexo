package main

import (
    "github.com/faiface/pixel/pixelgl"
    "github.com/faiface/pixel/text"
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
    Dt float64
    Atlas *text.Atlas
}

func NewEnvironment(ob *Observer, win *pixelgl.Window, cursor *Cursor, static []StaticShape, dynamic []DynamicShape, dt float64, atlas *text.Atlas) *Environment {
    return &Environment{ob, win, cursor, static, dynamic, dt, atlas}
}

