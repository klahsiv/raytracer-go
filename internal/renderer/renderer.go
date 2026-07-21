package renderer

import (
	"fmt"
	"ray-tracing/internal/scene/gpu"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type Renderer struct {
	texture             RenderTexture
	shader              ComputeShader
	renderFrame         int32
	accumulationTexture RenderTexture
}

func NewRenderer(width, height int) *Renderer {

	texture := CreateRenderTexture(width, height)
	accumulationTexture := CreateRenderTexture(width, height)
	shader := LoadComputeShader("shaders/raytracing.cs")

	renderer := Renderer{texture: *texture, shader: *shader, renderFrame: 0, accumulationTexture: *accumulationTexture}

	return &renderer
}

func (renderer *Renderer) Render(camera *Camera) {
	gl.UseProgram(renderer.shader.programID)

	renderer.renderFrame++
	iFrameLoc := renderer.shader.GetUniformLocation("iFrame")
	gl.Uniform1i(iFrameLoc, renderer.renderFrame)

	viewLoc := renderer.shader.GetUniformLocation("invViewMat")
	view := rl.MatrixToFloat(camera.GetInverseCameraViewMatrix())
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])

	projLoc := renderer.shader.GetUniformLocation("invProjMat")
	projection := rl.MatrixToFloat(camera.GetInverseCameraProjectionMatrix())
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	camPosLoc := renderer.shader.GetUniformLocation("camPos")
	gl.Uniform3f(camPosLoc, camera.Camera.Position.X, camera.Camera.Position.Y, camera.Camera.Position.Z)

	renderer.texture.Bind(0)
	renderer.accumulationTexture.Bind(1)
	gl.DispatchCompute(uint32((renderer.texture.Width+7)/8), uint32((renderer.texture.Height+7)/8), 1)
	gl.MemoryBarrier(gl.ALL_BARRIER_BITS)
}

func (renderer *Renderer) Draw() {
	renderer.texture.Draw()
	rl.DrawFPS(10, 10)
}

func (renderer *Renderer) UploadScene(scene *gpu.Scene) {

	gl.UseProgram(renderer.shader.programID)

	sphereSsbo := uploadSsbo(scene.Spheres, 2)
	if sphereSsbo > 0 {
		rl.UpdateShaderBuffer(sphereSsbo, unsafe.Pointer(&scene.Spheres[0]), uint32(len(scene.Spheres)*int(unsafe.Sizeof(scene.Spheres[0]))), 0)
	}

	materialSsbo := uploadSsbo(scene.Materials, 3)
	if materialSsbo > 0 {
		rl.UpdateShaderBuffer(materialSsbo, unsafe.Pointer(&scene.Materials[0]), uint32(len(scene.Materials)*int(unsafe.Sizeof(scene.Materials[0]))), 0)
	}

	triangleSsbo := uploadSsbo(scene.Triangles, 4)
	if triangleSsbo > 0 {
		rl.UpdateShaderBuffer(triangleSsbo, unsafe.Pointer(&scene.Triangles[0]), uint32(len(scene.Triangles)*int(unsafe.Sizeof(scene.Triangles[0]))), 0)
	}

	nodeSsbo := uploadSsbo(scene.Nodes, 5)
	if nodeSsbo > 0 {
		rl.UpdateShaderBuffer(nodeSsbo, unsafe.Pointer(&scene.Nodes[0]), uint32(len(scene.Nodes)*int(unsafe.Sizeof(scene.Nodes[0]))), 0)
	}

	primitiveSsbo := uploadSsbo(scene.Primitives, 6)
	if primitiveSsbo > 0 {
		rl.UpdateShaderBuffer(primitiveSsbo, unsafe.Pointer(&scene.Primitives[0]), uint32(len(scene.Primitives)*int(unsafe.Sizeof(scene.Primitives[0]))), 0)
	}

	primitiveIndicesSsbo := uploadSsbo(scene.PrimitiveIndices, 7)
	if primitiveIndicesSsbo > 0 {
		rl.UpdateShaderBuffer(primitiveIndicesSsbo, unsafe.Pointer(&scene.PrimitiveIndices[0]), uint32(len(scene.PrimitiveIndices)*int(unsafe.Sizeof(scene.PrimitiveIndices[0]))), 0)
	}

	sphereCount := len(scene.Spheres)
	sphereCountLoc := renderer.shader.GetUniformLocation("sphereCount")
	gl.Uniform1i(sphereCountLoc, int32(sphereCount))

	triangleCount := len(scene.Triangles)
	triangleCountLoc := renderer.shader.GetUniformLocation("triangleCount")
	gl.Uniform1i(triangleCountLoc, int32(triangleCount))

	fmt.Println("Upload Scene")
}

func uploadSsbo[T any](data []T, binding int32) uint32 {
	if len(data) == 0 {
		return 0
	}
	size := uint32(len(data) * int(unsafe.Sizeof(data[0])))
	ssbo := rl.LoadShaderBuffer(size, unsafe.Pointer(&data[0]), rl.StaticRead)

	rl.BindShaderBuffer(ssbo, uint32(binding))

	return ssbo
}

func (renderer *Renderer) ResetAccumulation() {

	renderer.renderFrame = 0
	gl.ClearTexImage(renderer.texture.ID, 0, gl.RGBA, gl.FLOAT, nil)
	gl.ClearTexImage(renderer.accumulationTexture.ID, 0, gl.RGBA, gl.FLOAT, nil)
}
