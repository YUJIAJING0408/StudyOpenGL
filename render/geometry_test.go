package render

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"testing"
)

func TestNewGeometry(t *testing.T) {
	var data = make(map[int]map[string]interface{}, 1)
	geo := &Geometry{
		Type: SPHERE,
		Name: "球壳",
	}
	innerMap := make(map[string]interface{}, 8) //一般参数不会大于8个的
	innerMap["center"] = mgl32.Vec3{0, 0, 0}
	innerMap["radio"] = 5.0
	data[geo.Type] = innerMap
	geometry := geo.NewGeometry(data)
	sphere := geometry.(Sphere)
	fmt.Printf("几何体信息\n名字：%s\n中心：x=%0.2f，y=%0.2f，z=%0.2f\n半径：%0.2f\n", sphere.Name, sphere.Center.X(), sphere.Center.Y(), sphere.Center.Z(), sphere.Radio)
}
