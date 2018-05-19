package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type planet struct {
	orb
	*player

	satellites    []*planet
	ships         []*ship
	shipsProduced float64
	shipAngleMod  float64
	radius        float64

	sprite *pixelgl.Canvas
}

func newPlanet(dist, radius, dir float64, vel pixel.Vec, anchor *pixel.Vec, player *player, sprite *pixelgl.Canvas) *planet {
	p := &planet{
		orb: orb{
			dist:   dist,
			anchor: anchor,
			vel:    vel,
			dir:    dir,
		},
		player:     player,
		radius:     radius,
		satellites: []*planet{},
		ships:      make([]*ship, int(radius/3)),
		sprite:     sprite,
	}

	if anchor != nil {
		p.pos.X = anchor.X + dist
	} else {
		p.pos.X = dist
	}

	for i := 0; i < len(p.ships); i++ {
		p.ships[i] = newShip(p, player)
	}
	p.setShips(0)

	objectCount++

	return p
}

// rotateGroup rotates the planet and adjusts the position of its satellites accordingly.
func (p *planet) rotateGroup(dt float64) {
	dvec := p.rotate(dt)
	for i := 0; i < len(p.satellites); i++ {
		p.satellites[i].pos.X += dvec.X
		p.satellites[i].pos.Y += dvec.Y
	}
}

func (p *planet) update(dt float64) {
	p.rotateGroup(dt)
	// Ship production depends on planet size: production = sqrt(radius)/5
	prod := math.Sqrt(p.radius) * productionFactor
	p.shipsProduced += prod * dt

	// Add new ships to slice.
	for i := 0; i < int(p.shipsProduced); i++ {
		added := false
		// Search a free spot and if there is none append.
		nship := newShip(p, p.player)
		for j := 0; j < len(p.ships); j++ {
			if p.ships[i] == nil {
				p.ships[i] = nship
				added = true
			}
		}
		if !added {
			p.ships = append(p.ships, nship)
		}
		p.shipsProduced--
	}

	p.setShips(dt)
}

func (p *planet) draw(translation pixel.Matrix) {
	// TODO magic numbers
	p.sprite.DrawColorMask(worldCanvas, pixel.IM.Moved(p.pos).Scaled(p.pos, p.radius/30), nil)

	// Draw all ships stationed at this planet.
	for _, s := range p.ships {
		s.draw()
	}
}

// distributeShips evenly distributes ships around a planet.
func (p *planet) setShips(dt float64) {
	amount := len(p.ships)
	step := (2 * math.Pi) / float64(amount)
	p.shipAngleMod += dt * 0.5
	if p.shipAngleMod > 2*math.Pi {
		p.shipAngleMod -= 2 * math.Pi
	}

	for i := 0; i < amount; i++ {
		p.ships[i].pos.X = p.pos.X + p.ships[i].dist
		p.ships[i].pos.Y = p.pos.Y

		omega := float64(i) * step
		rotatePoint(&p.pos, &p.ships[i].pos, omega+p.shipAngleMod)
	}
}

type ship struct {
	orb
	*player
}

func newShip(planet *planet, player *player) *ship {
	// Test if there is a ship available for recycling.
	var sp *ship
	i := -1
	for i, sp = range recycledShips {
		if sp != nil {
			break
		}
	}
	if sp != nil && i >= 0 {
		// Remove ship from recycled slice.
		recycledShips[i] = nil
	} else {
		// Create new ship.
		sp = &ship{}
	}

	objectCount++

	// TODO remove magic numbers
	sp.dist = planet.radius * 2
	sp.anchor = &planet.pos
	sp.vel = pixel.V(5, 5)
	sp.dir = 1
	sp.player = player

	return sp
}

func (s *ship) draw() {
	sprites.ship.Draw(batches.ships, pixel.IM.Moved(s.pos))
}
