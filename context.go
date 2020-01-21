package main

import (
    "fmt"
    "image/color"
    "log"
//    "math"
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
    HandleEscape() bool
}

type NullContext struct {
}

func (self *NullContext) HandleEscape() bool {
    return false
}

func (self *NullContext) Handle(env *Environment) int {
    return HANDLED
}

type MainContext struct {
}

func (self *MainContext) HandleEscape() bool {
    return false
}

func (self *MainContext) Handle(env *Environment) int {
    me := env.Observer
    win := env.Window
    cursor := env.Cursor
    static := env.Static
    dynamic := env.Dynamic
    dt := env.Dt
    atlas := DefaultAtlas()

    txt := text.New(pixel.ZV, atlas)
    txt.Color = colornames.Brown
    txt.WriteString("jump: [space]\n")
    txt.WriteString("move: [wasd] or arrow keys\n")
    txt.WriteString("toggle cursor: [right-click]\n")
    txt.WriteString("menu: [ESC]\n")

    //desired := me.Height/40
    //scale := desired/h
    //menuMat = menuMat.ScaledXY(pixel.ZV, pixel.V(scale,scale))
    //menuMat = menuMat.Moved(pixel.V(h, me.Height-2*h))

    mat := pixel.IM
    desired := me.Height/40
    scale := desired/atlas.LineHeight()
    mat = mat.ScaledXY(pixel.ZV, pixel.V(scale,scale))
    mat = mat.Moved(pixel.V(desired/2,me.Height-desired))

    scoreTxt := text.New(pixel.ZV, atlas)
    scoreTxt.Color = colornames.Black
    scoreMat := pixel.IM
    scoreMat = scoreMat.ScaledXY(pixel.ZV, pixel.V(scale,scale))
    scoreMat = scoreMat.Moved(pixel.V(desired/2, 4*desired))
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
        //cursor.Color.A = 150

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

func (self *MenuContext) HandleEscape() bool {
    menu := self.Menu
    if menu.Root != nil {
        self.Menu = menu.Root
        return true
    }
    return false
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

    /*
    if win.JustPressed(pixelgl.KeyEnter) {
        switch menu.Active {
        case 0:
            return HANDLED
        case 3:
            return EXIT
        default:
            log.Println(menu.Options[menu.Active])
        }

    }
    */

    if win.JustPressed(pixelgl.KeyUp) {
        menu.Up()
    }

    if win.JustPressed(pixelgl.KeyDown) {
        menu.Down()
    }

    menuMat := pixel.IM
    desired := me.Height/40
    scale := desired/menu.Atlas.LineHeight()
    menuMat = menuMat.ScaledXY(pixel.ZV, pixel.V(scale,scale))
    menuMat = menuMat.Moved(pixel.V(desired/2, me.Height-desired))

    menu.Draw(win, menuMat)
    if win.JustPressed(pixelgl.KeyEnter) {
        m, code := menu.Handle(menu.Active)
        self.Menu = m
        return code
    }
    return HANDLING
    //return HANDLING
}

type CraftingContext struct {
}

func (self *CraftingContext) HandleEscape() bool {
    return false
}

func (self *CraftingContext) Handle(env *Environment) int {
    //me := env.Observer
    win := env.Window
    cursor := env.Cursor
    //static := env.Static
    //dynamic := env.Dynamic
    dt := env.Dt
    cursor.SetActive()
    movable := env.Movable
    immov := env.Immovable
    templates := env.Templates

    //win.Clear(color.RGBA{0xff,0xf9,0xf9,255})
    win.Clear(color.RGBA{0x60,0x60,0x70,255})

    // align to mouse
    prev := win.MousePreviousPosition()
    pos := win.MousePosition()

    // compute mouse distance traveled
    dx := pos.X - prev.X
    dy := pos.Y - prev.Y

    old_x := cursor.X
    old_y := cursor.Y

    cursor.Move(dx, dy, dt)
    x := cursor.X
    y := cursor.Y

    for _,s := range templates {
        if win.JustPressed(pixelgl.MouseButtonLeft) && s.Contains(old_x, old_y) {
            //side := 50.0
            //sq := NewSquare(side, old_x-side/2, old_y-side/2, 5, colornames.Red)
            sq := s.Copy()
            movable = append(movable, sq)
            env.Movable = movable
        }
    }

    for i,s := range movable {
        if win.JustPressed(pixelgl.KeyR) && s.Contains(old_x, old_y) {
            s.Rotate()
        }

        if win.Pressed(pixelgl.MouseButtonLeft) {
            if s.Contains(old_x, old_y) {
                if win.JustPressed(pixelgl.MouseButtonLeft) {
                    s.SetDragging(true)
                }

                if s.GetDragging() {
                    s.SetLoc(s.GetX()+(x-old_x), s.GetY()+(y-old_y))
                }
            }
        } else if win.JustPressed(pixelgl.MouseButtonRight) {
            if s.Contains(old_x, old_y) {
                // problems when there's two shapes on top of each other
                log.Println(i, len(movable))
                movable = append(movable[:i], movable[i+1:]...)
                env.Movable = movable
                s = nil
            }
        } else {
            s.Snap()
            s.SetDragging(false)
        }

    }



    for _,s := range immov {
        s.Draw(win)
    }

    for _,s := range templates {
        s.Draw(win)
    }

    for _,s := range movable {
        s.Draw(win)
    }
    cursor.Draw(win)
    return HANDLING

}

