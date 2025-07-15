package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ghostTrailLength = 5 // Number of frames to keep in the trail
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
	// Color fading properties
	FadeStartColor color.Color
	FadeEndColor   color.Color
	FadeProgress   float64 // 0.0 to 1.0, where 0 is start color and 1 is end color
	FadeSpeed      float64 // How fast to fade (increment per frame)
	IsFading       bool    // Whether the object is currently fading
	drawCount      int

	Trail []drawablePolygon
}

type drawablePolygon []Vector2

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
		Vertices:       vertices,
		Position:       Vector2{X: 0, Y: 0},
		Velocity:       Vector2{X: 0, Y: 0},
		Rotation:       0,
		RotationSpeed:  0,
		Scale:          1.0,
		Color:          color.White,
		LineWidth:      1.0,
		FadeStartColor: color.White,
		FadeEndColor:   color.White,
		FadeProgress:   0.0,
		FadeSpeed:      0.0,
		IsFading:       false,
		Trail:          make([]drawablePolygon, 0, ghostTrailLength),
	}
}

// CreatePlayer creates a spaceship polygon with wings and a divet at the back
func CreatePlayer(size float64) *PolygonObject {
	vertices := []Vector2{
		{X: 0, Y: -size},                 // Nose (top vertex, pointing up)
		{X: -size * 0.3, Y: -size * 0.2}, // Left side of nose
		{X: -size * 0.7, Y: size * 0.3},  // Left wing tip
		{X: -size * 0.4, Y: size * 0.6},  // Left wing inner
		{X: -size * 0.2, Y: size * 0.8},  // Left back corner
		{X: 0, Y: size * 0.6},            // Center back (creates divet)
		{X: size * 0.2, Y: size * 0.8},   // Right back corner
		{X: size * 0.4, Y: size * 0.6},   // Right wing inner
		{X: size * 0.7, Y: size * 0.3},   // Right wing tip
		{X: size * 0.3, Y: -size * 0.2},  // Right side of nose
	}
	return &PolygonObject{
		Vertices:       vertices,
		Position:       Vector2{X: 0, Y: 0},
		Velocity:       Vector2{X: 0, Y: 0},
		Rotation:       0,
		RotationSpeed:  0,
		Scale:          1.0,
		Color:          color.White,
		LineWidth:      1.0,
		FadeStartColor: color.White,
		FadeEndColor:   color.White,
		FadeProgress:   0.0,
		FadeSpeed:      0.0,
		IsFading:       false,
		Trail:          make([]drawablePolygon, 0, ghostTrailLength),
	}
}

