package main

import (
    "fmt"
    "log"
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "golang.org/x/image/colornames"
    "github.com/faiface/pixel/text"
)

const (
    HANDLED = iota
    HANDLING
    EXIT
)

type Context interface {
    Handle(env *Environment) int
}

type NullContext struct {
}

func (self *NullContext) Handle(env *Environment) int {
    return HANDLED
}

type MainContext struct {
}

func (self *MainContext) Handle(env *Environment) int {
    me := env.Observer
    win := env.Window
    cursor := env.Cursor
    static := env.Static
    dynamic := env.Dynamic
    dt := env.Dt

    txt := text.New(pixel.ZV, env.Atlas)
    txt.Color = colornames.Brown
    txt.WriteString("jump: [space]\n")
    txt.WriteString("move: [wasd] or arrow keys\n")
    txt.WriteString("toggle cursor: [right-click]\n")
    txt.WriteString("menu: [ESC]\n")
    mat := pixel.IM
    mat = mat.ScaledXY(pixel.ZV, pixel.V(3,3))
    mat = mat.Moved(pixel.V(20,me.Height-2*20))

    scoreTxt := text.New(pixel.ZV, env.Atlas)
    scoreTxt.Color = colornames.Black
    scoreMat := pixel.IM
    scoreMat = scoreMat.ScaledXY(pixel.ZV, pixel.V(3,3))
    scoreMat = scoreMat.Moved(pixel.V(20, 4*20))
    dot := scoreTxt.Dot
    scoreTxt.Clear()
    scoreTxt.Dot = dot
    scoreTxt.WriteString(fmt.Sprintf("score: %d\r", me.Score))



    win.Clear(colornames.Aliceblue)

    if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
        me.Forward(dt)
    }

    if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
        me.Backward(dt)
    }

    if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
        me.Left(dt)
    }

    if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
        me.Right(dt)
    }

    if win.Pressed(pixelgl.KeyK) {
        me.Ascend(dt)
    }

    if win.Pressed(pixelgl.KeyJ) {
        me.Descend(dt)
    }

    if win.JustPressed(pixelgl.KeySpace) {
        me.Jump()
    }
    me.Freefall(dt)

    if win.JustPressed(pixelgl.MouseButtonLeft) {
        me.Score += 1
    }

    if win.JustPressed(pixelgl.MouseButtonRight) {
        if me.Locked {
            me.Locked = false
        } else {
            me.Locked = true
        }

    }

    // align to mouse
    prev := win.MousePreviousPosition()
    pos := win.MousePosition()
    //win.SetMousePosition(center)

    // compute mouse distance traveled
    dx := pos.X - prev.X
    dy := pos.Y - prev.Y
    if me.Locked {
        cursor.SetActive()
        //cursor.Color.A = 255
        cursor.Move(dx, dy, dt)
    } else {
        cursor.SetInactive()
        cursor.Color.A = 150

        me.Yaw(dx, dt)
        me.Pitch(dy, dt)

    }

    // draw things

    // static shapes
    for _,shape := range static {
        shape.Draw(win, me)
    }

    // dynamic shapes
    for _,shape := range dynamic {
        shape.Draw(win, me, dt)
    }

    // cursor
    cursor.Draw(win)

    // text
    txt.Draw(win, mat)
    scoreTxt.Draw(win, scoreMat)

    // main context always returns handling
    return HANDLING

}

type MenuContext struct {
    Menu *Menu
}

func NewMenuContext(menu *Menu) *MenuContext {
    return &MenuContext{menu}
}

func (self *MenuContext) Handle(env *Environment) int {

    me := env.Observer
    win := env.Window
    //cursor := env.Cursor
    //static := env.Static
    //dynamic := env.Dynamic
    //dt := env.Dt
    menu := self.Menu

    win.Clear(colornames.Black)
    menu.Write()

    if win.JustPressed(pixelgl.KeyEnter) {
        switch menu.Active {
        case 0:
            return HANDLED
        case 1:
            log.Println("Saving")
        case 2:
            log.Println("Options")
        case 3:
            return EXIT
        default:
            log.Println(menu.Active)
        }

    }
    if win.JustPressed(pixelgl.KeyUp) {
        menu.Up()
    }

    if win.JustPressed(pixelgl.KeyDown) {
        menu.Down()
    }

    menuMat := pixel.IM
    menuMat = menuMat.ScaledXY(pixel.ZV, pixel.V(2,2))
    menuMat = menuMat.Moved(pixel.V(20, me.Height-2*20))

    menu.Draw(win, menuMat)
    return HANDLING
}
