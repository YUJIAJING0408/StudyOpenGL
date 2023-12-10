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

func TestIntersectionWithSphere(t *testing.T) {
	var ray1 = NewRay(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{2, 1, 0})   //球外射线
	var ray2 = NewRay(mgl32.Vec3{4, 1, 0}, mgl32.Vec3{1, 0, 0})   //球心射线
	var ray3 = NewRay(mgl32.Vec3{4, 1.5, 0}, mgl32.Vec3{0, 1, 0}) //球内射线
	var ray4 = NewRay(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})   //球内射线（未命中）
	var g = Geometry{
		Type: SPHERE,
		Name: "球壳",
	}
	data := make(map[int]map[string]interface{}, 8)
	data[SPHERE] = make(map[string]interface{}, 8)
	data[SPHERE]["center"] = mgl32.Vec3{4, 1, 0}
	data[SPHERE]["radio"] = 2.0
	sphere := g.NewGeometry(data).(Sphere)
	if isIntersection, point := intersectionWithSphere(*ray1, sphere); isIntersection {
		fmt.Printf("Ray1交点为：\nx=%0.2f\ny=%0.2f\nz=%0.2f\n", point.X(), point.Y(), point.Z())
	}
	if isIntersection, point := intersectionWithSphere(*ray2, sphere); isIntersection {
		fmt.Printf("Ray2交点为：\nx=%0.2f\ny=%0.2f\nz=%0.2f\n", point.X(), point.Y(), point.Z())
	}
	if isIntersection, point := intersectionWithSphere(*ray3, sphere); isIntersection {
		fmt.Printf("Ray3交点为：\nx=%0.2f\ny=%0.2f\nz=%0.2f\n", point.X(), point.Y(), point.Z())
	}
	if isIntersection, point := intersectionWithSphere(*ray4, sphere); isIntersection {
		fmt.Printf("Ray4交点为：\nx=%0.2f\ny=%0.2f\nz=%0.2f\n", point.X(), point.Y(), point.Z())
	} else {
		fmt.Printf("Ray4未命中\n")
	}
}
