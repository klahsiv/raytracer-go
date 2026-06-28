package cpu

type ObjVertex struct {
	VertexIdx int
	NormalIdx int
}

type ObjFace struct {
	Vertices []ObjVertex
}
