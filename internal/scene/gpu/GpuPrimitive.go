package gpu

import rl "github.com/gen2brain/raylib-go/raylib"

type GpuPrimitive struct {
	SsboInfo rl.Vector4
	// SsboInfo : X -> Type, Y -> SSBO Index
}
