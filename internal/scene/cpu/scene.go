package cpu

import "ray-tracing/internal/renderer"

type Scene struct {
	Spheres   []Sphere
	Materials []Material
	Camera    renderer.Camera
}
