use crate::rlp::{self, Decodable, Encodable};
use crate::types::{
    messages::BtpMessage, messages::Message, messages::SerializedMessage, TransferableAsset,
};
use btp_common::errors::BshError;
use near_sdk::base64::{self, URL_SAFE_NO_PAD};
use near_sdk::serde::{Deserialize, Serialize};
use near_sdk::AccountId;
use std::convert::TryFrom;

#[derive(Clone, PartialEq, Eq, Debug, Deserialize, Serialize)]
pub enum TokenServiceType {
    RequestTokenTransfer {
        sender: String,
        receiver: String,
        assets: Vec<TransferableAsset>,
    },
    RequestTokenRegister,
    ResponseHandleService {
        code: u8,
        message: String,
    },
    RequestBlacklist {
        request_type: BlackListType,
        addresses: Vec<String>,
        network: String,
    },
    RequestChangeTokenLimit {
        token_names: Vec<String>,
        token_limits: Vec<u128>,
        network: String,
    },
    ResponseBlacklist {
        code: u128,
        message: String,
    },
    ResponseChangeTokenLimit {
        code: u128,
        message: String,
    },
    UnknownType,
    UnhandledType,
}
#[derive(Clone, PartialEq, Eq, Debug, Deserialize, Serialize)]
pub enum BlackListType {
    AddToBlacklist,
    RemoveFromBlacklist,
    UnhandledType,
}

impl Default for TokenServiceType {
    fn default() -> Self {
        Self::UnknownType
    }
}

impl TokenServiceMessage {
    pub fn new(service_type: TokenServiceType) -> Self {
        Self { service_type }
    }

    pub fn service_type(&self) -> &TokenServiceType {
        &self.service_type
    }
}

impl Encodable for BlackListType {
    fn rlp_append(&self, stream: &mut rlp::RlpStream) {
        stream.begin_unbounded_list();
        match self {
            BlackListType::AddToBlacklist => stream.append::<u128>(&0),
            BlackListType::RemoveFromBlacklist => stream.append::<u128>(&1),
            BlackListType::UnhandledType => stream.append::<u128>(&2),
        };
    }
}

impl Encodable for TokenServiceMessage {
    fn rlp_append(&self, stream: &mut rlp::RlpStream) {
        stream.begin_unbounded_list();
        match *self.service_type() {
            TokenServiceType::RequestTokenTransfer {
                ref sender,
                ref receiver,
                ref assets,
            } => {
                let mut params = rlp::RlpStream::new_list(3);
                params
                    .append::<String>(sender)
                    .append::<String>(receiver)
                    .append_list(assets);
                stream.append::<u128>(&0).append(&params.out());
            }
            TokenServiceType::ResponseHandleService {
                ref code,
                ref message,
            } => {
                let mut params = rlp::RlpStream::new_list(2);
                params.append::<u8>(code).append::<String>(message);
                stream.append::<u128>(&2).append(&params.out());
            }
            TokenServiceType::ResponseBlacklist {
                ref code,
                ref message,
            } => {
                let mut params = rlp::RlpStream::new_list(2);
                params.append::<u128>(code).append::<String>(message);
                stream.append::<u128>(&3).append(&params.out());
            }
            TokenServiceType::ResponseChangeTokenLimit {
                ref code,
                ref message,
            } => {
                let mut params = rlp::RlpStream::new_list(2);
                params.append::<u128>(code).append::<String>(message);
                stream.append::<u128>(&4).append(&params.out());
            }
            _ => (),
        }
        stream.finalize_unbounded_list();
    }
}

impl Decodable for TokenServiceMessage {
    fn decode(rlp: &rlp::Rlp) -> Result<Self, rlp::DecoderError> {
        Ok(Self {
            service_type: TokenServiceType::try_from((rlp.val_at::<u128>(0)?, &rlp.val_at(1)?))?,
        })
    }
}

impl Decodable for BlackListType {
    fn decode(rlp: &rlp::Rlp) -> Result<Self, rlp::DecoderError> {
        match rlp.as_val::<u128>()? {
            0 => Ok(Self::AddToBlacklist),
            1 => Ok(Self::RemoveFromBlacklist),
            _ => Ok(Self::UnhandledType),
        }
    }
}

