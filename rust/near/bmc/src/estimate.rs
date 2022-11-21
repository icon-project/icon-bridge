use super::*;

#[allow(dead_code)]
pub const SEND_MESSAGE: Gas = Gas(5_000_000_000_000);

#[allow(dead_code)]
pub const HANDLE_EXTERNAL_SERVICE_MESSAGE_CALLBACK: Gas = Gas(5_000_000_000_000);

#[allow(dead_code)]
pub const BSH_HANDLE_BTP_MESSAGE: Gas = Gas(80_000_000_000_000);

#[allow(dead_code)]
pub const GATHER_FEE: Gas = Gas(1_000_000_000_000);

#[allow(dead_code)]
pub const NO_DEPOSIT: Balance = 0;

#[allow(dead_code)]
pub const HANDLE_BTP_ERROR_CALLBACK: Gas = Gas(5_000_000_000_000);
