package render

import "github.com/go-gl/mathgl/mgl32"

const LIGHTTYPE = 4

const (
	POINTLIGHT = iota
	PARALLELLIGHT
	FACELIGHT
	AMBIENTLIGHT
)

// ParallelLight 平行光
type ParallelLight struct {
	Direction mgl32.Vec3
	Color     mgl32.Vec3
}

// PointLight 平行光
type PointLight struct {
	Position mgl32.Vec3
	Color    mgl32.Vec3
}

type FaceLight struct {
	Position  mgl32.Vec3
	Direction mgl32.Vec3
	Radio     float64
}

type AmbientLight struct {
	Color mgl32.Vec3
}

func NewParallelLight(dir, color mgl32.Vec3) *ParallelLight {
	return &ParallelLight{
		Direction: dir.Normalize(),
		Color:     color,
	}
}

func Reflect(v mgl32.Vec3, n mgl32.Vec3) (reflectNorm mgl32.Vec3) {
	vn := v.Normalize()
	nn := n.Normalize()
	reflectNorm = nn.Mul(vn.Dot(nn) * 2).Sub(vn)
	return reflectNorm
}
