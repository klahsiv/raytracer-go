package scene

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Sphere struct {
	center      rl.Vector3
	radius      float32
	materialIdx int
}
