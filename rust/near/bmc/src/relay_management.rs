use super::*;

#[near_bindgen]
impl BtpMessageCenter {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Relay Management  * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    pub fn add_relays(&mut self, link: BTPAddress, relays: Vec<AccountId>) {
        self.assert_have_permission();
        self.assert_link_exists(&link);

        if let Some(link_property) = self.links.get(&link).as_mut() {
            link_property.relays_mut().set(&relays);
            self.links.set(&link, link_property);
        }
    }

    pub fn add_relay(&mut self, link: BTPAddress, relay: AccountId) {
        self.assert_have_permission();
        self.assert_link_exists(&link);
        self.assert_relay_not_exists(&link, &relay);

        if let Some(link_property) = self.links.get(&link).as_mut() {
            link_property.relays_mut().add(&relay);
            self.links.set(&link, link_property)
        }
    }

    pub fn remove_relay(&mut self, link: BTPAddress, relay: AccountId) {
        self.assert_have_permission();
        self.assert_link_exists(&link);
        self.assert_relay_exists(&link, &relay);

        if let Some(link_property) = self.links.get(&link).as_mut() {
            link_property.relays_mut().remove(&relay);
            self.links.set(&link, link_property);
        }
    }

    pub fn get_relays(&self, link: BTPAddress) -> Value {
        self.assert_link_exists(&link);
        
        if let Some(link_property) = self.links.get(&link).as_mut() {
            to_value(link_property.relays().to_vec()).unwrap()
        } else {
            to_value(Vec::new() as Vec<String>).unwrap()
        }
    }
}
