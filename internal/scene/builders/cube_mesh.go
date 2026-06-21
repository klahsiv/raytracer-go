package builders

import (
	"fmt"
	"ray-tracing/internal/obj"
	"ray-tracing/internal/renderer"
	"ray-tracing/internal/scene/convert"
	"ray-tracing/internal/scene/cpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func CubeMeshScene() *cpu.Scene {

	cubeMeshPath := "assets/cube.obj"
	materials := []cpu.Material{
		{Colour: rl.NewVector3(1.0, 0.0, 0.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
	}

	cubeMesh, err := obj.LoadObj(cubeMeshPath)
	if err != nil {
		panic("Error Loading ObjMesh")
	}

	triangles := convert.MeshToTriangles(cubeMesh, 0)
	fmt.Println(len(triangles))
	fmt.Println("Triangle 0 : ", triangles[0])

	postion := rl.NewVector3(0, 0, -10)
	target := rl.NewVector3(0, 0, 0)
	fovy := float32(45)
	cam := renderer.NewCamera(postion, target, fovy, rl.CameraPerspective)
	defaultScene := cpu.Scene{Spheres: []cpu.Sphere(nil), Triangles: triangles, Materials: materials, Camera: *cam}

	return &defaultScene
}
