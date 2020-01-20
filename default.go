package main

import (
    "math"
    "image/color"
	"golang.org/x/image/colornames"
    "golang.org/x/image/font/basicfont"
    "github.com/faiface/pixel/text"
)

const FOV = math.Pi/3
const WIDTH = 1920.0
//const WIDTH = 800.0
const HEIGHT = 1080.0
//const HEIGHT = 600.0
const GRAVITY = -50
const SENSITIVITY = 0.1

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
        &Point{0.0, 0.0, 0.0}, // starting position
        false, // locked, for allowing cursor movement
    )
}

func DefaultCursor() *Cursor {
    return NewCursor(
        WIDTH/2, // starting x
        HEIGHT/2, // starting y
        WIDTH, // max x
        HEIGHT, // max y
        35, // multiplier on sensitivity
        10, // radius
        3, //thickness
        color.RGBA{128, 128, 128, 150}, // color
        255, // active alpha
        150, // inactive alpha
    )
}

func DefaultStaticShapes() []StaticShape {
    static := []StaticShape{}
    for y:=-10;y<11;y+=6 {
    //for y:=5;y<11;y+=6 {
        x := float64(20)
        for z:=-9;z<4;z+=6 {
        //for z:=0;z<4;z+=6 {
            p := Point{x, float64(y), float64(z)}
            //q := Point{float64(y), x, float64(z)}
            cube1 := NewCube(5, &p, colornames.Black)
            //cube2 := NewCube(5, &q)
            //cubes = append(cubes, cube1, cube2)
            static = append(static, cube1)
        }

    }
    circ := NewSphere(&Point{30.0, 0.0, 5.0}, 0.5, colornames.Orange)
    static = append(static, circ)
    return static
}

func DefaultDynamicShapes() []DynamicShape {
    dynamic := []DynamicShape{}
    bounce := NewBouncing(&Point{50.0, 30.0, 1.0}, 10, 0.0, 0.0, 50.0, 0, colornames.Blue)
    //fount := NewFountain(&Point{20.0, -20.0, 0.0}, 100, 0.0, 0.0, 5.0, 0.5, colornames.Navy)

    dynamic = append(dynamic, bounce)
    //dynamic = append(dynamic, fount)
    return dynamic
}

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

func DefaultAtlas() *text.Atlas {
    return text.NewAtlas(basicfont.Face7x13, text.ASCII)
}

func DefaultContexts() map[string]Context {
    contexts := make(map[string]Context)
    atlas := DefaultAtlas()

    // main context
    contexts["main"] = &MainContext{}

    // menu context
    menu := NewMenu(atlas, []string{"resume", "save", "options", "exit"}, colornames.White, colornames.Orange)
    contexts["menu"] = NewMenuContext(menu)

    //crafting context
    contexts["crafting"] = &CraftingContext{}
    return contexts
}
