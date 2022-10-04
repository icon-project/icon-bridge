use std::str::FromStr;

use libraries::types::Request;
use libraries::types::WrappedI128;

use super::*;

#[near_bindgen]
impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * *  Messaging  * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    pub fn handle_btp_message(&mut self, message: BtpMessage<SerializedMessage>) {
        self.assert_predecessor_is_bmc();
        self.assert_valid_service(message.service());
        let outcome = self.handle_service_message(message.clone().try_into());

        if outcome.is_err() {
            let error = outcome.clone().unwrap_err();
            self.send_response(
                message.serial_no(),
                message.source(),
                TokenServiceMessage::new(TokenServiceType::ResponseHandleService {
                    code: 1,
                    message: format!("{}", error),
                }),
            );
        } else {
            match outcome.clone().unwrap() {
                Some(service_message) => {
                    self.send_response(message.serial_no(), message.source(), service_message);
                }
                None => (),
            }
        }

        #[cfg(feature = "testable")]
        outcome.clone().unwrap();
    }

    pub fn handle_btp_error(
        &mut self,
        source: BTPAddress,
        service: String,
        serial_no: i128,
        code: u128,
        message: String,
    ) {
        self.assert_predecessor_is_bmc();
        self.assert_valid_service(&service);
        self.handle_response(
            &WrappedI128::new(serial_no),
            1,
            &format!(
                "[BTPError] source: {}, code: {} message: {}",
                source, code, message
            ),
        )
        .unwrap();
    }

    #[cfg(feature = "testable")]
    pub fn last_request(&self) -> Option<Request> {
        self.requests().get(self.serial_no())
    }

    #[private]
    pub fn send_service_message_callback(&mut self, message: TokenServiceMessage, serial_no: i128) {
        match message.service_type() {
            TokenServiceType::RequestTokenTransfer {
                sender,
                receiver,
                assets,
            } => match env::promise_result(0) {
                PromiseResult::Successful(_) => {
                    let mut assets_log: Vec<Value> = Vec::new();
                    assets.iter().for_each(|asset| {
                        assets_log.push(json!({
                        "token_name": asset.name(),
                        "amount": asset.amount().to_string(),
                        "fee": asset.fees().to_string(),
                        }))
                    });
                    let log = json!({
                      "event": "TransferStart",
                      "code": "0",
                      "sender_address": sender,
                      "serial_number": serial_no.to_string(),
                      "receiver_address": receiver,
                      "assets": assets_log
                    });

                    log!(near_sdk::serde_json::to_string(&log).unwrap())
                }
                PromiseResult::NotReady => log!("Not Ready"),
                PromiseResult::Failed => {
                    let mut assets_log: Vec<Value> = Vec::new();
                    assets.iter().for_each(|asset| {
                        assets_log.push(json!({
                        "token_name": asset.name(),
                        "amount": asset.amount().to_string(),
                        "fee": asset.fees().to_string(),
                        }))
                    });
                    let log = json!({
                      "event": "TransferStart",
                      "code": "1",
                      "sender_address": sender,
                      "serial_number": serial_no.to_string(),
                      "receiver_address": receiver,
                      "assets" : assets_log
                    });

                    log!(near_sdk::serde_json::to_string(&log).unwrap());
                    self.rollback_external_transfer(&AccountId::from_str(sender).unwrap(), assets)
                }
            },

            _ => {}
        }
    }
}

