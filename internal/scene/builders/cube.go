package builders

import (
	"ray-tracing/internal/renderer"
	"ray-tracing/internal/scene/cpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func CubeScene() *cpu.Scene {

	indices := []rl.Vector3{
		rl.NewVector3(-1, -1, 1),
		rl.NewVector3(1, -1, 1),
		rl.NewVector3(1, 1, 1),
		rl.NewVector3(-1, 1, 1),

		rl.NewVector3(-1, -1, -1),
		rl.NewVector3(1, -1, -1),
		rl.NewVector3(1, 1, -1),
		rl.NewVector3(-1, 1, -1),
	}

	normals := []rl.Vector3{
		rl.NewVector3(0, 0, 1),
		rl.NewVector3(0, 0, -1),
		rl.NewVector3(1, 0, 0),
		rl.NewVector3(-1, 0, 0),
		rl.NewVector3(0, 1, 0),
		rl.NewVector3(0, -1, 0),
	}

	materials := []cpu.Material{
		{Colour: rl.NewVector3(1.0, 0.0, 0.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
		{Colour: rl.NewVector3(0.0, 1.0, 0.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
		{Colour: rl.NewVector3(0.0, 0.0, 1.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
		{Colour: rl.NewVector3(0.0, 0.0, 0.0), EmissionColor: rl.NewVector3(1.0, 1.0, 1.0), EmissionStrength: 1.0},
		{Colour: rl.NewVector3(0.5, 1.0, 1.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
	}

	triangles := []cpu.Triangle{
		//+Z
		{
			PosA: indices[0], PosB: indices[1], PosC: indices[2],
			NormalA: normals[0], NormalB: normals[0], NormalC: normals[0], MaterialIdx: 0,
		},
		{
			PosA: indices[0], PosB: indices[2], PosC: indices[3],
			NormalA: normals[0], NormalB: normals[0], NormalC: normals[0], MaterialIdx: 0,
		},
		//-Z
		{
			PosA: indices[5], PosB: indices[4], PosC: indices[7],
			NormalA: normals[1], NormalB: normals[1], NormalC: normals[1], MaterialIdx: 1,
		},
		{
			PosA: indices[5], PosB: indices[7], PosC: indices[6],
			NormalA: normals[1], NormalB: normals[1], NormalC: normals[1], MaterialIdx: 1,
		},
		//+X
		{
			PosA: indices[1], PosB: indices[5], PosC: indices[6],
			NormalA: normals[2], NormalB: normals[2], NormalC: normals[2], MaterialIdx: 2,
		},
		{
			PosA: indices[1], PosB: indices[6], PosC: indices[2],
			NormalA: normals[2], NormalB: normals[2], NormalC: normals[2], MaterialIdx: 2,
		},
		//-X
		{
			PosA: indices[4], PosB: indices[0], PosC: indices[3],
			NormalA: normals[3], NormalB: normals[3], NormalC: normals[3], MaterialIdx: 0,
		},
		{
			PosA: indices[4], PosB: indices[3], PosC: indices[7],
			NormalA: normals[3], NormalB: normals[3], NormalC: normals[3], MaterialIdx: 0,
		},
		//+Y
		{
			PosA: indices[2], PosB: indices[6], PosC: indices[7],
			NormalA: normals[4], NormalB: normals[4], NormalC: normals[4], MaterialIdx: 1,
		},
		{
			PosA: indices[2], PosB: indices[7], PosC: indices[3],
			NormalA: normals[4], NormalB: normals[4], NormalC: normals[4], MaterialIdx: 1,
		},
		//-Y
		{
			PosA: indices[0], PosB: indices[4], PosC: indices[5],
			NormalA: normals[5], NormalB: normals[5], NormalC: normals[5], MaterialIdx: 2,
		},
		{
			PosA: indices[0], PosB: indices[5], PosC: indices[1],
			NormalA: normals[5], NormalB: normals[5], NormalC: normals[5], MaterialIdx: 2,
		},
	}

	postion := rl.NewVector3(0, 0, -10)
	target := rl.NewVector3(0, 0, 0)
	fovy := float32(45)
	cam := renderer.NewCamera(postion, target, fovy, rl.CameraPerspective)
	defaultScene := cpu.Scene{Spheres: []cpu.Sphere(nil), Triangles: triangles, Materials: materials, Camera: *cam}

	return &defaultScene
}
