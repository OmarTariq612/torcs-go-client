package controller

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	FocusSensorsNum     = 5
	TrackSensorsNum     = 19
	OpponentsSensorsNum = 36
	WheelsNum           = 4
)

type CarState struct {
	Gear                  int
	RacePosition          int
	Angle                 float64
	CurrentLapTime        float64
	Damage                float64
	DistanceFromStartLine float64
	DistanceRaced         float64
	FuelLevel             float64
	LastLapTime           float64
	RPM                   float64
	SpeedX                float64
	SpeedY                float64
	SpeedZ                float64
	TrackPosition         float64
	Z                     float64
	FocusSensors          [FocusSensorsNum]float64
	TrackEdgeSensors      [TrackSensorsNum]float64
	OpponentSensors       [OpponentsSensorsNum]float64
	WheenSpinVelocity     [WheelsNum]float64
}

func ScanProperties(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	i := bytes.IndexByte(data, '(')
	j := bytes.IndexByte(data, ')')
	if i >= 0 {
		if j > i {
			return j + 1, data[i+1 : j], nil
		} else {
			// TODO
			return 0, nil, errors.New("")
		}
	}
	return 0, nil, nil
}

func parseFloat64(properties []string) (float64, error) {
	if len(properties) != 1 {
		return 0, fmt.Errorf("1 != len(properties) = %v\n", len(properties))
	}
	return strconv.ParseFloat(properties[0], 64)
}

func parseInt(properties []string) (int, error) {
	if len(properties) != 1 {
		return 0, fmt.Errorf("1 != len(properties) = %v\n", len(properties))
	}
	intVal, err := strconv.ParseInt(properties[0], 10, 0)
	return int(intVal), err
}

func (cs *CarState) UnmarshalText(text []byte) error {
	scanner := bufio.NewScanner(bytes.NewBuffer(text))
	scanner.Split(ScanProperties)

	for scanner.Scan() {
		propertiesStr := scanner.Text()
		propertiesStr = strings.Trim(propertiesStr, " ")
		properties := strings.Split(propertiesStr, " ")
		if len(properties) < 2 {
			continue
		}
		var err error
		switch properties[0] {
		case "gear":
			cs.Gear, err = parseInt(properties[1:])

		case "racePos":
			cs.RacePosition, err = parseInt(properties[1:])

		case "angle":
			cs.Angle, err = parseFloat64(properties[1:])

		case "curLapTime":
			cs.CurrentLapTime, err = parseFloat64(properties[1:])

		case "damage":
			cs.Damage, err = parseFloat64(properties[1:])

		case "distFromStart":
			cs.DistanceFromStartLine, err = parseFloat64(properties[1:])

		case "distRaced":
			cs.DistanceRaced, err = parseFloat64(properties[1:])

		case "focus":
			for i := range cs.FocusSensors {
				cs.FocusSensors[i], err = parseFloat64(properties[i+1 : i+2])
				if err != nil {
					fmt.Printf("ERR (focus): %v\n", err)
					cs.FocusSensors[i] = 0
					err = nil
					// break
				}
			}

		case "fuel":
			cs.FuelLevel, err = parseFloat64(properties[1:])

		case "lastLapTime":
			cs.LastLapTime, err = parseFloat64(properties[1:])

		case "opponents":
			for i := range cs.OpponentSensors {
				cs.OpponentSensors[i], err = parseFloat64(properties[i+1 : i+2])
				if err != nil {
					fmt.Printf("ERR (opponents): %v\n", err)
					cs.OpponentSensors[i] = 0
					err = nil
					// break
				}
			}

		case "rpm":
			cs.RPM, err = parseFloat64(properties[1:])

		case "speedX":
			cs.SpeedX, err = parseFloat64(properties[1:])

		case "speedY":
			cs.SpeedY, err = parseFloat64(properties[1:])

		case "speedZ":
			cs.SpeedZ, err = parseFloat64(properties[1:])

		case "track":
			for i := range cs.TrackEdgeSensors {
				cs.TrackEdgeSensors[i], err = parseFloat64(properties[i+1 : i+2])
				if err != nil {
					fmt.Printf("ERR (track): %v\n", err)
					cs.TrackEdgeSensors[i] = 0
					err = nil
					// break
				}
			}

		case "trackPos":
			cs.TrackPosition, err = parseFloat64(properties[1:])

		case "wheelSpinVel":
			for i := range cs.WheenSpinVelocity {
				cs.WheenSpinVelocity[i], err = parseFloat64(properties[i+1 : i+2])
				if err != nil {
					fmt.Printf("ERR (wheelSpinVel): %v", err)
					cs.WheenSpinVelocity[i] = 0
					err = nil
					// break
				}
			}

		case "z":
			cs.Z, err = parseFloat64(properties[1:])

		}

		if err != nil {
			fmt.Printf("ERR (after): %v\n", err)
			// return err
		}
	}

	return nil
}
