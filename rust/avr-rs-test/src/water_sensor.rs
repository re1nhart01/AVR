use arduino_hal::{Adc, clock::MHz16, hal::{Atmega, port::PF4}, pac::adc::ADC, port::{Pin, mode::{Analog, Floating, Input}}};

struct WaterSensor {
    pin: Pin<Input<Floating>, PF4>,
}


impl WaterSensor {
    pub fn new(pin: Pin<Input<Floating>, PF4>) -> Self {
        return Self { pin }
    }

    pub fn read(self, mut adc: arduino_hal::adc::Adc) -> u16 {
        let analog_pin_internal = &self.pin.into_analog_input(&mut adc);

        return adc.read_blocking(analog_pin_internal)
    }
}