package cpu

import (
	"ray-tracing/internal/bvh"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Triangle struct {
	PosA, PosB, PosC          rl.Vector3
	NormalA, NormalB, NormalC rl.Vector3
	MaterialIdx               int
}

func (t *Triangle) Bounds() bvh.AABB {
	maxi := t.PosA
	mini := t.PosA

	maxi.X = max(maxi.X, t.PosB.X)
	mini.X = min(mini.X, t.PosB.X)

	maxi.Y = max(maxi.Y, t.PosB.Y)
	mini.Y = min(mini.Y, t.PosB.Y)

	maxi.Z = max(maxi.Z, t.PosB.Z)
	mini.Z = min(mini.Z, t.PosB.Z)

	maxi.X = max(maxi.X, t.PosC.X)
	mini.X = min(mini.X, t.PosC.X)

	maxi.Y = max(maxi.Y, t.PosC.Y)
	mini.Y = min(mini.Y, t.PosC.Y)

	maxi.Z = max(maxi.Z, t.PosC.Z)
	mini.Z = min(mini.Z, t.PosC.Z)

	bounds := bvh.AABB{Min: mini, Max: maxi}
	return bounds
}

func (t *Triangle) Centroid() rl.Vector3 {
	x := (t.PosA.X + t.PosB.X + t.PosC.X) / 3
	y := (t.PosA.Y + t.PosB.Y + t.PosC.Y) / 3
	z := (t.PosA.Z + t.PosB.Z + t.PosC.Z) / 3

	centroid := rl.NewVector3(x, y, z)
	return centroid
}
