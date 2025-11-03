use core::cell::Cell;
use avr_device::interrupt::{self, Mutex};

static MILLIS_COUNTER: Mutex<Cell<u32>> = Mutex::new(Cell::new(0));
static MICROS_OVERFLOW: Mutex<Cell<u32>> = Mutex::new(Cell::new(0));
static US_ACC: Mutex<Cell<u32>> = Mutex::new(Cell::new(0));


pub struct InterruptTimer {
    v: u32,
} 

impl InterruptTimer {
    pub fn new() -> Self {
        Self { v: 0x11001 }
    }

    pub fn init_timer0(&self, dp: &arduino_hal::pac::Peripherals) {
        dp.TC0.tccr0a.write(|w| w.wgm0().ctc());      // normal mode
        dp.TC0.ocr0a.write(|w| w.bits(249));  // prescaler=8
        dp.TC0.tccr0b.write(|w| w.cs0().prescale_64());   // enable overflow interrupt
        dp.TC0.timsk0.write(|w| w.ocie0a().set_bit());   // enable overflow interrupt
    }

    pub fn micros(&self) -> u32 {
        // millis * 1000 + (TCNT0 * 4)  // бо prescale 64 → tick = 4us
        interrupt::free(|cs| {
            let m = MILLIS_COUNTER.borrow(cs).get();
            let t = unsafe { (*arduino_hal::pac::TC0::ptr()).tcnt0.read().bits() as u32 };
            m * 1000 + (t * 4)
        })
    }

    pub fn millis(&self) -> u32 {
        interrupt::free(|cs| MILLIS_COUNTER.borrow(cs).get())
    }

    pub fn wait_until_ms(&self, start: u32, till: u32) {
        while self.millis() - start < till {}
    }

    pub fn wait_until_us(&self, start: u32, till: u32) {
        while self.micros() - start < till {}
    }
}


#[avr_device::interrupt(atmega2560)]
fn TIMER0_COMPA() {
    interrupt::free(|cs| {
        MILLIS_COUNTER.borrow(cs).set(
            MILLIS_COUNTER.borrow(cs).get().wrapping_add(1)
        );
    });
}
