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

        match self.handle_service_message(message.clone().try_into()) {
            Ok(outcome) => {
                if let Some(service_message) = outcome {
                    self.send_response(message.serial_no(), message.source(), service_message)
                }
            }
            Err(error) => {
                self.send_response(
                    message.serial_no(),
                    message.source(),
                    TokenServiceMessage::new(TokenServiceType::ResponseHandleService {
                        code: 1,
                        message: format!("{}", error),
                    }),
                );

                #[cfg(feature = "testable")]
                env::panic_str(error.to_string().as_str())
            }
        }
    }

    pub fn handle_btp_error(
        &mut self,
        source: BTPAddress,
        service: String,
        serial_no: i128,
        message: BtpMessage<SerializedMessage>,
    ) {
        self.assert_predecessor_is_bmc();
        self.assert_valid_service(&service);

        let error_message: BtpMessage<ErrorMessage> = message.try_into().unwrap();
        self.handle_response(
            &WrappedI128::new(serial_no).negate(),
            RC_ERROR,
            &format!(
                "[BTPError] source: {}, code: {} message: {:?}",
                source,
                RC_ERROR,
                error_message.message().clone().unwrap()
            ),
        )
        .unwrap();
    }

    #[cfg(feature = "testable")]
    pub fn last_request(&self) -> Option<Request> {
        self.requests().get(self.serial_no())
    }

    #[private]
    pub fn send_service_message_callback(
        &mut self,
        destination_network: String,
        message: TokenServiceMessage,
        serial_no: i128,
    ) {
        if let TokenServiceType::RequestTokenTransfer {
            sender,
            receiver,
            assets,
        } = message.service_type()
        {
            match env::promise_result(0) {
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
                      "receiver_address": format!("btp://{}/{}", destination_network, receiver),
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
                      "receiver_address": format!("btp://{}/{}", destination_network, receiver),
                      "assets" : assets_log
                    });

                    log!(near_sdk::serde_json::to_string(&log).unwrap());
                    self.rollback_external_transfer(&AccountId::from_str(sender).unwrap(), assets)
                }
            };
        }
    }
}

impl BtpTokenService {
    fn handle_service_message(
        &mut self,
        message: Result<BtpMessage<TokenServiceMessage>, BshError>,
    ) -> Result<Option<TokenServiceMessage>, BshError> {
        let btp_message = message?;

        if let Some(service_message) = btp_message.message() {
            match service_message.service_type() {
                TokenServiceType::RequestTokenTransfer {
                    sender: _,
                    ref receiver,
                    ref assets,
                } => self.handle_token_transfer(receiver, assets),

                TokenServiceType::ResponseHandleService {
                    ref code,
                    ref message,
                } => self.handle_response(btp_message.serial_no(), *code, message),

                TokenServiceType::RequestBlacklist {
                    request_type,
                    addresses,
                    #[allow(unused_variables)]
                    network,
                } => {
                    let mut non_valid_addresses: Vec<String> = Vec::new();
                    let mut valid_addresses: Vec<AccountId> = Vec::new();
                    addresses.iter().clone().for_each(|address| {
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
                    token_names,
                    token_limits,
                    #[allow(unused_variables)]
                    network,
                } => match self.set_token_limit(token_names.clone(), token_limits.clone()) {
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
                sender_id.into(),
                destination.contract_address().unwrap(),
                assets,
            ),
        );
        self.send_message(serial_no, destination.network_address().unwrap(), message);
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
            service_message,
        );
    }

    fn handle_response(
        &mut self,
        serial_no: &WrappedI128,
        code: u8,
        message: &str,
    ) -> Result<Option<TokenServiceMessage>, BshError> {
        if let Some(request) = self.requests().get(*serial_no.get()) {
            let sender_id = AccountId::try_from(request.sender().to_owned()).unwrap();
            if code == RC_OK {
                self.finalize_external_transfer(&sender_id, request.assets());
            } else if code == RC_ERROR {
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
        message: TokenServiceMessage,
    ) {
        #[cfg(feature = "testable")]
        {
            let service_message: SerializedMessage = message.clone().into();
            self.message.set(&(service_message.data().clone().into()));
        }

        ext_bmc::ext(self.bmc.clone())
            .send_service_message(
                serial_no,
                self.name.clone(),
                destination_network.clone(),
                message.clone().into(),
            )
            .then(
                Self::ext(env::current_account_id()).send_service_message_callback(
                    destination_network,
                    message,
                    serial_no,
                ),
            );
    }
}
