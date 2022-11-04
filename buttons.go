package main

import (
	"machine"
	"time"
)

type PinInterface interface {
	High()
	Low()
	Get() bool
}

const (
	R1_BUTTON_PINOUT = uint8_t(machine.PD3)
	R2_BUTTON_PINOUT = uint8_t(machine.PD4)
	R3_BUTTON_PINOUT = uint8_t(machine.PD5)
	R4_BUTTON_PINOUT = uint8_t(machine.PD6)
	C1_BUTTON_PINOUT = uint8_t(machine.PB0)
	C2_BUTTON_PINOUT = uint8_t(machine.PB1)
	C3_BUTTON_PINOUT = uint8_t(machine.PB2)
	C4_BUTTON_PINOUT = uint8_t(machine.PB3)
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

func (buttons *Buttons) Configure() {
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
	for {
		time.Sleep(time.Second * 3)
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
	if pin != nil {
		time.Sleep(time.Microsecond * 10)
		if flag {
			pin.High()
		} else {
			pin.Low()
		}
	}
}

func checkOnValidPackage(symbol1, symbol2, symbol3, symbol4 string) {
	row1Value := ROW_PIN_1.Get()
	row2Value := ROW_PIN_2.Get()
	row3Value := ROW_PIN_3.Get()
	row4Value := ROW_PIN_4.Get()
	if row1Value == true && row2Value == false && row3Value == false && row4Value == false {
		println(symbol1)
	} else if row1Value == false && row2Value == true && row3Value == false && row4Value == false {
		println(symbol2)
	} else if row1Value == false && row2Value == false && row3Value == true && row4Value == false {
		println(symbol3)
	} else if row1Value == false && row2Value == false && row3Value == false && row4Value == true {
		println(symbol4)
	}
	time.Sleep(time.Microsecond * 100)
}
