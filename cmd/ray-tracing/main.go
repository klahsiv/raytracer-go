package main

import (
	"ray-tracing/internal/renderer"
	"ray-tracing/internal/scene/builders"
	"ray-tracing/internal/scene/convert"

	rl "github.com/gen2brain/raylib-go/raylib"
	gl "github.com/go-gl/gl/v4.6-core/gl"
)

func main() {
	rl.InitWindow(960, 540, "Ray Tracing with Raylib-go")

	err := gl.Init()
	if err != nil {
		panic(err)
	}
	width := rl.GetScreenWidth()
	height := rl.GetScreenHeight()

	scene := builders.DefaultScene()
	gpuScene := convert.BuildGpuScene(scene)

	r := renderer.NewRenderer(width, height)

	r.UploadScene(gpuScene)
	rl.DisableCursor()

	for !rl.WindowShouldClose() {

		cameraUpdated := scene.Camera.Update()
		if cameraUpdated {
			r.ResetAccumulation()
		}

		r.Render(&scene.Camera)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		r.Draw()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
