package main

import (
    "fmt"
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

    // defualt context
    var context Context = contexts["main"]
    var savedContext Context = contexts["main"]

    // Create observer
    me := DefaultObserver()

    // Create cursor
    cursor := DefaultCursor()

    // Placeholders for static and dynamic shapes
    static := DefaultStaticShapes()
    dynamic := DefaultDynamicShapes()
    flat := DefaultFlatShapes()

    last := time.Now()

    // looping update code
	for !win.Closed() {
        // dt will be needed for many contexts
        dt := time.Since(last).Seconds()
        last = time.Now()

        // ESC is a context switcher
        if win.JustPressed(pixelgl.KeyEscape) {
            if context != contexts["menu"] {
                // saved previous context
                savedContext = context
                // start menu context
                context = contexts["menu"]
            } else {
                // load saved context
                context = savedContext
            }
        }
        if win.JustPressed(pixelgl.KeyC) {
            if context == contexts["main"] {
                savedContext = context
                context = contexts["crafting"]
            } else if context == contexts["crafting"] {
                context = savedContext
            }
        }

        // Gather up the objects that are collectively known as the "environment"
        // pass it to the current context
        e := NewEnvironment(me, win, cursor, static, dynamic, flat, dt)

        // the main context should always return HANDLING
        // if any other context returns HANDLED, go back to the main context
        // if any context returns EXIT, then exit
        code := context.Handle(e)
        if code == HANDLED {
            context = contexts["main"]
        } else if code == EXIT {
            return
        }

        // And update
        win.Update()
	}
}

func main() {
	pixelgl.Run(run)
    log.Println("Done")
}
