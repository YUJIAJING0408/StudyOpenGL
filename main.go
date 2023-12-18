package main

import (
	"StudyOpenGL/render"
	"StudyOpenGL/utils"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"go/build"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"runtime"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
)

var keys = make(map[int]bool, 1024)
var camera render.Camera
var firstMouse = true
var lastX float64
var lastY float64
var lastTime float64

// 初始化项目
func init() {
	// 锁定线程是为了GLFW更稳定的运行
	runtime.LockOSThread()
	//utils.TestingMC()

}

func main() {
	// GLFW初始化
	if err := glfw.Init(); err != nil {
		log.Fatalln("GLFW初始化失败:", err)
	}
	// 延迟关闭GLFW
	defer glfw.Terminate()
	// GLFW初始化
	glfw.WindowHint(glfw.Resizable, glfw.False)                 //窗口伸缩
	glfw.WindowHint(glfw.ContextVersionMajor, 3)                //OpenGL大版本
	glfw.WindowHint(glfw.ContextVersionMinor, 3)                //OpenGL小版本
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) //使用OpenGL的核心模式
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4) //抗锯齿采样 4
	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	//绑定上下文
	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	//注册事件
	window.SetKeyCallback(keyCallBack)
	window.SetCursorPosCallback(mouseCallBack)
	window.SetScrollCallback(scrollCallBack)
	//禁用指针
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// GL初始化
	if err = gl.Init(); err != nil {
		panic(err)
	} else {
		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Println("OpenGL版本：", version[0:5])
		fmt.Println("Nvidia版本：", version[13:])
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	//默认着色器（顶点+片元）
	var shader = &utils.CommonShader{
		FragPath: "glsl/shader.frag",
		VertPath: "glsl/shader.vert",
	}
	program, err := shader.BuildProgram() //构建着色器
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program) //使用着色器，以便后续VBO的使用

	//camera = utils.NewCamera(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 1, 0}, -90.0, 0, 10.0, 0.2, 45.0)
	camera = render.NewCamera(mgl32.Vec3{0, 0, 10}, mgl32.Vec3{0, 1, 0}, -90.0, 0)
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	view := camera.GetViewMatrix()
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	//投影矩阵，
	//projection := mgl32.Perspective(45.0, windowWidth/windowHeight, 0.1, 100.0)
	projection := camera.GetPerspective(WindowWidth / WindowHeight)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	//camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	path, err := utils.Relative2FullPath("texture/square.png")
	if err != nil {
		return
	}
	texture, err := newTexture(path)
	if err != nil {
		log.Fatalln(err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	//vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	//gl.EnableVertexAttribArray(vertAttrib)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)

	//texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	//gl.EnableVertexAttribArray(texCoordAttrib)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.5, 0.5, 0.5, 1.0)

	angle := 0.0
	lastTime = glfw.GetTime()
	fps := utils.NewFps(lastTime)
	for !window.ShouldClose() {

		// Update
		now := glfw.GetTime()
		if ok, f := fps.Get(now); ok {
			fmt.Printf("当前帧数为：%d\n", f) //输出
		}
		deltaTime := now - lastTime
		lastTime = now

		DoMovement(deltaTime)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		view := camera.GetViewMatrix()

		perspective := camera.GetPerspective(WindowWidth / WindowHeight)

		//angle += deltaTime
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Render
		gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])
		gl.UniformMatrix4fv(projectionUniform, 1, false, &perspective[0])
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}

func keyCallBack(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// ESC键被被按下
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
	// 缓存按键被按下
	if key >= 0 && key <= 1024 {
		if action == glfw.Press {
			keys[int(key)] = true
		}
		if action == glfw.Release {
			keys[int(key)] = false
		}
	}
}

func mouseCallBack(w *glfw.Window, xPos float64, yPos float64) {
	//第一次移动时进行更新
	if firstMouse {
		lastX, lastY, firstMouse = xPos, xPos, false
	}
	xOffSet := xPos - lastX
	yOffSet := lastY - yPos
	lastX, lastY = xPos, yPos
	//获得两次的差值，并提交给camera进行更新
	camera.ProcessMouseMovement(xOffSet, yOffSet, true)
	//更新last值

}

func scrollCallBack(w *glfw.Window, xOff float64, yOff float64) {
	//xOff表示鼠标横向滚动，大部分鼠标都只有竖向滚轮
	//将yOff传输给缩放
	camera.ProcessMouseScroll(yOff)
	//print(camera.Zoom)
}

func DoMovement(deltaTime float64) {
	if keys[int(glfw.KeyW)] {
		camera.ProcessKeyboard(render.FORWARD, deltaTime)
	}
	if keys[int(glfw.KeyS)] {
		camera.ProcessKeyboard(render.BACKWARD, deltaTime)
	}
	if keys[int(glfw.KeyA)] {
		camera.ProcessKeyboard(render.LEFT, deltaTime)
	}
	if keys[int(glfw.KeyD)] {
		camera.ProcessKeyboard(render.RIGHT, deltaTime)
	}
}
