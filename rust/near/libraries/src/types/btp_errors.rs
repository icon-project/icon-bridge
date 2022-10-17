use near_sdk::json_types::U128;
use near_sdk::{
    borsh::{self, BorshDeserialize, BorshSerialize},
    serde::{Deserialize, Serialize},
    serde_json::{from_value, to_value, Value},
};

#[derive(Serialize, Deserialize, BorshDeserialize, BorshSerialize)]
#[serde(crate = "near_sdk::serde")]
pub struct BTPError {
    service: String,
    sequence: U128,
    code: U128,
    message: String,
    btp_error_code: U128,
    btp_error_message: String,
}

impl BTPError {
    pub fn new(
        service: String,
        sequence: U128,
        code: U128,
        message: String,
        btp_error_code: U128,
        btp_error_message: String
    ) -> Self {
        Self {
            service,
            sequence,
            code,
            message,
            btp_error_code,
            btp_error_message,
        }
    }

    pub fn service(&self) ->&String{
        &self.service
    }
    pub fn sequence(&self) -> &U128{
        &self.sequence
    }

    pub fn code(&self) -> &U128{
        &self.code
    }
    pub fn message(&self) -> &String{
        &self.message
    }
    pub fn btp_error_code(&self) -> &U128{
        &self.btp_error_code
    }
    pub fn btp_error_message(&self) -> &String{
        &self.btp_error_message
    }
}
