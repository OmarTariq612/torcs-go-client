package controller

import (
	"bytes"
	"testing"
)

func TestMarshalText(t *testing.T) {
	action := CarAction{
		Accelerator: 5,
		Brake:       0.5,
		Clutch:      0.35,
		Gear:        5,
		Steering:    1,
		RestartRace: false,
		Focus:       45,
	}

	resultBytes, err := action.MarshalText()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expectedBytes := []byte("(accel 1) (brake 0.5) (clutch 0.35) (gear 5) (steer 1) (meta 0) (focus 45)")

	if !bytes.Equal(resultBytes, expectedBytes) {
		t.Fatalf("\noutput: %s\nexpected: %s", resultBytes, expectedBytes)
	}
}
