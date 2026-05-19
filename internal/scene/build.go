package scene

import (
	"ray-tracing/internal/gpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func BuildGpuScene(scene *Scene) *gpu.GpuScene {

	gpuScene := &gpu.GpuScene{}

	for _, sphere := range scene.Spheres {
		gpuScene.GpuSheres = append(gpuScene.GpuSheres, gpu.GpuSphere{
			CenterRadius: rl.NewVector4(sphere.center.X, sphere.center.Y, sphere.center.Z, sphere.radius),
			MaterialId:   rl.NewVector4(float32(sphere.materialIdx), 0, 0, 0)})
	}

	for _, material := range scene.Materials {
		gpuScene.GpuMaterials = append(gpuScene.GpuMaterials, gpu.GpuMaterial{Colour: rl.NewVector4(material.Color.X, material.Color.Y, material.Color.Z, 1.0),
			Emission: rl.NewVector4(material.EmissionColor.X, material.EmissionColor.Y, material.EmissionColor.Z, material.EmissionStrength)})
	}

	return gpuScene
}
