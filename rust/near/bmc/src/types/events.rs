use libraries::{
    rlp::{self, Decodable},
    types::{
        messages::{BtpMessage, SerializedMessage},
        BTPAddress,
    },
    BytesMut,
};
use near_sdk::base64::{self, URL_SAFE_NO_PAD};
use std::{convert::TryFrom, ops::Deref};

#[derive(Clone, PartialEq, Eq, Debug)]
pub struct Event {
    next: BTPAddress,
    sequence: u128,
    message: BtpMessage<SerializedMessage>,
}

impl Event {
    pub fn next(&self) -> &BTPAddress {
        &self.next
    }

    pub fn sequence(&self) -> u128 {
        self.sequence
    }

    pub fn message(&self) -> &BtpMessage<SerializedMessage> {
        &self.message
    }

    pub fn btp_message(&self) -> Option<BtpMessage<SerializedMessage>> {
        match BtpMessage::try_from(self.message.clone()) {
            // TODO : OPTIMIZE
            Ok(message) => Some(message),
            Err(_) => None,
        }
    }
}
#[derive(Clone, PartialEq, Eq, Debug)]
pub struct Events(Vec<Event>);

impl Deref for Events {
    type Target = Vec<Event>;
    fn deref(&self) -> &Self::Target {
        &self.0
    }
}

impl Decodable for Event {
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
