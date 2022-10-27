use bts::BtpTokenService;
use near_sdk::{env, json_types::U128, testing_env, AccountId, VMContext, test_utils::VMContextBuilder, Gas};
pub mod accounts;
use accounts::*;
use libraries::types::{
    messages::{BtpMessage, TokenServiceMessage, TokenServiceType},
    Account, AccountBalance, Asset, BTPAddress, Math, WrappedI128, WrappedNativeCoin,
};
mod token;
use libraries::types::{Request, TransferableAsset};
use std::convert::TryInto;
use token::*;

pub type Coin = Asset<WrappedNativeCoin>;

fn get_context(is_view: bool, signer_account_id: AccountId, attached_deposit: u128, account_balance: u128) -> VMContext {
    VMContextBuilder::new()
        .current_account_id(alice())
        .is_view(is_view)
        .signer_account_id(signer_account_id.clone())
        .predecessor_account_id(signer_account_id)
        .storage_usage(env::storage_usage())
        .prepaid_gas(Gas(10u64.pow(18)))
        .attached_deposit(attached_deposit)
        .account_balance(account_balance)
        .build()
}

#[test]
#[cfg(feature = "testable")]
fn handle_fee_gathering() {
    use libraries::types::AccumulatedAssetFees;

    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 1000));

    contract.deposit();
    contract.transfer(
        nativecoin.name().to_string(),
        destination.clone(),
        U128::from(999),
    );

    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();

    expected.deposit_mut().add(1).unwrap();
    expected.locked_mut().add(900).unwrap();
    expected.locked_mut().add(99).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), nativecoin.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
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

    let result = contract.balance_of(alice(), nativecoin.name().to_string());
    assert_eq!(result, U128::from(999));

    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1).unwrap();

    assert_eq!(result, Some(expected));

    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![AccumulatedAssetFees {
            name: nativecoin.name().to_string(),
            network: nativecoin.network().to_string(),
            accumulated_fees: 100
        }]
    );

    let fee_aggregator =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    contract.handle_fee_gathering(fee_aggregator, "nativecoin".to_string());

    let accumulted_fees = contract.accumulated_fees();

    assert_eq!(
        accumulted_fees,
        vec![AccumulatedAssetFees {
            name: nativecoin.name().to_string(),
            network: nativecoin.network().to_string(),
            accumulated_fees: 0
        }]
    );

    let result = contract.balance_of(alice(), nativecoin.name().to_string());
    assert_eq!(result, U128::from(899));
}

#[test]
fn get_fee() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(false, account_id, deposit, 0)
    };
    testing_env!(context(alice(), 0));

    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    testing_env!(context(alice(), 0));

    let result = contract.get_fee("NEAR".into(), U128(1000));
    assert_eq!(result, U128::from(101));

    contract.set_fee_ratio("NEAR".into(), 10.into(), 10.into());

    testing_env!(context(charlie(), 0));
    let result = contract.get_fee("NEAR".into(), U128(1000));
    assert_eq!(result, U128::from(11));

}
