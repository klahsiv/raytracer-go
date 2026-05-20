package scene

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Material struct {
	Colour           rl.Vector3
	EmissionColor    rl.Vector3
	EmissionStrength float32
}
