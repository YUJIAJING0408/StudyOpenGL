package render

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const (
	TRACEDEPTH = 5 //1次直接光照 + n-1次间接光照
	TRACEMINT  = 0.0001
	EPSILON    = 0.00001
)

type Scene struct {
	Name      string
	GeoCount  int
	GeoData   map[int]interface{}
	LightData map[int]interface{}
	Camera    Camera
	Ambient   AmbientLight
}

func NewScene(name string, c Camera, ambientColor mgl32.Vec3) Scene {
	lightData := make(map[int]interface{}, LIGHTTYPE-1) //环境光单独一份存在Scene中
	geoData := make(map[int]interface{}, GEOMETRYTYPE)
	for i := 0; i < LIGHTTYPE; i++ {
		switch i {
		case POINTLIGHT:
			lightData[i] = make([]PointLight, 0)
		case PARALLELLIGHT:
			lightData[i] = make([]ParallelLight, 0)
		case FACELIGHT:
			lightData[i] = make([]FaceLight, 0)
		case AMBIENTLIGHT:
			lightData[i] = make([]AmbientLight, 0)
		}
	}
	for i := 0; i < GEOMETRYTYPE; i++ {
		switch i {
		case SPHERE:
			geoData[i] = make([]Sphere, 0)
		case CUBE:
			geoData[i] = make([]Cube, 0)
		case BOX:
			geoData[i] = make([]Box, 0)
		case POINT:
			geoData[i] = make([]Point, 0)
		case OTHERS:
			//自定义几何体，包括assimp
			geoData[i] = make([]Other, 0)
		}
	}
	return Scene{
		Name:      name,
		GeoData:   geoData,
		LightData: lightData,
		Camera:    c,
		Ambient:   AmbientLight{Color: ambientColor},
	}
}

func (s Scene) AddLight(lightType int, light interface{}) {
	switch lightType {
	case POINTLIGHT:
		s.LightData[lightType] = append(s.LightData[lightType].([]PointLight), light.(PointLight))
	case PARALLELLIGHT:
		s.LightData[lightType] = append(s.LightData[lightType].([]ParallelLight), light.(ParallelLight))
	case FACELIGHT:
		s.LightData[lightType] = append(s.LightData[lightType].([]FaceLight), light.(FaceLight))
	}
}

func (s Scene) AddGeometry(geometryType int, geo interface{}) {
	switch geometryType {
	//构建几何数据
	case SPHERE:
		s.GeoData[geometryType] = append(s.GeoData[geometryType].([]Sphere), geo.(Sphere))
	case CUBE:
		s.GeoData[geometryType] = append(s.GeoData[geometryType].([]Cube), geo.(Cube))
	case BOX:
		s.GeoData[geometryType] = append(s.GeoData[geometryType].([]Box), geo.(Box))
	case POINT:
		s.GeoData[geometryType] = append(s.GeoData[geometryType].([]Point), geo.(Point))
	case OTHERS:
		s.GeoData[geometryType] = append(s.GeoData[geometryType].([]Other), geo.(Other))
	}
}

// SceneIntersect 判断指定观察光线是否与场景中物体相交，若相交，则存储与光线 首个相交的物体 的数组编号 、 法向量 与 交点坐标
func (s Scene) SceneIntersect(ray Ray) (isHit bool, hitGeoId int, hitPos, hitNorm mgl32.Vec3) {
	if s.GeoCount <= 0 {
		// 场景没有物体直接返回
		return
	} //如果没有几何体直接返回
	var closerPoint, closerHitNorm mgl32.Vec3
	var closerGeoId = -1
	var closerLen float64 = 10000000000.0
	//遍历每个几何体数据
	for _, sphere := range s.GeoData[SPHERE].([]Sphere) {
		if isIntersection, point, hitNormal, length := IntersectionWithSphere(ray, sphere); isIntersection {
			if !isHit {
				isHit = true
			}
			if closerLen >= length {
				closerGeoId = sphere.Id
				closerPoint = point
				closerLen = length
				closerHitNorm = hitNormal
			}
		}
	}
	if isHit {
		hitGeoId = closerGeoId
		hitPos = closerPoint
		hitNorm = closerHitNorm
	}
	return isHit, hitGeoId, hitPos, hitNorm
}

// ShadowRayIntersect 阴影检测
func (s Scene) ShadowRayIntersect(shadowRay Ray) (isShadow bool) {
	if s.GeoCount <= 0 {
		return false
	}
	for _, sphere := range s.GeoData[SPHERE].([]Sphere) {
		isIntersection, _, _, _ := IntersectionWithSphere(shadowRay, sphere)
		if isIntersection {
			return true
		}
	}
	return false
}

