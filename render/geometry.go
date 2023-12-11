package render

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const GEOMETRYTYPE = 5 //

const (
	SPHERE = iota
	BOX
	CUBE
	POINT
	//其他几何体
	OTHERS
)

// Geometry 几何体
type Geometry struct {
	Mat  Material //材质
	Type int
	Id   int // 不许是负数
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

// Other 其他几何体
type Other struct {
	Center mgl32.Vec3
	//？？其他属性
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

// Intersection 与射线求交点
func (g Geometry) Intersection(r Ray, Geo interface{}) (flag bool, point mgl32.Vec3) {
	switch g.Type {
	case SPHERE: //射线与球壳求交
		sphere, ok := Geo.(Sphere)
		if !ok {
			return
		} else {
			// 解法一（速度较慢）：射线方程 p' = p + t * dir ，表示t时刻位置，
			// 与球壳方程联立，解一元二次方程，将t用求根公式表示
			// r.Direction
			// 解法二（优化流程），先判读可不可能出现交点，部分情况可以直接退出计算
			flag, point, _, _ = IntersectionWithSphere(r, sphere)
		}
	case POINT:
		flag = false
	}
	return flag, point
}

func IntersectionWithSphere(ray Ray, sphere Sphere) (isIntersection bool, hitPoint, hitNorm mgl32.Vec3, rayTouchSphereLen float64) {
	isIntersection = true
	// 先判断是否在出射点是否在球内
	rayToSphere := sphere.Center.Sub(ray.Position)
	rayCenterToSphereLen := rayToSphere.Len() //出射点到球心距离
	projectOnRayDot := rayToSphere.Dot(ray.Direction)
	projectOnRayLenAbs := math.Abs(float64(projectOnRayDot)) // 投影长度
	if float64(rayCenterToSphereLen) < sphere.Radio {
		//在球内必有交点
		acrossFlats := math.Sqrt(math.Pow(float64(rayCenterToSphereLen), 2) - math.Pow(projectOnRayLenAbs, 2.0)) //夹角的对边
		path := math.Sqrt(math.Pow(sphere.Radio, 2) - math.Pow(acrossFlats, 2))                                  //球内平面圆通过垂径定理构建直角三角形的临边
		if projectOnRayDot > 0 {
			rayTouchSphereLen = projectOnRayLenAbs + path                              //等价时间t
			hitPoint = ray.Position.Add(ray.Direction.Mul(float32(rayTouchSphereLen))) //p' = p + t * dir
		} else if projectOnRayDot == 0 {
			hitPoint = ray.Position.Add(ray.Direction.Mul(float32(sphere.Radio)))
		} else {
			rayTouchSphereLen = path - projectOnRayLenAbs //等价时间t
			hitPoint = ray.Position.Add(ray.Direction.Mul(float32(rayTouchSphereLen)))
		}
	} else {
		//在球外
		//rayToSphereNor := rayToSphere.Normalize()
		//投影大于0说明夹角为锐角，小于0时不存在交点
		//x in (0,90) => cos(x) > 0
		//x in (90,180) => cos(x) < 0
		// 先求中垂线长度，勾股定理
		acrossFlats := math.Sqrt(math.Pow(float64(rayCenterToSphereLen), 2) - math.Pow(projectOnRayLenAbs, 2.0)) //夹角的对边
		path := math.Sqrt(math.Pow(sphere.Radio, 2) - math.Pow(acrossFlats, 2))                                  //球内平面圆通过垂径定理构建直角三角形的临边
		if projectOnRayDot > 0 && acrossFlats < sphere.Radio {
			//有且仅有这一种情况，夹角对边的长度小于球壳半径 且投影大于0
			rayTouchSphereLen = projectOnRayLenAbs - path                              //等价时间t
			hitPoint = ray.Position.Add(ray.Direction.Mul(float32(rayTouchSphereLen))) //p' = p + t * dir
		} else {
			isIntersection = false
		}
	}
	if isIntersection {
		hitNorm = hitPoint.Sub(sphere.Center).Normalize() //碰撞点法线
	}
	return isIntersection, hitPoint, hitNorm, rayTouchSphereLen
}
