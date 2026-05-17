package renderer

import (
	"os"
	"strings"

	gl "github.com/go-gl/gl/v4.6-core/gl"
)

type ComputeShader struct {
	programID uint32
}

func LoadComputeShader(path string) *ComputeShader {

	sourceBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	source := string(sourceBytes) + "\x00"
	shader := gl.CreateShader(gl.COMPUTE_SHADER)

	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()

	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))

		gl.GetShaderInfoLog(
			shader,
			logLength,
			nil,
			gl.Str(log),
		)

		panic(log)
	}
	program := gl.CreateProgram()
	gl.AttachShader(program, shader)
	gl.LinkProgram(program)

	var linkStatus int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &linkStatus)

	if linkStatus == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))

		gl.GetProgramInfoLog(
			program,
			logLength,
			nil,
			gl.Str(log),
		)

		panic(log)
	}
	computeShader := ComputeShader{programID: program}

	return &computeShader

}
