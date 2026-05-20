package gpu

import rl "github.com/gen2brain/raylib-go/raylib"

type Sphere struct {
	CenterRadius rl.Vector4
	MaterialId   rl.Vector4
}

type Triangle struct {
	centerRadius rl.Vector4
	materialId   rl.Vector4
}
type Material struct {
	Colour   rl.Vector4
	Emission rl.Vector4
}
