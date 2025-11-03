
use arduino_hal::{hal::port::{PE0, PE1}, port::Pin};


pub type SerialType = arduino_hal::Usart<
    arduino_hal::pac::USART0,
    Pin<arduino_hal::port::mode::Input, PE0>,
    Pin<arduino_hal::port::mode::Output, PE1>
>;

#[macro_export]
macro_rules! log {
    ($($arg:tt)*) => {
        avr_device::interrupt::free(|cs| {
            if let Some(ref mut logger) = crate::LOGGER.borrow(cs).borrow_mut().as_mut() {
                ufmt::uwriteln!(logger, $($arg)*).ok();
            }
        });
    };
}