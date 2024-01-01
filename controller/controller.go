package controller

type Stage byte

type AnglesInitializer interface {
	InitAngles() []float64
}

type AnglesInitializerFunc func() []float64

func (d AnglesInitializerFunc) InitAngles() []float64 {
	return d()
}

func DefaultInitializer() []float64 {
	angles := make([]float64, 19)
	for i := range angles {
		angles[i] = -90 + float64(i)*10
	}
	return angles
}

type CarController interface {
	AnglesInitializer
	Control(*CarState, *CarAction) // CarAction is an output parameter
}
