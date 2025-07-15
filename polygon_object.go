package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	// Rotation angle in radians
	Rotation float64
	// Scale factor
	Scale float64
	// Color for drawing
	Color color.Color
	// Line width for drawing
	LineWidth float32
}

// NewPolygonObject creates a new polygon object
func NewPolygonObject(vertices []Vector2) *PolygonObject {
	return &PolygonObject{
		Vertices:  vertices,
		Position:  Vector2{X: 0, Y: 0},
		Rotation:  0,
		Scale:     1.0,
		Color:     color.White,
		LineWidth: 1.0,
	}
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
	return NewPolygonObject(vertices)
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

// Draw renders the polygon to the screen
func (p *PolygonObject) Draw(screen *ebiten.Image) {
	if len(p.Vertices) < 3 {
		return // Can't draw a polygon with less than 3 vertices
	}

	transformedVertices := p.getTransformedVertices()

	// Draw lines between consecutive vertices
	for i := 0; i < len(transformedVertices); i++ {
		start := transformedVertices[i]
		end := transformedVertices[(i+1)%len(transformedVertices)] // Wrap to first vertex for last line

		ebitenutil.DrawLine(
			screen,
			start.X, start.Y,
			end.X, end.Y,
			p.Color,
		)
	}
}

// DrawFilled renders the polygon as a filled shape (simplified implementation)
func (p *PolygonObject) DrawFilled(screen *ebiten.Image) {
	// For now, just draw the outline with thicker lines
	// In a more advanced implementation, you could use triangulation
	// or the vector package's proper filling methods
	transformedVertices := p.GetTransformedVertices()

	if len(transformedVertices) < 3 {
		return
	}

	// Draw multiple overlapping lines to simulate a filled effect
	for thickness := 0; thickness < 3; thickness++ {
		for i := 0; i < len(transformedVertices); i++ {
			start := transformedVertices[i]
			end := transformedVertices[(i+1)%len(transformedVertices)]

			// Draw slightly offset lines to create thickness
			offset := float64(thickness) * 0.5
			ebitenutil.DrawLine(
				screen,
				start.X+offset, start.Y+offset,
				end.X+offset, end.Y+offset,
				p.Color,
			)
		}
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
