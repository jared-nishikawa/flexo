package main

import (
    "fmt"
    "image/color"
    "log"
    "time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
    // Create window
	cfg := pixelgl.WindowConfig{
		Title:  fmt.Sprintf("Flexo"),
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(err)
	}

    // Set cursor
    win.DisableCursor()
    win.SetMousePosition(win.Bounds().Center())

    // Create contexts
    contexts := DefaultContexts()

    // Context switcher
    cs := NewContextSwitcher(contexts)

    // Create observer
    me := DefaultObserver()

    // Create cursor
    cursor := DefaultCursor()

    // Placeholders for static and dynamic shapes
    static := DefaultStaticShapes()
    dynamic := DefaultDynamicShapes()
    movable := DefaultMovableShapes()
    immovable := DefaultImmovableShapes()
    templates := DefaultTemplates()
    //flat := DefaultFlatShapes()

    // Gather up the objects that are collectively known as the "environment"
    // will be passed to the context for handling
    env := NewEnvironment(me, win, cursor, static, dynamic, movable, immovable, templates, 0)

    pos := &Point{me.Pos[0]+10, me.Pos[1], me.Pos[2]-3}
    c := NewCube(2, pos, color.RGBA{0xff, 0, 0, 0x7f})

    last := time.Now()

    // looping update code
	for !win.Closed() {
        // dt will be needed for many contexts
        dt := time.Since(last).Seconds()
        last = time.Now()

        // handle context switching
        cs.Switch(win)

        // adjust env variables
        env.Cursor = cursor
        env.Dt = dt

        // the main context should always return HANDLING
        // if any other context returns HANDLED, go back to the main context
        // if any context returns EXIT, then exit
        code := cs.Current().Handle(env)
        if code == HANDLED {
            cs.PopMenu()
            //cs.Current = cs.Contexts["main"]
        } else if code == EXIT {
            return
        }

        c.Draw(win, me)

        // And update
        win.Update()
	}
}

func main() {
	pixelgl.Run(run)
    log.Println("Done")
}
