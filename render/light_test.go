package render

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"testing"
)

func TestReflect(t *testing.T) {
	v := mgl32.Vec3{-1, 1, 0}
	n := mgl32.Vec3{0, 1, 0}
	reflectNorm := Reflect(v, n)
	fmt.Printf("x = %0.2f,y = %0.2f,z = %0.2f\n", reflectNorm.X(), reflectNorm.Y(), reflectNorm.Z())
}
