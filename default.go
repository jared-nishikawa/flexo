package main

import (
    "log"
    "math"
    //"math/rand"
    "image/color"
	"golang.org/x/image/colornames"
    "golang.org/x/image/font/basicfont"
    "github.com/faiface/pixel/text"
    "github.com/faiface/pixel"
)

const FOV = math.Pi/3
const WIDTH = 1920.0
//const WIDTH = 800.0
const HEIGHT = 1080.0
//const HEIGHT = 600.0
const GRAVITY = -50
const SENSITIVITY = 0.1
const WSIZE = 101
const SNAP = 1.0

func DefaultObserver() *Observer {
    return NewObserver(
        FOV,
        WIDTH,
        HEIGHT,
        20.0, // speed
        0.0, // vertical speed
        GRAVITY,
        0.0, // theta
        math.Pi/2, //phi
        SENSITIVITY, // for mouse movement translated to 3d camera pan
        &Point{0.0, 0.0, 6.0}, // starting position
        false, // locked, for allowing cursor movement
    )
}

func DefaultCursor() *Cursor {
    return NewCursor(
        WIDTH/2, // starting x
        HEIGHT/2, // starting y
        WIDTH, // max x
        HEIGHT, // max y
        180, // multiplier on sensitivity
        10, // radius
        3, //thickness
        &color.RGBA{128, 128, 128, 150}, // color
        255, // active alpha
        150, // inactive alpha
    )
}

func DefaultWorld() *World {
    grid := make([][][]Object, WSIZE)
    for i := range grid {
        grid[i] = make([][]Object, WSIZE)
        for j := range grid[i] {
            grid[i][j] = make([]Object, WSIZE)
        }
    }

    pos := &Point{SNAP*30.0, 0.0, 0.0}
    cube := NewBorderedCube(SNAP, pos, &color.RGBA{0x0f, 0x0f, 0x0f, 0x7f}, &colornames.Black)
    w := &World{
        Grid: grid,
        Snap: SNAP,
    }

    w.Set(pos, cube)

    return w
}

/*
func DefaultStaticShapes() []StaticShape {
    static := []StaticShape{}
    //for y:=0;y<36;y+=6 {
    //    for z:=0;z<36;z+=6 {
    //        for x:=0;x<36;x+=6 {
    //            p1 := Point{float64(x), float64(y), float64(z)}

    //            r := uint8(rand.Int() % 256)
    //            g := uint8(rand.Int() % 256)
    //            b := uint8(rand.Int() % 256)
    //            a := uint8(128)

    //            cube1 := NewSolidCube(5, &p1, color.RGBA{r, g, b, a})

    //            static = append(static, cube1)
    //        }
    //    }

    //}
    //circ := NewSphere(&Point{30.0, 0.0, 5.0}, 0.5, colornames.Orange)
    //static = append(static, circ)

    //wall1 := NewWall(100.0, 100.0, 100.0, 120.0, 0.0)
    //wall2 := NewWall(100.0, 120.0, 120.0, 120.0, 0.0)
    //wall3 := NewWall(120.0, 120.0, 120.0, 100.0, 0.0)
    //wall4 := NewWall(300.0, 300.0, 400.0, 300.0, 100.0)
    //static = append(static, []StaticShape{wall1, wall2, wall3, wall4}...)
    return static
}
*/

/*
func DefaultDynamicShapes() []DynamicShape {
    dynamic := []DynamicShape{}
    //bounce := NewBouncing(&Point{50.0, 30.0, 0.0}, 10, 0.0, 0.0, 50.0, 0, colornames.Blue)
    //fount := NewFountain(&Point{20.0, -20.0, 0.0}, 100, 0.0, 0.0, 5.0, 0.5, colornames.Navy)

    //dynamic = append(dynamic, bounce)
    //dynamic = append(dynamic, fount)
    return dynamic
}
*/

