// This reads from UART1 and outputs to default serial, usually UART0 or USB.
// Example of how to work with UARTs other than the default.
package main

import (
	"machine"
	"runtime/volatile"
	"strconv"
	"time"
	"unsafe"
)

type RFID_ME struct {
	RX      machine.Pin
	TX      machine.Pin
	WithLed bool
	uart    *machine.UART
	spi     *machine.SPI
}

var (
	led  = machine.LED
	SPDR = (*volatile.Register8)(unsafe.Pointer(uintptr(0x4e)))
	SPSR = (*volatile.Register8)(unsafe.Pointer(uintptr(0x4d)))
	SPCR = (*volatile.Register8)(unsafe.Pointer(uintptr(0x4c)))
)

func (rfid *RFID_ME) StartUART() {
	if rfid.WithLed {
		println("Started!")
		led.Configure(machine.PinConfig{Mode: machine.PinOutput})
		led.Low()
		time.Sleep(time.Microsecond * 100)
		led.High()
	}
	rfid.uart = machine.UART0
	if rfid.uart != nil {
		rfid.uart.Configure(machine.UARTConfig{BaudRate: 9600, TX: rfid.TX, RX: rfid.RX})
		rfid.uart.Write([]byte("Start"))
	}
}

func (rfid *RFID_ME) Listen() {
	for {
		if rfid.uart.Buffered() > 0 {
			bytes := make([]byte, 10)
			_, _ = rfid.uart.Read(bytes)
			for _, v := range bytes {
				hex := strconv.FormatInt(int64(v), 16)
				print(hex, " ", " ")
			}
			led.Low()
			time.Sleep(time.Millisecond * 1000)
			led.High()
		}
		// time.Sleep(10 * time.Millisecond)
	}
}
