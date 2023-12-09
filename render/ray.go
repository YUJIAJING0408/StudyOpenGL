package render

import "github.com/go-gl/mathgl/mgl32"

// Ray 光线
type Ray struct {
	Position  mgl32.Vec3 //光源位置
	Direction mgl32.Vec3 //光源方向
}
