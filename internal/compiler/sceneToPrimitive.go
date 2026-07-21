package compiler

import (
	"ray-tracing/internal/bvh"
	"ray-tracing/internal/scene/cpu"
)

func SceneToPrimitive(scene *cpu.Scene) []bvh.Primitive {
	primitivies := make([]bvh.Primitive, 0, (len(scene.Spheres) + len(scene.Triangles)))

	for idx, sphere := range scene.Spheres {
		primitivies = append(primitivies, bvh.Primitive{
			Type:   bvh.PrimitiveSphere,
			Index:  idx,
			Bounds:  sphere.Bounds(),
			Centroid: sphere.Centroid(),
		})
	}

	for idx, triangle := range scene.Triangles {
		primitivies = append(primitivies, bvh.Primitive{
			Type:   bvh.PrimitiveTriangle,
			Index:  idx,
			Bounds:  triangle.Bounds(),
			Centroid: triangle.Centroid(),
		})
	}
	return primitivies
}