impl BtpTokenService {
    fn handle_service_message(
        &mut self,
        message: Result<BtpMessage<TokenServiceMessage>, BshError>,
    ) -> Result<Option<TokenServiceMessage>, BshError> {
        let btp_message = message.clone()?;

        if let Some(service_message) = btp_message.message() {
            match service_message.service_type() {
                TokenServiceType::RequestTokenTransfer {
                    sender: _,
                    ref receiver,
                    ref assets,
                } => self.handle_coin_transfer(btp_message.source(), receiver, assets),

                TokenServiceType::ResponseHandleService {
                    ref code,
                    ref message,
                } => self.handle_response(btp_message.serial_no(), *code, &message),

                TokenServiceType::RequestBlacklist {
                    request_type,
                    addresses,
                    network,
                } => {
                    let mut non_valid_addresses: Vec<String> = Vec::new();
                    let mut valid_addresses: Vec<AccountId> = Vec::new();
                    addresses.into_iter().clone().for_each(|address| {
                        match AccountId::try_from(address.clone()) {
                            Ok(account_id) => valid_addresses.push(account_id),
                            Err(_) => non_valid_addresses.push(address.to_string()),
                        }
                    });
                    if !non_valid_addresses.is_empty() {
                        return Err(BshError::InvalidAddress {
                            message: non_valid_addresses.join(", "),
                        });
                    }
                    match request_type {
                        BlackListType::AddToBlacklist => {
                            self.add_to_blacklist(valid_addresses);
                            let response =
                                TokenServiceMessage::new(TokenServiceType::ResponseBlacklist {
                                    code: 0,
                                    message: "AddedToBlacklist".to_string(),
                                });
                            self.send_response(
                                btp_message.serial_no(),
                                btp_message.source(),
                                response,
                            );

                            Ok(None)
                        }
                        BlackListType::RemoveFromBlacklist => {
                            match self.remove_from_blacklist(valid_addresses) {
                                Ok(()) => {
                                    let response = TokenServiceMessage::new(
                                        TokenServiceType::ResponseBlacklist {
                                            code: 0,
                                            message: "RemovedFromBlacklist".to_string(),
                                        },
                                    );

                                    self.send_response(
                                        btp_message.serial_no(),
                                        btp_message.source(),
                                        response,
                                    );

                                    Ok(None)
                                }
                                Err(err) => {
                                    let response = TokenServiceMessage::new(
                                        TokenServiceType::ResponseBlacklist {
                                            code: 1,
                                            message: err.to_string(),
                                        },
                                    );
                                    self.send_response(
                                        btp_message.serial_no(),
                                        btp_message.source(),
                                        response,
                                    );

                                    Ok(None)
                                }
                            }
                        }
                        BlackListType::UnhandledType => todo!(),
                    }
                }
                TokenServiceType::RequestChangeTokenLimit {
                    coin_names,
                    token_limits,
                    network,
                } => match self.set_token_limit(coin_names.clone(), token_limits.clone()) {
                    Ok(()) => {
                        let response =
                            TokenServiceMessage::new(TokenServiceType::ResponseChangeTokenLimit {
                                code: 0,
                                message: "ChangeTokenLimit".to_string(),
                            });

                        self.send_response(btp_message.serial_no(), btp_message.source(), response);

                        Ok(None)
                    }
                    Err(err) => {
                        let response =
                            TokenServiceMessage::new(TokenServiceType::ResponseChangeTokenLimit {
                                code: 1,
                                message: err.to_string(),
                            });

                        self.send_response(btp_message.serial_no(), btp_message.source(), response);

                        Ok(None)
                    }
                },

                TokenServiceType::UnknownType => {
                    log!(
                        "Unknown Response: from {} for serial_no {}",
                        btp_message.source(),
                        btp_message.serial_no().get()
                    );
                    Ok(None)
                }

                _ => Ok(Some(TokenServiceMessage::new(
                    TokenServiceType::UnknownType,
                ))),
            }
        } else {
            Err(BshError::Unknown)
        }
    }

    pub fn send_request(
        &mut self,
        sender_id: AccountId,
        destination: BTPAddress,
        assets: Vec<TransferableAsset>,
    ) {
        let serial_no = self.serial_no.checked_add(1).unwrap();
        self.serial_no.clone_from(&serial_no);

        let message = TokenServiceMessage::new(TokenServiceType::RequestTokenTransfer {
            sender: sender_id.clone().into(),
            receiver: destination.contract_address().unwrap(),
            assets: assets.clone(),
        });

        self.requests_mut().add(
            serial_no,
            &Request::new(
                sender_id.clone().into(),
                destination.contract_address().unwrap(),
                assets,
            ),
        );
        self.send_message(
            serial_no,
            destination.network_address().unwrap(),
            message.into(),
        );
    }

    pub fn send_response(
        &mut self,
        serial_no: &WrappedI128,
        destination: &BTPAddress,
        service_message: TokenServiceMessage,
    ) {
        self.send_message(
            *serial_no.get(),
            destination.network_address().unwrap(),
            service_message.into(),
        );
    }

    fn handle_response(
        &mut self,
        serial_no: &WrappedI128,
        code: u128,
        message: &str,
    ) -> Result<Option<TokenServiceMessage>, BshError> {
        if let Some(request) = self.requests().get(*serial_no.get()) {
            let sender_id = AccountId::try_from(request.sender().to_owned()).unwrap();
            if code == 0 {
                self.finalize_external_transfer(&sender_id, request.assets());
            } else if code == 1 {
                self.rollback_external_transfer(&sender_id, request.assets());
            }
            self.requests_mut().remove(*serial_no.get());

            let log = json!({
                "event": "TransferEnd",
                "code": code.to_string(),
                "serial_number": serial_no.get().to_string(),
                "message": message,
            });
            log!(near_sdk::serde_json::to_string(&log).unwrap())
        }
        Ok(None)
    }

    pub fn send_message(
        &mut self,
        serial_no: i128,
        destination_network: String,
        message: SerializedMessage,
    ) {
        ext_bmc::send_service_message(
            serial_no,
            self.name.clone(),
            destination_network,
            message.clone(),
            self.bmc.clone(),
            estimate::NO_DEPOSIT,
            estimate::GAS_FOR_SEND_SERVICE_MESSAGE,
        )
        .then(ext_self::send_service_message_callback(
            message.clone().try_into().unwrap(),
            serial_no,
            env::current_account_id(),
            estimate::NO_DEPOSIT,
            estimate::GAS_FOR_SEND_SERVICE_MESSAGE,
        ));

        #[cfg(feature = "testable")]
        self.message.set(&(message.data().clone().into()));
    }
}
