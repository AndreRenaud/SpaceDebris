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

// CreateTriangle creates a triangle polygon pointing upward
func CreateTriangle(size float64) *PolygonObject {
	vertices := []Vector2{
		{X: 0, Y: -size},          // Top vertex (pointing up)
		{X: -size * 0.5, Y: size}, // Bottom left
		{X: size * 0.5, Y: size},  // Bottom right
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

// GetBoundingBox returns the bounding box of the transformed polygon
func (p *PolygonObject) GetBoundingBox() (minX, minY, maxX, maxY float64) {
	transformedVertices := p.getTransformedVertices()
	if len(transformedVertices) == 0 {
		return 0, 0, 0, 0
	}

	minX, minY = transformedVertices[0].X, transformedVertices[0].Y
	maxX, maxY = minX, minY

	for _, vertex := range transformedVertices[1:] {
		if vertex.X < minX {
			minX = vertex.X
		}
		if vertex.X > maxX {
			maxX = vertex.X
		}
		if vertex.Y < minY {
			minY = vertex.Y
		}
		if vertex.Y > maxY {
			maxY = vertex.Y
		}
	}

	return minX, minY, maxX, maxY
}

// DrawWithWrapping draws the polygon with screen wrapping
func (p *PolygonObject) DrawWithWrapping(screen *ebiten.Image, screenWidth, screenHeight float64) {
	if len(p.Vertices) < 3 {
		return
	}

	// Get bounding box to check if we need wrapping
	minX, minY, maxX, maxY := p.GetBoundingBox()

	// Store original position
	originalPos := p.Position

	// Determine which edges the polygon is crossing
	drawOffsets := []Vector2{{0, 0}} // Always draw at original position

	// Check horizontal wrapping
	if minX < 0 && maxX > 0 {
		// Crossing left edge - also draw on right side
		drawOffsets = append(drawOffsets, Vector2{screenWidth, 0})
	} else if maxX > screenWidth && minX < screenWidth {
		// Crossing right edge - also draw on left side
		drawOffsets = append(drawOffsets, Vector2{-screenWidth, 0})
	}

	// Check vertical wrapping
	if minY < 0 && maxY > 0 {
		// Crossing top edge - also draw on bottom
		for i := len(drawOffsets) - 1; i >= 0; i-- {
			offset := drawOffsets[i]
			drawOffsets = append(drawOffsets, Vector2{offset.X, screenHeight})
		}
	} else if maxY > screenHeight && minY < screenHeight {
		// Crossing bottom edge - also draw on top
		for i := len(drawOffsets) - 1; i >= 0; i-- {
			offset := drawOffsets[i]
			drawOffsets = append(drawOffsets, Vector2{offset.X, -screenHeight})
		}
	}

	// Draw the polygon at each required position
	for _, offset := range drawOffsets {
		p.Position.X = originalPos.X + offset.X
		p.Position.Y = originalPos.Y + offset.Y
		p.Draw(screen)
	}

	// Restore original position
	p.Position = originalPos
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

// SetVelocity sets the velocity of the polygon
func (p *PolygonObject) SetVelocity(vx, vy float64) {
	p.Velocity.X = vx
	p.Velocity.Y = vy
}

// SetRotationSpeed sets the rotation speed in radians per frame
func (p *PolygonObject) SetRotationSpeed(speed float64) {
	p.RotationSpeed = speed
}

// UpdateWithWrapping updates the polygon and wraps position around screen edges
func (p *PolygonObject) Update(screenWidth, screenHeight float64) {
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

	// Wrap position around screen edges
	if p.Position.X < 0 {
		p.Position.X += screenWidth
	} else if p.Position.X > screenWidth {
		p.Position.X -= screenWidth
	}

	if p.Position.Y < 0 {
		p.Position.Y += screenHeight
	} else if p.Position.Y > screenHeight {
		p.Position.Y -= screenHeight
	}
}
