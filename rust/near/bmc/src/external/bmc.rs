use crate::{BTPAddress, BtpMessage, SerializedMessage,SerializedBtpMessages};
use near_sdk::ext_contract;

#[ext_contract(bmc_contract)]
pub trait BmcContract {
    fn emit_message(link: BTPAddress, btp_message: BtpMessage<SerializedMessage>);
    fn emit_error();
    fn handle_external_service_message_callback(source: BTPAddress, message: BtpMessage<SerializedMessage>);
}
