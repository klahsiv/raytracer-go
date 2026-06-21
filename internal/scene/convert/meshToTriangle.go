package convert

import (
	"ray-tracing/internal/scene/cpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func MeshToTriangles(mesh cpu.ObjMesh, materialIdx int) []cpu.Triangle {
	var triangles []cpu.Triangle
	for _, face := range mesh.Face {
		vertexA := mesh.Vertices[face.A]
		vertexB := mesh.Vertices[face.B]
		vertexC := mesh.Vertices[face.C]

		normal := calculateFaceNormal(vertexA, vertexB, vertexC)

		triangles = append(triangles, cpu.Triangle{
			PosA: vertexA,
			PosB: vertexB,
			PosC: vertexC,

			NormalA: normal,
			NormalB: normal,
			NormalC: normal,

			MaterialIdx: materialIdx,
		})
	}

	return triangles
}

func calculateFaceNormal(a, b, c rl.Vector3) rl.Vector3 {

	edgeAB := b.Subtract(a)
	edgeAC := c.Subtract(a)
	normal := rl.Vector3CrossProduct(edgeAB, edgeAC)

	return normal.Normalize()
}
