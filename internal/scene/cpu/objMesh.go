package cpu

import rl "github.com/gen2brain/raylib-go/raylib"

type ObjMesh struct {
	Vertices []rl.Vector3
	Face     []ObjFace
}
