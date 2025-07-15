package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Vector2 represents a 2D point or vector
type Vector2 struct {
	X, Y float64
}

// PolygonObject represents a closed polygon that can be drawn
type PolygonObject struct {
	// Vertices relative to the object's origin (0,0)
	Vertices []Vector2
	// Position of the object's origin in world space
	Position Vector2
	// Velocity in pixels per frame
	Velocity Vector2
	// Rotation angle in radians
	Rotation float64
	// Rotation speed in radians per frame
	RotationSpeed float64
	// Scale factor
	Scale float64
	// Color for drawing
	Color color.Color
	// Line width for drawing
	LineWidth float32
}

// CreateAsteroid creates an irregular asteroid-like polygon
func CreateAsteroid(baseRadius float64, irregularity float64, numVertices int) *PolygonObject {
	vertices := make([]Vector2, numVertices)
	angleStep := 2 * math.Pi / float64(numVertices)

	for i := 0; i < numVertices; i++ {
		angle := float64(i) * angleStep
		// Add some irregularity to the radius
		radius := baseRadius + (math.Sin(angle*3)+math.Cos(angle*5))*irregularity
		vertices[i] = Vector2{
			X: math.Cos(angle) * radius,
			Y: math.Sin(angle) * radius,
		}
	}
	return &PolygonObject{
		Vertices:      vertices,
		Position:      Vector2{X: 0, Y: 0},
		Velocity:      Vector2{X: 0, Y: 0},
		Rotation:      0,
		RotationSpeed: 0,
		Scale:         1.0,
		Color:         color.White,
		LineWidth:     1.0,
	}
}

// GetTransformedVertices returns the vertices transformed by position, rotation, and scale
func (p *PolygonObject) getTransformedVertices() []Vector2 {
	transformed := make([]Vector2, len(p.Vertices))
	cos := math.Cos(p.Rotation)
	sin := math.Sin(p.Rotation)

	for i, vertex := range p.Vertices {
		// Scale
		scaledX := vertex.X * p.Scale
		scaledY := vertex.Y * p.Scale

		// Rotate
		rotatedX := scaledX*cos - scaledY*sin
		rotatedY := scaledX*sin + scaledY*cos

		// Translate
		transformed[i] = Vector2{
			X: rotatedX + p.Position.X,
			Y: rotatedY + p.Position.Y,
		}
	}

	return transformed
}

// whiteImage is a 1x1 white image used for drawing colored shapes
var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

// Draw renders the polygon to the screen with antialiased lines
func (p *PolygonObject) Draw(screen *ebiten.Image) {
	if len(p.Vertices) < 3 {
		return // Can't draw a polygon with less than 3 vertices
	}

	transformedVertices := p.getTransformedVertices()

	// Draw the polygon outline using vector.StrokeLine for each edge
	for i := 0; i < len(transformedVertices); i++ {
		start := transformedVertices[i]
		end := transformedVertices[(i+1)%len(transformedVertices)] // Wrap to first vertex for last line

		vector.StrokeLine(
			screen,
			float32(start.X), float32(start.Y),
			float32(end.X), float32(end.Y),
			p.LineWidth,
			p.Color,
			true, // antialiasing
		)
	}
}

// SetPosition sets the world position of the polygon
func (p *PolygonObject) SetPosition(x, y float64) {
	p.Position.X = x
	p.Position.Y = y
}

// SetRotation sets the rotation angle in radians
func (p *PolygonObject) SetRotation(angle float64) {
	p.Rotation = angle
}

// SetScale sets the scale factor
func (p *PolygonObject) SetScale(scale float64) {
	p.Scale = scale
}

// SetColor sets the drawing color
func (p *PolygonObject) SetColor(c color.Color) {
	p.Color = c
}

// Rotate rotates the polygon by the given angle in radians
func (p *PolygonObject) Rotate(angle float64) {
	p.Rotation += angle
}

// Move translates the polygon by the given offset
func (p *PolygonObject) Move(dx, dy float64) {
	p.Position.X += dx
	p.Position.Y += dy
}

// Update updates the polygon's position and rotation based on velocity and rotation speed
func (p *PolygonObject) Update() {
	// Update position based on velocity
	p.Position.X += p.Velocity.X
	p.Position.Y += p.Velocity.Y

	// Update rotation based on rotation speed
	p.Rotation += p.RotationSpeed

	// Keep rotation in the range [0, 2π] for cleaner values
	if p.Rotation > 2*math.Pi {
		p.Rotation -= 2 * math.Pi
	} else if p.Rotation < 0 {
		p.Rotation += 2 * math.Pi
	}
}

// SetVelocity sets the velocity of the polygon
func (p *PolygonObject) SetVelocity(vx, vy float64) {
	p.Velocity.X = vx
	p.Velocity.Y = vy
}

// SetRotationSpeed sets the rotation speed in radians per frame
func (p *PolygonObject) SetRotationSpeed(speed float64) {
	p.RotationSpeed = speed
}
