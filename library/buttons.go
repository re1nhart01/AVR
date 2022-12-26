package main

import (
	"machine"
	"time"
)

type PinInterface interface {
	High()
	Low()
	Get() bool
	Set(bool)
}

const (
	R4_BUTTON_PINOUT = uint8_t(machine.PD6)
	R3_BUTTON_PINOUT = uint8_t(machine.PD7)
	R2_BUTTON_PINOUT = uint8_t(machine.PB0)
	R1_BUTTON_PINOUT = uint8_t(machine.PB1)
	C4_BUTTON_PINOUT = uint8_t(machine.PB2)
	C3_BUTTON_PINOUT = uint8_t(machine.PB3)
	C2_BUTTON_PINOUT = uint8_t(machine.PB4)
	C1_BUTTON_PINOUT = uint8_t(machine.PB5)
)

const (
	ROW_PIN_1    = machine.Pin(R1_BUTTON_PINOUT)
	ROW_PIN_2    = machine.Pin(R2_BUTTON_PINOUT)
	ROW_PIN_3    = machine.Pin(R3_BUTTON_PINOUT)
	ROW_PIN_4    = machine.Pin(R4_BUTTON_PINOUT)
	COLUMN_PIN_1 = machine.Pin(C1_BUTTON_PINOUT)
	COLUMN_PIN_2 = machine.Pin(C2_BUTTON_PINOUT)
	COLUMN_PIN_3 = machine.Pin(C3_BUTTON_PINOUT)
	COLUMN_PIN_4 = machine.Pin(C4_BUTTON_PINOUT)
)

type Buttons struct {
}

var currentValue = ""

func (buttons *Buttons) Configure() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ROW_PIN_1.Configure(machine.PinConfig{Mode: machine.PinInput})
	ROW_PIN_2.Configure(machine.PinConfig{Mode: machine.PinInput})
	ROW_PIN_3.Configure(machine.PinConfig{Mode: machine.PinInput})
	ROW_PIN_4.Configure(machine.PinConfig{Mode: machine.PinInput})

	COLUMN_PIN_1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	COLUMN_PIN_2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	COLUMN_PIN_3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	COLUMN_PIN_4.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func (buttons *Buttons) Checker() {
	setFlag(COLUMN_PIN_1, false)
	setFlag(COLUMN_PIN_2, false)
	setFlag(COLUMN_PIN_3, false)
	setFlag(COLUMN_PIN_4, false)
	for {
		setFlag(COLUMN_PIN_1, true)
		setFlag(COLUMN_PIN_2, false)
		setFlag(COLUMN_PIN_3, false)
		setFlag(COLUMN_PIN_4, false)
		checkOnValidPackage("1", "4", "7", "*")
		setFlag(COLUMN_PIN_1, false)
		setFlag(COLUMN_PIN_2, true)
		setFlag(COLUMN_PIN_3, false)
		setFlag(COLUMN_PIN_4, false)
		checkOnValidPackage("2", "5", "8", "0")
		setFlag(COLUMN_PIN_1, false)
		setFlag(COLUMN_PIN_2, false)
		setFlag(COLUMN_PIN_3, true)
		setFlag(COLUMN_PIN_4, false)
		checkOnValidPackage("3", "6", "9", "#")
		setFlag(COLUMN_PIN_1, false)
		setFlag(COLUMN_PIN_2, false)
		setFlag(COLUMN_PIN_3, false)
		setFlag(COLUMN_PIN_4, true)
		checkOnValidPackage("G", "A", "Y", "BAN")
	}
}

func setFlag(pin PinInterface, flag bool) {
	pin.Set(flag)
}

func led2() {
	led.Low()
	time.Sleep(time.Millisecond * 420)
	led.High()
	time.Sleep(time.Millisecond * 420)
}

func checkOnValidPackage(symbol1, symbol2, symbol3, symbol4 string) {
	if ROW_PIN_1.Get() == true && ROW_PIN_2.Get() == false && ROW_PIN_3.Get() == false && ROW_PIN_4.Get() == false {
		if currentValue == symbol1 {
			println(symbol1)
		}
		currentValue = symbol1
	} else if ROW_PIN_1.Get() == false && ROW_PIN_2.Get() == true && ROW_PIN_3.Get() == false && ROW_PIN_4.Get() == false {
		if currentValue == symbol2 {
			println(symbol2)
		}
		currentValue = symbol2
	} else if ROW_PIN_1.Get() == false && ROW_PIN_2.Get() == false && ROW_PIN_3.Get() == true && ROW_PIN_4.Get() == false {
		if currentValue == symbol3 {
			println(symbol3)
		}
		currentValue = symbol3
	} else if ROW_PIN_1.Get() == false && ROW_PIN_2.Get() == false && ROW_PIN_3.Get() == false && ROW_PIN_4.Get() == true {
		if currentValue == symbol4 {
			println(symbol4)
		}
		currentValue = symbol4
	}
}
