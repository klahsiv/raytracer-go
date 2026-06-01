package gpu

import rl "github.com/gen2brain/raylib-go/raylib"

type Sphere struct {
	CenterRadius rl.Vector4
	MaterialId   rl.Vector4
}

type Material struct {
	Colour   rl.Vector4
	Emission rl.Vector4
}

type Triangle struct {
	PosA, PosB, PosC          rl.Vector4
	NormalA, NormalB, NormalC rl.Vector4
	MaterialId                rl.Vector4
}
