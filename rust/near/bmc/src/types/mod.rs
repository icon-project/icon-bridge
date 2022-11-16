mod relay_message;
pub use relay_message::RelayMessage;
mod receipts;
pub(super) use receipts::Receipt;
mod events;
pub(super) use events::{Event, Events};

use btp_common::errors::BmcError;
use libraries::{
    rlp::{self, Decodable},
    types::{
        messages::{BtpMessage, SerializedMessage},
        BTPAddress,
    },
    BytesMut,
};

use near_sdk::{
    base64::{self, URL_SAFE_NO_PAD},
    serde::{de, Deserialize, Serialize},
};

use std::{convert::TryFrom, ops::Deref};
