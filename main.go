package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Game implements ebiten.Game interface.
type Game struct {
	asteroids    []*PolygonObject
	player       *PolygonObject
	screenWidth  float64
	screenHeight float64
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Handle player input
	g.handlePlayerInput()

	// Update player with wrapping
	g.player.Update(g.screenWidth, g.screenHeight)

	// Update all asteroids with wrapping
	for _, asteroid := range g.asteroids {
		asteroid.Update(g.screenWidth, g.screenHeight)
	}
	return nil
}

// handlePlayerInput processes keyboard input for player movement
func (g *Game) handlePlayerInput() {
	const rotationSpeed = 0.1 // radians per frame
	const acceleration = 0.2  // pixels per frame squared
	const maxSpeed = 5.0      // maximum speed
	const friction = 0.98     // velocity decay factor

	// Rotation controls
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.Rotation -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.Rotation += rotationSpeed
	}

	// Forward/backward thrust
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		// Accelerate in the direction the ship is facing
		thrustX := math.Sin(g.player.Rotation) * acceleration
		thrustY := -math.Cos(g.player.Rotation) * acceleration
		g.player.Velocity.X += thrustX
		g.player.Velocity.Y += thrustY
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		// Decelerate (reverse thrust)
		thrustX := math.Sin(g.player.Rotation) * acceleration * 0.5
		thrustY := -math.Cos(g.player.Rotation) * acceleration * 0.5
		g.player.Velocity.X -= thrustX
		g.player.Velocity.Y -= thrustY
	}

	// Apply friction to gradually slow down the ship
	g.player.Velocity.X *= friction
	g.player.Velocity.Y *= friction

	// Limit maximum speed
	speed := math.Sqrt(g.player.Velocity.X*g.player.Velocity.X + g.player.Velocity.Y*g.player.Velocity.Y)
	if speed > maxSpeed {
		g.player.Velocity.X = (g.player.Velocity.X / speed) * maxSpeed
		g.player.Velocity.Y = (g.player.Velocity.Y / speed) * maxSpeed
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw player ship
	g.player.DrawWithWrapping(screen, g.screenWidth, g.screenHeight)

	// Draw all asteroids with wrapping
	for _, asteroid := range g.asteroids {
		asteroid.DrawWithWrapping(screen, g.screenWidth, g.screenHeight)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.screenWidth), int(g.screenHeight)
}

// NewGame creates a new game instance with initialized asteroids and player
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		asteroids:    nil,
		screenWidth:  800,
		screenHeight: 600,
	}

	// Create player ship (triangle)
	game.player = CreateTriangle(15)                                 // 15 pixel triangle
	game.player.SetPosition(game.screenWidth/2, game.screenHeight/2) // Center of screen
	game.player.SetColor(color.White)

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
			50+rand.Float64()*(game.screenWidth-100),  // X between 50 and 750
			50+rand.Float64()*(game.screenHeight-100), // Y between 50 and 550
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

		game.asteroids = append(game.asteroids, asteroid)
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
