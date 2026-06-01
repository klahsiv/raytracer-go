package cpu

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Triangle struct {
	PosA, PosB, PosC          rl.Vector3
	NormalA, NormalB, NormalC rl.Vector3
	MaterialIdx               int
}
