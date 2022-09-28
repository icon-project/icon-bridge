use libraries::types::{messages::TokenServiceMessage, AssetId};
use near_sdk::ext_contract;

#[ext_contract(ext_self)]
pub trait BtpTokenService {
    fn on_withdraw(
        &mut self,
        account: AccountId,
        amount: u128,
        coin_name: String,
        coin_id: AssetId,
    );

    fn on_mint(
        &mut self,
        amount: u128,
        coin_id: AssetId,
        coin_symbol: String,
        receiver_id: AccountId,
    );

    fn on_burn(&mut self, amount: u128, coin_id: AssetId, coin_symbol: String);

    fn send_service_message_callback(&mut self, message: TokenServiceMessage, serial_no: i128);
}
