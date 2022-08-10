use libraries::{
    rlp::{self, Decodable},
    types::{
        messages::{BtpMessage, SerializedMessage},
        BTPAddress,
    },
};
use std::{ops::Deref, convert::TryFrom};

#[derive(Clone, PartialEq, Eq, Debug)]
pub struct Proof {
    next: BTPAddress,
    sequence: u128,
    message: BtpMessage<SerializedMessage>,
}

impl Proof {
    pub fn next(&self) -> &BTPAddress {
        &self.next
    }

    pub fn sequence(&self) -> u128 {
        self.sequence
    }

    pub fn message(&self) -> &BtpMessage<SerializedMessage> {
        &self.message
    }
    // pub fn btp_message(&self) -> Option<BtpMessage<SerializedMessage>> {
    //     match BtpMessage::try_from(self.message()) {
    //         Ok(message) => {
    //             return Some(message)
    //         },
    //         Err(_) => return None,
    //     }
    // }
}
#[derive(Clone, PartialEq, Eq, Debug)]
pub struct Proofs(Vec<Proof>);

impl Deref for Proofs {
    type Target = Vec<Proof>;
    fn deref(&self) -> &Self::Target {
        &self.0
    }
}

impl Decodable for Proof {
    fn decode(rlp: &rlp::Rlp) -> Result<Self, rlp::DecoderError> {
        // let data = rlp.as_val::<Vec<u8>>()?;
        // let rlp = rlp::Rlp::new(&data);
        Ok(Self {
            next: rlp.val_at(0)?,
            sequence: rlp.val_at(1)?,
            message: rlp.val_at(2)?,
        })
    }
}

impl Decodable for Proofs {
    fn decode(rlp: &rlp::Rlp) -> Result<Self, rlp::DecoderError> {
        let data = rlp.as_val::<Vec<u8>>()?;
        let rlp = rlp::Rlp::new(&data);
        Ok(Self(rlp.as_list()?))
    }
}
