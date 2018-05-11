package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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
	worldCanvas = pixelgl.NewCanvas(win.Bounds())
	// Set the camera to look at camPos.
	cam = pixel.IM.Moved(worldCanvas.Bounds().Center().Sub(camPos))
	worldCanvas.SetMatrix(cam)
}

func initPlayers(playerName string, ais int) {
	players = []player{
		// A pseudo player that represents 'no player'.
		player{
			id:    0,
			name:  "not occupied",
			ai:    false,
			color: colornames.Antiquewhite,
		},
		player{
			id:    1,
			name:  playerName,
			ai:    false,
			color: colornames.Skyblue,
		},
	}
}

// We distribute the planets homogeneously on the X axis starting inside of the given range (span).
func initSolarSystem(planetAmount, maxSatellites, minDist, maxDist int) {
	span := maxDist - minDist
	step := span / planetAmount
	current := minDist

	for i := 0; i < planetAmount; i++ {
		size, vel, dir := genPlanetParameters(planetSizes)
		p := newPlanet(float64(current), size, dir, pixel.V(vel, vel), origin, &players[0])
		// Add a little random adjustment to the planet's position to make
		// it look less static.
		shift := float64(rand.Intn(step/3)*2 - step/3)
		p.pos.X += shift
		// The planet is generated. Add it to our global planets slice.
		planets = append(planets, p)

		// Now we do more or less the same again as above. Just this time we are adding satellites
		// which orbit the previously generated planet.
		sats := rand.Intn(maxSatellites + 1)

		for s := 0; s < sats; s++ {
			size, vel, dir := genPlanetParameters(satelliteSizes)
			sat := newPlanet(float64((s+1)*20), size, dir, pixel.V(vel, vel), &p.pos, &players[0])
			sat.rotate(rand.Float64() * sat.dist)
			p.satellites = append(p.satellites, sat)
			planets = append(planets, sat)
		}

		// Now that the planet and its satellites exist we rotate them randomly
		// to achieve a nice distribution "on the clock".
		p.rotateGroup(rand.Float64() * p.dist)

		// Next planet please..
		current += step
	}
}

func makeBasicShapes() {
	imd = imdraw.New(nil)
	// Draw the sun.
	imd.Color = colornames.Gold
	imd.Push(pixel.ZV)
	imd.Circle(20, 0)
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
	for i := 0; i < len(planets); i++ {
		planets[i].update(dt)
	}
}

// draw is called after update and just draws
// everything visible to the screen.
func draw() {
	// Clear everything before drawing.
	window.Clear(colornames.Black)
	worldCanvas.Clear(colornames.Black)

	// Draw the game objects onto the canvas.
	makeBasicShapes()
	for _, p := range planets {
		p.draw()
	}
	imd.Draw(worldCanvas)

	// Draw the canvas onto the window.
	worldCanvas.Draw(window, cam)

	// Draw HUD to window not canvas so we can use screen coordinates directly.
	fpsText.Clear()
	fpsText.WriteString(fmt.Sprintf("FPS: %d", int(math.Round(fps))))
	fpsText.Draw(window, pixel.IM)

}

func run() {
	// First call all init functions to setup the game.
	initScreen()
	setFPS(0)
	initPlayers("RagingDave", 0)
	initSolarSystem(12, 3, 100, int(screenHeight/2))

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
