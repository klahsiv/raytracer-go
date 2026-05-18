package gpu

import rl "github.com/gen2brain/raylib-go/raylib"

type GpuSphere struct {
	CenterRadius rl.Vector4
	MaterialId   rl.Vector4
}

type GpuTriangle struct {
	centerRadius rl.Vector4
	materialId   rl.Vector4
}
type GpuMaterial struct {
	Colour   rl.Vector4
	Emission rl.Vector4
}
