package controller

import (
	"math"

	"github.com/OmarTariq612/torcs-go-client/pid"
)

var (
	gearUp      = [6]int{5000, 6000, 6000, 6500, 7000, 0}
	gearDown    = [6]int{0, 2500, 3000, 3000, 3500, 3500}
	wheelRadius = [4]float64{0.3179, 0.3179, 0.3276, 0.3276}
)

const (
	stuckTime              = 25
	stuckAngle             = 0.523598775
	maxSpeedDist           = 70
	maxSpeed               = 300
	sin5                   = 0.08716
	cos5                   = 0.99619
	steerLock              = 0.785398
	steerSensitivityOffset = 80.0
	wheelSensitivityCoeff  = 1
	absSlip                = 2.0
	absRange               = 3.0
	absMinSpeed            = 3.0
	clutchMax              = 0.5
	clutchDelta            = 0.05
	clutchRange            = 0.82
	clutchDeltaTime        = 0.02
	clutchDeltaRaced       = 10
	clutchDec              = 0.01
	clutchMaxModifier      = 1.3
	clutchMaxTime          = 1.5

	kp       = 4
	ki       = 1.5
	kd       = 0.01
	setPoint = 0
	dt       = 0.1
)

type SimpleDriver struct {
	*pid.PID
	stuck  int
	clutch float64
}

func NewSimpleDriver() *SimpleDriver {
	return &SimpleDriver{
		PID:    pid.NewPID(kp, ki, kd, setPoint, dt),
		stuck:  0,
		clutch: 0,
	}
}

func (sd *SimpleDriver) getGear(state *CarState) int {
	gear := state.Gear
	rpm := state.RPM
	switch {
	case gear < 1:
		gear = 1
	case gear < 6 && rpm >= float64(gearUp[gear-1]):
		gear++
	case gear > 1 && rpm <= float64(gearDown[gear-1]):
		gear--
	}

	return gear
}

func (sd *SimpleDriver) getSteer(state *CarState) float64 {
	measuredValue := (state.Angle / math.Pi) + (state.TrackEdgeSensors[6]-state.TrackEdgeSensors[12])/180 + (state.TrackEdgeSensors[7]-state.TrackEdgeSensors[11])/250 + (-state.TrackPosition)/7
	output := sd.Compute(measuredValue)
	return -output / 4
}

func (sd *SimpleDriver) getAccel(state *CarState) float64 {
	if math.Abs(state.TrackPosition) < 1 {
		rxSensor := state.TrackEdgeSensors[10]
		cSensor := state.TrackEdgeSensors[9]
		sxSensor := state.TrackEdgeSensors[8]

		var targetSpeed float64
		if cSensor > maxSpeedDist || (cSensor >= rxSensor && cSensor >= sxSensor) {
			targetSpeed = maxSpeed
		} else {
			// approaching a turn on right
			if rxSensor > sxSensor {
				// computing approximately the "angle" of turn
				h := cSensor * sin5
				b := rxSensor - cSensor*cos5
				sinAngle := b * b / (h*h + b*b)
				// estimate the target speed depending on turn and on how close it is
				targetSpeed = maxSpeed * (cSensor * sinAngle / maxSpeedDist)
			} else {
				// computing approximately the "angle" of turn
				h := cSensor * sin5
				b := sxSensor - cSensor*cos5
				sinAngle := b * b / (h*h + b*b)
				// estimate the target speed depending on turn and on how close it is
				targetSpeed = maxSpeed * (cSensor * sinAngle / maxSpeedDist)
			}
		}

		return 2/(1+math.Exp(state.SpeedX-targetSpeed)) - 1
	} else {
		return 0.3
	}
}

func (sd *SimpleDriver) Control(state *CarState, action *CarAction) {
	if math.Abs(state.Angle) > stuckAngle {
		sd.stuck++
	} else {
		sd.stuck = 0
	}

	if sd.stuck > stuckTime {
		steer := -state.Angle / steerLock
		gear := -1
		if state.Angle*state.TrackPosition > 0 {
			gear = 1
			steer = -steer
		}

		sd.clutching(state)

		action.Accelerator = 1
		action.Brake = 0
		action.Gear = gear
		action.Steering = steer
		action.Clutch = sd.clutch

	} else {

		accelAndBrake := sd.getAccel(state)
		gear := sd.getGear(state)
		steer := sd.getSteer(state)

		if steer < -1 {
			steer = -1
		}

		if steer > 1 {
			steer = 1
		}

		var accel, brake float64
		if accelAndBrake > 0 {
			accel = accelAndBrake
			brake = 0
		} else {
			accel = 0
			brake = sd.filterABS(state, -accelAndBrake)
		}
		sd.clutching(state)

		action.Accelerator = accel
		action.Brake = brake
		action.Gear = gear
		action.Steering = steer
		action.Clutch = sd.clutch
	}
}

func (sd *SimpleDriver) filterABS(state *CarState, brake float64) float64 {
	speed := state.SpeedX / 3.6
	if speed < absMinSpeed {
		return brake
	}

	slip := 0.0
	for i := 0; i < 4; i++ {
		slip += state.WheenSpinVelocity[i] * wheelRadius[i]
	}
	slip = speed - slip/4
	if slip > absSlip {
		brake = brake - (slip-absSlip)/absRange
	}

	if brake < 0 {
		return 0
	}

	return brake
}

func (sd *SimpleDriver) clutching(state *CarState) {
	maxClutch := clutchMax
	delta := clutchDelta
	if state.Gear < 2 {
		delta *= 50
		if state.CurrentLapTime < clutchMaxTime {
			sd.clutch = maxClutch
		}
	}
	if sd.clutch != maxClutch {
		sd.clutch -= delta
		sd.clutch = max(0.0, sd.clutch)
	} else {
		sd.clutch -= clutchDec
	}
}

func (sd *SimpleDriver) InitAngles() []float64 {
	angles := make([]float64, 19)
	for i := 0; i < 5; i++ {
		angles[i] = -90 + float64(i)*15
		angles[18-i] = 90 - float64(i)*15
	}
	for i := 5; i < 9; i++ {
		angles[i] = -20 + float64(i-5)*5
		angles[18-i] = 20 - float64(i-5)*5
	}
	angles[9] = 0
	return angles
}