impl TryFrom<(u128, &Vec<u8>)> for TokenServiceType {
    type Error = rlp::DecoderError;
    fn try_from((index, payload): (u128, &Vec<u8>)) -> Result<Self, Self::Error> {
        let payload = rlp::Rlp::new(payload as &[u8]);
        match index {
            0 => Ok(Self::RequestTokenTransfer {
                sender: payload.val_at(0)?,
                receiver: payload.val_at(1)?,
                assets: payload.list_at(2)?,
            }),
            2 => Ok(Self::ResponseHandleService {
                code: payload.val_at(0)?,
                message: payload.val_at(1)?,
            }),
            3 => Ok(Self::RequestBlacklist {
                request_type: payload.val_at::<BlackListType>(0)?,
                addresses: payload.list_at(1)?,
                network: payload.val_at::<String>(2)?,
            }),
            4 => Ok(Self::RequestChangeTokenLimit {
                token_names: payload.list_at(0)?,
                token_limits: payload.list_at(1)?,
                network: payload.val_at::<String>(2)?,
            }),
            5 => Ok(Self::UnknownType),
            _ => Ok(Self::UnhandledType),
        }
    }
}

impl Message for TokenServiceMessage {}

impl From<TokenServiceMessage> for SerializedMessage {
    fn from(message: TokenServiceMessage) -> Self {
        Self::new(rlp::encode(&message).to_vec())
    }
}

impl TryFrom<&Vec<u8>> for TokenServiceMessage {
    type Error = BshError;
    fn try_from(value: &Vec<u8>) -> Result<Self, Self::Error> {
        let rlp = rlp::Rlp::new(value as &[u8]);
        Self::decode(&rlp).map_err(|error| BshError::DecodeFailed {
            message: format!("rlp: {}", error),
        })
    }
}

impl TryFrom<SerializedMessage> for TokenServiceMessage {
    type Error = BshError;
    fn try_from(value: SerializedMessage) -> Result<Self, Self::Error> {
        Self::try_from(value.data())
    }
}

impl From<TokenServiceMessage> for Vec<u8> {
    fn from(service_message: TokenServiceMessage) -> Self {
        rlp::encode(&service_message).to_vec()
    }
}

impl TryFrom<BtpMessage<SerializedMessage>> for BtpMessage<TokenServiceMessage> {
    type Error = BshError;
    fn try_from(value: BtpMessage<SerializedMessage>) -> Result<Self, Self::Error> {
        Ok(Self::new(
            value.source().clone(),
            value.destination().clone(),
            value.service().clone(),
            value.serial_no().clone(),
            value.payload().clone(),
            Some(TokenServiceMessage::try_from(value.payload())?),
        ))
    }
}

impl TryFrom<&BtpMessage<TokenServiceMessage>> for BtpMessage<SerializedMessage> {
    type Error = BshError;
    fn try_from(value: &BtpMessage<TokenServiceMessage>) -> Result<Self, Self::Error> {
        Ok(Self::new(
            value.source().clone(),
            value.destination().clone(),
            value.service().clone(),
            value.serial_no().clone(),
            value
                .message()
                .clone()
                .ok_or(BshError::EncodeFailed {
                    message: "Encoding Failed".to_string(),
                })?
                .into(),
            None,
        ))
    }
}

impl TryFrom<String> for TokenServiceMessage {
    type Error = BshError;
    fn try_from(value: String) -> Result<Self, Self::Error> {
        let decoded = base64::decode_config(value, URL_SAFE_NO_PAD).map_err(|error| {
            BshError::DecodeFailed {
                message: format!("base64: {}", error),
            }
        })?;
        let rlp = rlp::Rlp::new(&decoded);
        Self::decode(&rlp).map_err(|error| BshError::DecodeFailed {
            message: format!("rlp: {}", error),
        })
    }
}

impl From<&TokenServiceMessage> for String {
    fn from(service_message: &TokenServiceMessage) -> Self {
        let rlp = rlp::encode(service_message);
        base64::encode_config(rlp, URL_SAFE_NO_PAD)
    }
}

#[derive(Clone, PartialEq, Eq, Debug, Deserialize, Serialize)]
pub struct TokenServiceMessage {
    service_type: TokenServiceType,
}

#[cfg(test)]
mod tests {
    use btp_common::errors::BshError;

    use super::{
        BlackListType, BtpMessage, SerializedMessage, TokenServiceMessage, TokenServiceType,
        TransferableAsset,
    };
    use std::convert::{TryFrom, TryInto};

