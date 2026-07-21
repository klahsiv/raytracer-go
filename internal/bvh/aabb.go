package bvh

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AABB struct {
	Min rl.Vector3
	Max rl.Vector3
}

func NewAABB(min, max rl.Vector3) *AABB {
	return &AABB{Min: min, Max: max}
}

func Empty() AABB {
	inf := float32(math.Inf(1))
	min := rl.NewVector3(inf, inf, inf)
	max := rl.NewVector3(-inf, -inf, -inf)

	return AABB{Min: min, Max: max}
}

func (a AABB) Union(b AABB) AABB {
	result := Empty()
	result.Max = rl.Vector3Max(a.Max, b.Max)
	result.Min = rl.Vector3Min(a.Min, b.Min)

	return result
}

func (a AABB) Extent() rl.Vector3 {
	extent := rl.Vector3Subtract(a.Max, a.Min)
	return extent
}

func (a AABB) LongestAxis() int {

	axis := 0
	extent := a.Extent()
	largest := extent.X

	if extent.Y > largest {
		axis = 1
		largest = extent.Y
	}
	if extent.Z > largest {
		axis = 2
		largest = extent.Z
	}
	return axis
}
