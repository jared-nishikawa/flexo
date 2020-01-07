package main

import (
    "math"
)

//const FOV = math.Pi/3
//const WIDTH = 1920.0
////const WIDTH = 800.0
//const HEIGHT = 1080.0
////const HEIGHT = 600.0
//const GRAVITY = -50
//const SENSITIVITY = 0.1

type Observer struct {
    HFov float64
    VFov float64
    Width float64
    Height float64
    Speed float64
    VerticalSpeed float64
    Gravity float64
    Theta float64
    Phi float64
    Pos *Point
    Locked bool
    Score int
}

func NewObserver(fov, w, h, s, vs, g, th, ph float64, p *Point, l bool) *Observer {
    vfov := h/w * fov
    return &Observer{fov, vfov, w, h, s, vs, g, th, ph, p, l, 0}
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
        self.VerticalSpeed += self.Gravity*dt
        if self.Pos[2] < 0 {
            self.Pos[2] = 0
        }
    } else {
        self.VerticalSpeed = 0
    }
}




func (self *Observer) Project(theta, phi float64) (float64, float64) {
    theta = math.Mod(theta, 2*math.Pi)
    if theta < -math.Pi {
        theta += 2*math.Pi
    }

    if theta > math.Pi {
        theta -= 2*math.Pi
    }

    x := self.Width  * (((-theta/(self.HFov/2))+1)/2)
    y := self.Height * (((-(phi-math.Pi/2)/(self.VFov/2))+1)/2)

    return x,y
}

func (self *Observer) Snap(P *Point) *Point {
    P1 := Translate(P, self.Pos)
    P2 := Rotate(P1, -self.Theta)
    return Tilt(P2, math.Pi/2-self.Phi)
}

func (self *Observer) PointInView (P *Point) bool {
    // translate observer to origin
    // align x-axis so that observer is looking straight down the x-axis
    Q := self.Snap(P)

    // calculate theta and phi
    _, th, ph := RecToSphere(Q)

    eps := 0.1

    if th + eps < -self.HFov/2 || th - eps > self.HFov/2 {
        return false
    }

    // translate observer's phi to pi/2

    if ph - math.Pi/2 + eps < -self.VFov/2 || ph - math.Pi/2 - eps > self.VFov/2 {
        return false
    }
    return true
}

