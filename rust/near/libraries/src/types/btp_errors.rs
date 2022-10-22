use near_sdk::json_types::{Base64VecU8, U128};
use near_sdk::{
    borsh::{self, BorshDeserialize, BorshSerialize},
    serde::{Deserialize, Serialize},
    serde_json::{from_value, to_value, Value},
};
#[derive(Serialize, Deserialize, BorshDeserialize, BorshSerialize, Debug)]
#[serde(crate = "near_sdk::serde")]
pub struct BtpError {
    service: String,
    serial_no: U128,
    code: u32,
    message: Base64VecU8,
    btp_error_code: u32,
    btp_error_message: String,
}

impl BtpError {
    pub fn new(
        service: String,
        serial_no: U128,
        code: u32,
        message: Base64VecU8,
        btp_error_code: u32,
        btp_error_message: String,
    ) -> Self {
        Self {
            service,
            serial_no,
            code,
            message,
            btp_error_code,
            btp_error_message,
        }
    }

    pub fn service(&self) -> &String {
        &self.service
    }
    pub fn serial_no(&self) -> &U128 {
        &self.serial_no
    }

    pub fn code(&self) -> u32 {
        self.code
    }
    pub fn message(&self) -> &Base64VecU8 {
        &self.message
    }
    pub fn btp_error_code(&self) -> u32 {
        self.btp_error_code
    }
    pub fn btp_error_message(&self) -> &String {
        &self.btp_error_message
    }
}
