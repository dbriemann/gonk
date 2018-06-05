package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
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

// func noiseStripe() (canvas *pixelgl.Canvas) {
// 	xsize, ysize := 400, 100
// 	canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(xsize), float64(ysize)))
// 	pixels := canvas.Pixels()

// 	for y := 0; y < ysize; y++ {
// 		for x := 0; x < xsize; x++ {
// 			noise := uint8(layerNoise(8, float64(x), float64(y), 0.5, 0.02, 0, 255))
// 			index := y*xsize*4 + x*4
// 			pixels[index] = noise
// 			pixels[index+1] = noise
// 			pixels[index+2] = noise
// 			pixels[index+3] = 255
// 		}
// 	}
// 	canvas.SetPixels(pixels)

// 	return
// }

func genSprites(planets int) {
	for i := 0; i < planets; i++ {
		sprite := genPlanet(30)
		sprites.planets = append(sprites.planets, sprite)
	}
	sprites.sun = genGradientDisc(30, 0.6, colornames.Gold)
	sprites.ship = genGradientDisc(16, 0.95, colornames.White)
	batches.ships = pixel.NewBatch(&pixel.TrianglesData{}, sprites.ship)
}

func initSolarSystem(planetAmount, maxSatellites, minDist, maxDist int) {
	// We distribute the planets homogeneously on the X axis inside the given range (span).
	span := maxDist - minDist
	step := span / planetAmount
	current := minDist

	for i := 0; i < planetAmount; i++ {
		size, vel, dir := genPlanetParameters(planetSizes)
		r := rand.Intn(len(sprites.planets))
		p := newPlanet(float64(current), size, dir, pixel.V(vel, vel), origin, &players[0], sprites.planets[r])
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
			r = rand.Intn(len(sprites.planets))
			sat := newPlanet(float64((s+1)*20), size, dir, pixel.V(vel, vel), &p.pos, &players[0], sprites.planets[r])
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
	worldCanvas.Clear(pixel.Alpha(0))
	batches.ships.Clear()

	// Draw the game objects onto the canvas.
	sprites.sun.Draw(worldCanvas, pixel.IM)
	for _, p := range planets {
		p.draw(cam)
	}

	batches.ships.Draw(worldCanvas)

	// // Draw the canvas onto the window.
	worldCanvas.Draw(window, cam)

	// Draw HUD to window not canvas so we can use screen coordinates directly.
	fpsText.Clear()
	fpsText.WriteString(fmt.Sprintf("FPS: %d", int(math.Round(fps))))
	fpsText.Draw(window, pixel.IM)
	objectsText.Clear()
	objectsText.WriteString(fmt.Sprintf("Objects: %d", objectCount))
	objectsText.Draw(window, pixel.IM)
}

func run() {
	rand.Seed(time.Now().UnixNano())

	// First call all init functions to setup the game.
	initScreen()
	setFPS(0)
	initPlayers("RagingDave", 0)
	genSprites(10)
	initSolarSystem(8, 3, 100, int(screenHeight/2))

	// TODO init texts in extra function at some point.
	fpsText = text.New(pixel.V(10, window.Bounds().H()-20), text.Atlas7x13)
	fpsText.Color = colornames.Antiquewhite
	objectsText = text.New(pixel.V(10, window.Bounds().H()-40), text.Atlas7x13)
	objectsText.Color = colornames.Antiquewhite

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
