package obj

import (
	"bufio"
	"errors"
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

		if strings.HasPrefix(line, "v") {
			vertex, err := parseVertex(line)
			if err != nil {
				return cpu.ObjMesh{}, err
			}
			mesh.Vertices = append(mesh.Vertices, vertex)
		}
		if strings.HasPrefix(line, "f") {
			face, err := parseFace(line)
			if err != nil {
				return cpu.ObjMesh{}, err
			}
			mesh.Face = append(mesh.Face, face)
		}
	}

	if err := scanner.Err(); err != nil {
		return cpu.ObjMesh{}, err
	}

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

func parseFace(line string) (cpu.ObjFace, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		err := errors.New("Incorrect face data : " + line)
		return cpu.ObjFace{}, err
	}

	a, err := strconv.Atoi(fields[1])
	if err != nil {
		return cpu.ObjFace{}, err
	}
	b, err := strconv.Atoi(fields[2])
	if err != nil {
		return cpu.ObjFace{}, err
	}
	c, err := strconv.Atoi(fields[3])
	if err != nil {
		return cpu.ObjFace{}, err
	}

	return cpu.ObjFace{A: a - 1, B: b - 1, C: c - 1}, nil
}
