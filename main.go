package main

import (
    "fmt"
    "image/color"
    "log"
    "math"
    "math/rand"
    "time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
    "github.com/faiface/pixel/imdraw"
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

type Menu struct {
    Options []string
    Active int
    Color color.RGBA
    ActiveColor color.RGBA
    Atlas *text.Atlas
    Text *text.Text
}

func NewMenu(atlas *text.Atlas, options []string) *Menu {
    return &Menu{options, 0, colornames.White, colornames.Orange, atlas, nil}
}

func (self *Menu) Write() {
    optText := text.New(pixel.ZV, self.Atlas)
    optText.Color = self.Color
    for i, opt := range self.Options {
        if i == self.Active {
            optText.Color = self.ActiveColor
        } else {
            optText.Color = self.Color
        }
        optText.WriteString(opt)
        optText.WriteString("\n")
    }
    self.Text = optText
}

func (self *Menu) Draw(win *pixelgl.Window, mat pixel.Matrix) {
    self.Text.Draw(win, mat)
}


func (self *Menu) Down() {
    self.Active = (self.Active + 1) % len(self.Options)
    if self.Active < 0 {
        self.Active += len(self.Options)
    }
}

func (self *Menu) Up() {
    self.Active = (self.Active - 1) % len(self.Options)
    if self.Active < 0 {
        self.Active += len(self.Options)
    }

}

type Observer struct {
    HFov float64
    VFov float64
    Width float64
    Height float64
    Speed float64
    VerticalSpeed float64
    Theta float64
    Phi float64
    Pos *Point
    Locked bool
    Curs *Cursor
    Score int
}

func (self *Observer) Forward(dt float64) {
    x := math.Cos(self.Theta)
    y := math.Sin(self.Theta)
    self.Pos[0] += x*dt*self.Speed
    self.Pos[1] += y*dt*self.Speed
}

func (self *Observer) Backward(dt float64) {
    x := math.Cos(self.Theta)
    y := math.Sin(self.Theta)
    self.Pos[0] -= x*dt*self.Speed
    self.Pos[1] -= y*dt*self.Speed
}

func (self *Observer) Left(dt float64) {
    x := -math.Sin(self.Theta)
    y := math.Cos(self.Theta)
    self.Pos[0] += x*dt*self.Speed
    self.Pos[1] += y*dt*self.Speed
}

func (self *Observer) Right(dt float64) {
    x := -math.Sin(self.Theta)
    y := math.Cos(self.Theta)
    self.Pos[0] -= x*dt*self.Speed
    self.Pos[1] -= y*dt*self.Speed
}

func (self *Observer) Ascend(dt float64) {
    self.Pos[2] += dt*self.Speed
}

func (self *Observer) Descend(dt float64) {
    self.Pos[2] -= dt*self.Speed
}


func (self *Observer) Jump() {
    self.VerticalSpeed += 30
}

func (self *Observer) Freefall(dt float64) {
    self.Pos[2] += self.VerticalSpeed*dt
    if self.Pos[2] > 0 {
        self.VerticalSpeed += GRAVITY*dt
        if self.Pos[2] < 0 {
            self.Pos[2] = 0
        }
    } else {
        self.VerticalSpeed = 0
    }
}


func NewObserver(fov, w, h, s, vs, th, ph float64, p *Point, l bool) *Observer {
    vfov := h/w * fov
    col := color.RGBA{128, 128, 128, 150}
    cur := NewCursor(WIDTH/2, HEIGHT/2, 10, 3, col)
    return &Observer{fov, vfov, w, h, s, vs, th, ph, p, l, cur, 0}
}

func SphereToRec(rho, theta, phi float64) *Point {
    x := rho*math.Cos(theta)*math.Sin(phi)
    y := rho*math.Sin(theta)*math.Sin(phi)
    z := rho*math.Cos(phi)
    return &Point{x, y, z}
}

func RecToSphere(P *Point) (float64, float64, float64) {
    // Atan2 returns a value in (-pi, pi)
    x := P[0]
    y := P[1]
    z := P[2]
    rho := math.Sqrt(math.Pow(x,2) + math.Pow(y,2) + math.Pow(z,2))
    // rho is guaranteed >= 0
    theta := math.Atan2(y, x)
    // theta is guaranteed in (-pi, pi)
    phi := 0.0
    if rho != 0.0 {
        // proj is in the x-y plane
        proj := math.Sqrt(math.Pow(x,2) + math.Pow(y,2))
        phi = math.Atan2(proj, z)
        // since proj is >= 0, phi is guaranteed in (0, pi)
    }
    return rho, theta, phi
}


func Project(theta, phi float64, ob *Observer) (float64, float64) {
    // Assumes viewer is on the origin, looking toward the x-axis
    // Theta increases from right to left
    // Phi increases from top to bottom
    // theta = -FOV/2 = right side of screen
    // theta = FOV/2 = left side of screen
    // phi = -FOV/2 = top of screen
    // phi = FOV/2 = bottom of screen
    // Linearly interpolate the rest

    //theta -= ob.Theta
    theta = math.Mod(theta, 2*math.Pi)
    if theta < -math.Pi {
        theta += 2*math.Pi
    }

    if theta > math.Pi {
        theta -= 2*math.Pi
    }

    //phi = phi - ob.Phi + math.Pi/2

    //below_horizon := ob.Phi - math.Pi/2
    //phi -= below_horizon

    x := ob.Width  * (((-theta/(ob.HFov/2))+1)/2)
    y := ob.Height * (((-(phi-math.Pi/2)/(ob.VFov/2))+1)/2)

    return x,y
}


type Point = [3]float64

// distance
func Distance(P, Q *Point) float64 {
    return math.Sqrt(
        math.Pow(P[0]-Q[0],2) +
        math.Pow(P[1]-Q[1],2) +
        math.Pow(P[2]-Q[2],2))
}

func Magnitude(P *Point) float64 {
    return Distance(P, &Point{0.0,0.0,0.0})
}

// translates Q to origin
func Translate(P, Q *Point) *Point {
    return &Point{P[0]-Q[0], P[1]-Q[1], P[2]-Q[2]}
}

func Add(P, Q *Point) *Point {
    return &Point{P[0]+Q[0], P[1]+Q[1], P[2]+Q[2]}
}

// rotate around z-axis
func Rotate(P *Point, theta float64) *Point {
    theta = -theta
    x := math.Cos(theta)*P[0] + math.Sin(theta)*P[1]
    y := -math.Sin(theta)*P[0] + math.Cos(theta)*P[1]
    z := P[2]
    return &Point{x,y,z}

    //rh, th, ph := RecToSphere(P)
    //th += theta
    //return SphereToRec(rh, th, ph)
}

// tilt down from vertical
func Tilt(P *Point, phi float64) *Point {
    phi = -phi
    x := math.Cos(phi)*P[0]-math.Sin(phi)*P[2]
    z := math.Sin(phi)*P[0]+math.Cos(phi)*P[2]
    y := P[1]
    return &Point{x,y,z}

    //rh, th, ph := RecToSphere(P)
    //ph += phi
    //return SphereToRec(rh, th, ph)
}

func Snap(P *Point, ob *Observer) *Point {
    P1 := Translate(P, ob.Pos)
    P2 := Rotate(P1, -ob.Theta)
    return Tilt(P2, math.Pi/2-ob.Phi)
}

// effectively translates Q to origin
func Relative(P *Point, Q *Point) *Point {
    return &Point{P[0]-Q[0], P[1]-Q[1], P[2]-Q[2]}
}

type Line struct {
    P *Point
    Q *Point
}

func PointInView (P *Point, ob *Observer) bool {
    // translate observer to origin
    // align x-axis so that observer is looking straight down the x-axis
    Q := Snap(P, ob)

    // calculate theta and phi
    _, th, ph := RecToSphere(Q)

    eps := 0.1

    if th + eps < -ob.HFov/2 || th - eps > ob.HFov/2 {
        return false
    }

    // translate observer's phi to pi/2

    if ph - math.Pi/2 + eps < -ob.VFov/2 || ph - math.Pi/2 - eps > ob.VFov/2 {
        return false
    }
    return true
}


func (self *Line) InView(ob *Observer) bool {
    if PointInView(self.P, ob) && PointInView(self.Q, ob) {
        return true
    }
    return false
}

func (self *Line) Draw(win *pixelgl.Window, ob *Observer) {
    if !self.InView(ob) {
        return
    }
    imd := imdraw.New(nil)
    imd.Color = colornames.Black

    relative_P := Snap(self.P, ob)
    relative_Q := Snap(self.Q, ob)

    rh1,th1,ph1 := RecToSphere(relative_P)
    x1,y1 := Project(th1,ph1,ob)

    rh2,th2,ph2 := RecToSphere(relative_Q)
    x2,y2 := Project(th2,ph2,ob)

    avg := (rh1+rh2) / 2
    imd.Push(pixel.V(x1,y1), pixel.V(x2,y2))
    imd.Line(200/avg)

    imd.Draw(win)
}

func NewLine(P,Q *Point) *Line {
    return &Line{P, Q}
}

type Cube struct {
    Lines []*Line
}

func NewCube(s float64, P *Point) *Cube {
    x := P[0]
    y := P[1]
    z := P[2]
    bottom := []*Point{
        &Point{x,y,z},
        &Point{x+s,y,z},
        &Point{x+s,y+s,z},
        &Point{x,y+s,z}}

    top := []*Point{
        &Point{x,y,z+s},
        &Point{x+s,y,z+s},
        &Point{x+s,y+s,z+s},
        &Point{x,y+s,z+s}}

    lines := []*Line{}
    for i:=0;i<4;i++ {
        v := NewLine(bottom[i], top[i])
        t := NewLine(top[i], top[(i+1)%4])
        b := NewLine(bottom[i], bottom[(i+1)%4])
        lines = append(lines, v)
        lines = append(lines, t)
        lines = append(lines, b)
    }

    return &Cube{lines}
}

func (self *Cube) Draw(win *pixelgl.Window, ob *Observer) {
    for _,line := range self.Lines {
        line.Draw(win, ob)
    }
}

type Circle struct {
    Center *Point
    Radius float64
    Active bool
}

func NewCircle(c *Point, r float64) *Circle {
    return &Circle{c, r, false}
}

func (self *Circle) Draw(win *pixelgl.Window, ob *Observer) {
    //if !self.InView(ob) {
    //    return
    //}
    imd := imdraw.New(nil)

    relative_C := Snap(self.Center, ob)
    orth := &Point{relative_C[1], -relative_C[0], 0.0}
    m := self.Radius/Magnitude(orth)
    orth = &Point{orth[0]*m, orth[1]*m}
    edge := Add(orth, relative_C)

    rh1,th1,ph1 := RecToSphere(relative_C)
    x1,y1 := Project(th1,ph1,ob)

    _,th2,ph2 := RecToSphere(edge)
    x2,y2 := Project(th2,ph2,ob)

    r := math.Sqrt(math.Pow(x1-x2,2) + math.Pow(y1-y2,2))

    x3,y3 := ob.Curs.X,ob.Curs.Y

    d := math.Sqrt(math.Pow(x1-x3,2) + math.Pow(y1-y3,2))

    if d < r {
        imd.Color = colornames.Lightblue
        self.Active = true
    } else {
        imd.Color = colornames.Navy
        self.Active = false
    }
    imd.Push(pixel.V(x1, y1))
    _ = rh1
    imd.Circle(r, 0)
    //imd.Circle(r, 200/rh1)
    //imd.Ellipse(pixel.V(400/rh1,400/rh1), 0)

    imd.Draw(win)

}

type Particle struct {
    Pos *Point
    Theta float64
    Phi float64
    V float64
    Decay float64
    Step float64
}

func NewParticle(pos *Point, theta, phi, v, decay, step float64) *Particle {
    return &Particle{pos, rand.Float64()*math.Pi*2, phi, v, decay, step}
}

func (self *Particle) Draw(win *pixelgl.Window, ob *Observer, dt float64) {
    dx := self.Step*math.Cos(self.Theta)*math.Sqrt(self.V)
    dy := self.Step*math.Sin(self.Theta)*math.Sqrt(self.V)
    dz := 0.001*GRAVITY*math.Pow(self.Step, 2) + self.V*self.Step + self.Pos[2]
    x := self.Pos[0]
    y := self.Pos[1]
    z := self.Pos[2]

    if self.Step > self.Decay {
        dx = 0
        dy = 0
        dz = 0
        self.Theta = rand.Float64()*math.Pi*2
        self.Step = 0
    }
    newPos := &Point{x+dx, y+dy, z+dz}
    c := NewCircle(newPos, 0.5)
    self.Step += dt
    //fmt.Println(newPos, self.Step)
    c.Draw(win, ob)
}

type Fountain struct {
    Pos *Point
    Parts int
    Particles []*Particle
    Theta float64
    Phi float64
    V float64
    Decay float64
    Step float64
}

func NewFountain(pos *Point, parts int, theta, phi, v float64, delay float64) *Fountain {
    decay := float64(parts) * delay
    particles := make([]*Particle, parts)
    for i:=0;i<parts;i++ {
        particles[i] = NewParticle(pos, theta, phi, v, decay, float64(i)*delay)
    }

    return &Fountain{pos, parts, particles, theta, phi, v, decay, 0}
}

func (self *Fountain) Draw(win *pixelgl.Window, ob *Observer, dt float64) {
    for _,p := range self.Particles {
        p.Draw(win, ob, dt)
    }

    //for i:=0;i<self.Parts;i++ {
    //    t := self.Step + float64(i)*0.01
    //    z := GRAVITY*math.Pow(t, 2) + self.V*t + self.Pos[2]
    //    if self.Step > self.Decay {
    //        z = 0
    //        self.Step = 0
    //    }
    //    newPos := &Point{self.Pos[0], self.Pos[1], z}
    //    c := NewCircle(newPos, 0.5)
    //    self.Step += dt
    //    //fmt.Println(newPos, self.Step)
    //    c.Draw(win, ob)
    //}
}

type Bouncing struct {
    Pos *Point
    Parts int
    Theta float64
    Phi float64
    V float64
    Decay int
    Step float64
}

func NewBouncing(pos *Point, parts int, theta, phi, v float64, decay int) *Bouncing {
    return &Bouncing{pos, parts, theta, phi, v, decay, 0}
}

func (self *Bouncing) Draw(win *pixelgl.Window, ob *Observer, dt float64) {
    z := GRAVITY*math.Pow(self.Step, 2) + self.V*self.Step + self.Pos[2]
    if z < 0 {
        z = 0
        self.Step = 0
    }
    newPos := &Point{self.Pos[0], self.Pos[1], z}
    c := NewCircle(newPos, 0.5)
    self.Step += dt
    c.Draw(win, ob)
}

type Cursor struct {
    X float64
    Y float64
    Radius float64
    Thickness float64
    Color color.RGBA
}

func NewCursor(x, y, r, t float64, col color.RGBA) *Cursor {
    return &Cursor{x, y, r, t, col}
}

func (self *Cursor) Draw(win *pixelgl.Window) {
    imd := imdraw.New(nil)
    imd.Color = self.Color
    imd.Push(pixel.V(self.X, self.Y))
    imd.Circle(self.Radius, self.Thickness)
    imd.Draw(win)
}

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
    win.SetMousePosition(win.Bounds().Center())
    win.SetCursorVisible(false)
    //col := color.RGBA{128, 128, 128, 150}
    //cur := NewCursor(WIDTH/2, HEIGHT/2, 10, 3, col)

    // start interesting code here

    origin := Point{0.0, 0.0, 0.0}
    me := NewObserver(
        FOV,
        WIDTH,
        HEIGHT,
        20.0,
        0.0,
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
            cube1 := NewCube(5, &p)
            //cube2 := NewCube(5, &q)
            //cubes = append(cubes, cube1, cube2)
            cubes = append(cubes, cube1)
        }

    }
    circ := NewCircle(&Point{30.0, 0.0, 5.0}, 0.5)
    bounce := NewBouncing(&Point{50.0, 30.0, 1.0}, 10, 0.0, 0.0, 50.0, 0)
    fount := NewFountain(&Point{20.0, -20.0, 0.0}, 100, 0.0, 0.0, 5.0, 0.5)

    alice := colornames.Aliceblue
    bg := alice
    inMenu := false
    center := win.Bounds().Center()

    atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
    txt := text.New(pixel.ZV, atlas)
    txt.Color = colornames.Brown
    txt.WriteString("jump: [space]\n")
    txt.WriteString("move: [wasd] or arrow keys\n")
    txt.WriteString("toggle cursor: [right-click]\n")
    txt.WriteString("menu: [ESC]\n")

    scoreTxt := text.New(pixel.ZV, atlas)
    scoreTxt.Color = colornames.Black


    menu := NewMenu(atlas, []string{"resume", "save", "exit"})

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
            win.SetMousePosition(center)
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
            if circ.Active {
                me.Score += 1
            }
        }


        if win.JustPressed(pixelgl.MouseButtonRight) {
            if me.Locked {
                me.Locked = false
            } else {
                me.Locked = true
            }

        }
        // align to mouse
        pos := win.MousePosition()
        win.SetMousePosition(center)

        // compute mouse distance traveled
        dx := pos.X - center.X
        dy := pos.Y - center.Y
        if me.Locked {
            me.Curs.Color.A = 255
            me.Curs.X += dx*dt*100
            me.Curs.Y += dy*dt*100
            if me.Curs.X > me.Width {
                me.Curs.X = me.Width
            }
            if me.Curs.X < 0 {
                me.Curs.X = 0
            }
            if me.Curs.Y > me.Height {
                me.Curs.Y = me.Height
            }
            if me.Curs.Y < 0 {
                me.Curs.Y = 0
            }
        } else {
            me.Curs.Color.A = 150

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
        me.Curs.Draw(win)

        for _,cube := range cubes {
            cube.Draw(win, me)
        }
        circ.Draw(win, me)
        bounce.Draw(win, me, dt)
        fount.Draw(win, me, dt)

        // update
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
    log.Println("Done")
}
