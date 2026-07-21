package gpu

type Scene struct {
	Spheres          []Sphere
	Triangles        []Triangle
	Materials        []Material
	Primitives       []GpuPrimitive
	Nodes            []GpuNode
	PrimitiveIndices []int
}
