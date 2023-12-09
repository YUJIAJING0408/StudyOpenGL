package render

import "github.com/go-gl/mathgl/mgl32"

// ParallelLight 平行光
type ParallelLight struct {
	Direction mgl32.Vec3
	Color     mgl32.Vec3
}

func NewParallelLight(dir, color mgl32.Vec3) *ParallelLight {
	return &ParallelLight{
		Direction: dir.Normalize(),
		Color:     color,
	}
}
