package obj

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"ray-tracing/internal/scene/cpu"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadObj(path string) (cpu.ObjMesh, error) {
	file, err := os.Open(path)
	if err != nil {
		return cpu.ObjMesh{}, err
	}
	defer file.Close()

	var mesh cpu.ObjMesh

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "vn ") {
			vertexNormal, err := parseVertexNormal(line)
			if err != nil {
				return cpu.ObjMesh{}, err
			}
			mesh.Normals = append(mesh.Normals, vertexNormal)
		}
		if strings.HasPrefix(line, "v ") {
			vertex, err := parseVertex(line)
			if err != nil {
				return cpu.ObjMesh{}, err
			}
			mesh.Vertices = append(mesh.Vertices, vertex)
		}
		if strings.HasPrefix(line, "f ") {
			face, err := parseFace(line)
			if err != nil {
				return cpu.ObjMesh{}, err
			}
			mesh.Faces = append(mesh.Faces, face)
		}
	}

	if err := scanner.Err(); err != nil {
		return cpu.ObjMesh{}, err
	}
	fmt.Println(len(mesh.Vertices))
	fmt.Println(len(mesh.Normals))
	fmt.Println(len(mesh.Faces))

	return mesh, nil
}

func parseVertex(line string) (rl.Vector3, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		err := errors.New("Incorrect vertex data : " + line)
		return rl.Vector3{}, err
	}

	x, err := strconv.ParseFloat(fields[1], 32)
	if err != nil {
		return rl.Vector3{}, err
	}
	y, err := strconv.ParseFloat(fields[2], 32)
	if err != nil {
		return rl.Vector3{}, err
	}
	z, err := strconv.ParseFloat(fields[3], 32)
	if err != nil {
		return rl.Vector3{}, err
	}
	return rl.NewVector3(float32(x), float32(y), float32(z)), nil
}

func parseVertexNormal(line string) (rl.Vector3, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		err := errors.New("Incorrect vertex normal data : " + line)
		return rl.Vector3{}, err
	}

	x, err := strconv.ParseFloat(fields[1], 32)
	if err != nil {
		return rl.Vector3{}, err
	}
	y, err := strconv.ParseFloat(fields[2], 32)
	if err != nil {
		return rl.Vector3{}, err
	}
	z, err := strconv.ParseFloat(fields[3], 32)
	if err != nil {
		return rl.Vector3{}, err
	}
	return rl.NewVector3(float32(x), float32(y), float32(z)), nil
}

func parseFace(line string) (cpu.ObjFace, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		err := errors.New("Incorrect face data : " + line)
		return cpu.ObjFace{}, err
	}

	face := cpu.ObjFace{}

	for i := 1; i < len(fields); i++ {
		fv, err := parseFaceVertex(fields[i])
		if err != nil {
			return cpu.ObjFace{}, err
		}

		face.Vertices = append(face.Vertices, fv)
	}

	return face, nil
}

func parseFaceVertex(token string) (cpu.ObjVertex, error) {
	parts := strings.Split(token, "//")
	if len(parts) != 2 {
		err := errors.New("Unsupported face token : " + token)
		return cpu.ObjVertex{}, err
	}
	vertexIdx, err := strconv.Atoi(parts[0])
	if err != nil {
		return cpu.ObjVertex{}, err
	}
	normalIdx, err := strconv.Atoi(parts[1])
	if err != nil {
		return cpu.ObjVertex{}, err
	}

	return cpu.ObjVertex{VertexIdx: vertexIdx - 1, NormalIdx: normalIdx - 1}, nil
}
