#![allow(unused_variables)]
#![allow(unused_imports)]
#![allow(unused_mut)]

use bts::{BtpTokenService, Token};
use near_sdk::{
    env, json_types::U128, test_utils::VMContextBuilder, testing_env, AccountId, Gas,
    PromiseResult, RuntimeFeesConfig, VMConfig, VMContext,
};
pub mod accounts;
use accounts::*;
use libraries::types::{
    messages::{BtpMessage, TokenServiceMessage, TokenServiceType},
    Account, AccountBalance, AccumulatedAssetFees, Asset, AssetItem, BTPAddress, Math,
    WrappedFungibleToken, WrappedI128,
};
mod token;
use libraries::types::{Request, TransferableAsset};
use std::convert::TryInto;
use token::*;

pub type TokenFee = AssetItem;

fn get_context(
    is_view: bool,
    signer_account_id: AccountId,
    attached_deposit: u128,
    account_balance: u128,
    prepaid_gas: u64,
) -> VMContext {
    VMContextBuilder::new()
        .current_account_id(alice())
        .is_view(is_view)
        .signer_account_id(signer_account_id.clone())
        .predecessor_account_id(signer_account_id)
        .storage_usage(env::storage_usage())
        .prepaid_gas(Gas(prepaid_gas))
        .attached_deposit(attached_deposit)
        .build()
}

#[test]
fn deposit_wnear() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };

    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let w_near = <Token>::new(WNEAR.to_owned());
    contract.register(w_near.clone());

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(100), "".to_string());

    testing_env!(context(chuck(), 0));
    let result = contract.balance_of(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));
}

#[test]
fn withdraw_wnear() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 1000, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let w_near = <Token>::new(WNEAR.to_owned());
    contract.register(w_near.clone());
    let token_id = contract.token_id(w_near.name()).unwrap();

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    testing_env!(context(chuck(), 0));
    let result = contract.balance_of(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1000).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));

    let storage_cost = contract
        .get_storage_balance(chuck(), w_near.name().to_string())
        .0
        + 1;
    testing_env!(context(chuck(), storage_cost));
    contract.withdraw(w_near.name().to_string(), U128::from(999));

    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.on_withdraw(chuck(), 999, w_near.name().to_string(), token_id.clone());

    let result = contract.balance_of(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn withdraw_wnear_higher_amount() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 1000, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let w_near = <Token>::new(WNEAR.to_owned());
    contract.register(w_near.clone());

    testing_env!(context(wnear(), 0));
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    testing_env!(context(chuck(), 0));
    let result = contract.balance_of(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1000).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));

    testing_env!(context(chuck(), 1));
    contract.withdraw(w_near.name().to_string(), U128::from(1001));
}

#[test]
#[cfg(feature = "testable")]
fn external_transfer() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let w_near = <Token>::new(WNEAR.to_owned());
    contract.register(w_near.clone());

    let token_id = contract.token_id(w_near.name()).unwrap();

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    let storage_cost = contract
        .get_storage_balance(chuck(), w_near.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.transfer(
        w_near.name().to_string(),
        destination.clone(),
        U128::from(999),
    );

    let result = contract.account_balance(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();

    expected.deposit_mut().add(1).unwrap();
    expected.locked_mut().add(900).unwrap();
    expected.locked_mut().add(99).unwrap();

    assert_eq!(result, Some(expected));

    let request = contract.last_request().unwrap();
    assert_eq!(
        request,
        Request::new(
            chuck().to_string(),
            destination.account_id().to_string(),
            vec![TransferableAsset::new(w_near.name().to_owned(), 899, 100)]
        )
    )
}

#[test]
#[cfg(feature = "testable")]
fn handle_success_response_wnear_external_transfer() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let w_near = <Token>::new(WNEAR.to_owned());

    contract.register(w_near.clone());
    let token_id = contract.token_id(w_near.name()).unwrap();

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    let storage_cost = contract
        .get_storage_balance(chuck(), w_near.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.transfer(
        w_near.name().to_string(),
        destination.clone(),
        U128::from(999),
    );

    let result = contract.account_balance(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();

    expected.deposit_mut().add(1).unwrap();
    expected.locked_mut().add(899).unwrap();
    expected.locked_mut().add(100).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), w_near.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "TokenBSH".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 0,
                message: "Transfer Success".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    let result = contract.balance_of(alice(), w_near.name().to_string());
    assert_eq!(result, U128::from(999));

    let result = contract.account_balance(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1).unwrap();

    assert_eq!(result, Some(expected));

    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![
            AccumulatedAssetFees {
                name: nativecoin.name().to_string(),
                network: nativecoin.network().to_string(),
                accumulated_fees: 0
            },
            AccumulatedAssetFees {
                name: w_near.name().to_string(),
                network: w_near.network().to_string(),
                accumulated_fees: 100
            }
        ]
    );
}

