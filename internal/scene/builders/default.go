package builders

import (
	"ray-tracing/internal/renderer"
	"ray-tracing/internal/scene/cpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DefaultScene() *cpu.Scene {

	spheres := []cpu.Sphere{
		{Center: rl.NewVector3(0.0, 0.0, -1.0), Radius: 8.0, MaterialIdx: 0},
		{Center: rl.NewVector3(20.0, 0.0, -1.0), Radius: 10.0, MaterialIdx: 1},
		{Center: rl.NewVector3(-20.0, 0.0, -1.0), Radius: 6.0, MaterialIdx: 2},
		{Center: rl.NewVector3(0.0, 0.0, -500.0), Radius: 450.0, MaterialIdx: 3},
		{Center: rl.NewVector3(0.0, -510.0, 0.0), Radius: 500.0, MaterialIdx: 4},
	}
	materials := []cpu.Material{
		{Colour: rl.NewVector3(1.0, 0.0, 0.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
		{Colour: rl.NewVector3(0.0, 1.0, 0.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
		{Colour: rl.NewVector3(0.0, 0.0, 1.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
		{Colour: rl.NewVector3(0.0, 0.0, 0.0), EmissionColor: rl.NewVector3(1.0, 1.0, 1.0), EmissionStrength: 1.0},
		{Colour: rl.NewVector3(0.5, 1.0, 1.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
	}

	camera := rl.Camera{}
	camera.Fovy = 45
	camera.Position = rl.NewVector3(-7.08, 35.97, 22.82)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Projection = rl.CameraPerspective
	camera.Position = rl.NewVector3(0.0, 0.0, 15.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)

	cam := renderer.Camera{Camera: camera}
	defaultScene := cpu.Scene{Spheres: spheres, Materials: materials, Camera: cam}
	/*
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
	*/
	return &defaultScene
}
