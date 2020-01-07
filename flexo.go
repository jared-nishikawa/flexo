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
const WIDTH = 1920.0
//const WIDTH = 800.0
const HEIGHT = 1080.0
//const HEIGHT = 600.0
const GRAVITY = -50
const SENSITIVITY = 0.1


func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Walking",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)

	if err != nil {
		panic(err)
	}

    // cursor options
    win.DisableCursor()
    win.SetMousePosition(win.Bounds().Center())
    //win.SetCursorVisible(false)

    col := color.RGBA{128, 128, 128, 150}
    cursor := NewCursor(WIDTH/2, HEIGHT/2, 10, 3, col)

    // start interesting code here

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
        &origin,
        false)


    cubes := []*Cube{}
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
            cubes = append(cubes, cube1)
        }

    }
    circ := NewCircle(&Point{30.0, 0.0, 5.0}, 0.5, colornames.Orange)
    bounce := NewBouncing(&Point{50.0, 30.0, 1.0}, 10, 0.0, 0.0, 50.0, 0, colornames.Blue)
    fount := NewFountain(&Point{20.0, -20.0, 0.0}, 100, 0.0, 0.0, 5.0, 0.5, colornames.Navy)

    alice := colornames.Aliceblue
    bg := alice
    inMenu := false
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

    last := time.Now()
    // looping update code
	for !win.Closed() {
		win.Clear(bg)
        scoreMat := pixel.IM
        scoreMat = scoreMat.ScaledXY(pixel.ZV, pixel.V(3,3))
        scoreMat = scoreMat.Moved(pixel.V(20, 4*20))
        dot := scoreTxt.Dot
        scoreTxt.Clear()
        scoreTxt.Dot = dot
        scoreTxt.WriteString(fmt.Sprintf("score: %d\r", me.Score))
        scoreTxt.Draw(win, scoreMat)

        if win.JustPressed(pixelgl.KeyEscape) {
            curMenu := inMenu
            if curMenu {
                bg = alice
                inMenu = false
            } else {
                bg = colornames.Black
                inMenu = true
            }
        }

        if inMenu {
            menu.Write()
            if win.JustPressed(pixelgl.KeyEnter) {
                switch menu.Active {
                case 0:
                    bg = alice
                    inMenu = false
                case 1:
                    log.Println("Saving")
                case 2:
                    return
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

            // necessary
            last = time.Now()
            //win.SetMousePosition(center)
            win.Update()
            continue
        }

        mat := pixel.IM
        mat = mat.ScaledXY(pixel.ZV, pixel.V(3,3))
        mat = mat.Moved(pixel.V(20,me.Height-2*20))
        txt.Draw(win, mat)
        dt := time.Since(last).Seconds()
        last = time.Now()

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
            cursor.Color.A = 255
            cursor.X += dx*dt*100
            cursor.Y += dy*dt*100
            if cursor.X > me.Width {
                cursor.X = me.Width
            }
            if cursor.X < 0 {
                cursor.X = 0
            }
            if cursor.Y > me.Height {
                cursor.Y = me.Height
            }
            if cursor.Y < 0 {
                cursor.Y = 0
            }
        } else {
            cursor.Color.A = 150

            me.Theta -= dx*dt*SENSITIVITY
            me.Phi -= dy*dt*SENSITIVITY
            if me.Phi > math.Pi {
                me.Phi = math.Pi
            }
            if me.Phi < 0 {
                me.Phi = 0
            }
        }

        // draw things

        // cursor
        cursor.Draw(win)

        // cubes
        for _,cube := range cubes {
            cube.Draw(win, me)
        }

        // static circle
        circ.Draw(win, me)

        // bouncing circle (dynamic)
        bounce.Draw(win, me, dt)

        // fountain (dynamic)
        fount.Draw(win, me, dt)

        // update
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
    log.Println("Done")
}
