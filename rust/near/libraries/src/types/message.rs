use super::*;
#[derive(Debug, Serialize, Deserialize, BorshDeserialize, BorshSerialize)]
#[serde(crate = "near_sdk::serde")]
pub struct Message {
    next: BTPAddress,
    sequence: U128,
    message: Base64VecU8,
}

impl Message {
    pub fn new(next: BTPAddress, sequence: U128, message: Base64VecU8) -> Self {
        Self {
            next,
            sequence,
            message,
        }
    }

    pub fn sequence(&self) -> U128 {
        self.sequence
    }

    pub fn next(&self) -> &BTPAddress {
        &self.next
    }

    pub fn message(&self) -> &Base64VecU8 {
        &self.message
    }
}
