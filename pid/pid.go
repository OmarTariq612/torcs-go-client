package pid

type PID struct {
	Kp            float64
	Ki            float64
	Kd            float64
	OpErr         float64
	LastOpErr     float64
	Integral      float64
	Derivative    float64
	SetPoint      float64
	Dt            float64
	SampleTime    float64
	MaxIntegral   float64
	MaxDerivative float64
}

func NewPID(kp, ki, kd, setPoint, dt float64) *PID {
	return &PID{
		Kp:            kp,
		Ki:            ki,
		Kd:            kd,
		SetPoint:      setPoint,
		Dt:            dt,
		MaxIntegral:   1,
		MaxDerivative: 1,
	}
}

func (pid *PID) Compute(input float64) (output float64) {
	opErr := pid.SetPoint - input
	pid.Integral = max(-pid.MaxIntegral, min(pid.MaxIntegral, pid.Integral+opErr*pid.Dt))
	pid.Derivative = max(-pid.MaxDerivative, min(pid.MaxDerivative, (opErr-pid.LastOpErr)/pid.Dt))

	output = (pid.Kp * opErr) + (pid.Ki * pid.Integral) + (pid.Kd * pid.Derivative)
	pid.LastOpErr = opErr
	return output
}
