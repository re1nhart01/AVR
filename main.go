package main

import (
	"machine"
	"time"
)

// tinygo flash -target=arduino-mega2560 -port=COM8 .\main.go
// 0B 01 22 00 72 AA 89 FF 51 15 - card hex
func main() {
	// time.Sleep(5 * time.Second)
	// lcd := LCD{intf: 4}
	// lcd.Initialize(true)
	// led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	// led.High()
	// uart.Configure(machine.UARTConfig{BaudRate: 9600, TX: tx, RX: rx})
	// uart.Write([]byte("Zalupa"))
	// led.Low()
	// time.Sleep(time.Millisecond * 500)
	// led.High()
	// for {
	// 	if uart.Buffered() > 0 {
	// 		bytes := make([]byte, 10)
	// 		_, _ = uart.Read(bytes)
	// 		if bytes[4] != 0 && bytes[5] != 0 {
	// 			lcd.ClearScreen()
	// 		}
	// 		for _, v := range bytes {
	// 			hex := strconv.FormatInt(int64(v), 16)
	// 			print(hex, " ", " ")
	// 			lcd.WriteString(hex)

	// 		}
	// 		led.Low()
	// 		time.Sleep(time.Millisecond * 1000)
	// 		led.High()
	// 	}
	// 	// time.Sleep(10 * time.Millisecond)
	// }

	// // lcd.DisplayOnOff(true)
	// buttons := Buttons{}
	// buttons.Configure()
	// buttons.Checker()
	// rfid := &RFID_ME{}
	// rfid.StartSPI()
	// stepper := New(uint8(machine.PB3), uint8(machine.PB2), uint8(machine.PB1), uint8(machine.PB0))
	// for i := 0; i < 9999; i++ {
	// 	if i >= 4000 {
	// 		stepper.Step()
	// 	} else {
	// 		stepper.ReversedStep()
	// 	}
	// }
	// stepper.Step()
	// pin2_micro := machine.PD2
	// pin2_micro.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	// pin3_led := machine.PD3
	// pin3_led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	// pin3_led.Low()
	// for {
	// 	pin3_led.Low()
	// 	if pin2_micro.Get() == true {
	// 		print(true)
	// 		pin3_led.High()
	// 		time.Sleep(time.Second * 1)
	// 		pin3_led.Low()
	// 	}
	// }
	// println("started!")
	// wifiModule := MN35_New(machine.PG5, *machine.UART3, machine.PJ0, machine.PJ1)
	// wifiModule.ConnectToPoint(ACCESS_POINT_NAME, ACCESS_POINT_PASS)
	// pkg := &HttpPackage{
	// 	HttpVersion:   "HTTP/1.1",
	// 	Method:        "POST",
	// 	Path:          "/bober/me?hello=true",
	// 	Host:          "localhost:8080",
	// 	ContentLength: 128,
	// 	Headers: map[string]string{
	// 		"Content-Type": "application/json",
	// 	},
	// }
	// packageStr := wifiModule.generateHTTPackage(pkg)
	// time.Sleep(time.Second * 5)
	// wifiModule.UART_Instance.Buffer.Clear()
	// println("initialized")
	// for {
	// 	time.Sleep(time.Second * 3)
	// 	wifiModule.Try(packageStr, true)
	// }
	joy := NewJoy(machine.PG5, machine.ADC0, machine.ADC1)
	time.Sleep(time.Second * 5)
	println("initialized!")
	joy.Test()
}

type RFID_ME struct {
	RX      machine.Pin
	TX      machine.Pin
	WithLed bool
	uart    *machine.UART
	spi     *machine.SPI
}

const SPI_PinMode = 0

func (rfid *RFID_ME) StartSPI() {
	rfid.spi = &machine.SPI0
	err := rfid.spi.Configure(machine.SPIConfig{Mode: SPI_PinMode, Frequency: 38000, LSBFirst: false})
	if err != nil {
		print(err)
	}
	read_content := make([]byte, 16)
	println("CPU:", machine.CPUFrequency())
	for {
		err := rfid.spi.Tx(nil, read_content)
		if err != nil {
			print(err)
		}
		for _, v := range read_content {
			if v != 0 {
				println(v)
			}
		}
		read_content = make([]byte, 16)
	}
}

// func rfidSPI() {
// 	spi := machine.SPI0.Configure(machine.SPIConfig{Frequency: 9600})
// }
