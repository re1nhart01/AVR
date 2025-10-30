use arduino_hal::{Adc, clock::MHz16, hal::{Atmega, Pins, port::PF4}, pac::adc::ADC, port::{Pin, mode::{Analog, Floating, Input, Output}}};
use oorandom::Rand32;

struct RGBASensor {
    red: Pin<Output>,
    green: Pin<Output>,
    blue: Pin<Output>,
}


impl RGBASensor {
    pub fn new(red: Pin<Output>, green: Pin<Output>, blue: Pin<Output>) -> Self {
        return Self { red, green, blue }
    }

    pub fn toggleDigital(self, mut adc: arduino_hal::adc::Adc, bit1: u8, bit2: u8, bit3: u8) {
        let mut red = self.red.into_output();
        let mut green = self.green.into_output();
        let mut blue = self.blue.into_output();

        if bit1 == 1 { red.set_high(); } else { red.set_low(); }
        if bit2 == 1 { green.set_high(); } else { green.set_low(); }
        if bit3 == 1 { blue.set_high(); } else { blue.set_low(); }
    }

    pub fn toggleAnalog(self, mut adc: arduino_hal::adc::Adc, bit1: u8, bit2: u8, bit3: u8) {
        let mut red = self.red.into_output();
        let mut green = self.green.into_output();
        let mut blue = self.blue.into_output();

        if bit1 == 1 { red.set_high(); } else { red.set_low(); }
        if bit2 == 1 { green.set_high(); } else { green.set_low(); }
        if bit3 == 1 { blue.set_high(); } else { blue.set_low(); }
    }


    pub fn toggleRand(self, mut adc: arduino_hal::adc::Adc) {
        let mut rng = Rand32::new(unsafe {
            core::ptr::read_volatile(0x46 as *const u8) as u64
        });

        let bit1 = (rng.rand_u32() & 1) as u8;
        let bit2 = (rng.rand_u32() & 1) as u8;
        let bit3 = (rng.rand_u32() & 1) as u8;

        self.toggleDigital(adc, bit1, bit2, bit3);
    }

}