#[test]
#[cfg(feature = "testable")]
fn handle_success_response_baln_external_transfer() {
    let context = |account_id: AccountId, deposit: u128, prepaid_gas: u64| {
        get_context(false, account_id, deposit, 0, prepaid_gas)
    };
    testing_env!(
        context(alice(), 1_000_000_000_000_000_000_000_000, 10u64.pow(18)),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let baln = <Token>::new(BALN.to_owned());
    contract.register(baln.clone());
    let token_id = env::sha256(baln.name().to_owned().as_bytes());
    contract.register_token_callback(baln.clone(), token_id.clone().try_into().unwrap());

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "TokenBSH".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: destination.account_id().to_string(),
                receiver: chuck().to_string(),
                assets: vec![TransferableAsset::new(baln.name().to_owned(), 900, 99)],
            },
        )),
    );

    testing_env!(context(bmc(), 0, 10u64.pow(18)));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(
        context(alice(), 0, 10u64.pow(18)),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.on_mint(
        900,
        token_id.clone().try_into().unwrap(),
        baln.symbol().to_string(),
        chuck(),
        Ok(U128::from(700000)),
    );

    let storage_cost = contract
        .get_storage_balance(chuck(), baln.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost, 10u64.pow(18)));
    contract.transfer(
        baln.name().to_string(),
        destination.clone(),
        U128::from(800),
    );

    let result = contract.account_balance(chuck(), baln.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    expected.locked_mut().add(800).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), baln.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "TokenBSH".to_string(),
        WrappedI128::new(contract.serial_no()),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 0,
                message: "Transfer Success".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0, 10u64.pow(18)));

    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(
        context(alice(), 0, 10u64.pow(18)),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.on_burn(
        719,
        token_id.clone().try_into().unwrap(),
        baln.symbol().to_string(),
    );

    let result = contract.balance_of(alice(), baln.name().to_string());
    assert_eq!(result, U128::from(81));

    let result = contract.account_balance(chuck(), baln.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();

    assert_eq!(result, Some(expected));

    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![
            AccumulatedAssetFees {
                name: nativecoin.name().to_string(),
                network: nativecoin.network().to_string(),
                accumulated_fees: 0
            },
            AccumulatedAssetFees {
                name: baln.name().to_string(),
                network: baln.network().to_string(),
                accumulated_fees: 81
            }
        ]
    );
}

