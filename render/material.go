package render

import "github.com/go-gl/mathgl/mgl32"

const (
	ROUGH = iota
	REFLECTIVE
	REFRACTIVE
)

type Material struct {
	Type            int
	Color           mgl32.Vec3
	Shininess       float64 //发光度
	RefractiveIndex float64 //材料折射率
}

func NewMaterial(t int, color mgl32.Vec3, s float64, r float64) Material {
	return Material{
		Type:            0,
		Color:           mgl32.Vec3{},
		Shininess:       0,
		RefractiveIndex: 0,
	}
}
