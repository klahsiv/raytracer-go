package builders

import (
	"ray-tracing/internal/renderer"
	"ray-tracing/internal/scene/cpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestScene() *cpu.Scene {

	spheres := []cpu.Sphere{
		{Center: rl.NewVector3(0.0, 0.0, -1.0), Radius: 8.0, MaterialIdx: 0},
	}
	materials := []cpu.Material{
		{Colour: rl.NewVector3(1.0, 0.0, 0.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
	}

	camera := rl.Camera{}
	camera.Fovy = 45
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Projection = rl.CameraPerspective
	camera.Position = rl.NewVector3(0.0, 0.0, 5.0)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)

	cam := renderer.Camera{Camera: camera}
	testScene := cpu.Scene{Spheres: spheres, Materials: materials, Camera: cam}

	return &testScene
}
