package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// VectorFont represents a font made of vector lines for drawing digits
type VectorFont struct {
	digitWidth  float32
	digitHeight float32
	lineWidth   float32
	color       color.Color
}

// NewVectorFont creates a new vector font with specified dimensions
func NewVectorFont(width, height, lineWidth float32, c color.Color) *VectorFont {
	return &VectorFont{
		digitWidth:  width,
		digitHeight: height,
		lineWidth:   lineWidth,
		color:       c,
	}
}

// SetColor changes the color of the font
func (vf *VectorFont) SetColor(c color.Color) {
	vf.color = c
}

// drawLine draws a line from (x1, y1) to (x2, y2) on the screen
func (vf *VectorFont) drawLine(screen *ebiten.Image, x1, y1, x2, y2, offsetX, offsetY float32) {
	vector.StrokeLine(screen,
		offsetX+x1, offsetY+y1,
		offsetX+x2, offsetY+y2,
		vf.lineWidth, vf.color, false)
}

// DrawDigit draws a single digit at the specified position
func (vf *VectorFont) DrawDigit(screen *ebiten.Image, digit int, x, y float32) {
	if digit < 0 || digit > 9 {
		return // Invalid digit
	}

	w := vf.digitWidth
	h := vf.digitHeight
	hh := h / 2 // half height

	// Define the seven-segment display positions
	// Segments are numbered as follows:
	//  AAA
	// F   B
	// F   B
	//  GGG
	// E   C
	// E   C
	//  DDD

	switch digit {
	case 0:
		// Segments: A, B, C, D, E, F (all except G)
		vf.drawLine(screen, 0, 0, w, 0, x, y)  // A (top)
		vf.drawLine(screen, w, 0, w, hh, x, y) // B (top right)
		vf.drawLine(screen, w, hh, w, h, x, y) // C (bottom right)
		vf.drawLine(screen, w, h, 0, h, x, y)  // D (bottom)
		vf.drawLine(screen, 0, h, 0, hh, x, y) // E (bottom left)
		vf.drawLine(screen, 0, hh, 0, 0, x, y) // F (top left)

	case 1:
		// Segments: B, C (right side)
		vf.drawLine(screen, w, 0, w, hh, x, y) // B (top right)
		vf.drawLine(screen, w, hh, w, h, x, y) // C (bottom right)

	case 2:
		// Segments: A, B, G, E, D
		vf.drawLine(screen, 0, 0, w, 0, x, y)   // A (top)
		vf.drawLine(screen, w, 0, w, hh, x, y)  // B (top right)
		vf.drawLine(screen, w, hh, 0, hh, x, y) // G (middle)
		vf.drawLine(screen, 0, hh, 0, h, x, y)  // E (bottom left)
		vf.drawLine(screen, 0, h, w, h, x, y)   // D (bottom)

	case 3:
		// Segments: A, B, G, C, D
		vf.drawLine(screen, 0, 0, w, 0, x, y)   // A (top)
		vf.drawLine(screen, w, 0, w, hh, x, y)  // B (top right)
		vf.drawLine(screen, w, hh, 0, hh, x, y) // G (middle)
		vf.drawLine(screen, w, hh, w, h, x, y)  // C (bottom right)
		vf.drawLine(screen, w, h, 0, h, x, y)   // D (bottom)

	case 4:
		// Segments: F, G, B, C
		vf.drawLine(screen, 0, 0, 0, hh, x, y)  // F (top left)
		vf.drawLine(screen, 0, hh, w, hh, x, y) // G (middle)
		vf.drawLine(screen, w, 0, w, hh, x, y)  // B (top right)
		vf.drawLine(screen, w, hh, w, h, x, y)  // C (bottom right)

	case 5:
		// Segments: A, F, G, C, D
		vf.drawLine(screen, 0, 0, w, 0, x, y)   // A (top)
		vf.drawLine(screen, 0, 0, 0, hh, x, y)  // F (top left)
		vf.drawLine(screen, 0, hh, w, hh, x, y) // G (middle)
		vf.drawLine(screen, w, hh, w, h, x, y)  // C (bottom right)
		vf.drawLine(screen, w, h, 0, h, x, y)   // D (bottom)

	case 6:
		// Segments: A, F, G, E, D, C
		vf.drawLine(screen, 0, 0, w, 0, x, y)   // A (top)
		vf.drawLine(screen, 0, 0, 0, hh, x, y)  // F (top left)
		vf.drawLine(screen, 0, hh, w, hh, x, y) // G (middle)
		vf.drawLine(screen, 0, hh, 0, h, x, y)  // E (bottom left)
		vf.drawLine(screen, 0, h, w, h, x, y)   // D (bottom)
		vf.drawLine(screen, w, hh, w, h, x, y)  // C (bottom right)

	case 7:
		// Segments: A, B, C
		vf.drawLine(screen, 0, 0, w, 0, x, y)  // A (top)
		vf.drawLine(screen, w, 0, w, hh, x, y) // B (top right)
		vf.drawLine(screen, w, hh, w, h, x, y) // C (bottom right)

	case 8:
		// Segments: A, B, C, D, E, F, G (all segments)
		vf.drawLine(screen, 0, 0, w, 0, x, y)   // A (top)
		vf.drawLine(screen, w, 0, w, hh, x, y)  // B (top right)
		vf.drawLine(screen, w, hh, w, h, x, y)  // C (bottom right)
		vf.drawLine(screen, w, h, 0, h, x, y)   // D (bottom)
		vf.drawLine(screen, 0, h, 0, hh, x, y)  // E (bottom left)
		vf.drawLine(screen, 0, hh, 0, 0, x, y)  // F (top left)
		vf.drawLine(screen, 0, hh, w, hh, x, y) // G (middle)

	case 9:
		// Segments: A, B, C, D, F, G
		vf.drawLine(screen, 0, 0, w, 0, x, y)   // A (top)
		vf.drawLine(screen, w, 0, w, hh, x, y)  // B (top right)
		vf.drawLine(screen, w, hh, w, h, x, y)  // C (bottom right)
		vf.drawLine(screen, w, h, 0, h, x, y)   // D (bottom)
		vf.drawLine(screen, 0, 0, 0, hh, x, y)  // F (top left)
		vf.drawLine(screen, 0, hh, w, hh, x, y) // G (middle)
	}
}

// DrawNumber draws a multi-digit number at the specified position
func (vf *VectorFont) DrawNumber(screen *ebiten.Image, number int, x, y float32) {
	if number == 0 {
		vf.DrawDigit(screen, 0, x, y)
		return
	}

	// Convert number to string to get individual digits
	digits := []int{}
	temp := number
	if temp < 0 {
		temp = -temp // Handle negative numbers by making them positive
	}

	for temp > 0 {
		digits = append([]int{temp % 10}, digits...) // Prepend to reverse order
		temp /= 10
	}

	// Draw each digit with spacing
	spacing := vf.digitWidth + 4 // Small gap between digits
	currentX := x

	for _, digit := range digits {
		vf.DrawDigit(screen, digit, currentX, y)
		currentX += spacing
	}
}

// GetTextWidth calculates the width of a number when drawn
func (vf *VectorFont) GetTextWidth(number int) float32 {
	if number == 0 {
		return vf.digitWidth
	}

	digitCount := 0
	temp := number
	if temp < 0 {
		temp = -temp
	}

	for temp > 0 {
		digitCount++
		temp /= 10
	}

	spacing := vf.digitWidth + 4
	return float32(digitCount-1)*spacing + vf.digitWidth
}
