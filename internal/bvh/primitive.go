package bvh

import rl "github.com/gen2brain/raylib-go/raylib"

type PrimitiveType uint32

const (
	PrimitiveSphere PrimitiveType = iota
	PrimitiveTriangle
)

type Primitive struct {
	Type     PrimitiveType
	Index    int
	Bounds   AABB
	Centroid rl.Vector3
}
