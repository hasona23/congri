package utils

import "math"

type Vec2 struct {
	X, Y float32
}

func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(((v.X * v.X) + (v.Y * v.Y)))))
}

func (v *Vec2) NormalizeDir() {
	if v.X != 0 && v.Y != 0 {
		factor := 1 / v.Length()
		v.X *= factor
		v.Y *= factor
	}
}
func V2(x, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}
