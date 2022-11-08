use std::ops::Deref;

use crate::types::{Bts, Context, Contract};
use crate::{invoke_call, invoke_view};
use duplicate::duplicate;

#[duplicate(
    contract_type;
    [ Bts ];
)]

impl Contract<'_, contract_type> {
    pub fn register(&self, mut context: Context) -> Context {
        invoke_call!(
            self,
            context,
            "register",
            method_params,
            Some(10_000_000_000_000_000_000_000_000),
            Some(300_000_000_000_000)
        );
        context
    }

    pub fn storage_balance(&self, mut context: Context) -> Context {
        invoke_view!(self, context, "get_storage_balance", method_params);
        context
    }

    pub fn withdraw(&self, mut context: Context, deposit: u128) -> Context {
        invoke_call!(
            self,
            context,
            "withdraw",
            method_params,
            Some(deposit),
            Some(300_000_000_000_000)
        );
        context
    }
}
