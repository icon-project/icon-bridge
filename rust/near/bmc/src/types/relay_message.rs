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
        BmcServiceMessage, BtpMessage, ErrorMessage, SerializedBtpMessages, TokenServiceMessage,
    };
    use near_sdk::serde_json::{self, json};

    use super::*;

    #[test]
    fn deserialize_relay_message1() {
        let message = "-QGU-QGRuQGO-QGLAbkBgvkBf_kBfLhPYnRwOi8vMHgyLm5lYXIvNzI3MGE3OWJlNzg5ZDc3MGYyZGUwMTUwNDc2ODRlMjgwNjU5N2VlZWU5NmVlM2NhODdiMTc5YzYzOTlkZWFhZgO5ASf5ASS4OWJ0cDovLzB4Ny5pY29uL2N4MWFkNmZjYzQ2NWQxYjg2NDRjYTM3NWY5ZTEwYmFiZWVhNGMzODMxNbhPYnRwOi8vMHgyLm5lYXIvNzI3MGE3OWJlNzg5ZDc3MGYyZGUwMTUwNDc2ODRlMjgwNjU5N2VlZWU5NmVlM2NhODdiMTc5YzYzOTlkZWFhZoNidHMBuJH4jwC4jPiKqmh4NTRkOWJhMjIxZmJlOGE0NzVhOGJmMzhjN2QwNDg2NzViNWQ3Yjg1YbhAMGFlNzgzY2I1MmUzMTZiZTdjNWRjZjc2M2Q4NjQ1ZmI0MGQ4ZmUzODZjNWQ4YjliZWQ1MDJkNGY2ZGNhNjk1NdzbkGJ0cC0weDcuaWNvbi1JQ1iJAMHD-J9UD3CjhADKzNQ=";
        let btp_message: BtpMessage<SerializedMessage> =
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
