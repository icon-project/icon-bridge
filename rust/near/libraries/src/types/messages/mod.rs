#![allow(unused)]
mod token_service;
pub use token_service::*;
pub mod bmc_service;
pub use bmc_service::*;
mod btp_message;
pub use btp_message::*;
mod error_message;
pub use error_message::*;
pub trait Message {}

use crate::{
    rlp::{self, Decodable, Encodable},
    types::{BTPAddress, Nullable, TransferableAsset, WrappedI128},
};

use btp_common::errors::{BmcError, BshError};

use near_sdk::{
    base64::{self, URL_SAFE_NO_PAD}, // TODO: Confirm
    borsh::{self, maybestd::io, BorshDeserialize, BorshSerialize},
    serde::{de, ser, Deserialize, Deserializer, Serialize, Serializer},
};
use std::{
    convert::{TryFrom, TryInto},
    vec::IntoIter,
};
