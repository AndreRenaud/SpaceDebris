package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// GameState represents the current state of the game
type GameState int

const (
	GameStatePlaying GameState = iota
	GameStateGameOver
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
	score          int
	vectorFont     *VectorFont
	state          GameState
	gameOverReason string
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	switch g.state {
	case GameStatePlaying:
		return g.updatePlaying()
	case GameStateGameOver:
		return g.updateGameOver()
	}
	return nil
}

// updatePlaying handles the game logic when playing
func (g *Game) updatePlaying() error {
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

	// Check win condition (all asteroids destroyed)
	if len(g.asteroids) == 0 {
		g.state = GameStateGameOver
		g.gameOverReason = "YOU WIN!"
	}

	return nil
}

// updateGameOver handles the game logic when in game over state
func (g *Game) updateGameOver() error {
	// Check for restart input
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.Restart()
		g.state = GameStatePlaying
		g.gameOverReason = ""
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
		Trail:         make([]drawablePolygon, 0, ghostTrailLength),
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

				// Increment score for hitting an asteroid
				g.score++

				// Split the asteroid or remove it if too small
				g.splitAsteroid(j)

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
			// Set game over state
			g.state = GameStateGameOver
			g.gameOverReason = "GAME OVER"

			// Start a red flash fade effect for 1 second (60 frames)
			redFlash := color.RGBA{255, 50, 50, 255}
			blue := color.RGBA{0, 0, 255, 255} // Blue color
			g.player.SetColor(redFlash)
			g.player.StartFade(blue, 60)

			break
		}
	}
}

// splitAsteroid splits an asteroid into two smaller ones or removes it if too small
func (g *Game) splitAsteroid(asteroidIndex int) {
	asteroid := g.asteroids[asteroidIndex]

	// Calculate current size (approximate radius)
	bbox := asteroid.GetBoundingBox()
	currentSize := (bbox.MaxX - bbox.MinX + bbox.MaxY - bbox.MinY) / 4 // Average of width and height, divided by 2

	const minSize = 15.0 // Minimum size threshold

	if currentSize < minSize {
		// Remove asteroid if too small
		g.asteroids = append(g.asteroids[:asteroidIndex], g.asteroids[asteroidIndex+1:]...)
		return
	}

	// Create two smaller asteroids
	newSize := currentSize * 0.6    // Make them 60% of original size
	irregularity := newSize * 0.3   // Proportional irregularity
	numVertices := 6 + rand.Intn(5) // 6-10 vertices

	// Create first smaller asteroid
	asteroid1 := CreateAsteroid(newSize, irregularity, numVertices)
	asteroid1.SetPosition(asteroid.Position.X-newSize*0.5, asteroid.Position.Y-newSize*0.5)
	asteroid1.SetColor(asteroid.Color)

	// Give it some velocity based on original velocity plus some random spread
	vel1X := asteroid.Velocity.X + (rand.Float64()-0.5)*2
	vel1Y := asteroid.Velocity.Y + (rand.Float64()-0.5)*2
	asteroid1.SetVelocity(vel1X, vel1Y)
	asteroid1.SetRotationSpeed((rand.Float64() - 0.5) * 0.15)

	// Start a fade from white to red over 2 seconds (120 frames at 60 FPS)
	redColor := color.RGBA{255, 100, 100, 255}
	asteroid1.SetColor(redColor)
	asteroid1.StartFade(color.White, 120)

	// Create second smaller asteroid
	asteroid2 := CreateAsteroid(newSize, irregularity, numVertices)
	asteroid2.SetPosition(asteroid.Position.X+newSize*0.5, asteroid.Position.Y+newSize*0.5)
	asteroid2.SetColor(asteroid.Color)

	// Give it velocity in roughly opposite direction
	vel2X := asteroid.Velocity.X + (rand.Float64()-0.5)*2
	vel2Y := asteroid.Velocity.Y + (rand.Float64()-0.5)*2
	asteroid2.SetVelocity(vel2X, vel2Y)
	asteroid2.SetRotationSpeed((rand.Float64() - 0.5) * 0.15)

	// Start Pulse red
	asteroid2.SetColor(redColor)
	asteroid2.StartFade(color.White, 120)

	// Remove the original asteroid
	g.asteroids = append(g.asteroids[:asteroidIndex], g.asteroids[asteroidIndex+1:]...)

	// Add the two new asteroids
	g.asteroids = append(g.asteroids, asteroid1, asteroid2)
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

	// Draw score in top-right corner
	scoreStr := fmt.Sprintf("%d", g.score)
	scoreWidth := g.vectorFont.GetWidth(scoreStr)
	scoreX := float32(g.screenWidth) - scoreWidth - 20 // 20 pixels from right edge
	scoreY := float32(20)                              // 20 pixels from top
	g.vectorFont.DrawString(screen, scoreStr, scoreX, scoreY)

	// Draw game over screen if in game over state
	if g.state == GameStateGameOver {
		g.drawGameOverScreen(screen)
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
		screenWidth:    800,
		screenHeight:   600,
		bulletCooldown: 100 * time.Millisecond,                // 100ms cooldown
		vectorFont:     NewVectorFont(16, 24, 3, color.White), // 16x24 digit size, 2px line width, white color
	}

	// Use Restart to initialize the game state
	game.Restart()

	return game
}

// Restart resets the game state to initial conditions
func (g *Game) Restart() {
	// Reset game state
	g.state = GameStatePlaying
	g.gameOverReason = ""

	// Reset score
	g.score = 0

	// Clear all bullets and asteroids
	g.bullets = nil
	g.asteroids = nil

	// Reset bullet timing
	g.lastBulletTime = time.Now()

	// Create player ship
	g.player = CreatePlayer(20)
	g.player.SetPosition(g.screenWidth/2, g.screenHeight/2) // Center of screen
	blue := color.RGBA{0, 0, 255, 255}                      // Blue color
	g.player.SetColor(blue)

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
			50+rand.Float64()*(g.screenWidth-100),  // X between 50 and 750
			50+rand.Float64()*(g.screenHeight-100), // Y between 50 and 550
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

		g.asteroids = append(g.asteroids, asteroid)
	}
}

// drawGameOverScreen draws the game over screen with score and restart instruction
func (g *Game) drawGameOverScreen(screen *ebiten.Image) {
	centerX := float32(g.screenWidth / 2)
	centerY := float32(g.screenHeight / 2)

	// Draw game over reason (GAME OVER or YOU WIN!)
	reasonWidth := g.vectorFont.GetWidth(g.gameOverReason)
	reasonX := centerX - (reasonWidth / 2)
	reasonY := centerY - 60
	g.vectorFont.DrawString(screen, g.gameOverReason, reasonX, reasonY)

	// Draw final score
	scoreText := fmt.Sprintf("SCORE: %d", g.score)
	scoreWidth := g.vectorFont.GetWidth(scoreText)
	scoreX := centerX - (scoreWidth / 2)
	scoreY := centerY - 20
	g.vectorFont.DrawString(screen, scoreText, scoreX, scoreY)

	// Draw restart instruction
	restartText := "PRESS ENTER TO RESTART"
	restartWidth := g.vectorFont.GetWidth(restartText)
	restartX := centerX - (restartWidth / 2)
	restartY := centerY + 40
	g.vectorFont.DrawString(screen, restartText, restartX, restartY)
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Asteroids Game")

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
