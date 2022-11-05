use crate::invoke_call;
use crate::types::{Context, Contract, Bts};
use duplicate::duplicate;

#[duplicate(
    contract_type;
    [ Bts ];
)]

impl Contract<'_, contract_type> {
    pub fn register(&self, mut context: Context) -> Context {
        invoke_call!(self, context, "register", method_params, Some(10_000_000_000_000_000_000_000_000), Some(300_000_000_000_000));
        context
    }
}
