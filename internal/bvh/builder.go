package bvh

import rl "github.com/gen2brain/raylib-go/raylib"

const MaxLeafSize = 2

type Builder struct {
	primitives []Primitive
}

func (b *Builder) Build(primitives []Primitive) *Tree {
	b.primitives = primitives
	N := len(primitives)
	primitiveIdx := make([]int, N)

	for i := range N {
		primitiveIdx[i] = i
	}
	tree := Tree{PrimitiveIndices: primitiveIdx}
	tree.Nodes = append(tree.Nodes, Node{})

	root := &tree.Nodes[0]
	root.FirstPrimitive = 0
	root.PrimitiveCount = N
	root.Left = 0
	root.Right = 0

	b.ComputeBounds(&tree, 0)
	b.SubDivide(&tree, 0)

	return &tree
}

func (b *Builder) ComputeBounds(tree *Tree, idx int) {

	node := &tree.Nodes[idx]
	node.Bounds = b.primitives[tree.PrimitiveIndices[node.FirstPrimitive]].Bounds

	for first, i := node.FirstPrimitive, 0; i < node.PrimitiveCount; i++ {
		p := b.primitives[tree.PrimitiveIndices[first+i]]
		node.Bounds = node.Bounds.Union(p.Bounds)
	}
}

func (b *Builder) SubDivide(tree *Tree, idx int) {
	node := &tree.Nodes[idx]
	if node.PrimitiveCount <= MaxLeafSize {
		return
	}

	axis := node.Bounds.LongestAxis()
	splitPos := Component(node.Bounds.Min, axis) + Component(node.Bounds.Extent(), axis)*0.5

	i := node.FirstPrimitive
	j := i + node.PrimitiveCount - 1

	for i <= j {
		pi := b.primitives[tree.PrimitiveIndices[i]]

		if Component(pi.Centroid, axis) < splitPos {
			i++
		} else {
			tree.PrimitiveIndices[i], tree.PrimitiveIndices[j] = tree.PrimitiveIndices[j], tree.PrimitiveIndices[i]
			j--
		}
	}

	leftCount := i - node.FirstPrimitive
	if leftCount == 0 || leftCount == node.PrimitiveCount {
		return
	}
	leftIdx := len(tree.Nodes)
	tree.Nodes = append(tree.Nodes, Node{})
	rightIdx := len(tree.Nodes)
	tree.Nodes = append(tree.Nodes, Node{})
	
	node = &tree.Nodes[idx]

	tree.Nodes[leftIdx].FirstPrimitive = node.FirstPrimitive
	tree.Nodes[leftIdx].PrimitiveCount = leftCount

	tree.Nodes[rightIdx].FirstPrimitive = i
	tree.Nodes[rightIdx].PrimitiveCount = node.PrimitiveCount - leftCount

	node.Left = leftIdx
	node.Right = rightIdx
	node.PrimitiveCount = 0

	b.ComputeBounds(tree, leftIdx)
	b.ComputeBounds(tree, rightIdx)

	b.SubDivide(tree, leftIdx)
	b.SubDivide(tree, rightIdx)

}

func Component(v rl.Vector3, axis int) float32 {
	switch axis {
	case 0:
		return v.X
	case 1:
		return v.Y
	case 2:
		return v.Z
	default:
		return v.X
	}
}
