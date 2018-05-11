package main

import (
	"image/color"

	"github.com/faiface/pixel"
)

type player struct {
	id    int
	name  string
	ai    bool
	color color.Color
}

type orb struct {
	anchor *pixel.Vec
	vel    pixel.Vec
	pos    pixel.Vec
	dir    float64
	dist   float64
}

// rotate rotates an orb around its anchor and returns the shift vector.
// The angle is destined by the orb's velocity and distance to the anchor.
func (o *orb) rotate(dt float64) (shift pixel.Vec) {
	mat := pixel.IM
	len := o.vel.Len()
	omega := o.dir * len / o.dist

	mat = mat.Rotated(*o.anchor, omega*dt)
	npos := mat.Project(o.pos)
	shift.X, shift.Y = npos.X-o.pos.X, npos.Y-o.pos.Y
	o.pos.X, o.pos.Y = npos.XY()
	return
}