// Trace 光追
func (s Scene) Trace(ray Ray, depth int) (c mgl32.Vec3) {
	//超过递归深度最后返回环境光，避免返回零与前面的颜色乘完后整体颜色消失
	if depth > TRACEDEPTH {
		return s.Ambient.Color
	}
	var (
		isHit                  = false
		hitId                  = -1
		hitPos, hitNorm, color mgl32.Vec3
		hitMat                 Material
	)
	// 判断光线是否打到物体，否的话直接返回环境光
	if isHit, hitId, hitPos, hitNorm = s.SceneIntersect(ray); !isHit {
		return s.Ambient.Color
	}
	// 获取打击球体的材质
	hitSphere := s.GeoData[SPHERE].([]Sphere)[hitId]
	hitMat = hitSphere.Mat
	// 不同材质的渲染
	// ROUGH粗糙的材质
	if hitMat.Type == ROUGH {
		// 阴影计算
		// phong模型
		// 环境光
		color = s.Ambient.Color.Cross(hitMat.Color)
		// 测试阴影
		for _, light := range s.LightData[PARALLELLIGHT].([]ParallelLight) {
			ray = Ray{
				Position:  hitPos.Add(light.Direction.Mul(EPSILON)),
				Direction: light.Direction,
			}
			if !s.ShadowRayIntersect(ray) {
				// 漫反射
				diff := max(light.Direction.Dot(hitNorm), 0.0)
				diffuse := light.Color.Cross(hitMat.Color).Mul(diff)
				color = color.Add(diffuse)
				// 镜面高光
				viewDir := ray.Position.Sub(hitPos)
				reflectDir := Reflect(light.Direction.Mul(-1.0), hitNorm)
				if delta := viewDir.Dot(reflectDir); delta > 0 {
					specular := math.Pow(float64(delta), hitMat.Shininess)
					color = color.Mul(float32(1 + specular))
				}
			}
		}
		return color
	}

	// mat.type == REFLECTIVE
	// 计算反射光线方向（入射光线方向与法线方向点乘结果为负，故使用减法）
	reflectDir := ray.Direction.Sub(hitNorm.Mul(2 * ray.Direction.Dot(hitNorm)))
	// 构建反射光线
	reflectRay := Ray{
		Position:  hitPos.Add(reflectDir.Mul(EPSILON)),
		Direction: reflectDir.Normalize(),
	}
	// 使用菲涅尔方程的shlick近似方程计算反射光线颜色
	var cs mgl32.Vec3
	cosTheta := -ray.Direction.Dot(hitNorm)
	r0 := 0.0
	if hitMat.Type == REFLECTIVE {
		cs = hitMat.Color
	} else if hitMat.Type == REFRACTIVE {
		r0 = math.Pow((hitMat.RefractiveIndex-1)/(hitMat.RefractiveIndex+1), 2)
		R0 := mgl32.Vec3{float32(r0), float32(r0), float32(r0)}
		cs = R0.Add(mgl32.Vec3{1.0, 1.0, 1.0}.Sub(R0).Mul(float32(math.Pow(float64(1.0-cosTheta), 5.0))))
	}
	// 递归调用光追，直到最大深度
	color = cs.Cross(s.Trace(reflectRay, depth+1))
	if hitSphere.Mat.Type == REFRACTIVE {
		cos2Phi := 1 - (1-math.Pow(float64(cosTheta), 2))/math.Pow(hitMat.RefractiveIndex, 2)
		if cos2Phi >= 0 {
			reflectDir = ray.Direction.Sub(hitNorm.Mul(ray.Direction.Dot(hitNorm))).Mul(float32(1 / hitMat.RefractiveIndex)).Sub(hitNorm.Mul(float32(math.Sqrt(cos2Phi))))
			reflectRay = Ray{
				Position:  hitPos.Add(reflectDir.Mul(-EPSILON)),
				Direction: reflectDir.Normalize(),
			}
			color = color.Add(mgl32.Vec3{1.0, 1.0, 1.0}.Sub(cs).Cross(s.Trace(reflectRay, depth+1)))
		}
	}
	return color
}

// RenderScene 对整个场景进行光追
func (s Scene) RenderScene(viewWidth, viewHeight float64, function func(float64, float64, mgl32.Vec3)) {
	fov := max(viewWidth, viewHeight)
	for viewX := 0.0; viewX < fov; viewX++ {
		var cameraCoord mgl32.Vec2
		var pixelColor mgl32.Vec3
		for viewY := 0.0; viewY < fov; viewY++ {
			viewRay := s.Camera.GetRay(viewX, viewY, fov)
			pixelColor = s.Trace(viewRay, 0)
			// 计算视野中各像素点的位置，并调用函数renderFunc将像素点的坐标及颜色插入主程序的端点向量中
			cameraCoord = mgl32.Vec2{1, 0}.Mul(float32(2.0*(viewX+0.5)/fov - 1.0)).Add(mgl32.Vec2{0, 1}.Mul(float32(2.0*(viewY+0.5)/fov - 1.0)))
		}
		var renderX, renderY float64
		if flag := fov == viewWidth; flag {
			renderX = float64(cameraCoord.X())
		} else {
			renderX = float64(cameraCoord.X()) / viewWidth * viewHeight
		}
		if flag := fov == viewHeight; flag {
			renderX = float64(cameraCoord.Y())
		} else {
			renderX = float64(cameraCoord.Y()) / viewHeight * viewWidth
		}
		function(renderX, renderY, pixelColor)
	}
	return
}
