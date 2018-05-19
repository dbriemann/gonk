package main

import (
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	opensimplex "github.com/ojrac/opensimplex-go"
	"golang.org/x/image/colornames"
)

// see https://cmaher.github.io/posts/working-with-simplex-noise/
func octaveNoise(iterations int, x, y, persistence, scale, low, high float64) (result float64) {
	maxAmp := 0.0
	amp := 1.0
	freq := scale

	for i := 0; i < iterations; i++ {
		result += noise.Eval2(x*freq, y*freq) * amp
		maxAmp += amp
		amp *= persistence
		freq *= 2
	}

	result /= maxAmp

	result = result*(high-low)/2 + (high+low)/2
	return
}

func brighten(val uint8, factor float64) uint8 {
	r := float64(val) * factor
	if uint8(r) < val {
		return 255
	}
	return uint8(r)
}

func genGradientDisc(radius, density float64, c color.Color) (canvas *pixelgl.Canvas) {
	cr, cg, cb, ca := c.RGBA()
	size := int(radius*2 + 1)
	canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(size), float64(size)))
	pixels := canvas.Pixels()

	ncol := pixel.RGBA{
		R: float64(cr) / 0xffff,
		G: float64(cg) / 0xffff,
		B: float64(cb) / 0xffff,
		A: float64(ca) / 0xffff,
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dist := pixel.V(float64(x), float64(y)).Sub(pixel.V(radius, radius)).Len()
			factor := (dist - radius*density) / (radius * (1 - density))
			factor = math.Min(1, math.Max(0, factor)) // clamp

			index := y*size*4 + x*4
			pixels[index] = uint8(ncol.R * (1 - factor) * 255)
			pixels[index+1] = uint8(ncol.G * (1 - factor) * 255)
			pixels[index+2] = uint8(ncol.B * (1 - factor) * 255)
			pixels[index+3] = uint8(ncol.A * (1 - factor) * 255)
		}
	}

	canvas.SetPixels(pixels)

	return
}

func genMoon(radius float64) (canvas *pixelgl.Canvas) {
	noise = opensimplex.NewWithSeed(time.Now().UnixNano())
	size := int(radius*2 + 1)
	canvas = genGradientDisc(radius, 0.95, colornames.White)
	pixels := canvas.Pixels()

	scale := radius / (1000 * (radius / 40) * (radius / 40))

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			index := y*size*4 + x*4
			r, g, b, a := float64(pixels[index]), float64(pixels[index+1]), float64(pixels[index+2]), float64(pixels[index+3])

			if a > 0 {
				n := octaveNoise(16, float64(x), float64(y), 0.5, scale, 0, 1)

				pixels[index] = brighten(uint8(r*n), 1.5)
				pixels[index+1] = brighten(uint8(g*n), 1.5)
				pixels[index+2] = brighten(uint8(b*n), 1.5)
				pixels[index+3] = 255 // Make the planet opaque
			}
		}
	}

	canvas.SetPixels(pixels)

	return
}
