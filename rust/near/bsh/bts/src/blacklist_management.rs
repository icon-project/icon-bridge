use super::*;

#[near_bindgen]
impl BtpTokenService {
    /// Returns the blacklisted users
    pub fn get_blacklisted_users(&self) -> Vec<AccountId> {
        self.blacklisted_accounts.to_vec()
    }
    /// Checks whether the user is black listed or not
    pub fn is_user_black_listed(&self, user: AccountId) -> bool {
        self.blacklisted_accounts.contains(&user)
    }
}

impl BtpTokenService {
    /// We can add a user to blacklist
    /// caller should be a owner
    /// # Arguments
    /// * `users` - user accountId should be given
    /// 
    pub fn add_to_blacklist(&mut self, users: Vec<AccountId>) {
        users
            .iter()
            .for_each(|user| self.blacklisted_accounts.add(user));
    }

    
    /// We can remove the user from blacklist
    /// caller should be the owner
    /// # Arguments
    /// * `users` - user accountid should be given
    /// 
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
