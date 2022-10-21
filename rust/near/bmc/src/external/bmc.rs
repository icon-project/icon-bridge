use crate::{BTPAddress, BtpMessage, SerializedBtpMessages, SerializedMessage};
use libraries::types::WrappedI128;
use near_sdk::{ext_contract, json_types::U128};
#[ext_contract(bmc_contract)]
pub trait BmcContract {
    fn emit_message(link: BTPAddress, btp_message: BtpMessage<SerializedMessage>);
    fn handle_external_service_message_callback(
        source: BTPAddress,
        message: BtpMessage<SerializedMessage>,
    );
    fn handle_btp_error_callback(message: BtpMessage<SerializedMessage>);
}
