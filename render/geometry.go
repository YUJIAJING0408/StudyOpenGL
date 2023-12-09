package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	SPHERE = iota
	BOX
	CUBE
	POINT
	//其他几何体
)

// Geometry 几何体
type Geometry struct {
	Type int
	Name string
}

// Sphere 球体（球壳）
type Sphere struct {
	Geometry
	Center mgl32.Vec3
	Radio  float64
}

// Box 立方体
type Box struct {
	Geometry
	Center mgl32.Vec3
	Length float64
	Width  float64
	Height float64
}

// Cube 正方体
type Cube struct {
	Geometry
	Center mgl32.Vec3
	Size   float64
}

// Point 点
type Point struct {
	Geometry
	Center mgl32.Vec3 //等价Position
}

// NewGeometry 构造几何体 传入数据data为多重map，外层map的key为几何体类型int，内层为另一个map[string]interface{}
func (g Geometry) NewGeometry(data map[int]map[string]interface{}) interface{} {
	geometryData := data[g.Type]
	switch g.Type {
	case SPHERE:
		return Sphere{
			Geometry: g,
			Center:   geometryData["center"].(mgl32.Vec3),
			Radio:    geometryData["radio"].(float64),
		}
	case CUBE:
		return Cube{
			Geometry: g,
			Center:   geometryData["center"].(mgl32.Vec3),
			Size:     geometryData["size"].(float64),
		}
	case BOX:
		return Box{
			Geometry: g,
			Center:   geometryData["center"].(mgl32.Vec3),
			Length:   geometryData["length"].(float64),
			Width:    geometryData["width"].(float64),
			Height:   geometryData["height"].(float64),
		}
	default:
		return Point{
			Geometry: g,
			Center:   geometryData["center"].(mgl32.Vec3),
		}
	}
}

// IntersectionPoint 与射线求交点
func (g Geometry) IntersectionPoint(r Ray, Geo interface{}) {
	switch g.Type {
	case SPHERE: //射线与球壳求交
		//sphere := Geo.(Sphere)
		//print(sphere)
	}
}
