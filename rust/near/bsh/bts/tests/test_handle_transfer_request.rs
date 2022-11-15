#![allow(unused_variables)]
#![allow(unused_imports)]
#![allow(unused_mut)]

use bts::BtpTokenService;
use near_sdk::{
    env, json_types::U128, test_utils::VMContextBuilder, testing_env, AccountId, Gas,
    PromiseResult, RuntimeFeesConfig, VMConfig, VMContext,
};
use std::{collections::HashSet, convert::TryInto};
pub mod accounts;
use accounts::*;
use libraries::types::{
    messages::{BtpMessage, SerializedMessage, TokenServiceMessage, TokenServiceType},
    Account, AccountBalance, Asset, BTPAddress, Math, TransferableAsset, WrappedI128,
    WrappedNativeCoin,
};
mod token;
use std::convert::TryFrom;
use token::*;

pub type Token = Asset<WrappedNativeCoin>;

fn get_context(is_view: bool, signer_account_id: AccountId, attached_deposit: u128) -> VMContext {
    VMContextBuilder::new()
        .current_account_id(alice())
        .is_view(is_view)
        .signer_account_id(signer_account_id.clone())
        .predecessor_account_id(signer_account_id)
        .storage_usage(env::storage_usage())
        .prepaid_gas(Gas(10u64.pow(18)))
        .attached_deposit(attached_deposit)
        .build()
}

#[test]
#[cfg(feature = "testable")]

fn handle_transfer_mint_registered_icx() {
    use std::vec;

    let context = |account_id: AccountId, deposit: u128| get_context(false, account_id, deposit);

    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );

    let nativecoin = Token::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    let icx_coin = <Token>::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let token_id = env::sha256(icx_coin.name().to_owned().as_bytes());
    contract.register_token_callback(icx_coin.clone(), token_id.try_into().unwrap());

    let token_id = contract.token_id(icx_coin.name()).unwrap();

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: chuck().to_string(),
                receiver: destination.account_id().to_string(),
                assets: vec![TransferableAsset::new(icx_coin.name().to_owned(), 900, 99)],
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.on_mint(
        900,
        token_id,
        destination.account_id(),
        Ok(U128::from(700000)),
    );

    let result = contract.balance_of(destination.account_id(), icx_coin.name().to_string());
    assert_eq!(result, U128::from(900));
}

#[test]
fn message_generator() {
    let icx = <Token>::new(ICON_COIN.to_owned());

    let btp_message = BtpMessage::new(
        BTPAddress::new("btp://0x7.icon/cx1ad6fcc465d1b8644ca375f9e10babeea4c38315".to_string()),
        BTPAddress::new(
            "btp://0x2.near/7270a79be789d770f2de015047684e2806597eeee96ee3ca87b179c6399deaaf"
                .to_string(),
        ),
        "bts".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: "cx1ad6fcc465d1b8644ca375f9e10babeea4c38315".to_string(),
                receiver: "chuck".to_string(),
                assets: vec![TransferableAsset::new(
                    "btp-0x7.icon-icx".to_string(),
                    100000000000000,
                    99,
                )],
            },
        )),
    );
    let btp_messages = <BtpMessage<SerializedMessage>>::try_from(&btp_message).unwrap();
    println!("{:?}", String::from(&btp_messages));

    let token_s = <BtpMessage<SerializedMessage>>::try_from(String::from(&btp_messages)).unwrap();

    println!("{:?}", token_s)
}
