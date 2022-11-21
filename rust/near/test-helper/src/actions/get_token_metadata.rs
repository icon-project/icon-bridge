use crate::types::{Context, Contract, Nep141};
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

    pub fn ft_transfer(&self, mut context: Context) -> Context {
        invoke_call!(self, context, "ft_transfer", method_params);
        context
    }

    pub fn ft_transfer_call(&self, mut context: Context) -> Context {
        invoke_call!(self, context, "ft_transfer_call", method_params);
        context
    }
    pub fn ft_balance_of(&self, mut context: Context) -> Context {
        invoke_view!(self, context, "ft_balance_of", method_params);
        context
    }
}