/*
func DefaultFlatShapes() []FlatShape {
    flat := []FlatShape{}
    //s := NewSquare(100, 400, 500, 5, colornames.Red)
    //c := NewCircle(500, 600, 200, 5, colornames.Green)
    //flat = append(flat, s, c)
    return flat
}

func DefaultImmovableShapes() []FlatShape {
    immov := []FlatShape{}
    left := NewRectangle(300, 1080, 0, 0, 0, color.RGBA{200, 200, 200, 255})
    immov = append(immov, left)
    return immov
}

func DefaultTemplates() []FlatShape {
    templates := []FlatShape{}
    hydrogen := NewHydrogen(50, 150, 800, 0, colornames.White)
    oxygenA := NewOxygenA(50, 150, 700, 0, colornames.Green)
    oxygenB := NewOxygenB(50, 150, 600, 0, colornames.Green)
    nitrogen := NewNitrogen(50, 150, 500, 0, colornames.Red)
    carbon := NewCarbon(50, 150, 400, 0, colornames.Blue)

    hydrogen.Snap()
    oxygenA.Snap()
    oxygenB.Snap()
    nitrogen.Snap()
    carbon.Snap()

    templates = append(templates, hydrogen, oxygenA, oxygenB, nitrogen, carbon)
    return templates
}

func DefaultMovableShapes() []FlatShape {
    return []FlatShape{}
}
*/

func DefaultAtlas() *text.Atlas {
    return text.NewAtlas(basicfont.Face7x13, text.ASCII)
}

func DefaultContexts() map[string]Context {
    contexts := make(map[string]Context)
    atlas := DefaultAtlas()

    // main context
    contexts["main"] = &MainContext{}

    // menu context
    menu := NewMenu(nil, "root", atlas, []string{"resume", "save", "options", "exit"}, &colornames.White, &colornames.Orange)
    menu.Handle = func(num int, env *Environment) (*Menu, int) {
        switch num {
        case 0:
            return menu, HANDLED
        case 2:
            return menu.Children[menu.Options[num]], HANDLING

        case 3:
            return nil, EXIT
        default:
            log.Println(menu.Options[num])
            return menu, HANDLING
        }}
    opts := NewMenu(menu, "options", atlas, []string{"resolution", "vsync"}, &colornames.White, &colornames.Orange)
    opts.Handle = func(num int, env *Environment) (*Menu, int) {
        switch num {
        default:
            log.Println(opts.Options[num])
            return opts.Children[opts.Options[num]], HANDLING
        }}
    vs := NewMenu(opts, "vsync", atlas, []string{"off", "on"}, &colornames.White, &colornames.Orange)
    vs.Handle = func(num int, env *Environment) (*Menu, int) {
        win := env.Window
        switch num {
        case 0:
            win.SetVSync(false)
            return vs, HANDLING
        case 1:
            win.SetVSync(true)
            return vs, HANDLING
        default:
            log.Println(vs.Options[num])
            return vs.Children[vs.Options[num]], HANDLING
        }}
    res := NewMenu(opts, "resolution", atlas, []string{"640x480", "800x600", "1024x768", "1920x1080"}, &colornames.White, &colornames.Orange)
    res.Handle = func(num int, env *Environment) (*Menu, int) {
        win := env.Window
        ob := env.Observer
        cur := env.Cursor
        h := 0.0
        w := 0.0
        switch num {
        case 0:
            h = 480.0
            w = 640.0
            //return res, HANDLING
        case 1:
            h = 600
            w = 800
        case 2:
            h = 768
            w = 1024
        case 3:
            h = 1080
            w = 1920
        default:
            log.Println(res.Options[num])
            return res, HANDLING
        }
        win.SetBounds(pixel.R(0,0,w, h))
        ob.Width = w
        ob.Height = h
        ob.VFov = h/w * ob.HFov
        cur.X = w/2
        cur.Y = h/2
        cur.MaxX = w
        cur.MaxY = h
        return res, HANDLING

        }

    contexts["menu"] = NewMenuContext(menu)

    //crafting context
    //contexts["crafting"] = &CraftingContext{}
    return contexts
}

