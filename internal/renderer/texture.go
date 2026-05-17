package renderer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	gl "github.com/go-gl/gl/v4.6-core/gl"
)

type RenderTexture struct {
	ID      uint32
	Width   int
	Height  int
	Texture rl.Texture2D
}

func CreateRenderTexture(width, height int) *RenderTexture {
	var textureID uint32

	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, int32(width), int32(height), 0, gl.RGBA, gl.FLOAT, nil)

	rlTexture := rl.Texture2D{
		ID:     textureID,
		Width:  int32(width),
		Height: int32(height),
		Format: rl.UncompressedR32g32b32a32,
	}

	renderTexture := RenderTexture{
		ID:      textureID,
		Width:   width,
		Height:  height,
		Texture: rlTexture,
	}

	return &renderTexture
}

func (renderTexture *RenderTexture) Bind(binding uint32) {
	gl.BindImageTexture(binding, renderTexture.ID, 0, false, 0, gl.WRITE_ONLY, gl.RGBA32F)
}

func (renderTexture *RenderTexture) Draw() {

	rl.DrawTexturePro(
		renderTexture.Texture,
		rl.NewRectangle(0, 0, float32(renderTexture.Width), -float32(renderTexture.Height)),
		rl.NewRectangle(0, 0, float32(renderTexture.Width), float32(renderTexture.Height)),
		rl.Vector2{},
		0,
		rl.White,
	)
}

/*
	var textureId uint32
	gl.GenTextures(1, &textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureId)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(width), int32(height), 0, gl.RGBA, gl{}.UNSIGNED_BYTE, nil)
*/
