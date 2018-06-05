package main

import (
	"time"

	opensimplex "github.com/ojrac/opensimplex-go"

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

	// Our 'camera' targets (0,0) which will be the center of the screen.
	camPos = pixel.ZV
	cam    pixel.Matrix

	planets []*planet
	players []player

	sprites struct {
		// TODO planets -> one canvas -> spritesheet -> batch
		planets []*pixelgl.Canvas
		sun     *pixelgl.Canvas
		ship    *pixelgl.Canvas
	}

	batches struct {
		ships *pixel.Batch
	}

	origin         = &pixel.Vec{X: 0, Y: 0}
	planetSizes    = []int{9, 10, 11}
	satelliteSizes = []int{5, 6, 7}
	recycledShips  = []*ship{}

	productionFactor = 0.1

	frames      uint64
	fpsText     *text.Text
	objectsText *text.Text
	objectCount uint64 = 1 // Includes the sun at the start.

	noise *opensimplex.Noise
)
