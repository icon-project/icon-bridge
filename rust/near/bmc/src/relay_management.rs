use super::*;
use near_sdk::BlockHeight;

#[near_bindgen]
impl BtpMessageCenter {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Relay Management  * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

/// Adding the relays to the bmc
/// Caller should be a owner
/// 
/// # Arguments
/// * link - should be in the form of btp://0x1.near/account.testnet
/// * relays - Relays can be added by using vectors
/// Multiple relays can be added 

    pub fn add_relays(&mut self, link: BTPAddress, relays: Vec<AccountId>) {
        self.assert_have_permission();
        self.assert_link_exists(&link);
        if let Some(link_property) = self.links.get(&link).as_mut() {
            link_property.relays_mut().set(&relays);
            self.links.set(&link, &link_property);
        }
    }
/// Adding the relay to the bmc
///  #Arguments
/// * link - should be in the form of btp://0x1.near/account.testnet
/// * relay - Account id should be given to add that inside the link 
/// 
    pub fn add_relay(&mut self, link: BTPAddress, relay: AccountId) {
        self.assert_have_permission();
        self.assert_link_exists(&link);
        self.assert_relay_not_exists(&link, &relay);
        if let Some(link_property) = self.links.get(&link).as_mut(){
            link_property.relays_mut().add(&relay);
            self.links.set(&link, link_property)
        }
    }

/// Removing the relay from the bmc
/// 
/// # Arguments
/// * link - Should be in the form of btp://0x1.near/account.testnet
/// * relay - Account id should be given.

    pub fn remove_relay(&mut self, link: BTPAddress, relay: AccountId) {
        self.assert_have_permission();
        self.assert_link_exists(&link);
        self.assert_relay_exists(&link,&relay);
        if let Some(link_property) = self.links.get(&link).as_mut() {
            link_property.relays_mut().remove(&relay);
            self.links.set(&link, &link_property);
        }
    }


    /// Getting the list of relays prsent inside the bmc
    /// Caller can be any
    /// #Arguments
    /// *link - should be in the form of btp://0x1.near/account.testnet
    /// Fetch the relays respective to the given link
    pub fn get_relays(&self, link: BTPAddress) -> Value {
        self.assert_link_exists(&link);
        if let Some(link_property) = self.links.get(&link).as_mut() {
            to_value(link_property.relays().to_vec()).unwrap()
        } else {
            to_value(Vec::new() as Vec<String>).unwrap()
        }
    }
}
