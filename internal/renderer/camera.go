package renderer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	Camera rl.Camera

	Yaw   float32
	Pitch float32

	WorldUp rl.Vector3

	Forward rl.Vector3
	Right   rl.Vector3
	Up      rl.Vector3

	MoveSpeed   float32
	Sensitivity float32
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

func NewCamera(position, target rl.Vector3, fovy float32, perspective rl.CameraProjection) *Camera {
	worldUp := rl.NewVector3(0.0, 1.0, 0.0)

	rlCamera := rl.Camera{}
	rlCamera.Fovy = fovy
	rlCamera.Position = position
	rlCamera.Target = target
	rlCamera.Up = worldUp
	rlCamera.Projection = perspective

	forward := rl.Vector3Normalize(rl.Vector3Subtract(target, position))

	pitch := float32(math.Asin(float64(forward.Y))) * rl.Rad2deg
	yaw := float32(math.Atan2(float64(forward.X), float64(forward.Z))) * rl.Rad2deg

	cam := Camera{Camera: rlCamera, Forward: forward, WorldUp: worldUp, Yaw: yaw, Pitch: pitch, MoveSpeed: 20.0, Sensitivity: 0.1}
	cam.UpdateVectors()

	return &cam
}

func (c *Camera) UpdateVectors() {

	yawRad := float64(c.Yaw * rl.Deg2rad)
	pitchRad := float64(c.Pitch * rl.Deg2rad)

	c.Forward = rl.NewVector3(
		float32(math.Cos(pitchRad)*math.Sin(yawRad)),
		float32(math.Sin(pitchRad)),
		float32(math.Cos(pitchRad)*math.Cos(yawRad)),
	)

	c.Forward = rl.Vector3Normalize(c.Forward)

	c.Right = rl.Vector3Normalize(
		rl.Vector3CrossProduct(c.Forward, c.WorldUp),
	)

	c.Up = rl.Vector3Normalize(
		rl.Vector3CrossProduct(c.Right, c.Forward),
	)

	c.Camera.Target = rl.Vector3Add(
		c.Camera.Position,
		c.Forward,
	)

	c.Camera.Up = c.Up
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
	dt := rl.GetFrameTime()
	//dt := float32(1.0)

	if rl.IsKeyDown(rl.KeyW) {
		offset := rl.Vector3Scale(c.Forward, c.MoveSpeed*dt)
		c.Camera.Position = rl.Vector3Add(c.Camera.Position, offset)
		cameraMoved = true
	}
	if rl.IsKeyDown(rl.KeyS) {
		offset := rl.Vector3Scale(c.Forward, c.MoveSpeed*dt)
		c.Camera.Position = rl.Vector3Subtract(c.Camera.Position, offset)
		cameraMoved = true
	}
	if rl.IsKeyDown(rl.KeyA) {
		offset := rl.Vector3Scale(c.Right, c.MoveSpeed*dt)
		c.Camera.Position = rl.Vector3Subtract(c.Camera.Position, offset)
		cameraMoved = true
	}
	if rl.IsKeyDown(rl.KeyD) {
		offset := rl.Vector3Scale(c.Right, c.MoveSpeed*dt)
		c.Camera.Position = rl.Vector3Add(c.Camera.Position, offset)
		cameraMoved = true
	}
	if rl.IsKeyDown(rl.KeyUp) {
		offset := rl.Vector3Scale(c.Up, c.MoveSpeed*dt)
		c.Camera.Position = rl.Vector3Add(c.Camera.Position, offset)
		cameraMoved = true
	}
	if rl.IsKeyDown(rl.KeyDown) {
		offset := rl.Vector3Scale(c.Up, c.MoveSpeed*dt)
		c.Camera.Position = rl.Vector3Subtract(c.Camera.Position, offset)
		cameraMoved = true
	}

	if cameraMoved {
		c.UpdateVectors()
	}

	if rl.IsCursorHidden() {
		mouseDelta := rl.GetMouseDelta()

		if mouseDelta.X != 0 || mouseDelta.Y != 0 {
			c.Rotate(-mouseDelta.X*c.Sensitivity, -mouseDelta.Y*c.Sensitivity)
			cameraMoved = true
		}
	}

	return cameraMoved
}

func (c *Camera) Rotate(yawDelta, pitchDelta float32) {
	c.Yaw += yawDelta
	c.Pitch += pitchDelta

	if c.Pitch > 89 {
		c.Pitch = 89
	}

	if c.Pitch < -89 {
		c.Pitch = -89
	}

	c.UpdateVectors()
}
