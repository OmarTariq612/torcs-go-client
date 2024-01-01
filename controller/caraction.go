package controller

import "fmt"

type CarAction struct {
	Accelerator float64 // 0..=1
	Brake       float64 // 0..=1
	Clutch      float64 // 0..=1
	Gear        int     // -1..=6
	Steering    float64 // -1..=1
	RestartRace bool
	Focus       int // -90..=90
}

func (ca *CarAction) ApplyLimits() {
	ca.Accelerator = max(0, min(1, ca.Accelerator))
	ca.Brake = max(0, min(1, ca.Brake))
	ca.Clutch = max(0, min(1, ca.Clutch))
	ca.Steering = max(-1, min(1, ca.Steering))
	ca.Gear = max(-1, min(6, ca.Gear))
}

func (ca *CarAction) MarshalText() (text []byte, err error) {
	ca.ApplyLimits()

	restartRace := 0
	if ca.RestartRace {
		restartRace = 1
	}
	return []byte(fmt.Sprintf("(accel %v) (brake %v) (clutch %v) (gear %v) (steer %v) (meta %v) (focus %v)", ca.Accelerator, ca.Brake, ca.Clutch, ca.Gear, ca.Steering, restartRace, ca.Focus)), nil
}
