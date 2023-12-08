package utils

import "math/rand"

// 常见颜色
var (
	WHITE = Color{
		R: 1.0,
		G: 1.0,
		B: 1.0,
		A: 1.0,
	}
	BLACK = Color{
		R: 0.0,
		G: 0.0,
		B: 0.0,
		A: 1.0,
	}
	BLUE = Color{
		R: 0.0,
		G: 0.0,
		B: 1.0,
		A: 1.0,
	}
	RED = Color{
		R: 1.0,
		G: 0.0,
		B: 0.0,
		A: 1.0,
	}
	GREEN = Color{
		R: 0.0,
		G: 1.0,
		B: 0.0,
		A: 1.0,
	}
)

type Color struct {
	R, G, B, A float64
}

type Color32 struct {
	R, G, B, A float32
}

func (c Color) To32() Color32 {
	return Color32{
		R: float32(c.R),
		G: float32(c.G),
		B: float32(c.B),
		A: float32(c.A),
	}
}

// RandColor 随机颜色带透明通道
func RandColor(seed int64) Color {
	r := rand.New(rand.NewSource(seed))
	return Color{
		R: r.Float64(),
		G: r.Float64(),
		B: r.Float64(),
		A: 1.0,
	}
}

// Interpolation 颜色插值
func Interpolation(c1, c2 Color, factor float64) *Color {
	return &Color{
		R: c1.R + (c2.R-c1.R)*factor,
		G: c1.G + (c2.G-c1.G)*factor,
		B: c1.B + (c2.B-c1.B)*factor,
		A: c1.A + (c2.A-c1.A)*factor,
	}
}

// RandColorA 随机颜色带透明通道
func RandColorA(seed int64) Color {
	r := rand.New(rand.NewSource(seed))
	return Color{
		R: r.Float64(),
		G: r.Float64(),
		B: r.Float64(),
		A: r.Float64(),
	}
}
