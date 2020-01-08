package main

import (
    "fmt"
    "image/color"
    "log"
    "math"
    "time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
    "golang.org/x/image/font/basicfont"
    "github.com/faiface/pixel/text"
)

const FOV = math.Pi/3
//const WIDTH = 1920.0
const WIDTH = 800.0
//const HEIGHT = 1080.0
const HEIGHT = 600.0
const GRAVITY = -50
const SENSITIVITY = 0.1


func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Walking",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		//VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(err)
	}

    // cursor options
    win.DisableCursor()
    win.SetMousePosition(win.Bounds().Center())
    //win.SetCursorVisible(false)

    // start interesting code here

    mainContext := &MainContext{}
    menuContext := &MenuContext{}
    var context Context = mainContext
    //context := []Context{&NullContext{}}

    origin := Point{0.0, 0.0, 0.0}
    me := NewObserver(
        FOV,
        WIDTH,
        HEIGHT,
        20.0,
        0.0,
        GRAVITY,
        0.0,
        math.Pi/2,
        SENSITIVITY,
        &origin,
        false)

    col := color.RGBA{128, 128, 128, 150}
    cursor := NewCursor(WIDTH/2, HEIGHT/2, WIDTH, HEIGHT, 35, 10, 3, col, 255, 150)

    static := []StaticShape{}
    dynamic := []DynamicShape{}
    //meta := []MetaShape{}
    //meta = append(meta, cursor)
    //for y:=-10;y<11;y+=6 {
    for y:=5;y<11;y+=6 {
        x := float64(20)
        //for z:=-9;z<4;z+=6 {
        for z:=0;z<4;z+=6 {
            p := Point{x, float64(y), float64(z)}
            //q := Point{float64(y), x, float64(z)}
            cube1 := NewCube(5, &p, colornames.Black)
            //cube2 := NewCube(5, &q)
            //cubes = append(cubes, cube1, cube2)
            static = append(static, cube1)
        }

    }

    circ := NewCircle(&Point{30.0, 0.0, 5.0}, 0.5, colornames.Orange)
    static = append(static, circ)

    bounce := NewBouncing(&Point{50.0, 30.0, 1.0}, 10, 0.0, 0.0, 50.0, 0, colornames.Blue)
    fount := NewFountain(&Point{20.0, -20.0, 0.0}, 1, 0.0, 0.0, 5.0, 0.5, colornames.Navy)

    dynamic = append(dynamic, bounce)
    dynamic = append(dynamic, fount)

    alice := colornames.Aliceblue
    bg := alice
    //inMenu := false
    //center := win.Bounds().Center()

    atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
    txt := text.New(pixel.ZV, atlas)
    txt.Color = colornames.Brown
    txt.WriteString("jump: [space]\n")
    txt.WriteString("move: [wasd] or arrow keys\n")
    txt.WriteString("toggle cursor: [right-click]\n")
    txt.WriteString("menu: [ESC]\n")

    scoreTxt := text.New(pixel.ZV, atlas)
    scoreTxt.Color = colornames.Black


    menu := NewMenu(atlas, []string{"resume", "save", "exit"}, colornames.White, colornames.Orange)

    var (
        frames = 0
        second = time.Tick(time.Second)
    )
    last := time.Now()
    // looping update code
	for !win.Closed() {
		win.Clear(bg)
        dt := time.Since(last).Seconds()
        last = time.Now()


        //scoreMat := pixel.IM
        //scoreMat = scoreMat.ScaledXY(pixel.ZV, pixel.V(3,3))
        //scoreMat = scoreMat.Moved(pixel.V(20, 4*20))
        //dot := scoreTxt.Dot
        //scoreTxt.Clear()
        //scoreTxt.Dot = dot
        //scoreTxt.WriteString(fmt.Sprintf("score: %d\r", me.Score))
        //scoreTxt.Draw(win, scoreMat)

        if win.JustPressed(pixelgl.KeyEscape) {
            if context == mainContext {
                bg = colornames.Black
                context = menuContext
            } else {
                bg = alice
                context = mainContext
            }
        }

        //mat := pixel.IM
        //mat = mat.ScaledXY(pixel.ZV, pixel.V(3,3))
        //mat = mat.Moved(pixel.V(20,me.Height-2*20))
        //txt.Draw(win, mat)

        e := NewEnvironment(me, win, cursor, menu, static, dynamic, dt)
        code := context.Handle(e)
        if code == HANDLED {
            context = mainContext
        } else if code == EXIT {
            return
        }

        win.Update()
        frames++
        select {
        case <-second:
            win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
            frames = 0
        default:
        }

	}
}

func main() {
	pixelgl.Run(run)
    log.Println("Done")
}
