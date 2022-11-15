use near_sdk::{Balance, Gas};
/// These constant variables can create gas for resolve transfer, transfer call, mint, on_mint, send service message, no_deposit, burn, token storage deposit and one yocto
/// 
pub const GAS_FOR_RESOLVE_TRANSFER: Gas = Gas(10_000_000_000_000);
pub const GAS_FOR_FT_TRANSFER_CALL: Gas = Gas(25_000_000_000_000);
pub const GAS_FOR_MINT: Gas = Gas(10_000_000_000_000);
pub const GAS_FOR_ON_MINT: Gas = Gas(10_000_000_000_000);
pub const GAS_FOR_SEND_SERVICE_MESSAGE: Gas = Gas(15_000_000_000_000);
pub const NO_DEPOSIT: Balance = 0;
pub const GAS_FOR_BURN: Gas = Gas(1_000_000_000_000);
pub const GAS_FOR_TOKEN_STORAGE_DEPOSIT: Gas = Gas(8_000_000_000_000);
pub const ONE_YOCTO: Balance = 1;
