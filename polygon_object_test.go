package main

import (
	"image/color"
	"math"
	"testing"
)

func TestVector2(t *testing.T) {
	v := Vector2{X: 10, Y: 20}
	if v.X != 10 || v.Y != 20 {
		t.Errorf("Expected {10, 20}, got {%v, %v}", v.X, v.Y)
	}
}

func TestPolygonObject_SetPosition(t *testing.T) {
	p := &PolygonObject{}
	p.SetPosition(100, 200)
	if p.Position.X != 100 || p.Position.Y != 200 {
		t.Errorf("Expected position {100, 200}, got {%v, %v}", p.Position.X, p.Position.Y)
	}
}

func TestPolygonObject_Rotate(t *testing.T) {
	p := &PolygonObject{Rotation: 0}
	p.Rotate(math.Pi / 2)
	if math.Abs(p.Rotation-math.Pi/2) > 1e-9 {
		t.Errorf("Expected rotation %v, got %v", math.Pi/2, p.Rotation)
	}
}

func TestBoundingBoxOverlap(t *testing.T) {
	box1 := BoundingBox{MinX: 0, MinY: 0, MaxX: 10, MaxY: 10}
	box2 := BoundingBox{MinX: 5, MinY: 5, MaxX: 15, MaxY: 15}
	box3 := BoundingBox{MinX: 20, MinY: 20, MaxX: 30, MaxY: 30}

	if !box1.Overlaps(box2) {
		t.Errorf("Expected box1 and box2 to overlap")
	}
	if box1.Overlaps(box3) {
		t.Errorf("Expected box1 and box3 to not overlap")
	}
}

func TestPointInPolygon(t *testing.T) {
	// Square polygon
	vertices := []Vector2{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	insidePoint := Vector2{X: 5, Y: 5}
	outsidePoint := Vector2{X: 15, Y: 5}

	if !PointInPolygon(insidePoint, vertices) {
		t.Errorf("Expected point {%v, %v} to be inside polygon", insidePoint.X, insidePoint.Y)
	}
	if PointInPolygon(outsidePoint, vertices) {
		t.Errorf("Expected point {%v, %v} to be outside polygon", outsidePoint.X, outsidePoint.Y)
	}
}

func TestInterpolateColor(t *testing.T) {
	c1 := color.RGBA{0, 0, 0, 255}
	c2 := color.RGBA{255, 255, 255, 255}

	mid := interpolateColor(c1, c2, 0.5)
	r, _, _, a := mid.RGBA()

	// RGBA() returns 0-65535
	if uint8(r>>8) < 120 || uint8(r>>8) > 135 {
		t.Errorf("Expected red around 127, got %v", uint8(r>>8))
	}
	if uint8(a>>8) != 255 {
		t.Errorf("Expected alpha 255, got %v", uint8(a>>8))
	}
}

func TestCollisionDetection(t *testing.T) {
	polygon1 := &PolygonObject{
		Vertices: []Vector2{
			{X: 0, Y: 0},
			{X: 10, Y: 0},
			{X: 10, Y: 10},
			{X: 0, Y: 10},
		},
		Position: Vector2{X: 0, Y: 0},
		Scale:    1.0,
	}

	polygon2 := &PolygonObject{
		Vertices: []Vector2{
			{X: 5, Y: 5},
			{X: 15, Y: 5},
			{X: 15, Y: 15},
			{X: 5, Y: 15},
		},
		Position: Vector2{X: 0, Y: 0},
		Scale:    1.0,
	}

	polygon3 := &PolygonObject{
		Vertices: []Vector2{
			{X: 20, Y: 20},
			{X: 30, Y: 20},
			{X: 30, Y: 30},
			{X: 20, Y: 30},
		},
		Position: Vector2{X: 0, Y: 0},
		Scale:    1.0,
	}

	if !PolygonsCollide(polygon1, polygon2) {
		t.Errorf("Expected polygon1 and polygon2 to collide")
	}
	if PolygonsCollide(polygon1, polygon3) {
		t.Errorf("Expected polygon1 and polygon3 to not collide")
	}
}
