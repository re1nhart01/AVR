package main

import "machine"

type RGBOR struct {
	R machine.Pin
	G machine.Pin
	B machine.Pin
}

func New1(p1, p2, p3 uint8) *RGBOR {
	rgbor := RGBOR{
		R: machine.Pin(p1),
		G: machine.Pin(p2),
		B: machine.Pin(p3),
	}
	rgbor.R.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rgbor.G.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rgbor.B.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return &rgbor
}

func (rgb *RGBOR) configurePins(r, g, b bool) {
	rgb.R.Set(r)
	rgb.G.Set(g)
	rgb.B.Set(b)
}

func (rgb *RGBOR) SetColor(r, g, b bool) {
	rgb.configurePins(r, g, b)
}
