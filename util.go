package main

import "math"

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
