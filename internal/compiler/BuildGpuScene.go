package compiler

import (
	"ray-tracing/internal/bvh"
	"ray-tracing/internal/scene/cpu"
	"ray-tracing/internal/scene/gpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func BuildGpuScene(scene *cpu.Scene, primitives []bvh.Primitive, tree *bvh.Tree) *gpu.Scene {

	gpuScene := &gpu.Scene{}

	for _, sphere := range scene.Spheres {
		gpuScene.Spheres = append(gpuScene.Spheres, gpu.Sphere{
			CenterRadius: rl.NewVector4(sphere.Center.X, sphere.Center.Y, sphere.Center.Z, sphere.Radius),
			MaterialId:   rl.NewVector4(float32(sphere.MaterialIdx), 0.0, 0.0, 0.0)})
	}

	for _, triangle := range scene.Triangles {
		gpuScene.Triangles = append(gpuScene.Triangles, gpu.Triangle{
			PosA:       rl.NewVector4(triangle.PosA.X, triangle.PosA.Y, triangle.PosA.Z, 0.0),
			PosB:       rl.NewVector4(triangle.PosB.X, triangle.PosB.Y, triangle.PosB.Z, 0.0),
			PosC:       rl.NewVector4(triangle.PosC.X, triangle.PosC.Y, triangle.PosC.Z, 0.0),
			NormalA:    rl.NewVector4(triangle.NormalA.X, triangle.NormalA.Y, triangle.NormalA.Z, 0.0),
			NormalB:    rl.NewVector4(triangle.NormalB.X, triangle.NormalB.Y, triangle.NormalB.Z, 0.0),
			NormalC:    rl.NewVector4(triangle.NormalC.X, triangle.NormalC.Y, triangle.NormalC.Z, 0.0),
			MaterialId: rl.NewVector4(float32(triangle.MaterialIdx), 0.0, 0.0, 0.0)})
	}

	for _, material := range scene.Materials {
		gpuScene.Materials = append(gpuScene.Materials, gpu.Material{Colour: rl.NewVector4(material.Colour.X, material.Colour.Y, material.Colour.Z, 1.0),
			Emission: rl.NewVector4(material.EmissionColor.X, material.EmissionColor.Y, material.EmissionColor.Z, material.EmissionStrength)})
	}

	for _, primitive := range primitives {
		gpuScene.Primitives = append(gpuScene.Primitives, gpu.GpuPrimitive{
			SsboInfo: rl.NewVector4(float32(primitive.Type), float32(primitive.Index), 0.0, 0.0),
		})
	}

	for _, node := range tree.Nodes {
		gpuScene.Nodes = append(gpuScene.Nodes, gpu.GpuNode{
			Min:      rl.NewVector4(node.Bounds.Min.X, node.Bounds.Min.Y, node.Bounds.Min.Z, 0.0),
			Max:      rl.NewVector4(node.Bounds.Max.X, node.Bounds.Max.Y, node.Bounds.Max.Z, 0.0),
			NodeInfo: rl.NewVector4(float32(node.Left), float32(node.Right), float32(node.FirstPrimitive), float32(node.PrimitiveCount)),
		})
	}

	//gpuScene.PrimitiveIndices = append(gpuScene.PrimitiveIndices, tree.PrimitiveIndices...)

	for _, idx := range tree.PrimitiveIndices {
		gpuScene.PrimitiveIndices = append(gpuScene.PrimitiveIndices, int32(idx))
	}
	return gpuScene
}
