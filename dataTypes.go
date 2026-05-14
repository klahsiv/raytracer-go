package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Sphere struct {
	// pos : XYZ -> Center, Z -> Radius
	// material : X -> Material Index
	pos      rl.Vector4
	material rl.Vector4
}

type Triangle struct {
	// material : X -> Material Index
	posA, posB, posC          rl.Vector4
	normalA, normalB, normalC rl.Vector4
	material                  rl.Vector4
}

type Material struct {
	// emissive : XYZ -> Colour, W -> Strength
	colour   rl.Vector4
	emissive rl.Vector4
}

type AccumPixel struct {
	R float32
	G float32
	B float32
	A float32
}
