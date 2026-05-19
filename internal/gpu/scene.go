package gpu

type GpuScene struct {
	GpuSheres    []GpuSphere
	GpuMaterials []GpuMaterial
}

func BuildGpuSceneFromCpuScene(scene *Scene) *GpuScene {}
