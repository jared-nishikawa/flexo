package main

import (
    "fmt"
    "log"
    "time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
    "golang.org/x/image/font/basicfont"
    "github.com/faiface/pixel/text"
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

    // Atlas
    atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

    // Create contexts
    // -- main context
    mainContext := &MainContext{}

    // -- menu context
    menu := NewMenu(atlas, []string{"resume", "save", "options", "exit"}, colornames.White, colornames.Orange)
    menuContext := NewMenuContext(menu)

    // -- defualt context
    var context Context = mainContext

    // Create observer
    me := DefaultObserver()

    // Create cursor
    cursor := DefaultCursor()

    // Placeholders for static and dynamic shapes
    static := DefaultStaticShapes()
    dynamic := DefaultDynamicShapes()

    last := time.Now()

    // looping update code
	for !win.Closed() {
        // dt will be needed for many contexts
        dt := time.Since(last).Seconds()
        last = time.Now()

        // ESC is a context switcher
        if win.JustPressed(pixelgl.KeyEscape) {
            if context == mainContext {
                context = menuContext
            } else {
                context = mainContext
            }
        }

        // Gather up the objects that are collectively known as the "environment"
        // pass it to the current context
        e := NewEnvironment(me, win, cursor, static, dynamic, dt, atlas)

        // the main context should always return HANDLING
        // if any other context returns HANDLED, go back to the main context
        // if any context returns EXIT, then exit
        code := context.Handle(e)
        if code == HANDLED {
            context = mainContext
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
