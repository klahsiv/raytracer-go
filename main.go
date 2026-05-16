package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	gl "github.com/go-gl/gl/v4.6-core/gl"
)

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetTraceLogLevel(rl.LogAll)
	rl.InitWindow(960, 540, "Ray Tracing with Raylib-go")

	err := gl.Init()
	if err != nil {
		panic(err)
	}
	width := rl.GetScreenWidth()
	height := rl.GetScreenHeight()

	var textureId uint32
	gl.GenTextures(1, &textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureId)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	program := loadComputeShader("raytracing.cs")

	/*
		pixels := make([]uint8, width*height*4)

		for i := 0; i < len(pixels); i += 4 {
			pixels[i] = 255
			pixels[i+1] = 0
			pixels[i+2] = 0
			pixels[i+3] = 255
		}

		gl.BindTexture(gl.TEXTURE_2D, textureId)

		gl.TexSubImage2D(
			gl.TEXTURE_2D,
			0,
			0,
			0,
			int32(width),
			int32(height),
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			unsafe.Pointer(&pixels[0]),
		)
	*/
	tex := rl.Texture2D{
		ID:      textureId,
		Width:   int32(width),
		Height:  int32(height),
		Mipmaps: 1,
		Format:  rl.UncompressedR8g8b8a8,
	}

	var fps float64
	var frameCounter int = 0
	var elapsed float64 = 0.0
	for !rl.WindowShouldClose() {
		start := time.Now()
		gl.UseProgram(program)
		gl.BindImageTexture(0, textureId, 0, false, 0, gl.WRITE_ONLY, gl.RGBA8)
		gl.DispatchCompute(uint32((width+7)/8), uint32((height+7)/8), 1)
		gl.MemoryBarrier(gl.ALL_BARRIER_BITS)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		//rl.DrawTexture(tex, 0, 0, rl.White)
		rl.DrawTexturePro(
			tex,
			rl.NewRectangle(0, 0, float32(width), -float32(height)),
			rl.NewRectangle(0, 0, float32(width), float32(height)),
			rl.Vector2{},
			0,
			rl.White,
		)
		rl.DrawText(fmt.Sprintf("FPS: %.2f", fps), 10, 10, 20, rl.Green)

		rl.EndDrawing()

		dt := time.Since(start).Seconds()
		elapsed += dt
		frameCounter++
		if frameCounter >= 100 {
			fps = float64(frameCounter) / elapsed
			frameCounter = 0
			elapsed = 0.0
		}
	}

	rl.CloseWindow()
}

func loadComputeShader(path string) uint32 {

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
		//log := "Error complining shader"

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

	return program

}
