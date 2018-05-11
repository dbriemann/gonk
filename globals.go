package main

import (
	"time"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var (
	primaryMonitor *pixelgl.Monitor
	window         *pixelgl.Window

	frameTick *time.Ticker
	fps       float64

	screenWidth  = 1200
	screenHeight = 800
	title        = "Gonk"

	worldCanvas *pixelgl.Canvas
	imd         *imdraw.IMDraw

	// Our 'camera' targets (0,0) which will be the center of the screen.
	camPos = pixel.ZV
	cam    pixel.Matrix

	planets []*planet
	players []player

	origin         = &pixel.Vec{X: 0, Y: 0}
	planetSizes    = []int{7, 8, 9}
	satelliteSizes = []int{3, 4, 5}
	recycledShips  = []*ship{}

	productionFactor = 0.1

	frames  uint64
	fpsText *text.Text
)
