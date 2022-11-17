//! BTP Message Center
mod assertion;
mod estimate;
mod external;
mod internal_service;
mod link_management;
mod messaging;
mod owner_management;
mod relay_management;
mod route_management;
mod service_management;
mod types;

#[cfg(feature = "testable")]
use libraries::types::{BtpError, HashedCollection};

#[cfg(feature = "testable")]
use near_sdk::{
    collections::LazyOption,
    json_types::{Base64VecU8, U128, U64},
};

use btp_common::errors::{BmcError, BshError, BtpException, Exception};
use std::convert::TryInto;

use libraries::{
    emit_error, emit_message,
    types::{
        messages::{
            BmcServiceMessage, BmcServiceType, BtpMessage, ErrorMessage, SerializedBtpMessages,
            SerializedMessage,
        },
        Address, BTPAddress, BmcEvent, Connection, Connections, Link, LinkStatus, Links, Math,
        Owners, RelayStatus, Routes, Service, Services, WrappedI128,
    },
};

use near_sdk::{
    borsh::{self, BorshDeserialize, BorshSerialize},
    env, near_bindgen, require, serde_json,
    serde_json::{to_value, Value},
    AccountId, Balance, Gas, PanicOnDefault, PromiseResult,
};

use crate::{
    external::*,
    types::{Event, RelayMessage},
};

const SERVICE: &str = "bmc";

#[near_bindgen]
#[derive(BorshDeserialize, BorshSerialize, PanicOnDefault)]
pub struct BtpMessageCenter {
    block_interval: u64,
    btp_address: BTPAddress,
    owners: Owners,
    services: Services,
    links: Links,
    routes: Routes,
    connections: Connections,
    event: BmcEvent,
}

#[near_bindgen]
impl BtpMessageCenter {
    #[init]
    pub fn new(network: String, block_interval: u64) -> Self {
        require!(!env::state_exists(), "Already initialized");
        
        let mut owners = Owners::new();
        owners.add(&env::current_account_id());

        Self {
            block_interval,
            btp_address: BTPAddress::new(format!(
                "btp://{}/{}",
                network,
                env::current_account_id()
            )),
            owners,
            services: Services::new(),
            links: Links::new(),
            routes: Routes::new(),
            connections: Connections::new(),
            event: BmcEvent::new(),
        }
    }
}
