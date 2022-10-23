use near_sdk::{Balance, Gas};

pub const SEND_MESSAGE: Gas = Gas(5_000_000_000_000);
pub const HANDLE_EXTERNAL_SERVICE_MESSAGE_CALLBACK: Gas = Gas(5_000_000_000_000);
pub const BSH_HANDLE_BTP_MESSAGE: Gas = Gas(80_000_000_000_000);
pub const GATHER_FEE: Gas = Gas(1_000_000_000_000);
pub const NO_DEPOSIT: Balance = 0;
pub const HANDLE_BTP_ERROR_CALLBACK: Gas = Gas(5_000_000_000_000);
