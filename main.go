package main

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func initScreen() {
	primaryMonitor = pixelgl.PrimaryMonitor()
	cfg := pixelgl.WindowConfig{
		Title:   title,
		Bounds:  pixel.R(0, 0, float64(screenWidth), float64(screenHeight)),
		Monitor: nil,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	window = win
}

// setFPS allows us to set max frames per second.
// Disable any maximum by passing 0.
func setFPS(fps int) {
	if fps <= 0 {
		frameTick = nil
	} else {
		frameTick = time.NewTicker(time.Second / time.Duration(fps))
	}
}

// update handles all logic changes in the game. This
// includes moving objects or handling input.
func update(dt float64) {

}

// draw is called after update and just draws
// everything visible to the screen.
func draw() {
	window.Clear(colornames.Black)

	// Draw HUD
	fpsText.Clear()
	fpsText.WriteString(fmt.Sprintf("FPS: %d", int(math.Round(fps))))
	fpsText.Draw(window, pixel.IM)
}

func run() {
	// First call all init functions to setup the game.
	initScreen()
	setFPS(0)

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	fpsText = text.New(pixel.V(10, window.Bounds().H()-20), atlas)
	fpsText.Color = colornames.Antiquewhite

	start := time.Now()
	now := start
	for !window.Closed() {
		last := now
		now = time.Now()
		dt := now.Sub(last).Seconds()

		fps = float64(frames) / now.Sub(start).Seconds()

		update(dt)
		draw()

		frames++
		window.Update()

		if frameTick != nil {
			<-frameTick.C
		}
	}
}

func main() {
	pixelgl.Run(run)
}
