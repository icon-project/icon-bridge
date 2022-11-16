mod bmc_event;
pub use bmc_event::BmcEvent;

use crate::types::messages::{BtpMessage, SerializedMessage};
use crate::types::{BTPAddress, BtpError, Message};
use near_sdk::json_types::U128;
use near_sdk::serde_json::from_str;
use near_sdk::{
    borsh::{self, BorshDeserialize, BorshSerialize},
    collections::LazyOption,
    serde_json::to_value,
};
use std::convert::TryInto;

#[macro_export]
macro_rules! emit_message {
    ($self: ident, $event: ident, $($opt:expr),+) => {
        $self.$event.amend_event($($opt),+)
    }
}

#[macro_export]
macro_rules! emit_error {
    ($self: ident, $event: ident, $($opt:expr),+) => {
        $self.$event.amend_error($($opt),+)
    }
}
