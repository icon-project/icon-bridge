use super::Receipt;
use btp_common::errors::BmcError;
use libraries::{
    rlp::{self, Decodable},
    types::messages::SerializedMessage,
    types::messages::TokenServiceMessage,
};
use near_sdk::{
    base64::{self, URL_SAFE_NO_PAD},
    serde::{de, Deserialize, Serialize},
};
use std::convert::TryFrom;
#[derive(Clone, PartialEq, Eq, Debug)]
pub struct RelayMessage {
    receipts: Vec<Receipt>,
}

impl RelayMessage {
    pub fn receipts(&self) -> &Vec<Receipt> {
        &self.receipts
    }
}

impl Serialize for RelayMessage {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: near_sdk::serde::Serializer,
    {
        unimplemented!()
    }
}

impl Decodable for RelayMessage {
    fn decode(rlp: &rlp::Rlp) -> Result<Self, rlp::DecoderError> {
        Ok(Self {
            receipts: rlp.list_at(0)?,
        })
    }
}

impl TryFrom<String> for RelayMessage {
    type Error = BmcError;
    fn try_from(value: String) -> Result<Self, Self::Error> {
        let decoded = base64::decode_config(value, URL_SAFE_NO_PAD).map_err(|error| {
            BmcError::DecodeFailed {
                message: format!("base64: {}", error),
            }
        })?;
        let rlp = rlp::Rlp::new(&decoded);
        Self::decode(&rlp).map_err(|error| BmcError::DecodeFailed {
            message: format!("rlp: {}", error),
        })
    }
}

impl TryFrom<Vec<u8>> for RelayMessage {
    type Error = BmcError;
    fn try_from(value: Vec<u8>) -> Result<Self, Self::Error> {
        let decoded = base64::decode_config(value, URL_SAFE_NO_PAD).map_err(|error| {
            BmcError::DecodeFailed {
                message: format!("base64: {}", error),
            }
        })?;
        let rlp = rlp::Rlp::new(&decoded);
        Self::decode(&rlp).map_err(|error| BmcError::DecodeFailed {
            message: format!("rlp: {}", error),
        })
    }
}

impl<'de> Deserialize<'de> for RelayMessage {
    fn deserialize<D>(deserializer: D) -> Result<Self, <D as de::Deserializer<'de>>::Error>
    where
        D: de::Deserializer<'de>,
    {
        <String as Deserialize>::deserialize(deserializer)
            .and_then(|s| Self::try_from(s).map_err(de::Error::custom))
    }
}

#[cfg(test)]
mod tests {
    use std::convert::TryInto;

    use libraries::types::messages::{
        BmcServiceMessage, BtpMessage, SerializedBtpMessages, TokenServiceMessage,
    };
    use near_sdk::serde_json::{self, json};

    use super::*;

    #[test]
    fn deserialize_relay_message1() {
        let message = "-QEz-QEwuQEt-QEqGrkBIfkBHvkBG7g5YnRwOi8vMHgxLmljb24vY3gyM2E5MWVlM2RkMjkwNDg2YTkxMTNhNmE0MjQyOTgyNWQ4MTNkZTUzLbjd-Nu4OWJ0cDovLzB4MzguYnNjLzB4MDM0QWFERTg2QkY0MDJGMDIzQWExN0U1NzI1ZkFCQzRhYjlFOTc5OLg5YnRwOi8vMHgxLmljb24vY3gyM2E5MWVlM2RkMjkwNDg2YTkxMTNhNmE0MjQyOTgyNWQ4MTNkZTUzg2J0cxW4XvhcALhZ-FeqMHg3QTQzNDFBZjQ5OTU4ODQ1NDZCY2Y3ZTA5ZUI5OGJlRDNlRDI2RDI4qmh4OTM3NTE3YWMwNDJkMGExNGYwOWQ0Njc3ZDMwMmJiMjExMTg0YWM1ZsCEAT1bUQ==";
        let btp_message: BtpMessage<TokenServiceMessage> =
            RelayMessage::try_from(message.to_string())
                .unwrap()
                .receipts[0]
                .events()[0]
                .message()
                .clone()
                .try_into()
                .unwrap();
    }
}
