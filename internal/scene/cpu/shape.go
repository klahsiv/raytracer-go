package cpu

import (
	"ray-tracing/internal/bvh"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Shape interface {
	Bounds() bvh.AABB
	Centroid() rl.Vector3
}
