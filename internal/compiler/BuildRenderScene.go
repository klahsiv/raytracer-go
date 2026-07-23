package compiler

import (
	"fmt"
	"ray-tracing/internal/bvh"
	"ray-tracing/internal/scene/cpu"
)

func BuildRenderScene(cpuScene *cpu.Scene) RenderScene {

	primitives := SceneToPrimitive(cpuScene)

	builder := bvh.Builder{}
	tree := builder.Build(primitives)
	maxLeaf := 0
	sumLeaf := 0
	leafCount := 0

	for _, node := range tree.Nodes {
		if node.IsLeaf() {
			leafCount++
			sumLeaf += node.PrimitiveCount
			if node.PrimitiveCount > maxLeaf {
				maxLeaf = node.PrimitiveCount
			}
		}
	}

	fmt.Println("Leaf count:", leafCount)
	fmt.Println("Average leaf size:", float64(sumLeaf)/float64(leafCount))
	fmt.Println("Maximum leaf size:", maxLeaf)
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