#[test]
#[cfg(feature = "testable")]
fn handle_failure_response_wnear_external_transfer() {
    use libraries::types::{messages::ErrorMessage, AccumulatedAssetFees};

    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let w_near = <Token>::new(WNEAR.to_owned());
    contract.register(w_near.clone());
    let token_id = contract.token_id(w_near.name()).unwrap();

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    let storage_cost = contract
        .get_storage_balance(chuck(), w_near.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.transfer(
        w_near.name().to_string(),
        destination.clone(),
        U128::from(999),
    );

    let result = contract.account_balance(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();

    expected.deposit_mut().add(1).unwrap();
    expected.locked_mut().add(899).unwrap();
    expected.locked_mut().add(100).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), w_near.name().to_string());
    assert_eq!(result, U128::from(0));

    let source = BTPAddress::new("btp://0x1.icon/0x12345678".to_string());
    let btp_message = BtpMessage::new(
        source.clone(),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "TokenBSH".to_string(),
        WrappedI128::new(-1),
        vec![],
        Some(ErrorMessage::new(1, "Transfer Failed".to_string())),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_error(
        source.clone(),
        "TokenBSH".to_string(),
        -1,
        btp_message.try_into().unwrap(),
    );

    let result = contract.balance_of(alice(), w_near.name().to_string());
    assert_eq!(result, U128::from(100));

    let result = contract.account_balance(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1).unwrap();
    expected.refundable_mut().add(899).unwrap();

    assert_eq!(result, Some(expected));

    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![
            AccumulatedAssetFees {
                name: nativecoin.name().to_string(),
                network: nativecoin.network().to_string(),
                accumulated_fees: 0
            },
            AccumulatedAssetFees {
                name: w_near.name().to_string(),
                network: w_near.network().to_string(),
                accumulated_fees: 100
            }
        ]
    );
}

#[test]
#[cfg(feature = "testable")]
fn handle_failure_response_baln_coin_external_transfer() {
    use libraries::types::AccumulatedAssetFees;

    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 1_000_000_000_000_000_000_000_000),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let baln = <Token>::new(BALN.to_owned());
    contract.register(baln.clone());
    let token_id = env::sha256(baln.name().to_owned().as_bytes());
    contract.register_token_callback(baln.clone(), token_id.clone().try_into().unwrap());

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "TokenBSH".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: destination.account_id().to_string(),
                receiver: chuck().to_string(),
                assets: vec![TransferableAsset::new(baln.name().to_owned(), 900, 99)],
            },
        )),
    );

    testing_env!(
        context(bmc(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.handle_btp_message(btp_message.try_into().unwrap());

    contract.on_mint(
        900,
        token_id.clone().try_into().unwrap(),
        baln.symbol().to_string(),
        chuck(),
        Ok(U128::from(700000)),
    );

    let storage_cost = contract
        .get_storage_balance(chuck(), baln.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.transfer(
        baln.name().to_string(),
        destination.clone(),
        U128::from(800),
    );

    let result = contract.account_balance(chuck(), baln.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    expected.locked_mut().add(800).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), baln.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "TokenBSH".to_string(),
        WrappedI128::new(contract.serial_no()),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 1,
                message: "Transfer Failed".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    let result = contract.balance_of(alice(), baln.name().to_string());
    assert_eq!(result, U128::from(81));

    let result = contract.account_balance(chuck(), baln.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    expected.refundable_mut().add(719).unwrap();

    assert_eq!(result, Some(expected));
    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![
            AccumulatedAssetFees {
                name: nativecoin.name().to_string(),
                network: nativecoin.network().to_string(),
                accumulated_fees: 0
            },
            AccumulatedAssetFees {
                name: baln.name().to_string(),
                network: baln.network().to_string(),
                accumulated_fees: 81
            }
        ]
    );
}

#[test]
#[cfg(feature = "testable")]
fn reclaim_baln_coin() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 1_000_000_000_000_000_000_000_000),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let baln = <Token>::new(BALN.to_owned());
    contract.register(baln.clone());
    let token_id: [u8; 32] = env::sha256(baln.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(baln.clone(), token_id.clone());

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "TokenBSH".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: destination.account_id().to_string(),
                receiver: chuck().to_string(),
                assets: vec![TransferableAsset::new(baln.name().to_owned(), 899, 100)],
            },
        )),
    );

    testing_env!(
        context(bmc(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.handle_btp_message(btp_message.try_into().unwrap());
    contract.on_mint(
        899,
        token_id.clone(),
        baln.symbol().to_string(),
        chuck(),
        Ok(U128::from(700000)),
    );

    let storage_cost = contract
        .get_storage_balance(chuck(), baln.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), storage_cost));
    contract.transfer(
        baln.name().to_string(),
        destination.clone(),
        U128::from(800),
    );

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "TokenBSH".to_string(),
        WrappedI128::new(contract.serial_no()),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::ResponseHandleService {
                code: 1,
                message: "Transfer Failed".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(context(chuck(), 0));
    contract.reclaim(baln.name().to_string(), U128::from(700));

    let result = contract.account_balance(chuck(), baln.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(799).unwrap();
    expected.refundable_mut().add(19).unwrap();

    assert_eq!(result, Some(expected));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_higher_amount() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };

    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let w_near = <Token>::new(WNEAR.to_owned());
    contract.register(w_near.clone());

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    testing_env!(context(chuck(), 1000));
    contract.transfer(w_near.name().to_string(), destination, U128::from(1001));
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsToken: WNEAR")]
fn external_transfer_unregistered_coin() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let w_near = <Token>::new(WNEAR.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 0));
    contract.transfer(w_near.name().to_string(), destination, U128::from(1001));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_nil_balance() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };

    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let w_near = <Token>::new(WNEAR.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "TokenBSH".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    contract.register(w_near.clone());

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    testing_env!(context(chuck(), 0));
    contract.transfer(w_near.name().to_string(), destination, U128::from(1001));
}

