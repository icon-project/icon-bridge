use lazy_static::lazy_static;
pub use std::collections::HashSet;
use test_helper::types::{
    Context, Contract, Bts, BtsContract,
};

lazy_static! {
    pub static ref BTS_CONTRACT: Contract<'static, Bts> =
    BtsContract::new("bsh", "res/BTS_CONTRACT.wasm");
}

pub static NEW_CONTEXT: fn() -> Context = || Context::new();

pub static BTS_CONTRACT_IS_DEPLOYED: fn(Context) -> Context =
    |context: Context| BTS_CONTRACT.deploy(context);
