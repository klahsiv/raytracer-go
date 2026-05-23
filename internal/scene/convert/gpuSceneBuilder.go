package convert

import (
	"ray-tracing/internal/scene/cpu"
	"ray-tracing/internal/scene/gpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func BuildGpuScene(scene *cpu.Scene) *gpu.Scene {

	gpuScene := &gpu.Scene{}

	for _, sphere := range scene.Spheres {
		gpuScene.Spheres = append(gpuScene.Spheres, gpu.Sphere{
			CenterRadius: rl.NewVector4(sphere.Center.X, sphere.Center.Y, sphere.Center.Z, sphere.Radius),
			MaterialId:   rl.NewVector4(float32(sphere.MaterialIdx), 0.0, 0.0, 0.0)})
	}

	for _, material := range scene.Materials {
		gpuScene.Materials = append(gpuScene.Materials, gpu.Material{Colour: rl.NewVector4(material.Colour.X, material.Colour.Y, material.Colour.Z, 1.0),
			Emission: rl.NewVector4(material.EmissionColor.X, material.EmissionColor.Y, material.EmissionColor.Z, material.EmissionStrength)})
	}

	return gpuScene
}
