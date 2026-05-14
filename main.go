package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	gl "github.com/go-gl/gl/v4.6-core/gl"
	"unsafe"
)

func main() {

	rl.InitWindow(960, 540, "Ray Tracing with Raylib-go")

	err := gl.Init()
	if err != nil {
		panic(err)
	}
	camera := rl.Camera{}
	camera.Fovy = 45
	camera.Position = rl.NewVector3(-7.08, 35.97, 22.82)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Projection = rl.CameraPerspective

	delta := float32(0.01)

	shader := rl.LoadShader("", "raytracing.fs")

	resolutionLoc := rl.GetShaderLocation(shader, "resolution")
	viewLoc := rl.GetShaderLocation(shader, "viewMat")
	projLoc := rl.GetShaderLocation(shader, "projMat")
	camPosLoc := rl.GetShaderLocation(shader, "camPos")

	resolution := make([]float32, 2)
	resolution[0] = float32(rl.GetScreenWidth())
	resolution[1] = float32(rl.GetScreenHeight())
	rl.SetShaderValue(shader, resolutionLoc, resolution, rl.ShaderUniformVec2)

	spheres := []Sphere{
		{pos: rl.NewVector4(0.0, 0.0, -1.0, 8.0), material: rl.NewVector4(0.0, 0.0, 0.0, 0.0)},
		{pos: rl.NewVector4(20.0, 0.0, -1.0, 10.0), material: rl.NewVector4(1.0, 0.0, 0.0, 0.0)},
		{pos: rl.NewVector4(-20.0, 0.0, -1.0, 6.0), material: rl.NewVector4(2.0, 0.0, 0.0, 0.0)},
		{pos: rl.NewVector4(0.0, 0.0, -500.0, 450.0), material: rl.NewVector4(3.0, 0.0, 0.0, 0.0)},
		{pos: rl.NewVector4(0.0, -510.0, 0.0, 500.0), material: rl.NewVector4(4.0, 0.0, 0.0, 0.0)},
	}
	materials := []Material{
		{colour: rl.NewVector4(1.0, 0.0, 0.0, 1.0), emissive: rl.NewVector4(0.0, 0.0, 0.0, 0.0)},
		{colour: rl.NewVector4(0.0, 1.0, 0.0, 1.0), emissive: rl.NewVector4(0.0, 0.0, 0.0, 0.0)},
		{colour: rl.NewVector4(0.0, 0.0, 1.0, 1.0), emissive: rl.NewVector4(0.0, 0.0, 0.0, 0.0)},
		{colour: rl.NewVector4(0.0, 0.0, 0.0, 1.0), emissive: rl.NewVector4(1.0, 1.0, 1.0, 1.0)},
		{colour: rl.NewVector4(0.5, 1.0, 1.0, 1.0), emissive: rl.NewVector4(0.0, 0.0, 0.0, 0.0)},
	}

	triangles := []Triangle{
		{
			posA: rl.NewVector4(-10.0, 0.0, -10.0, 0.0),
			posB: rl.NewVector4(10.0, 0.0, -10.0, 0.0),
			posC: rl.NewVector4(0.0, 10.0, -10.0, 0.0),

			// Flat upward-facing normals
			normalA: rl.NewVector4(0.0, 0.0, 1.0, 0.0),
			normalB: rl.NewVector4(0.0, 0.0, 1.0, 0.0),
			normalC: rl.NewVector4(0.0, 0.0, 1.0, 0.0),

			// Material index = 0
			material: rl.NewVector4(0.0, 0.0, 0.0, 0.0),
		},
	}

	sphereSsbo := rl.LoadShaderBuffer(uint32(len(spheres)*int(unsafe.Sizeof(spheres[0]))),
		unsafe.Pointer(&spheres[0]),
		rl.DynamicDraw)
	rl.BindShaderBuffer(sphereSsbo, 0)
	rl.UpdateShaderBuffer(sphereSsbo, unsafe.Pointer(&spheres[0]), uint32(len(spheres)*int(unsafe.Sizeof(spheres[0]))), 0)

	materialSsbo := rl.LoadShaderBuffer(uint32(len(materials)*int(unsafe.Sizeof(materials[0]))),
		unsafe.Pointer(&materials[0]),
		rl.DynamicDraw)
	rl.BindShaderBuffer(materialSsbo, 1)
	rl.UpdateShaderBuffer(materialSsbo, unsafe.Pointer(&materials[0]), uint32(len(materials)*int(unsafe.Sizeof(materials[0]))), 0)

	triangleSsbo := rl.LoadShaderBuffer(uint32(len(triangles)*int(unsafe.Sizeof(triangles[0]))),
		unsafe.Pointer(&triangles[0]),
		rl.DynamicDraw)
	rl.BindShaderBuffer(triangleSsbo, 2)
	rl.UpdateShaderBuffer(triangleSsbo, unsafe.Pointer(&triangles[0]), uint32(len(triangles)*int(unsafe.Sizeof(triangles[0]))), 0)

	sphereCount := len(spheres)
	sphereCountLoc := rl.GetShaderLocation(shader, "sphereCount")
	rl.SetShaderValue(shader, sphereCountLoc, []float32{float32(sphereCount)}, rl.ShaderUniformFloat)

	triangleCount := len(triangles)
	triangleCountLoc := rl.GetShaderLocation(shader, "triangleCount")
	rl.SetShaderValue(shader, triangleCountLoc, []float32{float32(triangleCount)}, rl.ShaderUniformFloat)

	iFrame := 0
	iFrameLoc := rl.GetShaderLocation(shader, "iFrame")

	width := rl.GetScreenWidth()
	height := rl.GetScreenHeight()

	accumulation := make([]AccumPixel, width*height)
	accumulationSsbo := rl.LoadShaderBuffer(uint32(len(accumulation)*int(unsafe.Sizeof(accumulation[0]))),
		unsafe.Pointer(&accumulation[0]),
		rl.DynamicDraw)
	rl.BindShaderBuffer(accumulationSsbo, 3)
	rl.UpdateShaderBuffer(accumulationSsbo, unsafe.Pointer(&accumulation[0]), uint32(len(accumulation)*int(unsafe.Sizeof(accumulation[0]))), 0)

	var fps float64
	var frameCounter int = 0
	var elapsed float64 = 0.0
	for !rl.WindowShouldClose() {
		start := time.Now()
		iFrame++
		if rl.IsKeyDown(rl.KeyUp) {
			camera.Position.Y += delta
		}
		if rl.IsKeyDown(rl.KeyDown) {
			camera.Position.Y -= delta
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			camera.Position.X -= delta
		}
		if rl.IsKeyDown(rl.KeyRight) {
			camera.Position.X += delta
		}
		if rl.IsKeyDown(rl.KeyW) {
			camera.Position.Z += delta
		}
		if rl.IsKeyDown(rl.KeyS) {
			camera.Position.Z -= delta
		}

		if rl.IsKeyPressed(rl.KeyP) {
			println(camera.Position.X, camera.Position.Y, camera.Position.Z, camera.Fovy)
		}
		view := rl.GetCameraViewMatrix(&camera)
		proj := rl.GetCameraProjectionMatrix(&camera, float32(rl.GetScreenWidth())/float32(rl.GetScreenHeight()))
		camPos := []float32{camera.Position.X, camera.Position.Y, camera.Position.Z}

		rl.SetShaderValueMatrix(shader, viewLoc, rl.MatrixInvert(view))
		rl.SetShaderValueMatrix(shader, projLoc, rl.MatrixInvert(proj))
		rl.SetShaderValue(shader, camPosLoc, camPos, rl.ShaderUniformVec3)
		rl.SetShaderValue(shader, iFrameLoc, []float32{float32(iFrame)}, rl.ShaderUniformFloat)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Green)
		rl.BeginShaderMode(shader)
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Green)
		rl.EndShaderMode()

		rl.DrawText(fmt.Sprintf("FPS: %.2f", fps), 10, 10, 20, rl.Green)

		//rl.DrawFPS(10, 10)
		rl.EndDrawing()

		gl.Finish()

		dt := time.Since(start).Seconds()
		elapsed += dt
		frameCounter++
		if frameCounter >= 100 {
			fps = float64(frameCounter) / elapsed
			frameCounter = 0
			elapsed = 0.0
		}
	}

	rl.UnloadShader(shader)
	rl.CloseWindow()
}

func handleCameraMovement(cam *rl.Camera) {
}
