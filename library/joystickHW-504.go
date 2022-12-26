package main

import (
	"machine"
	"time"
)

type JOYSTICK_HW504 struct {
	SWITCH_PIN   machine.Pin
	ANALOG_PIN_X machine.ADC
	ANALOG_PIN_Y machine.ADC
}

var GLOBAL_X_CHANGER uint16 = 32767
var GLOBAL_Y_CHANGER uint16 = 32767

func NewJoy(swPin, x, y machine.Pin) JOYSTICK_HW504 {
	machine.InitADC()
	dummy := JOYSTICK_HW504{
		SWITCH_PIN:   swPin,
		ANALOG_PIN_X: machine.ADC{Pin: x},
		ANALOG_PIN_Y: machine.ADC{Pin: y},
	}
	dummy.SWITCH_PIN.Configure(machine.PinConfig{Mode: machine.PinInput})
	dummy.ANALOG_PIN_X.Configure(machine.ADCConfig{})
	dummy.ANALOG_PIN_Y.Configure(machine.ADCConfig{})

	return dummy
}

func (joy *JOYSTICK_HW504) Test() {
	for {
		time.Sleep(time.Microsecond * 50)
		x := joy.ANALOG_PIN_X.Get()
		y := joy.ANALOG_PIN_Y.Get()
		if y > 30000 && y < 36000 && x > 30000 && x < 36000 {
			continue
		}
		if x != GLOBAL_X_CHANGER || y != GLOBAL_Y_CHANGER {
			GLOBAL_X_CHANGER = x
			GLOBAL_Y_CHANGER = y
			println("X Value: ", GLOBAL_X_CHANGER, "Y Value: ", GLOBAL_Y_CHANGER)
			print("\n")
		}
	}
}

func (joy *JOYSTICK_HW504) ExecuteIfChange(timeBetween time.Duration, callback func(x, y uint16), isDebug bool) {
	time.Sleep(timeBetween)
	x := joy.ANALOG_PIN_X.Get()
	y := joy.ANALOG_PIN_Y.Get()
	if y > 30000 && y < 36000 && x > 30000 && x < 36000 {
		return
	}
	if x != GLOBAL_X_CHANGER || y != GLOBAL_Y_CHANGER {
		GLOBAL_X_CHANGER = x
		GLOBAL_Y_CHANGER = y
		if isDebug {
			println("X Value: ", GLOBAL_X_CHANGER, "Y Value: ", GLOBAL_Y_CHANGER)
			print("\n")
		}
		callback(GLOBAL_X_CHANGER, GLOBAL_Y_CHANGER)
	}
}

func normaliseInput(inputValue uint16) float32 {
	return float32(inputValue) / float32(0xffff) // ADC ranges from 0..0xffff
}