    #[test]
    fn deserialize_change_token_limit() {
        let service_message: BtpMessage<SerializedMessage> =  BtpMessage::try_from("-L64OWJ0cDovLzB4Mi5pY29uL2N4NjE5M2U2OTI3NzIzZWNiMzJkYWNiMGExMjVhOTg2NjMzNzY4N2IwM7hPYnRwOi8vMHgxLm5lYXIvN2ZlN2VkMGY4YjI2MTdmYjRlMTA4NWY3YzQzYTM0OWFjZDNmMzMwMGVlYTZiODgxODc2NDZhNDU4ZWNhYzIwY4NidHMKrOsEqejSkWJ0cC0weDEubmVhci1ORUFSy4oCHhngybqyQAAAiDB4MS5uZWFy".to_string()).unwrap();
        let service: BtpMessage<TokenServiceMessage> = service_message.clone().try_into().unwrap();
        let change_token_limit = TokenServiceMessage {
            service_type: TokenServiceType::RequestChangeTokenLimit {
                token_names: vec!["btp-0x1.near-NEAR".to_string()],
                token_limits: vec![10000000000000000000000],
                network: "0x1.near".to_string(),
            },
        };
        let result = service.message();
        assert_eq!(result, &Some(change_token_limit))
    }
    #[test]
    fn deserialize_blacklist_request_message() {
        let service_message:BtpMessage<SerializedMessage> = BtpMessage::try_from(
            "-K-4OWJ0cDovLzB4Mi5pY29uL2N4NjE5M2U2OTI3NzIzZWNiMzJkYWNiMGExMjVhOTg2NjMzNzY4N2IwM7hPYnRwOi8vMHgxLm5lYXIvN2ZlN2VkMGY4YjI2MTdmYjRlMTA4NWY3YzQzYTM0OWFjZDNmMzMwMGVlYTZiODgxODc2NDZhNDU4ZWNhYzIwY4NidHMIndwDmtkAzo1hbGljZS50ZXN0bmV0iDB4MS5uZWFy".to_string(),
        )
        .unwrap();
        let service: BtpMessage<TokenServiceMessage> = service_message.clone().try_into().unwrap();

        let blacklist_request = TokenServiceMessage {
            service_type: TokenServiceType::RequestBlacklist {
                request_type: BlackListType::AddToBlacklist,
                addresses: vec!["alice.testnet".to_string()],
                network: "0x1.near".to_string(),
            },
        };
        let result = service.message();
        assert_eq!(result, &Some(blacklist_request))
    }
    #[test]
    fn deserialize_transfer_request_message() {
        let service_message = TokenServiceMessage::try_from(
            "-FoAuFf4VaoweGMyOTRiMUE2MkU4MmQzZjEzNUE4RjliMmY5Y0FFQUEyM2ZiRDZDZjWKMHgxMjM0NTY3ON7JhElDT06DAwqEyYRUUk9OgwSQwMmEUEFSQYMBhEg".to_string(),
        )
        .unwrap();

        assert_eq!(
            service_message,
            TokenServiceMessage {
                service_type: TokenServiceType::RequestTokenTransfer {
                    sender: "0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5".to_string(),
                    receiver: "0x12345678".to_string(),
                    assets: vec![
                        TransferableAsset::new("ICON".to_string(), 199300, 0),
                        TransferableAsset::new("TRON".to_string(), 299200, 0),
                        TransferableAsset::new("PARA".to_string(), 99400, 0)
                    ]
                },
            },
        );
    }

    #[ignore] // TODO
    #[test]
    fn serialize_transfer_request_message() {
        let service_message = TokenServiceMessage {
            service_type: TokenServiceType::RequestTokenTransfer {
                sender: "0xc294b1A62E82d3f135A8F9b2f9cAEAA23fbD6Cf5".to_string(),
                receiver: "0x12345678".to_string(),
                assets: vec![
                    TransferableAsset::new("ICON".to_string(), 199300, 0),
                    TransferableAsset::new("TRON".to_string(), 299200, 0),
                    TransferableAsset::new("PARA".to_string(), 99400, 0),
                ],
            },
        };

        assert_eq!(
            String::from(&service_message),
            "-FoAuFf4VaoweGMyOTRiMUE2MkU4MmQzZjEzNUE4RjliMmY5Y0FFQUEyM2ZiRDZDZjWKMHgxMjM0NTY3ON7JhElDT06DAwqEyYRUUk9OgwSQwMmEUEFSQYMBhEg".to_string(),
        );
    }

    #[test]
    fn deserialize_transfer_request_message_btp_message() {
        let btp_message = <BtpMessage<SerializedMessage>>::try_from("-QEfuDlidHA6Ly8weDIuaWNvbi9jeDY3ZTIzOGQ4YjFiY2Q4MTVmMzI0ZWQ3ZjU4NWExYWExODhkZDYzZDm4T2J0cDovLzB4MS5uZWFyLzQzODI3ZGZjMDZiZTZiNmQ2ZGEwNDhlYWIyOTdmOWEzNWJhZjIyYjJjZmE1MTFiNDZiMzk1ZTRiODkxNTFmZjODYnRzB7iM-IoAuIf4hapoeGRlYjY5Yjg0YjhjNGY5ZmZhNmU2MWJiNmM3MzMxMTBmMjRlODU0NGO4QDE2YTEyNTVhMDViMTAyZGQzMTY1NDJiODc4MGNmZWRmN2M2M2ZhYjJlODNjYzEyMmIwMWJmYzRiYzQyMDBlMWbX1otidHAtMHgyLUlDWIkAwcP4n1QPcKM".to_string()).unwrap();
    }
}
