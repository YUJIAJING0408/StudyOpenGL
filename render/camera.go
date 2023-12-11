package render

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

// 相机移动事件
const (
	FORWARD = iota
	BACKWARD
	LEFT
	RIGHT
)

const (
	YAW         = -90.0
	PITCH       = 0.0
	SPEED       = 0.5
	SENSITIVITY = 0.005
	ZOOM        = 45.0
)

// Camera 相机类
type Camera struct {
	//属性
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Up       mgl32.Vec3
	Right    mgl32.Vec3
	WorldUp  mgl32.Vec3
	//欧拉角
	Yaw   float64
	Pitch float64
	//其他
	MoveSpeed        float64
	MouseSensitivity float64
	Zoom             float64
}

func NewCamera(pos, worldUp mgl32.Vec3, yaw, pitch float64) Camera {
	c := Camera{
		Position:         pos,
		WorldUp:          worldUp,
		Front:            mgl32.Vec3{0.0, 0.0, -1.0},
		Yaw:              yaw,
		Pitch:            pitch,
		MoveSpeed:        SPEED,
		MouseSensitivity: SENSITIVITY,
		Zoom:             ZOOM,
	}
	c.UpdateCameraVectors()
	return c
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	eye := c.Position
	center := c.Position.Add(c.Front)
	up := c.Up
	return mgl32.LookAtV(eye, center, up)
}

func (c *Camera) GetPerspective(aspect float32) mgl32.Mat4 {
	//fmt.Printf("zoom:%f\n", float32(c.Zoom))
	return mgl32.Perspective(float32(c.Zoom), aspect, 0.1, 100.0)
}

func (c *Camera) ProcessKeyboard(moveType int, deltaTime float64) {
	velocity := float32(c.MoveSpeed * deltaTime)
	println(moveType, velocity)
	if moveType == FORWARD {
		print("f")
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	}
	if moveType == BACKWARD {
		print("b")
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	}
	if moveType == LEFT {
		print("l")
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	}
	if moveType == RIGHT {
		print("r")
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	}
	fmt.Printf("x:%f,y:%f,z:%f\n", c.Position.X(), c.Position.Y(), c.Position.Z())
}

func (c *Camera) ProcessMouseMovement(xOffset float64, yOffset float64, constrainPitch bool) {
	xOffset *= c.MouseSensitivity
	yOffset *= c.MouseSensitivity
	c.Yaw += xOffset
	c.Pitch += yOffset
	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}

	c.UpdateCameraVectors()
}

func (c *Camera) ProcessMouseScroll(yOffset float64) {
	if c.Zoom >= 44.0 && c.Zoom <= 45.0 {
		c.Zoom -= yOffset * 0.01
	}
	if c.Zoom <= 44.0 {
		c.Zoom = 44.0
	}
	if c.Zoom >= 45.0 {
		c.Zoom = 45.0
	}
	//c.Zoom = 1.0
	//println(c.Zoom)
}

func (c *Camera) UpdateCameraVectors() {
	front := mgl32.Vec3{
		float32(math.Cos(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch))), //x
		float32(math.Sin(mgl64.DegToRad(c.Pitch))),                                   //y
		float32(math.Sin(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch))), //z
	}
	//重新构建相机的坐标系
	c.Front = front.Normalize()
	c.Right = front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}

func (c *Camera) GetRay(viewX, viewY, fov float64) (r Ray) {
	r.Direction = c.Front.Add(c.Right.Mul(float32(2*(viewX+0.5)/fov - 1))).Add(c.Up.Mul(float32(2*(viewY+0.5)/fov - 1))).Sub(c.Position)
	r.Position = c.Position
	return r
}
