package main

import "math"

type Point = [3]float64

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

func Negate(P *Point) *Point {
    return &Point{-P[0], -P[1], -P[2]}
}

func Subtract(P *Point, Q *Point) *Point {
    return Add(P, Negate(Q))
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

// translate Q to the origin, face the x-axis
func Snap(P, Q *Point, theta, phi float64) *Point {
    P1 := Translate(P, Q)
    P2 := Rotate(P1, -theta)
    return Tilt(P2, math.Pi/2-phi)
}

// align a to the nearest b (whether up or down)
func Align(a, b float64) (float64,int) {
    sign := 1.0
    if a < 0 {
        sign = -1.0
    }
    a = math.Abs(a)
    lowIndex := int(a/b)
    highIndex := int((a+b)/b)
    low := float64(lowIndex) * b
    high := float64(highIndex) * b
    highDiff := math.Abs(high - a)
    lowDiff := math.Abs(a - low)
    if highDiff <= lowDiff {
        return sign*high, int(sign)*highIndex
    }
    return sign*low, int(sign)*lowIndex
}

func Min(nums ...float64) float64 {
    if len(nums) == 0 {
        panic(nil)
    }
    m := nums[0]
    for _,num := range nums {
        if num <= m {
            m = num
        }
    }
    return m
}

func Max(nums ...float64) float64 {
    if len(nums) == 0 {
        panic(nil)
    }
    m := nums[0]
    for _,num := range nums {
        if num >= m {
            m = num
        }
    }
    return m
}
