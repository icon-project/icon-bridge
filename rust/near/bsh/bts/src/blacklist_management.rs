use super::*;

#[near_bindgen]
impl BtpTokenService {
    pub fn get_blacklisted_users(&self) -> Vec<AccountId> {
        self.blacklisted_accounts.to_vec()
    }

    pub fn is_user_black_listed(&self, user: AccountId) -> bool {
        self.blacklisted_accounts.contains(&user)
    }
}

impl BtpTokenService {
    pub fn add_to_blacklist(&mut self, users: Vec<AccountId>) {
        users
            .iter()
            .for_each(|user| self.blacklisted_accounts.add(user));
    }

    pub fn remove_from_blacklist(&mut self, users: Vec<AccountId>) -> Result<(), BshError> {
        let mut non_blacklisted_user: Vec<String> = Vec::new();

        users
            .iter()
            .for_each(|user| match self.ensure_user_blacklisted(user) {
                Ok(()) => {
                    self.blacklisted_accounts.remove(user);
                }
                Err(_) => non_blacklisted_user.push(user.to_string()),
            });

        if !non_blacklisted_user.is_empty() {
            return Err(BshError::NonBlacklistedUsers {
                message: non_blacklisted_user.join(", "),
            });
        }

        Ok(())
    }

    pub fn handle_request_blacklist(
        &mut self,
        request_type: &BlackListType,
        addresses: &[String],
        network: &str,
    ) -> Result<Option<TokenServiceMessage>, BshError> {
        let mut non_valid_addresses: Vec<String> = Vec::new();
        let mut valid_addresses: Vec<AccountId> = Vec::new();

        if network == self.network {
            addresses.iter().clone().for_each(|address| {
                match AccountId::try_from(address.clone()) {
                    Ok(account_id) => valid_addresses.push(account_id),
                    Err(_) => non_valid_addresses.push(address.to_string()),
                }
            });
        } else {
            return Err(BshError::InvalidAddress {
                message: network.to_string(),
            });
        }

        if !non_valid_addresses.is_empty() {
            return Err(BshError::InvalidAddress {
                message: non_valid_addresses.join(", "),
            });
        }

        match request_type {
            BlackListType::AddToBlacklist => {
                self.add_to_blacklist(valid_addresses);

                Ok(Some(TokenServiceMessage::new(
                    TokenServiceType::ResponseBlacklist {
                        code: 0,
                        message: "AddedToBlacklist".to_string(),
                    },
                )))
            }
            BlackListType::RemoveFromBlacklist => {
                self.remove_from_blacklist(valid_addresses)?;

                Ok(Some(TokenServiceMessage::new(
                    TokenServiceType::ResponseBlacklist {
                        code: 0,
                        message: "RemovedFromBlacklist".to_string(),
                    },
                )))
            }
            BlackListType::UnhandledType => todo!(),
        }
    }
}
