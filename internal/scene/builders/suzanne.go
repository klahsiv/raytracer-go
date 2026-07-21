package builders

import (
	"fmt"
	"ray-tracing/internal/compiler"
	"ray-tracing/internal/obj"
	"ray-tracing/internal/renderer"
	"ray-tracing/internal/scene/cpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func SuzanneScene() *cpu.Scene {

	cubeMeshPath := "assets/suzanne.obj"
	materials := []cpu.Material{
		{Colour: rl.NewVector3(1.0, 0.0, 0.0), EmissionColor: rl.NewVector3(0.0, 0.0, 0.0), EmissionStrength: 0.0},
	}

	cubeMesh, err := obj.LoadObj(cubeMeshPath)
	if err != nil {
		panic("Error Loading ObjMesh")
	}

	triangles := compiler.MeshToTriangles(cubeMesh, 0)
	fmt.Println("Triangles : ", len(triangles))
	fmt.Println("Triangles 0: ", triangles[0])
	fmt.Println("Triangles 100: ", triangles[100])
	fmt.Println("Triangles 500: ", triangles[500])

	postion := rl.NewVector3(0, 0, -10)
	target := rl.NewVector3(0, 0, 0)
	fovy := float32(45)
	cam := renderer.NewCamera(postion, target, fovy, rl.CameraPerspective)
	defaultScene := cpu.Scene{Spheres: []cpu.Sphere(nil), Triangles: triangles, Materials: materials, Camera: *cam}

	return &defaultScene
}
