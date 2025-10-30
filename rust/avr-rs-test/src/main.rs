#![no_std]
#![no_main]

mod water_sensor;
mod rgba_sensor;
mod rc522;

use core::{cell::RefCell, pin};

use arduino_hal::{hal::{port::PB0, spi}, port::{Pin, mode::Output}, prelude::_embedded_hal_blocking_spi_Transfer};
use embedded_hal::spi::{Mode, Phase, Polarity};
use embedded_hal::digital::OutputPin;
use heapless::String;
use panic_halt as _;

use crate::rc522::Rc522;

trait OrphanClone {
    fn clone(&self) -> Self;
}





// #[arduino_hal::entry]
fn main() -> ! {
    let dp: arduino_hal::Peripherals = arduino_hal::Peripherals::take().unwrap();
    let pins = arduino_hal::pins!(dp);
    // let mut adc = arduino_hal::Adc::new(dp.ADC, Default::default());
    // let mut serial = arduino_hal::default_serial!(dp, pins, 57600);

    // let mut red = pins.d40.into_output();
    // let mut green = pins.d46.into_output();
    // let mut blue = pins.d52.into_output();



    // let mut led = pins.d13.into_output();
    // let analog_pin1 = &pins.a0.into_analog_input(&mut adc);
    // let analog_pin2 = &pins.a1.into_analog_input(&mut adc);
    // let analog_pin3 = &pins.a2.into_analog_input(&mut adc);

    // let analog_pin = &pins.a4.into_analog_input(&mut adc);

    loop {
        // led.toggle();
        // arduino_hal::delay_ms(300);
        
        // // let noise1 = adc.read_blocking(analog_pin1);
        // // let noise2 = adc.read_blocking(analog_pin2);
        // // let noise3 = adc.read_blocking(analog_pin3);

        // let data = adc.read_blocking(analog_pin);



        // let bit0 = (noise1 & 1) as u8;
        // let bit1 = (noise2 & 1) as u8;
        // let bit2 = (noise3 & 1) as u8;

        // let mut rng = Rand32::new(unsafe {
        //     core::ptr::read_volatile(0x46 as *const u8) as u64
        // });
        // let bit0 = (rng.rand_u32() & 1) as u8;
        // let bit1 = (rng.rand_u32() & 1) as u8;
        // let bit2 = (rng.rand_u32() & 1) as u8;

        // ufmt::uwriteln!(&mut serial, "{} {} {}", data, data, data).unwrap();


        // if bit0 == 1 { red.set_high(); } else { red.set_low(); }
        // if bit1 == 1 { green.set_high(); } else { green.set_low(); }
        // if bit2 == 1 { blue.set_high(); } else { blue.set_low(); }

    }
}


#[arduino_hal::entry]
fn spi_test() -> ! {
    let dp = arduino_hal::Peripherals::take().unwrap();
    let pins = arduino_hal::pins!(dp);
    let mut serial = arduino_hal::default_serial!(dp, pins, 57600);

    // --- Настройка SPI ---
    let sclk = pins.d52.into_output();
    let mosi = pins.d51.into_output();
    let miso = pins.d50.into_pull_up_input();
    let mut cs_pin = pins.d53.into_output();
    cs_pin.set_high();

    let (mut spi, mut cs) = spi::Spi::new(
        dp.SPI,
        sclk,
        mosi,
        miso,
        cs_pin,
        spi::Settings {
            data_order: spi::DataOrder::MostSignificantFirst,
            clock: spi::SerialClockRate::OscfOver64, // <= медленно, но надёжно
            mode: Mode {
                polarity: Polarity::IdleLow,
                phase: Phase::CaptureOnFirstTransition,
            },
        },
    );

    ufmt::uwriteln!(&mut serial, "RC522 initialized").unwrap();

    let mut rc522 = Rc522::new(spi, cs);

    rc522.rc522_init();

    let mut atqa = [0u8; 2];

    loop {

        let UUID = rc522.rc_read(&mut atqa, &mut serial);
        ufmt::uwriteln!(&mut serial, "{}", UUID.as_str()).unwrap();
    }
}
