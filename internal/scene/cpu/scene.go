package cpu

import "ray-tracing/internal/renderer"

type Scene struct {
	Spheres   []Sphere
	Triangles []Triangle
	Materials []Material
	Camera    renderer.Camera
}
