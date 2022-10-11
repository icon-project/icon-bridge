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
}
