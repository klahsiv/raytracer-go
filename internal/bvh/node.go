package bvh

type Node struct {
	Bounds         AABB
	Left           int
	Right          int
	FirstPrimitive int
	PrimitiveCount int
}

func (n *Node) IsLeaf() bool {
	return (n.PrimitiveCount > 0)
}
