package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Bullet represents a projectile fired by the player
type Bullet struct {
	polygon *PolygonObject
}

// Game implements ebiten.Game interface.
type Game struct {
	asteroids      []*PolygonObject
	player         *PolygonObject
	bullets        []*Bullet
	screenWidth    float64
	screenHeight   float64
	lastBulletTime time.Time
	bulletCooldown time.Duration
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Handle player input
	g.handlePlayerInput()

	// Update player with wrapping
	g.player.Update(g.screenWidth, g.screenHeight, true)

	// Update all asteroids with wrapping
	for _, asteroid := range g.asteroids {
		asteroid.Update(g.screenWidth, g.screenHeight, true)
	}

	// Update bullets
	g.updateBullets()

	// Check collisions
	g.checkCollisions()

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

	// Shooting
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		now := time.Now()
		if now.Sub(g.lastBulletTime) > g.bulletCooldown {
			g.createBullet()
			g.lastBulletTime = now
		}
	}
}

// createBullet creates a new bullet at the tip of the player ship
func (g *Game) createBullet() {
	// Calculate the tip position of the player triangle
	tipOffset := 15.0 // Same as triangle size
	tipX := g.player.Position.X + math.Sin(g.player.Rotation)*tipOffset
	tipY := g.player.Position.Y - math.Cos(g.player.Rotation)*tipOffset

	// Create a small rectangle for the bullet (2x2)
	bulletPolygon := &PolygonObject{
		Vertices: []Vector2{
			{X: -1, Y: -1}, // Top left
			{X: 1, Y: -1},  // Top right
			{X: 1, Y: 1},   // Bottom right
			{X: -1, Y: 1},  // Bottom left
		},
		Position:      Vector2{X: tipX, Y: tipY},
		Velocity:      Vector2{X: 0, Y: 0},
		Rotation:      0,
		RotationSpeed: 0,
		Scale:         1.0,
		Color:         color.White,
		LineWidth:     1.0,
	}

	// Set bullet velocity in the direction the player is facing
	const bulletSpeed = 8.0
	bulletPolygon.Velocity.X = math.Sin(g.player.Rotation) * bulletSpeed
	bulletPolygon.Velocity.Y = -math.Cos(g.player.Rotation) * bulletSpeed

	// Add player's velocity to bullet (inherit momentum)
	bulletPolygon.Velocity.X += g.player.Velocity.X
	bulletPolygon.Velocity.Y += g.player.Velocity.Y

	bullet := &Bullet{polygon: bulletPolygon}
	g.bullets = append(g.bullets, bullet)
}

// updateBullets updates all bullets and removes those that have left the screen
func (g *Game) updateBullets() {
	// Update bullet positions
	for _, bullet := range g.bullets {
		bullet.polygon.Update(g.screenWidth, g.screenHeight, false)
	}

	// Remove bullets that are off-screen (with some margin for safety)
	margin := 50.0
	var activeBullets []*Bullet
	for _, bullet := range g.bullets {
		pos := bullet.polygon.Position
		if pos.X >= -margin && pos.X <= g.screenWidth+margin &&
			pos.Y >= -margin && pos.Y <= g.screenHeight+margin {
			activeBullets = append(activeBullets, bullet)
		}
	}
	g.bullets = activeBullets
}

// checkCollisions handles all collision detection in the game
func (g *Game) checkCollisions() {
	// Check bullet-asteroid collisions
	for i := len(g.bullets) - 1; i >= 0; i-- {
		bullet := g.bullets[i]
		bulletHit := false

		for j := len(g.asteroids) - 1; j >= 0; j-- {
			asteroid := g.asteroids[j]

			if PolygonsCollide(bullet.polygon, asteroid) {
				// Remove the bullet
				g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
				// Remove the asteroid
				g.asteroids = append(g.asteroids[:j], g.asteroids[j+1:]...)
				bulletHit = true
				break
			}
		}

		if bulletHit {
			break // Move to next bullet since this one was removed
		}
	}

	// Check player-asteroid collisions
	for _, asteroid := range g.asteroids {
		if PolygonsCollide(g.player, asteroid) {
			// For now, just reset player position to center
			// In a real game, you might handle lives, explosions, etc.
			g.player.SetPosition(g.screenWidth/2, g.screenHeight/2)
			g.player.Velocity = Vector2{X: 0, Y: 0}
			break
		}
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw player ship
	g.player.Draw(screen)

	// Draw all asteroids
	for _, asteroid := range g.asteroids {
		asteroid.Draw(screen)
	}

	// Draw all bullets
	for _, bullet := range g.bullets {
		bullet.polygon.Draw(screen)
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
		asteroids:      nil,
		bullets:        nil,
		screenWidth:    800,
		screenHeight:   600,
		lastBulletTime: time.Now(),
		bulletCooldown: 100 * time.Millisecond, // 100ms cooldown
	}

	// Create player ship (triangle)
	game.player = CreateTriangle(15)                                 // 15 pixel triangle
	game.player.SetPosition(game.screenWidth/2, game.screenHeight/2) // Center of screen
	blue := color.RGBA{0, 0, 255, 255}                               // Blue color
	game.player.SetColor(blue)

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
