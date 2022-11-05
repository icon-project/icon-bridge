use crate::types::{Nep141, Context, Contract};
use crate::{invoke_call, invoke_view};
use duplicate::duplicate;
use near_primitives::types::Gas;

#[duplicate(
    contract_type;
    [ Nep141 ];
)]

impl Contract<'_, contract_type> {
    pub fn get_metadata(&self, mut context: Context) -> Context {
        invoke_view!(self, context, "ft_metadata");
        context
    }
}