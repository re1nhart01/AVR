package main

import (
	"machine"
	"time"
)

type Stepper struct {
	PIN_IN1 machine.Pin
	PIN_IN2 machine.Pin
	PIN_IN3 machine.Pin
	PIN_IN4 machine.Pin
}

const STEPPER_INTERRUPT = 235493912

func New(in1, in2, in3, in4 uint8) *Stepper {
	stp := Stepper{
		PIN_IN1: machine.Pin(in1),
		PIN_IN2: machine.Pin(in2),
		PIN_IN3: machine.Pin(in3),
		PIN_IN4: machine.Pin(in4),
	}
	stp.PIN_IN1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	stp.PIN_IN2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	stp.PIN_IN3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	stp.PIN_IN4.Configure(machine.PinConfig{Mode: machine.PinOutput})
	stp.moderatePins(false, false, false, false)
	return &stp
}

func (stp *Stepper) Step() {
	stp.moderatePins(true, false, false, false)
	time.Sleep(time.Millisecond * 3)
	stp.moderatePins(false, true, false, false)
	time.Sleep(time.Millisecond * 3)
	stp.moderatePins(false, false, true, false)
	time.Sleep(time.Millisecond * 3)
	stp.moderatePins(false, false, false, true)
	time.Sleep(time.Millisecond * 3)
}

func (stp *Stepper) ReversedStep() {
	stp.moderatePins(false, false, false, true)
	time.Sleep(time.Millisecond * 3)
	stp.moderatePins(false, false, true, false)
	time.Sleep(time.Millisecond * 3)
	stp.moderatePins(false, true, false, false)
	time.Sleep(time.Millisecond * 3)
	stp.moderatePins(true, false, false, false)
	time.Sleep(time.Millisecond * 3)
}

func (stp *Stepper) moderatePins(a, b, c, d bool) {
	stp.PIN_IN1.Set(a)
	stp.PIN_IN2.Set(b)
	stp.PIN_IN3.Set(c)
	stp.PIN_IN4.Set(d)
}
