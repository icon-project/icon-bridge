use crate::invoke_call;
use crate::types::{Bmc, Context, Contract, Bts};
use duplicate::duplicate;

#[duplicate(
    contract_type;
    [ Bmc ];
    [ Bts ];
)]
impl Contract<'_, contract_type> {
    pub fn initialize(&self, mut context: Context) -> Context {
        invoke_call!(self, context, "new", method_params);
        context
    }
}
