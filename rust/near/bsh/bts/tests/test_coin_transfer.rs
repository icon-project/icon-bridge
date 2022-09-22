use bts::BtpTokenService;
use near_sdk::{env, json_types::U128, testing_env, AccountId, PromiseResult, VMContext};
pub mod accounts;
use accounts::*;
use libraries::types::{
    messages::{BtpMessage, TokenServiceMessage, TokenServiceType},
    Account, AccountBalance, AccumulatedAssetFees, Asset, BTPAddress, Math, WrappedI128,
    WrappedNativeCoin,
};
mod token;
use libraries::types::{Request, TransferableAsset};
use std::convert::TryInto;
use token::*;
pub type Coin = Asset<WrappedNativeCoin>;

fn get_context(
    input: Vec<u8>,
    is_view: bool,
    signer_account_id: AccountId,
    attached_deposit: u128,
    storage_usage: u64,
    account_balance: u128,
) -> VMContext {
    VMContext {
        current_account_id: alice().to_string(),
        signer_account_id: signer_account_id.to_string(),
        signer_account_pk: vec![0, 1, 2],
        predecessor_account_id: signer_account_id.to_string(),
        input,
        block_index: 0,
        block_timestamp: 0,
        account_balance,
        account_locked_balance: 0,
        storage_usage,
        attached_deposit,
        prepaid_gas: 10u64.pow(18),
        random_seed: vec![0, 1, 2],
        is_view,
        output_data_receivers: vec![],
        epoch_height: 19,
    }
}

#[test]
fn deposit_native_coin() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    testing_env!(context(chuck(), 100));

    contract.deposit();

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    assert_eq!(result, U128::from(expected.deposit()))
}

#[test]
fn withdraw_native_coin() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(
            vec![],
            false,
            account_id,
            deposit,
            env::storage_usage(),
            1000,
        )
    };
    testing_env!(context(alice(), 0));
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let coin_id = contract.coin_id(nativecoin.name()).unwrap();
    testing_env!(context(chuck(), 1000));

    contract.deposit();

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1000).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));

    testing_env!(context(chuck(), 1));
    contract.withdraw(nativecoin.name().to_string(), U128::from(999));

    testing_env!(
        context(alice(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.on_withdraw(chuck(), 999, nativecoin.name().to_string(), coin_id.clone());

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn withdraw_native_coin_higher_amount() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(
            vec![],
            false,
            account_id,
            deposit,
            env::storage_usage(),
            1000,
        )
    };
    testing_env!(context(alice(), 0));
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );
    let coin_id = contract.coin_id(nativecoin.name()).unwrap();
    testing_env!(context(chuck(), 100));

    contract.deposit();

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    assert_eq!(result, U128::from(expected.deposit()));

    testing_env!(context(chuck(), 1));
    contract.withdraw(nativecoin.name().to_string(), U128::from(1000));

    let result = contract.balance_of(chuck(), nativecoin.name().to_string());
    let expected = AccountBalance::default();
    assert_eq!(result, U128::from(expected.deposit()));
}

#[test]
#[cfg(feature = "testable")]
fn external_transfer() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
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

    let request = contract.last_request().unwrap();
    assert_eq!(
        request,
        Request::new(
            chuck().to_string(),
            destination.account_id().to_string(),
            vec![TransferableAsset::new(
                nativecoin.name().to_owned(),
                899,
                100
            )]
        )
    )
}

#[test]
#[cfg(feature = "testable")]
fn handle_success_response_native_coin_external_transfer() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
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
}

