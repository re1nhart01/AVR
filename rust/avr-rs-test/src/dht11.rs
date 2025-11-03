use core::{error::Error, ptr::read};

use arduino_hal::{delay_ms, delay_us, hal::port::Dynamic, port::{Pin, mode::{Floating, Input, PullUp}}, simple_pwm::Timer0Pwm};
use arduino_hal::hal::delay::Delay;
use heapless::{Vec, vec};

use crate::timer::InterruptTimer;
/**
 * 
    - мікроконтролер опускає DATA в LOW ~18мс щоб дати датчику сигнал що треба передати дані

    - мікроконтролер після цього відпускає лінію в HIGH і чекає

    - датчик відповідає LOW ~80мкс + HIGH ~80мкс

    - після цього датчик починає передавати 40 біт даних

    - кожен біт починається з LOW ~50мкс

    - якщо після цього HIGH короткий (~26–28мкс) → це 0

    - якщо HIGH довгий (~70мкс) → це 1

    - отримані 40 біт розбиваються на 5 байтів (вологість ціла, вологість дробова, температура ціла, температура дробова, контрольна сума)

    - контрольна сума = сума перших 4 байтів і повинна дорівнювати 5-му

    - якщо контрольна сума співпала → дані валідні

    - результати: перший байт — вологість %, третій байт — температура °C
 * 
 */



pub struct DHT11<'a> {
    pin: Option<Pin<Input<PullUp>, Dynamic>>,
    timer: &'a InterruptTimer,
}


impl <'a> DHT11<'a> {
    pub fn new(pin: Pin<Input<PullUp>, Dynamic>, timer: &'a InterruptTimer) -> Self {
        Self { pin: Some(pin), timer }
    }


    pub fn init<F>(&mut self, mut logger: F) -> Result<bool, ()> 
    where
        F: FnMut(u32),
    {
        let mut out_pin = self.pin.take().unwrap().into_output();
        out_pin.set_low();
        delay_ms(18);
        out_pin.set_high();

        let timer = self.timer;
        let in_pin = out_pin.into_pull_up_input();

        while in_pin.is_high() {}

        let start = timer.micros();

        while in_pin.is_low() {}

        let end = timer.micros();

        let dur1 = end - start;

        // logger(dur1);

        if !(dur1 >= 70 && dur1 <= 150) { return Err(()); }

        let start = timer.micros();

        while in_pin.is_high() {}


        let end = timer.micros();

        let dur2 = end - start;

        // logger(dur2);

        if !(dur2 >= 70 && dur2 <= 150) { return Err(()); }

        self.pin = Some(in_pin);

        return Ok(true);
    }


    pub fn read_bits<F>(&mut self, mut logger: F) -> Result<(u8, u8, u8, u8, u8), ()> 
    where
        F: FnMut(u32, &str)
    {
        let pin = self.pin.take().unwrap();

        let mut data: [u8; 5] = [0;5];

        for byte_index in 0..5 {
            let mut byte: u8 = 0;
            for _ in 0..8 {
                let bit = self.read_bit(&pin, &mut logger).unwrap();

                byte = (byte << 1) | bit;
            }

            data[byte_index] = byte;
        }

        let checksum = data[0].wrapping_add(data[1]).wrapping_add(data[2]).wrapping_add(data[3]);

        if checksum != data[4] {
            return Err(());
        }
       
       self.pin = Some(pin);

        return Ok((data[0], data[1], data[2], data[3], data[4]));
    } 


    pub fn read_bit<F>(&mut self, pin: &Pin<Input<PullUp>, Dynamic>, logger: &mut F) -> Result<u8, ()>
    where
        F: FnMut(u32, &str)
    {
        let timer = self.timer;
        while pin.is_low() {}

        let start = timer.micros();
        while pin.is_high() {}
        let dur = timer.micros() - start;

        if dur <= 40 {
            return Ok(0)
        } else if dur >= 40 {
            return Ok(1);
        }

        return Ok(1);
    }

}