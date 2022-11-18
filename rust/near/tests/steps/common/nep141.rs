use test_helper::types::{Context, Contract, Nep141, Nep141Contract};

pub fn nep141_contract(name: &'static str) -> Contract<'_, Nep141> {
    Nep141Contract::new(name, "")
}

pub fn register_token_account(name: &'static str, account_id: &str, context: Context) -> Context {
    nep141_contract(name).setup(context, account_id)
}
