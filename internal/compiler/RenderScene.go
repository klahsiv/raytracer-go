package compiler

import (
	"ray-tracing/internal/bvh"
	"ray-tracing/internal/scene/gpu"
)

type RenderScene struct {
	GPUScene   gpu.Scene
	Primitives []bvh.Primitive
	BVH        *bvh.Tree
}
