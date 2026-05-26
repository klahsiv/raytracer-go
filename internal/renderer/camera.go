package renderer

import rl "github.com/gen2brain/raylib-go/raylib"

type Camera struct {
	Camera    rl.Camera
	MoveSpeed float32
}

func CreateDefaultCamera() *Camera {
	camera := rl.Camera{}
	camera.Fovy = 45
	camera.Position = rl.NewVector3(-7.08, 35.97, 22.82)
	camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
	camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	camera.Projection = rl.CameraPerspective

	cam := Camera{Camera: camera}
	return &cam
}

func (camera *Camera) GetInverseCameraViewMatrix() rl.Matrix {

	view := rl.GetCameraViewMatrix(&camera.Camera)
	invCamViewMatrix := rl.MatrixInvert(view)

	return invCamViewMatrix
}

func (camera *Camera) GetInverseCameraProjectionMatrix() rl.Matrix {

	proj := rl.GetCameraProjectionMatrix(&camera.Camera, float32(rl.GetScreenWidth())/float32(rl.GetScreenHeight()))
	invCamProjectionMatrix := rl.MatrixInvert(proj)

	return invCamProjectionMatrix
}

func (c *Camera) Update() bool {

	var cameraMoved bool = false

	forward := rl.GetCameraForward(&c.Camera)
	right := rl.GetCameraRight(&c.Camera)
	up := rl.GetCameraUp(&c.Camera)

	if rl.IsKeyDown(rl.KeyW) {
		offset := rl.Vector3Scale(forward, c.MoveSpeed)
		c.Camera.Position = rl.Vector3Add(c.Camera.Position, offset)
		c.Camera.Target = rl.Vector3Add(c.Camera.Target, offset)
		cameraMoved = true
	} else if rl.IsKeyDown(rl.KeyS) {
		offset := rl.Vector3Scale(forward, c.MoveSpeed)
		c.Camera.Position = rl.Vector3Subtract(c.Camera.Position, offset)
		c.Camera.Target = rl.Vector3Subtract(c.Camera.Target, offset)
		cameraMoved = true
	} else if rl.IsKeyDown(rl.KeyA) {
		offset := rl.Vector3Scale(right, c.MoveSpeed)
		c.Camera.Position = rl.Vector3Subtract(c.Camera.Position, offset)
		c.Camera.Target = rl.Vector3Subtract(c.Camera.Target, offset)
		cameraMoved = true
	} else if rl.IsKeyDown(rl.KeyD) {
		offset := rl.Vector3Scale(right, c.MoveSpeed)
		c.Camera.Position = rl.Vector3Add(c.Camera.Position, offset)
		c.Camera.Target = rl.Vector3Add(c.Camera.Target, offset)
		cameraMoved = true
	} else if rl.IsKeyDown(rl.KeyUp) {
		offset := rl.Vector3Scale(up, c.MoveSpeed)
		c.Camera.Position = rl.Vector3Add(c.Camera.Position, offset)
		c.Camera.Target = rl.Vector3Add(c.Camera.Target, offset)
		cameraMoved = true
	} else if rl.IsKeyDown(rl.KeyDown) {
		offset := rl.Vector3Scale(up, c.MoveSpeed)
		c.Camera.Position = rl.Vector3Subtract(c.Camera.Position, offset)
		c.Camera.Target = rl.Vector3Subtract(c.Camera.Target, offset)
		cameraMoved = true
	}
	return cameraMoved
}
