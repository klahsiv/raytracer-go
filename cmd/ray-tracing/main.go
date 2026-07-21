package main

import (
	"ray-tracing/internal/compiler"
	"ray-tracing/internal/renderer"
	"ray-tracing/internal/scene/builders"

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

	//scene := builders.CubeMeshScene()
	scene := builders.SuzanneScene()
	renderScene := compiler.BuildRenderScene(scene)

	r := renderer.NewRenderer(width, height)

	r.UploadScene(&renderScene.GPUScene)
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
