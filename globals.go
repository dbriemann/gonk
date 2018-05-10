package main

import (
	"time"

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

	frames  uint64
	fpsText *text.Text
)
