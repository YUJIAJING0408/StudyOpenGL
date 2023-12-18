package render

import (
	"StudyOpenGL/utils"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"os"
	"strings"
)

type CommonShader struct {
	FragPath string
	VertPath string
}

func (cs CommonShader) BuildProgram() (program uint32, err error) {
	var (
		vertShader uint32
		fragShader uint32
	)
	fullVertPath, err := utils.Relative2FullPath(cs.VertPath)
	if err != nil {
		return 0, err
	}
	if vertBytes, err := os.ReadFile(fullVertPath); err != nil {
		return 0, err
	} else {
		if vertShader, err = CompileShader(string(vertBytes), gl.VERTEX_SHADER); err != nil {
			return 0, err
		}
	}
	fullFragPath, err := utils.Relative2FullPath(cs.FragPath)
	if err != nil {
		return 0, err
	}
	if fragBytes, err := os.ReadFile(fullFragPath); err != nil {
		return 0, err
	} else {
		if fragShader, err = CompileShader(string(fragBytes), gl.FRAGMENT_SHADER); err != nil {
			return 0, err
		}
	}
	program = gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}
	//清除顶点着色器和片段着色器
	gl.DeleteShader(vertShader)
	gl.DeleteShader(fragShader)
	return program, nil
}

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}
