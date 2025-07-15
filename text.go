package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// VectorFont represents a font made of vector lines for drawing digits
type VectorFont struct {
	runeWidth  float32
	runeHeight float32
	lineWidth  float32
	color      color.Color
}

type LineSegment struct {
	X1, Y1, X2, Y2 float32 // Start and end points of the line segment
}

// NewVectorFont creates a new vector font with specified dimensions
func NewVectorFont(width, height, lineWidth float32, c color.Color) *VectorFont {
	return &VectorFont{
		runeWidth:  width,
		runeHeight: height,
		lineWidth:  lineWidth,
		color:      c,
	}
}

// SetColor changes the color of the font
func (vf *VectorFont) SetColor(c color.Color) {
	vf.color = c
}

// drawLine draws a line from (x1, y1) to (x2, y2) on the screen
func (vf *VectorFont) drawLine(screen *ebiten.Image, x1, y1, x2, y2 float32) {
	vector.StrokeLine(screen,
		x1, y1,
		x2, y2,
		vf.lineWidth, vf.color, true)
}

var charMaps = map[rune][]LineSegment{
	// Define the seven-segment display positions
	// Segments are numbered as follows:
	//  AAA
	// F   B
	// F   B
	//  GGG
	// E   C
	// E   C
	//  DDD

	'0': {
		{0, 0, 1, 0},   // A (top)
		{1, 0, 1, 0.5}, // B (top right)
		{1, 0.5, 1, 1}, // C (bottom right)
		{1, 1, 0, 1},   // D (bottom)
		{0, 1, 0, 0.5}, // E (bottom left)
		{0, 0.5, 0, 0}, // F (top left)
	},
	'1': {
		{1, 0, 1, 0.5}, // B (top right)
		{1, 0.5, 1, 1}, // C (bottom right)
	},
	'2': {
		{0, 0, 1, 0},     // A (top)
		{1, 0, 1, 0.5},   // B (top right)
		{1, 0.5, 0, 0.5}, // G (middle)
		{0, 0.5, 0, 1},   // E (bottom left)
		{0, 1, 1, 1},     // D (bottom)
	},
	'3': {
		{0, 0, 1, 0},     // A (top)
		{1, 0, 1, 0.5},   // B (top right)
		{1, 0.5, 0, 0.5}, // G (middle)
		{1, 0.5, 1, 1},   // C (bottom right)
		{1, 1, 0, 1},     // D (bottom)
	},
	'4': {
		{0, 0, 0, 0.5},   // F (top left)
		{0, 0.5, 1, 0.5}, // G (middle)
		{1, 0, 1, 0.5},   // B (top right)
		{1, 0.5, 1, 1},   // C (bottom right)
	},
	'5': {
		{0, 0, 1, 0},     // A (top)
		{0, 0, 0, 0.5},   // F (top left)
		{0, 0.5, 1, 0.5}, // G (middle)
		{1, 0.5, 1, 1},   // C (bottom right)
		{1, 1, 0, 1},     // D (bottom)
	},
	'6': {
		{0, 0, 1, 0},     // A (top)
		{0, 0, 0, 0.5},   // F (top left)
		{0, 0.5, 1, 0.5}, // G (middle)
		{0, 0.5, 0, 1},   // E (bottom left)
		{0, 1, 1, 1},     // D (bottom)
		{1, 0.5, 1, 1},   // C (bottom right)
	},
	'7': {
		{0, 0, 1, 0},   // A (top)
		{1, 0, 1, 0.5}, // B (top right)
		{1, 0.5, 1, 1}, // C (bottom right)
	},
	'8': {
		{0, 0, 1, 0},     // A (top)
		{1, 0, 1, 0.5},   // B (top right)
		{1, 0.5, 1, 1},   // C (bottom right)
		{1, 1, 0, 1},     // D (bottom)
		{0, 1, 0, 0.5},   // E (bottom left)
		{0, 0.5, 0, 0},   // F (top left)
		{0, 0.5, 1, 0.5}, // G (middle)
	},
	'9': {
		{0, 0, 1, 0},     // A (top)
		{1, 0, 1, 0.5},   // B (top right)
		{1, 0.5, 1, 1},   // C (bottom right)
		{1, 1, 0, 1},     // D (bottom)
		{0, 0, 0, 0.5},   // F (top left)
		{0, 0.5, 1, 0.5}, // G (middle)
	},

	// Alphabet characters
	'A': {
		{0, 1, 0.5, 0},
		{1, 1, 0.5, 0},
		{0.2, 0.6, 0.8, 0.6},
	},
	'C': {
		{0, 0, 1, 0},   // A (top)
		{0, 0, 0, 0.5}, // F (top left)
		{0, 0.5, 0, 1}, // E (bottom left)
		{0, 1, 1, 1},   // D (bottom)
	},
	'E': {
		{0, 0, 1, 0},       // A (top)
		{0, 0, 0, 0.5},     // F (top left)
		{0, 0.5, 0.7, 0.5}, // G (middle, shortened)
		{0, 0.5, 0, 1},     // E (bottom left)
		{0, 1, 1, 1},       // D (bottom)
	},
	'G': {
		{0, 0, 1, 0},       // A (top)
		{0, 0, 0, 0.5},     // F (top left)
		{0, 0.5, 0, 1},     // E (bottom left)
		{0, 1, 1, 1},       // D (bottom)
		{1, 0.5, 1, 1},     // C (bottom right)
		{0.5, 0.5, 1, 0.5}, // G (middle, from center to right)
	},
	'M': {
		{0, 1, 0, 0},     // Left vertical (full height)
		{0, 0, 0.5, 0.5}, // Left diagonal to center
		{0.5, 0.5, 1, 0}, // Right diagonal from center
		{1, 0, 1, 1},     // Right vertical (full height)
	},
	'N': {
		{0, 0, 0, 1}, // Left vertical (full height)
		{0, 0, 1, 0}, // Top horizontal
		{1, 0, 1, 1}, // Right vertical (full height)
	},
	'O': {
		{0, 0, 1, 0},   // A (top)
		{1, 0, 1, 0.5}, // B (top right)
		{1, 0.5, 1, 1}, // C (bottom right)
		{1, 1, 0, 1},   // D (bottom)
		{0, 1, 0, 0.5}, // E (bottom left)
		{0, 0.5, 0, 0}, // F (top left)
	},
	'P': {
		{0, 0, 1, 0},     // A (top)
		{1, 0, 1, 0.5},   // B (top right)
		{1, 0.5, 0, 0.5}, // G (middle)
		{0, 0, 0, 0.5},   // F (top left)
		{0, 0.5, 0, 1},   // E (bottom left)
	},
	'R': {
		{0, 0, 1, 0},     // A (top)
		{1, 0, 1, 0.5},   // B (top right)
		{1, 0.5, 0, 0.5}, // G (middle)
		{0, 0, 0, 0.5},   // F (top left)
		{0, 0.5, 0, 1},   // E (bottom left)
		{0.5, 0.5, 1, 1}, // Diagonal from middle to bottom right
	},
	'S': {
		{0, 0, 1, 0},     // A (top)
		{0, 0, 0, 0.5},   // F (top left)
		{0, 0.5, 1, 0.5}, // G (middle)
		{1, 0.5, 1, 1},   // C (bottom right)
		{1, 1, 0, 1},     // D (bottom)
	},
	'T': {
		{0, 0, 1, 0},     // A (top horizontal)
		{0.5, 0, 0.5, 1}, // Center vertical line
	},
	'U': {
		{0, 0, 0, 0.5}, // F (top left)
		{0, 0.5, 0, 1}, // E (bottom left)
		{0, 1, 1, 1},   // D (bottom)
		{1, 1, 1, 0.5}, // C (bottom right)
		{1, 0.5, 1, 0}, // B (top right)
	},
	'V': {
		{0, 0, 0.5, 1}, // Left diagonal from top-left to bottom-center
		{1, 0, 0.5, 1}, // Right diagonal from top-right to bottom-center
	},
	'W': {
		{0, 0, 0, 1},           // Left vertical (full height)
		{0, 1, 0.33, 0.5},      // Left diagonal to first center point
		{0.33, 0.5, 0.67, 0.5}, // Center horizontal connection
		{0.67, 0.5, 1, 1},      // Right diagonal from second center
		{1, 1, 1, 0},           // Right vertical (full height)
	},
	'Y': {
		{0, 0, 0.5, 0.5},   // Left diagonal from top-left to center
		{1, 0, 0.5, 0.5},   // Right diagonal from top-right to center
		{0.5, 0.5, 0.5, 1}, // Center vertical from middle to bottom
	},
	'I': {
		{0, 0, 1, 0},     // A (top horizontal)
		{0.5, 0, 0.5, 1}, // Center vertical line
		{0, 1, 1, 1},     // D (bottom horizontal)
	},
	'!': {
		{0.5, 0, 0.5, 0.7},   // Vertical line (top part)
		{0.4, 0.8, 0.6, 0.8}, // Dot (top part)
		{0.4, 0.9, 0.6, 0.9}, // Dot (bottom part)
	},
	':': {
		{0.4, 0.3, 0.6, 0.3}, // Top dot (top part)
		{0.4, 0.4, 0.6, 0.4}, // Top dot (bottom part)
		{0.4, 0.6, 0.6, 0.6}, // Bottom dot (top part)
		{0.4, 0.7, 0.6, 0.7}, // Bottom dot (bottom part)
	},
}

// DrawDigit draws a single digit at the specified position
func (vf *VectorFont) DrawRune(screen *ebiten.Image, ch rune, x, y float32) {
	segments, ok := charMaps[ch]
	if !ok {
		return // No segments defined for this digit
	}
	// Draw each segment of the digit
	for _, seg := range segments {
		vf.drawLine(screen, x+seg.X1*vf.runeWidth, y+seg.Y1*vf.runeHeight,
			x+seg.X2*vf.runeWidth, y+seg.Y2*vf.runeHeight)
	}
}

// DrawNumber draws a multi-digit number at the specified position
func (vf *VectorFont) DrawString(screen *ebiten.Image, str string, x, y float32) {
	// Draw each digit with spacing
	spacing := vf.runeWidth + 4 // Small gap between digits
	currentX := x

	for _, ch := range str {
		vf.DrawRune(screen, ch, currentX, y)
		currentX += spacing
	}
}

// GetTextWidth calculates the width of a number when drawn
func (vf *VectorFont) GetWidth(str string) float32 {
	return vf.runeWidth*float32(len(str)) + 4*float32(len(str)-1)
}
