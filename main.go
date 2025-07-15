package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements ebiten.Game interface.
type Game struct {
	asteroids []*PolygonObject
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Update all asteroids (movement and rotation)
	for _, asteroid := range g.asteroids {
		asteroid.Update()
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw all asteroids
	for _, asteroid := range g.asteroids {
		asteroid.Draw(screen)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

// NewGame creates a new game instance with initialized asteroids
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		asteroids: make([]*PolygonObject, 3),
	}

	// Create 3 random asteroids
	for i := 0; i < 3; i++ {
		// Random base radius between 20 and 50
		baseRadius := 20.0 + rand.Float64()*30.0
		// Random irregularity between 5 and 15
		irregularity := 5.0 + rand.Float64()*10.0
		// Random number of vertices between 6 and 12
		numVertices := 6 + rand.Intn(7)

		asteroid := CreateAsteroid(baseRadius, irregularity, numVertices)

		// Random position within the screen bounds (with some margin)
		asteroid.SetPosition(
			50+rand.Float64()*(800-100), // X between 50 and 750
			50+rand.Float64()*(600-100), // Y between 50 and 550
		)

		// Random rotation
		asteroid.SetRotation(rand.Float64() * 6.28) // 0 to 2π radians

		// Random velocity (pixels per frame)
		vx := (rand.Float64() - 0.5) * 4 // -2 to 2 pixels per frame
		vy := (rand.Float64() - 0.5) * 4 // -2 to 2 pixels per frame
		asteroid.SetVelocity(vx, vy)

		// Random rotation speed (radians per frame)
		rotSpeed := (rand.Float64() - 0.5) * 0.1 // -0.05 to 0.05 radians per frame
		asteroid.SetRotationSpeed(rotSpeed)

		// Set color to white
		asteroid.SetColor(color.White)

		game.asteroids[i] = asteroid
	}

	return game
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Asteroids Game")

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
