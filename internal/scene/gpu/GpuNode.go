package gpu

import rl "github.com/gen2brain/raylib-go/raylib"

type GpuNode struct {
	Min      rl.Vector4
	Max      rl.Vector4
	NodeInfo rl.Vector4
	// NodeInfo :  X -> Left, Y -> Right, Z -> First Primitive, W -> Primitive Count
}
