use crate::types::messages::{BtpMessage, SerializedMessage};
use crate::types::{BTPAddress, BTPError, Message};
use near_sdk::json_types::{Base64VecU8, U128};
use near_sdk::serde_json::from_str;
use near_sdk::{
    borsh::{self, BorshDeserialize, BorshSerialize},
    collections::LazyOption,
    serde::{Deserialize, Serialize},
    serde_json::{from_value, to_value, Value},
};
use std::convert::TryInto;

#[derive(BorshDeserialize, BorshSerialize)]
pub struct BmcEvent {
    message: LazyOption<String>,
    error: LazyOption<String>,
}

impl BmcEvent {
    pub fn new() -> Self {
        Self {
            message: LazyOption::new(b"message".to_vec(), None),
            error: LazyOption::new(b"error".to_vec(), None),
        }
    }

    pub fn amend_event(
        &mut self,
        sequence: u128,
        next: BTPAddress,
        message: BtpMessage<SerializedMessage>,
    ) {
        self.message.set(
            &to_value(Message::new(
                next,
                sequence.into(),
                <Vec<u8>>::from(message).into(),
            ))
            .unwrap()
            .to_string(),
        );
    }

    pub fn amend_error(
        &mut self,
        service: String,
        sequence: U128,
        code: U128,
        message: String,
        btp_error_code: U128,
        btp_error_message: String,
    ) {
        self.error.set(&to_value(BTPError::new(
            service,
            sequence,
            code,
            message,
            btp_error_code,
            btp_error_message,
        )).unwrap().to_string()
        );
    }

    pub fn get_message(&self) -> Result<BtpMessage<SerializedMessage>, String> {
        let message: Message =
            from_str(&self.message.get().ok_or("Not Found")?).map_err(|e| format!(""))?;
        message
            .message()
            .0
            .clone()
            .try_into()
            .map_err(|e| format!("{}", e))
    }

    pub fn get_error(&self) {}
}
