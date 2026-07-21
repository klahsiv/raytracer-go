package compiler

import (
	"ray-tracing/internal/bvh"
	"ray-tracing/internal/scene/cpu"
)

func BuildRenderScene(cpuScene *cpu.Scene) RenderScene {

	primitives := SceneToPrimitive(cpuScene)

	builder := bvh.Builder{}
	tree := builder.Build(primitives)

	VerifyNode(tree, primitives, 0)

	gpuScene := BuildGpuScene(cpuScene, primitives, tree)

	renderScene := RenderScene{GPUScene: *gpuScene, Primitives: primitives, BVH: tree}

	return renderScene
}

func VerifyNode(tree *bvh.Tree, primitives []bvh.Primitive, idx int) {
	node := tree.Nodes[idx]

	if node.IsLeaf() {
		bounds := primitives[tree.PrimitiveIndices[node.FirstPrimitive]].Bounds

		for i := 1; i < node.PrimitiveCount; i++ {
			p := primitives[tree.PrimitiveIndices[node.FirstPrimitive+i]]
			bounds = bounds.Union(p.Bounds)
		}

		if bounds != node.Bounds {
			panic("leaf bounds wrong")
		}

		return
	}

	VerifyNode(tree, primitives, node.Left)
	VerifyNode(tree, primitives, node.Right)

	expected :=
		tree.Nodes[node.Left].Bounds.
			Union(tree.Nodes[node.Right].Bounds)

	if expected != node.Bounds {
		panic("internal bounds wrong")
	}
}