#[test]
#[cfg(feature = "testable")]
fn handle_success_response_icx_coin_external_transfer() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(
        context(alice(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());

    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let icx_coin = Coin::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let coin_id = env::sha256(icx_coin.name().to_owned().as_bytes());
    contract.register_coin_callback(icx_coin.clone(), coin_id);

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: destination.account_id().to_string(),
                receiver: chuck().to_string(),
                assets: vec![TransferableAsset::new(icx_coin.name().to_owned(), 900, 99)],
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(
        context(chuck(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let coin_id = contract.coin_id(icx_coin.name()).unwrap();

    contract.on_mint(
        900,
        coin_id.clone(),
        icx_coin.symbol().to_string(),
        chuck().clone(),
    );

    contract.transfer(
        icx_coin.name().to_string(),
        destination.clone(),
        U128::from(800),
    );

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    expected.locked_mut().add(800).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), icx_coin.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(contract.serial_no()),
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

    testing_env!(
        context(chuck(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.on_burn(720, coin_id.clone(), icx_coin.symbol().to_string());

    let result = contract.balance_of(alice(), icx_coin.name().to_string());
    assert_eq!(result, U128::from(80));

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
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
                name: icx_coin.name().to_string(),
                network: icx_coin.network().to_string(),
                accumulated_fees: 81
            }
        ]
    );
}

#[test]
#[cfg(feature = "testable")]
fn handle_failure_response_native_coin_external_transfer() {
    use libraries::types::AccumulatedAssetFees;

    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
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
    let coin_id = contract.coin_id(nativecoin.name()).unwrap();

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
                code: 1,
                message: "Transfer Failed".to_string(),
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    let result = contract.balance_of(alice(), nativecoin.name().to_string());
    assert_eq!(result, U128::from(100));

    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(1).unwrap();
    expected.refundable_mut().add(899).unwrap();

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
}

#[test]
#[cfg(feature = "testable")]
fn handle_failure_response_icx_coin_external_transfer() {
    use libraries::types::AccumulatedAssetFees;

    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(
        context(alice(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let icx_coin = Coin::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let coin_id = env::sha256(icx_coin.name().to_owned().as_bytes());
    contract.register_coin_callback(icx_coin.clone(), coin_id);

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: destination.account_id().to_string(),
                receiver: chuck().to_string(),
                assets: vec![TransferableAsset::new(icx_coin.name().to_owned(), 900, 99)],
            },
        )),
    );
    testing_env!(
        context(bmc(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    contract.handle_btp_message(btp_message.try_into().unwrap());

    let coin_id = contract.coin_id(icx_coin.name()).unwrap();

    contract.on_mint(
        900,
        coin_id.clone(),
        icx_coin.symbol().to_string(),
        chuck().clone(),
    );

    testing_env!(context(chuck(), 0));
    contract.transfer(
        icx_coin.name().to_string(),
        destination.clone(),
        U128::from(800),
    );

    testing_env!(context(chuck(), 0));

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(100).unwrap();
    expected.locked_mut().add(800).unwrap();

    assert_eq!(result, Some(expected));

    let result = contract.balance_of(alice(), icx_coin.name().to_string());
    assert_eq!(result, U128::from(0));

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
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

    let result = contract.balance_of(alice(), icx_coin.name().to_string());
    assert_eq!(result, U128::from(81));

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
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
                name: icx_coin.name().to_string(),
                network: icx_coin.network().to_string(),
                accumulated_fees: 81
            }
        ]
    );
}

#[test]
#[cfg(feature = "testable")]
fn reclaim_icx_coin() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(
        context(alice(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    let icx_coin = Coin::new(ICON_COIN.to_owned());
    contract.register(icx_coin.clone());
    let coin_id = env::sha256(icx_coin.name().to_owned().as_bytes());
    contract.register_coin_callback(icx_coin.clone(), coin_id);

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
        WrappedI128::new(1),
        vec![],
        Some(TokenServiceMessage::new(
            TokenServiceType::RequestTokenTransfer {
                sender: destination.account_id().to_string(),
                receiver: chuck().to_string(),
                assets: vec![TransferableAsset::new(icx_coin.name().to_owned(), 900, 99)],
            },
        )),
    );

    testing_env!(context(bmc(), 0));
    contract.handle_btp_message(btp_message.try_into().unwrap());

    testing_env!(
        context(chuck(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let coin_id = contract.coin_id(icx_coin.name()).unwrap();
    contract.on_mint(900, coin_id.clone(), icx_coin.symbol().to_string(), chuck());

    testing_env!(context(chuck(), 0));
    contract.transfer(
        icx_coin.name().to_string(),
        destination.clone(),
        U128::from(800),
    );

    let btp_message = &BtpMessage::new(
        BTPAddress::new("btp://0x1.icon/0x12345678".to_string()),
        BTPAddress::new("btp://1234.iconee/0x12345678".to_string()),
        "nativecoin".to_string(),
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
    contract.reclaim(icx_coin.name().to_string(), U128::from(700));

    let result = contract.account_balance(chuck(), icx_coin.name().to_string());
    let mut expected = AccountBalance::default();
    expected.deposit_mut().add(800).unwrap();
    expected.refundable_mut().add(19).unwrap();

    assert_eq!(result, Some(expected));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_higher_amount() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
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
    contract.transfer(nativecoin.name().to_string(), destination, U128::from(1001));
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsToken: ICON")]
fn external_transfer_unregistered_coin() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let icx_coin = Coin::new(ICON_COIN.to_owned());
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
    contract.transfer(icx_coin.name().to_string(), destination, U128::from(1001));
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_nil_balance() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(
        context(alice(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let icx_coin = Coin::new(ICON_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    contract.register(icx_coin.clone());
    let coin_id = env::sha256(icx_coin.name().to_owned().as_bytes());
    contract.register_coin_callback(icx_coin.clone(), coin_id);
    testing_env!(context(chuck(), 1000));

    contract.deposit();
    contract.transfer(icx_coin.name().to_string(), destination, U128::from(1001));
}

#[test]
#[cfg(feature = "testable")]
fn external_transfer_batch() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
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
    let coin_id = contract.coin_id(nativecoin.name()).unwrap();

    contract.deposit();
    contract.transfer_batch(
        vec![nativecoin.name().to_string()],
        destination,
        vec![U128::from(999)],
    );
    // TODO: Add other tokens
    let result = contract.account_balance(chuck(), nativecoin.name().to_string());
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
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
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
    contract.transfer_batch(
        vec![nativecoin.name().to_string()],
        destination,
        vec![U128::from(1001)],
    );
}

#[test]
#[should_panic(expected = "BSHRevertNotExistsToken")]
fn external_transfer_batch_unregistered_coin() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(context(alice(), 0));
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let icx_coin = Coin::new(ICON_COIN.to_owned());
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
    contract.transfer_batch(
        vec![nativecoin.name().to_string(), icx_coin.name().to_string()],
        destination,
        vec![U128::from(900), U128::from(1)],
    );
}

#[test]
#[should_panic(expected = "BSHRevertNotMinimumDeposit")]
fn external_transfer_batch_nil_balance() {
    let context = |account_id: AccountId, deposit: u128| {
        get_context(vec![], false, account_id, deposit, env::storage_usage(), 0)
    };
    testing_env!(
        context(alice(), 0),
        Default::default(),
        Default::default(),
        Default::default(),
        vec![PromiseResult::Successful(vec![1_u8])]
    );
    let nativecoin = Coin::new(NATIVE_COIN.to_owned());
    let icx_coin = Coin::new(ICON_COIN.to_owned());
    let destination =
        BTPAddress::new("btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
    let mut contract = BtpTokenService::new(
        "nativecoin".to_string(),
        bmc(),
        "0x1.near".into(),
        nativecoin.clone(),
    );

    contract.register(icx_coin.clone());
    let coin_id = env::sha256(icx_coin.name().to_owned().as_bytes());
    contract.register_coin_callback(icx_coin.clone(), coin_id);
    testing_env!(context(chuck(), 1000));

    contract.deposit();
    contract.transfer_batch(
        vec![nativecoin.name().to_string(), icx_coin.name().to_string()],
        destination,
        vec![U128::from(900), U128::from(1)],
    );
}
