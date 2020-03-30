package main

import (
    "fmt"
    "image/color"
//    "log"
//    "sort"
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
    bat := env.Batch
    cursor := env.Cursor
    //static := env.Static
    world := env.World
    //dynamic := env.Dynamic
    dt := env.Dt
    atlas := DefaultAtlas()

    txt := text.New(pixel.ZV, atlas)
    txt.Color = &colornames.Brown
    txt.WriteString("jump: [space]\n")
    txt.WriteString("move: [wasd] or arrow keys\n")
    txt.WriteString("place cube: [left-click]\n")
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
    scoreTxt.Color = &colornames.Black
    scoreMat := pixel.IM
    scoreMat = scoreMat.ScaledXY(pixel.ZV, pixel.V(scale,scale))
    scoreMat = scoreMat.Moved(pixel.V(desired/2, 4*desired))
    dot := scoreTxt.Dot
    scoreTxt.Clear()
    scoreTxt.Dot = dot
    scoreTxt.WriteString(fmt.Sprintf("score: %d\r", me.Score))

    win.Clear(&colornames.Aliceblue)

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

    if win.Pressed(pixelgl.KeyV) {
        me.Ascend(dt)
    }

    if win.Pressed(pixelgl.KeyC) {
        me.Descend(dt)
    }

    if win.JustPressed(pixelgl.KeySpace) {
        me.Jump()
    }
    //me.Freefall(dt)

    hand := env.Cast()

    if win.JustPressed(pixelgl.MouseButtonLeft) && hand != nil {
        me.Score += 1
        cube := NewBorderedCube(world.Snap, hand, &color.RGBA{0xff, 0, 0, 0xff}, &colornames.Black)
        world.Set(hand, cube)
        //env.Static = append(env.Static, cube)
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
    // take into account observer's FOV and position
    // draw 4 rays
    d := 40.0
    hray1 := me.Theta - me.HFov
    hray2 := me.Theta + me.HFov
    p1 := world.Ray(me.Pos, d, hray1, me.Phi)
    p2 := world.Ray(me.Pos, d, hray2, me.Phi)

    vray1 := me.Phi - me.VFov
    vray2 := me.Phi + me.VFov
    q1 := world.Ray(me.Pos, d, me.Theta, vray1)
    q2 := world.Ray(me.Pos, d, me.Theta, vray2)

    minX := Min(p1[0], p2[0], q1[0], q2[0], me.Pos[0])
    maxX := Max(p1[0], p2[0], q1[0], q2[0], me.Pos[0])
    minY := Min(p1[1], p2[1], q1[1], q2[1], me.Pos[1])
    maxY := Max(p1[1], p2[1], q1[1], q2[1], me.Pos[1])
    minZ := Min(p1[2], p2[2], q1[2], q2[2], me.Pos[2])
    maxZ := Max(p1[2], p2[2], q1[2], q2[2], me.Pos[2])

    for i := minX; i <= maxX; i++ {
        for j := minY; j <= maxY; j++ {
            for k := minZ; k <= maxZ; k++ {
                p := &Point{i,j,k}
                if !world.Exists(p) {
                    continue
                }
                obj := world.Get(p)
                if obj != nil {
                    obj.Draw(bat, me, dt)
                }
            }
        }
    }

    // draw template
    if hand != nil {
        templateCube := NewBorderedCube(world.Snap, hand, &color.RGBA{0x0, 0xff, 0, 0x5f}, &colornames.Black)
        templateCube.Draw(bat, me, dt)
    }

    // cursor
    cursor.Draw(bat)

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
    menu := self.Menu

    win.Clear(colornames.Black)
    menu.Write()

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
        m, code := menu.Handle(menu.Active, env)
        self.Menu = m
        return code
    }
    return HANDLING
    //return HANDLING
}

/*
type CraftingContext struct {
}

func (self *CraftingContext) HandleEscape() bool {
    return false
}

func (self *CraftingContext) Handle(env *Environment) int {
    //me := env.Observer
    win := env.Window
    bat := env.Batch
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
        s.Draw(bat)
    }

    for _,s := range templates {
        s.Draw(bat)
    }

    for _,s := range movable {
        s.Draw(bat)
    }
    cursor.Draw(bat)
    return HANDLING

}

*/
