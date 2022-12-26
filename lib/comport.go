package main

import (
	"machine"
	"time"
)

func ConfigureLogger() {
	machine.UART0.Configure(machine.UARTConfig{ // init UART
		BaudRate: 9600,
		TX:       machine.PE1,
		RX:       machine.PE0,
	})
	time.Sleep(time.Millisecond * 500)
}
