package renderer

import rl "github.com/gen2brain/raylib-go/raylib"

type Camera struct {
	Camera rl.Camera
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
