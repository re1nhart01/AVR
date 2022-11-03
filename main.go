package main

import (
	"machine"
	"strconv"
	"time"
)

var (
	led  = machine.LED
	uart = machine.UART0
	tx   = machine.PD1
	rx   = machine.PD0
)

// 0B 01 22 00 72 AA 89 FF 51 15 - card hex
func main() {
	time.Sleep(5 * time.Second)
	lcd := LCD{intf: 4}
	lcd.Initialize(true)
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()
	uart.Configure(machine.UARTConfig{BaudRate: 9600, TX: tx, RX: rx})
	uart.Write([]byte("Zalupa"))
	led.Low()
	time.Sleep(time.Millisecond * 500)
	led.High()
	for {
		if uart.Buffered() > 0 {
			bytes := make([]byte, 10)
			_, _ = uart.Read(bytes)
			if bytes[4] != 0 && bytes[5] != 0 {
				lcd.ClearScreen()
			}
			for _, v := range bytes {
				hex := strconv.FormatInt(int64(v), 16)
				print(hex, " ", " ")
				lcd.WriteString(hex)

			}
			led.Low()
			time.Sleep(time.Millisecond * 1000)
			led.High()
		}
		// time.Sleep(10 * time.Millisecond)
	}

	// lcd.DisplayOnOff(true)
}

// func rfidSPI() {
// 	spi := machine.SPI0.Configure(machine.SPIConfig{Frequency: 9600})
// }
