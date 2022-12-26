package main

import (
	"machine"
	"time"
)

type uint8_t uint8
type char_t = string

const intf = 8

var emptylcd = LCD{intf: 8}

var CHARS_MAPPING_1 = " !\"#$%&'()*+,-./"
var CHARS_MAPPING_2 = "0123456789:;<=>?"
var CHARS_MAPPING_3 = "@ABCDEFGHIJKLMNO"
var CHARS_MAPPING_4 = "PQRSTUVWXYZ[!]^_"
var CHARS_MAPPING_5 = "`abcdefghijklmno"
var CHARS_MAPPING_6 = "pqrstuvwxyz{|}~"

/*
* RS   PD7
* R/W  GND
* DB7  PD2
* DB6  PD3
* DB5  PD4
* DB4  PD5
* DB3  -- not using
* DB2  -- not using
* DB1  -- not using
* DB0  -- not using
 */

const (
	PD2_LCD_PIN_17 = uint8_t(machine.PD2)
	PD3_LCD_PIN_18 = uint8_t(machine.PD3)
	PD4_LCD_PIN_19 = uint8_t(machine.PD4)
	PD5_LCD_PIN_20 = uint8_t(machine.PD5)
	PD6_LCD_PIN_21 = uint8_t(machine.PD6)
	PD7_LCD_PIN_22 = uint8_t(machine.PD7)
)

const (
	RS_PIN     = machine.Pin(PD7_LCD_PIN_22)
	ENABLE_PIN = machine.Pin(PD6_LCD_PIN_21)
	DB7_PIN    = machine.Pin(PD2_LCD_PIN_17)
	DB6_PIN    = machine.Pin(PD3_LCD_PIN_18)
	DB5_PIN    = machine.Pin(PD4_LCD_PIN_19)
	DB4_PIN    = machine.Pin(PD5_LCD_PIN_20)
)

type LCD struct {
	RS   uint8_t
	intf uint8_t
}

func configurePins() {
	RS_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ENABLE_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	DB7_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	DB6_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	DB5_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	DB4_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func (lcd *LCD) Initialize(withConfigure bool) {

	if withConfigure {
		configurePins()
	}

	time.Sleep(time.Millisecond * 100) // Wait for more than 40 ms after VDD rises to 4.5 V

	setPinsMode(false, false, false, true, true)
	setPinsMode(false, false, false, true, true)
	setPinsMode(false, false, false, true, true)
	setPinsMode(false, false, false, true, false)

	setPinsMode(false, false, false, true, false)
	setPinsMode(false, true, false, false, false)

	setPinsMode(false, false, false, false, false)
	setPinsMode(false, true, true, false, false)

	setPinsMode(false, false, false, false, false)
	setPinsMode(false, false, false, false, true)
	time.Sleep(time.Millisecond * 10)
	setPinsMode(false, false, false, false, false)
	setPinsMode(false, false, true, true, false)

	print("LCD WINSTAR 1602A initialized!")
}

func getCharBits(decimal int) [4]bool {
	binary := [4]bool{false, false, false, false}
	counter := 3
	for decimal != 0 {
		delim := decimal % 2
		if delim == 0 {
			binary[counter] = false
		} else {
			binary[counter] = true
		}
		decimal = decimal / 2
		counter--
	}
	return binary
}

func getCharByMapping(charCode int) [2][4]bool {
	if charCode < 32 && charCode > 128 {
		return [2][4]bool{}
	}
	currentGroupVertical := [4]bool{false, false, false, false}
	currentGroupHorizontal := [4]bool{false, false, true, false}
	if charCode >= 32 && charCode <= 47 {
		currentGroupHorizontal = [4]bool{false, false, true, false}
		currentGroupVertical = getCharBits(charCode - 32)
	} else if charCode >= 48 && charCode <= 63 {
		currentGroupHorizontal = [4]bool{false, false, true, true}
		currentGroupVertical = getCharBits(charCode - 48)
	} else if charCode >= 64 && charCode <= 79 {
		currentGroupHorizontal = [4]bool{false, true, false, false}
		currentGroupVertical = getCharBits(charCode - 64)
	} else if charCode >= 80 && charCode <= 95 {
		currentGroupHorizontal = [4]bool{false, true, false, true}
		currentGroupVertical = getCharBits(charCode - 80)
	} else if charCode >= 96 && charCode <= 111 {
		currentGroupHorizontal = [4]bool{false, true, true, false}
		currentGroupVertical = getCharBits(charCode - 96)
	} else if charCode >= 112 && charCode <= 126 {
		currentGroupHorizontal = [4]bool{false, true, true, true}
		currentGroupVertical = getCharBits(charCode - 112)
	}
	return [2][4]bool{currentGroupVertical, currentGroupHorizontal}
}

func (lcd *LCD) ClearScreen() {
	time.Sleep(time.Millisecond * 15)
	setPinsMode(false, false, false, false, false)
	setPinsMode(false, false, false, false, true)
}

func (lcd *LCD) DisplayOnOff(flag bool) {
	if flag {
		time.Sleep(time.Millisecond * 100)
		setPinsMode(false, true, false, true, true)
	}
}

func (lcd *LCD) WriteString(str string) {
	for _, v := range str {
		responder := getCharByMapping(int(v))
		if len(responder) > 1 {
			setPinsModeArray(true, responder[1])
			setPinsModeArray(true, responder[0])
		}
	}
}

func setPinsMode(rs, db7, db6, db5, db4 bool) {
	time.Sleep(time.Microsecond * 120)
	ENABLE_PIN.High()
	time.Sleep(time.Microsecond * 40)
	RS_PIN.Set(rs)
	DB7_PIN.Set(db7)
	DB6_PIN.Set(db6)
	DB5_PIN.Set(db5)
	DB4_PIN.Set(db4)
	ENABLE_PIN.Low()
	// args for 8 bit interface db3, db2, db1, db0
	// DB3_PIN.Set(rs) --not in use
	// DB2_PIN.Set(rs) --not in use
	// DB1_PIN.Set(rs) --not in use
	// DB0_PIN.Set(rs) --not in use

}

func setPinsModeArray(rs bool, pins [4]bool) {
	if len(pins) == 4 {
		time.Sleep(time.Microsecond * 120)
		ENABLE_PIN.High()
		time.Sleep(time.Microsecond * 40)
		RS_PIN.Set(rs)
		DB7_PIN.Set(pins[0])
		DB6_PIN.Set(pins[1])
		DB5_PIN.Set(pins[2])
		DB4_PIN.Set(pins[3])
		ENABLE_PIN.Low()
	}
}