#[test]
#[cfg(feature = "testable")]
fn external_transfer_batch() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );

    let w_near = <Token>::new(WNEAR.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    contract.register(w_near.clone());
    let token_id = contract.token_id(w_near.name()).unwrap();

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    let stroage_cost = contract
        .get_storage_balance(chuck(), w_near.name().to_string())
        .0
        + 1;

    testing_env!(context(chuck(), stroage_cost));
    contract.transfer_batch(
        vec![w_near.name().to_string()],
        destination,
        vec![U128::from(999)],
    );
    // TODO: Add other tokens
    let result = contract.account_balance(chuck(), w_near.name().to_string());
    let mut expected = AccountBalance::default();

    expected.deposit_mut().add(1).unwrap();
    expected.locked_mut().add(900).unwrap();
    expected.locked_mut().add(99).unwrap();

    assert_eq!(result, Some(expected));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_batch_higher_amount() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };

    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );

    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let w_near = <Token>::new(WNEAR.to_owned());
    contract.register(w_near.clone());

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    testing_env!(context(chuck(), 0));
    contract.transfer_batch(
        vec![w_near.name().to_string()],
        destination,
        vec![U128::from(1001)],
    );
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsToken")]
fn external_transfer_batch_unregistered_coin() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };
    testing_env!(
        context(alice(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );

    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let w_near = <Token>::new(WNEAR.to_owned());
    let baln = <Token>::new(BALN.to_owned());
    contract.register(w_near.clone());

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    testing_env!(context(chuck(), 0));
    contract.transfer_batch(
        vec![w_near.name().to_string(), baln.name().to_string()],
        destination,
        vec![U128::from(900), U128::from(1)],
    );
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_batch_nil_balance() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0, 10u64.pow(18))
    };

    testing_env!(
        context(alice(), 1_000_000_000_000_000_000_000_000),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let nativecoin = <Token>::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let w_near = <Token>::new(WNEAR.to_owned());
    let baln = <Token>::new(BALN.to_owned());
    contract.register(w_near.clone());

    testing_env!(
        context(alice(), 1_000_000_000_000_000_000_000_000),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.register(baln.clone());
    let token_id: [u8; 32] = env::sha256(baln.name().to_owned().as_bytes())
        .try_into()
        .unwrap();
    contract.register_token_callback(baln.clone(), token_id.clone());

    testing_env!(
        context(wnear(), 0),
        VMConfig::test(),
        RuntimeFeesConfig::test(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.ft_on_transfer(chuck(), U128::from(1000), "".to_string());

    testing_env!(context(chuck(), 0));
    contract.transfer_batch(
        vec![w_near.name().to_string(), baln.name().to_string()],
        destination,
        vec![U128::from(900), U128::from(1)],
    );
}
