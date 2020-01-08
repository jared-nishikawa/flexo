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
    Menu *Menu
    Static []StaticShape
    Dynamic []DynamicShape
    Dt float64
}

func NewEnvironment(ob *Observer, win *pixelgl.Window, cursor *Cursor, menu *Menu, static []StaticShape, dynamic []DynamicShape, dt float64) *Environment {
    return &Environment{ob, win, cursor, menu, static, dynamic, dt}
}