// GetTransformedVertices returns the vertices transformed by position, rotation, and scale
func (p *PolygonObject) getTransformedVertices() drawablePolygon {
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

func (d drawablePolygon) Draw(screen *ebiten.Image, lineWidth float32, color color.Color) {
	if len(d) < 3 {
		return // Can't draw a polygon with less than 3 vertices
	}

	// Draw the polygon outline using vector.StrokeLine for each edge
	for i := 0; i < len(d); i++ {
		start := d[i]
		end := d[(i+1)%len(d)] // Wrap to first vertex for last line

		vector.StrokeLine(
			screen,
			float32(start.X), float32(start.Y),
			float32(end.X), float32(end.Y),
			lineWidth,
			color,
			true, // antialiasing
		)
	}
}

// Draw renders the polygon to the screen with antialiased lines
func (p *PolygonObject) Draw(screen *ebiten.Image) {
	if len(p.Vertices) < 3 {
		return // Can't draw a polygon with less than 3 vertices
	}
	p.drawCount++

	max := len(p.Trail)
	for i, trail := range p.Trail {
		r, g, b, _ := p.Color.RGBA()
		ratio := 1 - float64(i)/float64(max)
		newCol := color.RGBA{
			R: uint8(int(float64(r)*ratio) >> 8),
			G: uint8(int(float64(g)*ratio) >> 8),
			B: uint8(int(float64(b)*ratio) >> 8),
			A: 0xff,
		}
		trail.Draw(screen, p.LineWidth, newCol)
	}
	transformedVertices := p.getTransformedVertices()
	transformedVertices.Draw(screen, p.LineWidth, p.Color)

	// Don't add everything to the trail
	if p.drawCount%4 == 0 {
		p.Trail = append([]drawablePolygon{transformedVertices}, p.Trail...)
		if len(p.Trail) > ghostTrailLength {
			p.Trail = p.Trail[:ghostTrailLength] // Limit trail length
		}
	}

}

// BoundingBox represents a rectangular bounding box
type BoundingBox struct {
	MinX, MinY, MaxX, MaxY float64
}

// GetBoundingBox returns the axis-aligned bounding box of the polygon
func (p *PolygonObject) GetBoundingBox() BoundingBox {
	transformedVertices := p.getTransformedVertices()
	if len(transformedVertices) == 0 {
		return BoundingBox{0, 0, 0, 0}
	}

	minX, minY := transformedVertices[0].X, transformedVertices[0].Y
	maxX, maxY := minX, minY

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

	return BoundingBox{minX, minY, maxX, maxY}
}

// BoundingBoxesOverlap checks if two bounding boxes overlap (fast check)
func BoundingBoxesOverlap(box1, box2 BoundingBox) bool {
	return box1.MinX <= box2.MaxX && box1.MaxX >= box2.MinX &&
		box1.MinY <= box2.MaxY && box1.MaxY >= box2.MinY
}

// PointInPolygon checks if a point is inside a polygon using ray casting algorithm
func PointInPolygon(point Vector2, vertices []Vector2) bool {
	if len(vertices) < 3 {
		return false
	}

	inside := false
	j := len(vertices) - 1

	for i := 0; i < len(vertices); i++ {
		xi, yi := vertices[i].X, vertices[i].Y
		xj, yj := vertices[j].X, vertices[j].Y

		if ((yi > point.Y) != (yj > point.Y)) &&
			(point.X < (xj-xi)*(point.Y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// LineSegmentsIntersect checks if two line segments intersect
func LineSegmentsIntersect(p1, p2, p3, p4 Vector2) bool {
	// Calculate the direction vectors
	d1 := Vector2{p2.X - p1.X, p2.Y - p1.Y}
	d2 := Vector2{p4.X - p3.X, p4.Y - p3.Y}
	d3 := Vector2{p1.X - p3.X, p1.Y - p3.Y}

	// Calculate cross products
	cross1 := d1.X*d2.Y - d1.Y*d2.X
	cross2 := d3.X*d2.Y - d3.Y*d2.X
	cross3 := d3.X*d1.Y - d3.Y*d1.X

	// Check if lines are parallel
	if math.Abs(cross1) < 1e-10 {
		return false // Parallel lines
	}

	// Calculate intersection parameters
	t1 := cross2 / cross1
	t2 := cross3 / cross1

	// Check if intersection point lies within both line segments
	return t1 >= 0 && t1 <= 1 && t2 >= 0 && t2 <= 1
}

// PolygonsCollide checks if two polygons collide
func PolygonsCollide(poly1, poly2 *PolygonObject) bool {
	// Fast bounding box check first
	box1 := poly1.GetBoundingBox()
	box2 := poly2.GetBoundingBox()

	if !BoundingBoxesOverlap(box1, box2) {
		return false // No collision possible if bounding boxes don't overlap
	}

	// Get transformed vertices for both polygons
	vertices1 := poly1.getTransformedVertices()
	vertices2 := poly2.getTransformedVertices()

	if len(vertices1) < 3 || len(vertices2) < 3 {
		return false
	}

	// Check if any vertex of polygon1 is inside polygon2
	for _, vertex := range vertices1 {
		if PointInPolygon(vertex, vertices2) {
			return true
		}
	}

	// Check if any vertex of polygon2 is inside polygon1
	for _, vertex := range vertices2 {
		if PointInPolygon(vertex, vertices1) {
			return true
		}
	}

	// Check if any edge of polygon1 intersects any edge of polygon2
	for i := 0; i < len(vertices1); i++ {
		edge1Start := vertices1[i]
		edge1End := vertices1[(i+1)%len(vertices1)]

		for j := 0; j < len(vertices2); j++ {
			edge2Start := vertices2[j]
			edge2End := vertices2[(j+1)%len(vertices2)]

			if LineSegmentsIntersect(edge1Start, edge1End, edge2Start, edge2End) {
				return true
			}
		}
	}

	return false
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
	p.IsFading = false // Stop fading when color is set directly
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
func (p *PolygonObject) Update(screenWidth, screenHeight float64, withWrapping bool) {
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

	// Update color fading
	p.updateFade()

	if withWrapping {
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
}

// interpolateColor interpolates between two colors based on progress (0.0 to 1.0)
func interpolateColor(startColor, endColor color.Color, progress float64) color.Color {
	// Clamp progress to [0, 1]
	if progress < 0 {
		progress = 0
	} else if progress > 1 {
		progress = 1
	}

	// Convert colors to RGBA
	sr, sg, sb, sa := startColor.RGBA()
	er, eg, eb, ea := endColor.RGBA()

	// Interpolate each component
	r := uint8((float64(sr>>8)*(1-progress) + float64(er>>8)*progress))
	g := uint8((float64(sg>>8)*(1-progress) + float64(eg>>8)*progress))
	b := uint8((float64(sb>>8)*(1-progress) + float64(eb>>8)*progress))
	a := uint8((float64(sa>>8)*(1-progress) + float64(ea>>8)*progress))

	return color.RGBA{r, g, b, a}
}

// StartFade begins a color fade from current color to target color
func (p *PolygonObject) StartFade(targetColor color.Color, duration float64) {
	p.FadeStartColor = p.Color
	p.FadeEndColor = targetColor
	p.FadeProgress = 0.0
	p.FadeSpeed = 1.0 / duration // duration in frames (60 FPS)
	p.IsFading = true
}

// updateFade updates the fade progress and color (called internally by Update)
func (p *PolygonObject) updateFade() {
	if !p.IsFading {
		return
	}

	p.FadeProgress += p.FadeSpeed

	if p.FadeProgress >= 1.0 {
		// Fade complete
		p.FadeProgress = 1.0
		p.Color = p.FadeEndColor
		p.IsFading = false
	} else {
		// Update color based on current progress
		p.Color = interpolateColor(p.FadeStartColor, p.FadeEndColor, p.FadeProgress)
	}
}
