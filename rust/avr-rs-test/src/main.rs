#![no_std]
#![no_main]
#![feature(abi_avr_interrupt)]


mod water_sensor;
mod rgba_sensor;
mod rc522;
mod dht11;
mod timer;

use core::{cell::{Cell, RefCell}, pin};

use arduino_hal::{Usart, adc::AdcOps, clock::MHz16, delay_ms, delay_us, hal::{Atmega, delay, port::{PB0, PE0, PE1}, spi}, pac::USART0, port::{Pin, mode::{Input, Output}}, prelude::_embedded_hal_blocking_spi_Transfer};
use avr_device::interrupt::Mutex;
use embedded_hal::spi::{Mode, Phase, Polarity};
use panic_halt as _;


use crate::{dht11::DHT11, rc522::Rc522, timer::InterruptTimer};

pub type SerialType = arduino_hal::Usart<
    arduino_hal::pac::USART0,
    Pin<arduino_hal::port::mode::Input, PE0>,
    Pin<arduino_hal::port::mode::Output, PE1>
>;

static LOGGER: Mutex<RefCell<Option<SerialType>>> =
    Mutex::new(RefCell::new(None));


#[arduino_hal::entry]
fn main() -> ! {
    let dp = arduino_hal::Peripherals::take().unwrap();

    let timer0 = InterruptTimer::new();

    timer0.init_timer0(&dp);

    unsafe { avr_device::interrupt::enable() };

    let mut pins = arduino_hal::pins!(dp);
    
    let mut serial = arduino_hal::default_serial!(dp, pins, 57600);

    // ufmt::uwriteln!(&mut serial, "{} {} {}", end - start, start, end).unwrap();

    let dig_pin = pins.d52.into_pull_up_input().downgrade();

    let mut dht11 = DHT11::new(dig_pin, &timer0);

    // let start = dp.TC0.tcnt0.read().bits();
    // let end = dp.TC0.tcnt0.read().bits();
    // let ticks = end.wrapping_sub(start);

    // let mut adc = arduino_hal::adc::Adc::new(dp.ADC, Default::default());

    // let mut serial = arduino_hal::default_serial!(dp, pins, 57600);

    // let sound_dig = pins.a0.into_analog_input(&mut adc);    



    loop {
       if dht11.init(|data| {
            ufmt::uwriteln!(&mut serial, "{}", data).unwrap();
       }).is_ok() {
        match dht11.read_bits(|data, b| {
            ufmt::uwriteln!(&mut serial, "{} {}", data, b).unwrap();
       }) {
            Ok((h1,h2,t1,t2,chk)) => {
                ufmt::uwriteln!(
                    &mut serial,
                    "H={}.% T={}.C chk={}", 
                    h1, t1, chk
                ).unwrap();
            }
            Err(_) => {
                ufmt::uwriteln!(&mut serial, "ERR").unwrap();
            }
        }
    } else {
        ufmt::uwriteln!(&mut serial, "INIT FAIL").unwrap();
    }

    delay_ms(1500); // <- DHT11 не можна читати швидше ніж 1 раз/сек



        // arduino_hal::delay_ms(500);

        // let mut vec: heapless::Vec<u16, 100> = heapless::Vec::new();

        // for i in 0..100 {
        //     let reading= sound_dig.analog_read(&mut adc);
        //     vec.push(reading).unwrap();
        //     arduino_hal::delay_ms(5);
        // }

        // let res: u16 = vec.iter().sum();

        // ufmt::uwriteln!(&mut serial, "{}", res  / 100).unwrap();

        // if false {
        //     ufmt::uwriteln!(&mut serial, "There is sound").unwrap();
        // } else {
        //     ufmt::uwriteln!(&mut serial, "There is no sound").unwrap();
        // }
    }
}


fn led_button() {
    let dp: arduino_hal::Peripherals = arduino_hal::Peripherals::take().unwrap();
    let pins = arduino_hal::pins!(dp);
    let mut led = pins.d13.into_output();

    let button_pin = pins.d52.into_pull_up_input();
  
    loop {
        let state = button_pin.is_high();
        if state {
            led.set_high();
        } else {
            led.set_low();
        }
    }
}


// #[arduino_hal::entry]
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
