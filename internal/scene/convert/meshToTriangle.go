package convert

import (
	"ray-tracing/internal/scene/cpu"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func MeshToTriangles(mesh cpu.ObjMesh, materialIdx int) []cpu.Triangle {
	var triangles []cpu.Triangle
	for _, face := range mesh.Faces {
		for i := 1; i < len(face.Vertices)-1; i++ {
			v0 := face.Vertices[0]
			v1 := face.Vertices[i]
			v2 := face.Vertices[i+1]

			va := mesh.Vertices[v0.VertexIdx]
			vb := mesh.Vertices[v1.VertexIdx]
			vc := mesh.Vertices[v2.VertexIdx]

			na := mesh.Normals[v0.NormalIdx]
			nb := mesh.Normals[v1.NormalIdx]
			nc := mesh.Normals[v2.NormalIdx]

			triangles = append(triangles, cpu.Triangle{
				PosA: va,
				PosB: vb,
				PosC: vc,

				NormalA: na,
				NormalB: nb,
				NormalC: nc,

				MaterialIdx: materialIdx,
			})
		}

	}

	return triangles
}

func calculateFaceNormal(a, b, c rl.Vector3) rl.Vector3 {

	edgeAB := b.Subtract(a)
	edgeAC := c.Subtract(a)
	normal := rl.Vector3CrossProduct(edgeAB, edgeAC)

	return normal.Normalize()
}
