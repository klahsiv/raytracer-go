package cpu

import (
	"ray-tracing/internal/bvh"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Sphere struct {
	Center      rl.Vector3
	Radius      float32
	MaterialIdx int
}

func (s *Sphere) Bounds() bvh.AABB {
	min := rl.Vector3SubtractValue(s.Center, s.Radius)
	max := rl.Vector3AddValue(s.Center, s.Radius)

	bound := bvh.AABB{Min: min, Max: max}
	return bound
}

func (s *Sphere) Centroid() rl.Vector3 {
	return s.Center
}
