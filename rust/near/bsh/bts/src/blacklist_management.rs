use super::*;

impl BtpTokenService {
    pub fn add_to_blacklist(&mut self, users: Vec<BTPAddress>) -> Result<(), BshError> {
        self.assert_have_permission();
        for user in users.iter() {
            match user.is_valid() {
                Ok(valid) => {
                    if valid {
                        match self.ensure_user_not_blacklisted(&user.account_id()) {
                            Ok(()) => {
                                self.blacklisted_accounts.add(&user.account_id());
                            }
                            Err(err) => return Err(err),
                        }
                    } else {
                        return Err(BshError::InvalidAddress {
                            message: user.account_id().to_string(),
                        });
                    }
                }
                Err(err) => return Err(BshError::InvalidAddress { message: err }),
            }
        }
        Ok(())
    }

    pub fn remove_from_blacklist(&mut self, users: Vec<BTPAddress>) -> Result<(), BshError> {
        self.assert_have_permission();
        for user in users.iter() {
            match user.is_valid() {
                Ok(valid) => {
                    if valid {
                        match self.ensure_user_blacklisted(&user.account_id()) {
                            Ok(()) => {
                                self.blacklisted_accounts.remove(&user.account_id());
                            }
                            Err(err) => return Err(err),
                        }
                    } else {
                        return Err(BshError::InvalidAddress {
                            message: user.account_id().to_string(),
                        });
                    }
                }
                Err(err) => return Err(BshError::InvalidAddress { message: err }),
            }
        }
        Ok(())
    }

    pub fn get_blacklisted_user(&self) -> Vec<AccountId> {
        self.blacklisted_accounts.to_vec()
    }
}
