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
    // VSync caps FPS at refresh rate of monitor
	cfg := pixelgl.WindowConfig{
		Title:  fmt.Sprintf("Flexo"),
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		//VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
    batch := pixel.NewBatch(&pixel.TrianglesData{}, nil)

	if err != nil {
		panic(err)
	}

    // For measuring FPS
    var (
        frames = 0
        second = time.Tick(time.Second)
    )

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
    env := NewEnvironment(me, win, batch, cursor, static, dynamic, movable, immovable, templates, 0)


    last := time.Now()

    // looping update code
	for !win.Closed() {
        frames++
        select {
        case <-second:
            win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
            frames = 0
        default:
        }
        // dt will be needed for many contexts
        dt := time.Since(last).Seconds()
        last = time.Now()

        // adjust env variables
        env.Cursor = cursor
        env.Dt = dt

        // Clear the batch right before letting the current context draw
        batch.Clear()

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

        // handle context switching
        cs.Switch(win)

        // Draw the batch to the window
        batch.Draw(win)

        // And update
        win.Update()
	}
}

func main() {
	pixelgl.Run(run)
    log.Println("Done")
}
