use arduino_hal::{clock::MHz16, hal::{Usart, port::{PE0, PE1}}, pac::USART0, port::{Pin, mode::{Input, Output}}};
use embedded_hal::{digital::OutputPin, spi::SpiBus};
use heapless::{String};


// Команды контроллера MFRC522 (PCD)
const PCD_IDLE: u8 = 0x00;
const PCD_MEM: u8 = 0x01;
const PCD_GENERATE_RANDOM_ID: u8 = 0x02;
const PCD_TRANSCIEVE: u8 = 0x0C;
const PCD_SOFTRESET: u8 = 0x0F;

// Регистры (не все, только нужные)
const COMMAND_REG: u8 = 0x01;
const FIFO_DATA_REG: u8 = 0x09;
const ERROR_REG: u8 = 0x06;
const DIV_IRQ_REG: u8 = 0x05;
const VERSION_REG: u8 = 0x37;

const COMM_IE_N_REG: u8 = 0x02;
const DIVLEN_REG: u8 = 0x03;
const COMM_IRQ_REG: u8 = 0x04;
const FIFO_LEVEL_REG: u8 = 0x0A;
const CONTROL_REG: u8 = 0x0C;
const BIT_FRAMING_REG: u8 = 0x0D;
const COLL_REG: u8 = 0x0E;
const MODE_REG: u8 = 0x11;
const TX_MODE_REG: u8 = 0x12;
const RX_MODE_REG: u8 = 0x13;
const TX_CONTROL_REG: u8 = 0x14;
const RFCFG_REG: u8 = 0x26;
const TMODE_REG: u8 = 0x2A;
const TPRESCALER_REG: u8 = 0x2B;
const TRELOAD_REG_L: u8 = 0x2C;
const TRELOAD_REG_H: u8 = 0x2D;
const CRC_RESULT_REG_L: u8 = 0x22;
const CRC_RESULT_REG_M: u8 = 0x21;

// Команды:
const PCD_TRANSCEIVE: u8 = 0x0C;
const PCD_CALCCRC: u8 = 0x03;

// Карточные команды:
const PICC_REQA: u8 = 0x26;
const PICC_SEL_CL1: u8 = 0x93;
const PICC_ANTICOLL_CL1: u8 = 0x20;
const PICC_SEL_CMD: u8 = 0x70;



pub struct Rc522<SPI, CS>
where
    SPI: SpiBus<u8, Error = core::convert::Infallible>,
    CS: OutputPin<Error = core::convert::Infallible>,
{
    spi: SPI,
    cs: CS,
}


