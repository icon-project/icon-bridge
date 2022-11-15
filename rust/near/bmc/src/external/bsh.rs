use crate::{BTPAddress, BtpMessage, SerializedMessage};
use near_sdk::ext_contract;

#[ext_contract(bsh_contract)]
pub trait BshContract {
    /// Btp messages are handled
    fn handle_btp_message(message: BtpMessage<SerializedMessage>);
    fn handle_btp_error(
        source: BTPAddress,
        service: String,
        serial_no: i128,
        message: BtpMessage<SerializedMessage>,
    );

    /// Handling the fee gathering
    /// # Arguments
    /// * `fee_aggregator` - should be in the form of btp://0x1.near/account.testnet
    /// * `service` - name of the existence service should be given
    /// 
    fn handle_fee_gathering(fee_aggregator: BTPAddress, service: String);
}
