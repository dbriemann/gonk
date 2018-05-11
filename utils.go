package main

import (
	"math/rand"

	"github.com/faiface/pixel"
)

// genPlanetParameters generates random numbers for all parameters of a planet
// the valid values / ranges are passed in as arrays.
func genPlanetParameters(sizes []int) (size, vel, dir float64) {
	size = float64(sizes[rand.Intn(len(sizes))])
	vel = (rand.Float64() * 7) + 3
	// dir is just 1 or -1, which determines if a planet moves
	// clockwise or counter-clockwise.
	dir = float64(rand.Intn(2)*2 - 1)
	return
}

// rotatePoint rotates point around anchor by angle omega (rad).
func rotatePoint(anchor, point *pixel.Vec, omega float64) {
	mat := pixel.IM.Rotated(*anchor, omega)
	npos := mat.Project(*point)

	point.X = npos.X
	point.Y = npos.Y
}
