use arduino_hal::{hal::port::Dynamic, port::{Pin, mode::{Floating, Input, PullUp}}};


pub struct Servo<'a> {
    angle: u8,
    tc3: &'a arduino_hal::pac::TC3,
}

impl <'a>Servo<'a> {
    pub fn new(t: &'a arduino_hal::pac::TC3) -> Self {
        return Self { angle: 0, tc3: t };
    }

    pub fn init(&mut self, pin: Pin<Input<Floating>, Dynamic>) {
        pin.into_output();

        let tc3: &arduino_hal::pac::TC3 = self.tc3;

        tc3.icr3.write(|w| unsafe { w.bits(40000) });

        tc3.tccr3a.write(|w| {
            w.wgm3().bits(2)
            .com3a().match_clear()
        });

        tc3.tccr3b.write(|w| {
            w.wgm3().bits(3) 
            .cs3().prescale_8()
        });

    }

    pub fn angle(&mut self, v: u16) {
        let tc3 = self.tc3;

        let ocr_from_angle: u16 = 2000 + (v as u32 * 2000 / 180) as u16;

        tc3.ocr3a.write(|w| unsafe { w.bits(ocr_from_angle) });

        self.angle = v as u8;
    }
}

