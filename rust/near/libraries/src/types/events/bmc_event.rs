use super::*;
use crate::types::{BmcEventType, StorageKey};

#[derive(BorshDeserialize, BorshSerialize)]
pub struct BmcEvent {
    message: LazyOption<String>,
    error: LazyOption<String>,
}

impl BmcEvent {
    pub fn new() -> Self {
        Self {
            message: LazyOption::new(StorageKey::BmcEvent(BmcEventType::Message), None),
            error: LazyOption::new(StorageKey::BmcEvent(BmcEventType::Error), None),
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
        serial_no: U128,
        code: u32,
        message: String,
        btp_error_code: u32,
        btp_error_message: String,
    ) {
        self.error.set(
            &to_value(BtpError::new(
                service,
                serial_no,
                code,
                <Vec<u8>>::from(message).into(),
                btp_error_code,
                btp_error_message,
            ))
            .unwrap()
            .to_string(),
        );
    }

    pub fn get_message(&self) -> Result<BtpMessage<SerializedMessage>, String> {
        let message: Message =
            from_str(&self.message.get().ok_or("Not Found")?).map_err(|e| format!("{}", e))?;
        message.message().0.clone().try_into()
    }

    pub fn get_error(&self) -> Result<BtpError, String> {
        from_str(&self.error.get().ok_or("Not Found")?).map_err(|e| format!("{}", e))
    }
}

impl Default for BmcEvent {
    fn default() -> Self {
        Self::new()
    }
}