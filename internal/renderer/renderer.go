package renderer

import "github.com/go-gl/gl/v4.6-core/gl"

type Renderer struct {
	texture RenderTexture
	shader  ComputeShader
}

func NewRenderer(width, height int) *Renderer {

	texture := CreateRenderTexture(width, height)
	shader := LoadComputeShader("shaders/raytracing.cs")

	renderer := Renderer{texture: *texture, shader: *shader}

	return &renderer
}

func (renderer *Renderer) Render() {
	gl.UseProgram(renderer.shader.programID)
	renderer.texture.Bind(0)
	gl.DispatchCompute(uint32((renderer.texture.Width+7)/8), uint32((renderer.texture.Height+7)/8), 1)
	gl.MemoryBarrier(gl.ALL_BARRIER_BITS)
}

func (renderer *Renderer) Draw() {
	renderer.texture.Draw()
}