impl<SPI, CS> Rc522<SPI, CS>
where
    SPI: SpiBus<u8, Error = core::convert::Infallible>,
    CS: OutputPin<Error = core::convert::Infallible>,
{
    pub fn new(spi: SPI, cs: CS) -> Self {
        Self { spi, cs }
    }

    pub fn get_version(&mut self) -> Result<u8, ()> {
        let version = self.rc522_read_reg(VERSION_REG);
        return Ok(version)
    }

    fn rc522_read_reg(&mut self, reg: u8) -> u8 {
        let mut buf = [(reg << 1) | 0x80, 0];
        self.cs.set_low().unwrap();
        self.spi.transfer_in_place(&mut buf).unwrap();
        self.cs.set_high().unwrap();
        buf[1]
    }

    fn rc522_write_reg(&mut self, reg: u8, val: u8) {
        let mut buf = [(reg << 1) & 0x7E, val];
        self.cs.set_low().unwrap();
        self.spi.transfer_in_place(&mut buf).unwrap();
        self.cs.set_high().unwrap();
    }

    pub fn rc522_init(&mut self) {
        self.rc522_soft_reset();
        arduino_hal::delay_ms(50);

        self.rc522_write_reg(TMODE_REG,       0x8D);
        self.rc522_write_reg(TPRESCALER_REG,  0x3E);
        self.rc522_write_reg(TRELOAD_REG_L,   30);
        self.rc522_write_reg(TRELOAD_REG_H,   0);

        self.rc522_write_reg(TX_MODE_REG,     0x00);
        self.rc522_write_reg(RX_MODE_REG,     0x00);
        self.rc522_write_reg(MODE_REG,        0x3D);
        self.rc522_write_reg(0x15,        0x40);

        self.rc522_write_reg(RFCFG_REG,       0x7F);
        self.rc522_write_reg(TX_CONTROL_REG,  0x83); 
    }

    fn rc522_soft_reset(&mut self) {
        self.rc522_write_reg(COMMAND_REG, PCD_SOFTRESET);
        arduino_hal::delay_ms(50);
    }

    pub fn picc_reqa(&mut self, atqa: &mut [u8;2]) -> bool {
        match self.rc522_transceive(&[PICC_REQA], atqa, 0x07) {
            Ok(n) if n == 2 => true,
            _ => false,
        }
    }

    fn rc522_transceive(
        &mut self,
        send: &[u8],
        rx_buf: &mut [u8],
        bit_framing: u8,
    ) -> Result<usize, u8> {
        self.rc522_write_reg(COMM_IRQ_REG, 0x7F);      
        self.rc522_write_reg(FIFO_LEVEL_REG, 0x80);   
        self.rc522_write_reg(BIT_FRAMING_REG, bit_framing);

        for &b in send {
            self.rc522_write_reg(FIFO_DATA_REG, b);
        }

        self.rc522_write_reg(COMMAND_REG, PCD_TRANSCEIVE);
        self.rc522_set_bit_mask(BIT_FRAMING_REG, 0x80);

        let mut i = 2000;
        while i > 0 {
            let irq = self.rc522_read_reg(COMM_IRQ_REG);
            if irq & 0x30 != 0 { break; }
            if self.rc522_read_reg(ERROR_REG) & 0x1B != 0 {
                return Err(0xE1);
            }
            i -= 1;
        }
        self.rc522_clear_bit_mask(BIT_FRAMING_REG, 0x80);

        let n = (self.rc522_read_reg(FIFO_LEVEL_REG) as usize).min(rx_buf.len());
        for i in 0..n {
            rx_buf[i] = self.rc522_read_reg(FIFO_DATA_REG);
        }
        Ok(n)
    }

    fn rc522_set_bit_mask(&mut self, reg: u8, mask: u8) {
        let tmp = self.rc522_read_reg(reg);
        self.rc522_write_reg(reg, tmp | mask);
    }

    fn rc522_clear_bit_mask(&mut self, reg: u8, mask: u8) {
        let tmp = self.rc522_read_reg(reg);
        self.rc522_write_reg(reg, tmp & !mask);
    }

    fn picc_anticoll(&mut self) -> Result<[u8; 5], ()> {
        self.rc522_write_reg(0x0E, 0x80);
        self.rc522_write_reg(0x0D, 0x00);

        let cmd = [0x93, 0x20];
        let mut back = [0u8; 10];

        let n = self.rc522_transceive(&cmd, &mut back, 0x00).map_err(|_| ())?;
        if n < 5 {
            return Err(());
        }

        let uid0 = back[0];
        let uid1 = back[1];
        let uid2 = back[2];
        let uid3 = back[3];
        let bcc  = back[4];

        if (uid0 ^ uid1 ^ uid2 ^ uid3) != bcc {
            return Err(());
        }

        Ok([uid0, uid1, uid2, uid3, bcc])
    }

    fn picc_select(&mut self, uid_cl1: [u8;5]) -> Result<u8, ()> {
        let mut frame = [0u8; 9];
        frame[0] = PICC_SEL_CL1;
        frame[1] = PICC_SEL_CMD;
        frame[2..7].copy_from_slice(&uid_cl1);

        let (crc_l, crc_m) = self.rc522_calc_crc(&frame[0..7]);
        frame[7] = crc_l;
        frame[8] = crc_m;

        let mut back = [0u8; 3]; // SAK (+CRC)
        let n = self.rc522_transceive(&frame, &mut back, 0x00).map_err(|_| ())?;
        if n < 1 { return Err(()); }

        let sak = back[0];
        Ok(sak)
    }

    fn rc522_calc_crc(&mut self, data: &[u8]) -> (u8, u8) {
        self.rc522_write_reg(COMMAND_REG, PCD_IDLE);
        self.rc522_write_reg(DIV_IRQ_REG, 0x04);    // Clear CRCIRq
        self.rc522_write_reg(FIFO_LEVEL_REG, 0x80); // Flush FIFO

        for &b in data {
            self.rc522_write_reg(FIFO_DATA_REG, b);
        }
        self.rc522_write_reg(COMMAND_REG, PCD_CALCCRC);

        let mut i = 5000;
        while i > 0 {
            if self.rc522_read_reg(DIV_IRQ_REG) & 0x04 != 0 { break; }
            i -= 1;
        }

        let l = self.rc522_read_reg(CRC_RESULT_REG_L);
        let m = self.rc522_read_reg(CRC_RESULT_REG_M);
        (l, m)
    }

    pub fn rc_read(
    &mut self,
    atqa: &mut [u8; 2],
    serial: &mut Usart<USART0, Pin<Input, PE0>, Pin<Output, PE1>, MHz16>,
    ) -> String<8> {
        if self.picc_reqa(atqa) {
            ufmt::uwriteln!(serial, "Card detected! ATQA: {:02X} {:02X}", atqa[0], atqa[1]).unwrap();

            match self.picc_anticoll() {
                Ok(uid5) => {
                    let uuid_string = heapless::format!(8; "{:02X}{:02X}{:02X}{:02X}", uid5[0], uid5[1], uid5[2], uid5[3]).unwrap();
                    
                    ufmt::uwriteln!(serial,"UID: {}", uuid_string.as_str()).unwrap();

                    if let Ok(sak) = self.picc_select(uid5) {
                        ufmt::uwriteln!(serial, "SAK: 0x{:02X}", sak).unwrap();
                    } else {
                        ufmt::uwriteln!(serial, "Select failed").unwrap();
                    }

                    return uuid_string;
                }
                Err(_) => {
                    ufmt::uwriteln!(serial, "Anticollision failed").unwrap();
                }
            }
        } else {
            ufmt::uwriteln!(serial, "No card").unwrap();
        }

        arduino_hal::delay_ms(500);

        return String::new()
    }
}