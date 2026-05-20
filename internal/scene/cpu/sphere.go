package cpu

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Sphere struct {
	Center      rl.Vector3
	Radius      float32
	MaterialIdx int
}
