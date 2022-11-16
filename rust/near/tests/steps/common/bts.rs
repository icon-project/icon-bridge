use lazy_static::lazy_static;
use serde_json::json;
use test_helper::types::{Context, Contract, BtsContract, Bts};

lazy_static! {
    pub static ref BTS_CONTRACT: Contract<'static, Bts> =
        BtsContract::new("bts", "res/BTS_CONTRACT.wasm");
}

pub static BTS_CONTRACT_IS_DEPLOYED: fn(Context) -> Context =
    |context| BTS_CONTRACT.deploy(context);

pub static BTS_CONTRACT_IS_INITIALZIED: fn(Context) -> Context = |mut context: Context| {
    context.add_method_params(
        "new",
        json!({
            "service_name": "bts",
            "bmc": context.contracts().get("bmc").id(),
            "network": super::NEAR_NETWORK,
            "native_coin": {
                "metadata": {
                    "name": "NEAR",
                    "label": "Native NEAR Token",
                    "symbol": "NEAR",
                    "fee_numerator": "100",
                    "fixed_fee": "155000000000000000000000",
                    "network": super::NEAR_NETWORK
                }
            },
        }),
    );

    BTS_CONTRACT.initialize(context)
};

pub static BTS_CONTRACT_IS_DEPLOYED_AND_INITIALIZED: fn(Context) -> Context = |context: Context| {
    context
        .pipe(BTS_CONTRACT_IS_DEPLOYED)
        .pipe(BTS_CONTRACT_IS_INITIALZIED)
};
