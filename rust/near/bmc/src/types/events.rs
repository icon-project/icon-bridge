use libraries::{
    rlp::{self, Decodable},
    types::{
        messages::{BtpMessage, SerializedMessage},
        BTPAddress,
    }, BytesMut,
};
use std::{ops::Deref, convert::TryFrom};
use near_sdk::base64::{self, URL_SAFE_NO_PAD};

#[derive(Clone, PartialEq, Eq, Debug)]
pub struct Event {
    next: BTPAddress,
    sequence: u128,
    message: BtpMessage<SerializedMessage>,
}

impl Event {
    ///Returns the address
    pub fn next(&self) -> &BTPAddress {
        &self.next
    }
    /// returns the sequence
    pub fn sequence(&self) -> u128 {
        self.sequence
    }
    /// Returns the btp serialised message
    pub fn message(&self) -> &BtpMessage<SerializedMessage> {
        &self.message
    }

    pub fn btp_message(&self) -> Option<BtpMessage<SerializedMessage>> {
        match BtpMessage::try_from(self.message.clone()) { // TODO : OPTIMIZE
            Ok(message) => {
                return Some(message)
            },
            Err(_) => return None,
        }
    }
}

/// An events to do work on it
#[derive(Clone, PartialEq, Eq, Debug)]
pub struct Events(Vec<Event>);

impl Deref for Events {
    /// Event - To take the particular target
    type Target = Vec<Event>;
    /// Returns the target
    fn deref(&self) -> &Self::Target {
        &self.0
    }
}

impl Decodable for Event {
    /// Decodes the result
    fn decode(rlp: &rlp::Rlp) -> Result<Self, rlp::DecoderError> {
        let data = rlp.val_at::<BytesMut>(2).unwrap();
        
        Ok(Self {
            next: rlp.val_at(0)?,
            sequence: rlp.val_at(1)?,
            message: rlp::Rlp::new(&data).as_val()?,
        })
    }
}

impl Decodable for Events {
    fn decode(rlp: &rlp::Rlp) -> Result<Self, rlp::DecoderError> {
        let data = rlp.as_val::<Vec<u8>>()?;
        let rlp = rlp::Rlp::new(&data);
        Ok(Self(rlp.as_list()?))
    }
}
