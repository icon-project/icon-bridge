use crate::types::Event;

use super::*;

impl BtpMessageCenter {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Internal Validations  * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    /// Check whether signer account id is an owner
    pub fn assert_have_permission(&self) {
        require!(
            self.owners.contains(&env::predecessor_account_id()),
            format!("{}", BmcError::PermissionNotExist)
        );
    }


    ///check whether the given link is exists or not
    pub fn assert_link_exists(&self, link: &BTPAddress) {
        require!(
            self.links.contains(link),
            format!("{}", BmcError::LinkNotExist)
        );
    }
///check whether the given link is present or not
/// #Arguments
/// * `link` - Should be in the form of btp://0x1.near/account.testnet
/// 
    pub fn assert_link_does_not_exists(&self, link: &BTPAddress) {
        require!(
            !self.links.contains(link),
            format!("{}", BmcError::LinkExist)
        );
    }
    /// checks whether the link contains the route connection or not
    ///  #Arguments
    /// * `link` - Should be in the form of btp://0x1.near/account.testnet
    /// 
    pub fn assert_link_does_not_have_route_connection(&self, link: &BTPAddress) {
        require!(
            !self
                .connections
                .contains(&Connection::Route(link.network_address().unwrap())),
            format!("{}", BmcError::LinkRouteExist)
        )
    }
    /// checks whether the owner is exists or not by providing the accountId
    /// # Arguments 
    /// * `account` - Accountid should be provided (bmc.testnet)
    /// 
    pub fn assert_owner_exists(&self, account: &AccountId) {
        require!(
            self.owners.contains(&account),
            format!("{}", BmcError::OwnerNotExist)
        );
    }

    /// Checks that the given owner does not exist
    /// # Arguments
    /// * `account` - Accountid should be provided (bmc.testnet)
    /// 
    pub fn assert_owner_does_not_exists(&self, account: &AccountId) {
        require!(
            !self.owners.contains(account),
            format!("{}", BmcError::OwnerExist)
        );
    }
///Checks the present owner is not the last owner
    pub fn assert_owner_is_not_last_owner(&self) {
        require!(self.owners.len() > 1, format!("{}", BmcError::LastOwner));
    }
/// Checks if the relay is registered or not by giving the link.
/// # Arguments
/// * `link` - Should be in the form of btp://0x1.near/account.testnet
/// 
    pub fn assert_relay_is_registered(&self, link: &BTPAddress) {
        let link = self.links.get(link).unwrap();
        require!(
            link.relays().contains(&env::predecessor_account_id()),
            format!(
                "{}",
                BmcError::Unauthorized {
                    message: "not registered relay"
                }
            )
        )
    }
/// Checks whether the given relay is valid or not
/// # arguments
/// * `accepted_relay` - Accountid should be provided (bmc.testnet)
/// * `relay` - Accountid should be provided (bmc.testnet)
/// 
    pub fn assert_relay_is_valid(&self, accepted_relay: &AccountId, relay: &AccountId) {
        require!(
            relay == accepted_relay,
            format!(
                "{}",
                BmcError::Unauthorized {
                    message: "invalid relay"
                }
            )
        );
    }
///checks the given route is exists.
/// # Arguments
/// * `destination` - Should be in the form of btp://0x1.near/account.testnet
/// 
    pub fn assert_route_exists(&self, destination: &BTPAddress) {
        require!(
            self.routes.contains(destination),
            format!("{}", BmcError::RouteNotExist)
        );
    }
/// Checks whether the given route does not exists
/// # Arguments
/// * `destination` - Should be in the form of btp://0x1.near/account.testnet
/// 
    pub fn assert_route_does_not_exists(&self, destination: &BTPAddress) {
        require!(
            !self.routes.contains(destination),
            format!("{}", BmcError::RouteExist)
        );
    }
///Checks whether the sender service is authorised
/// # Arguments
/// * `service` - should be given as string.
    pub fn assert_sender_is_authorized_service(&self, service: &str) {
        require!(
            self.services.get(service) == Some(env::predecessor_account_id()),
            format!("{}", BmcError::PermissionNotExist)
        );
    }
/// Checks whether the service exists
/// # Arguments
/// * `name` - Service name should be given in the form of string
/// 
    pub fn assert_service_exists(&self, name: &str) {
        require!(
            self.services.contains(name),
            format!("{}", BmcError::ServiceNotExist)
        );
    }
/// Checks whether the service does not exists
/// # Arguments
/// * `name` - Service name should be given in the form of string
/// 
    pub fn assert_service_does_not_exists(&self, name: &str) {
        require!(
            !self.services.contains(name),
            format!("{}", BmcError::ServiceExist)
        );
    }


    ///Setting the link
    /// # Arguments
    /// * `max_aggregation` - should be in unsigned number
    /// * `delay_limit` - should be in the unsigned number
    /// 
    pub fn assert_valid_set_link_param(&self, max_aggregation: u64, delay_limit: u64) {
        require!(
            max_aggregation >= 1 && delay_limit >= 1,
            format!("{}", BmcError::InvalidParam)
        );
    }

    /// Checks that the relay does not exist
    /// # Arguments
    /// * `link` - It should be in the form of btp://0x1.near/account.testnet
    /// * `relay` - AccountId should be given
    /// 
    pub fn assert_relay_not_exists(&self, link: &BTPAddress, relay: &AccountId) {
        if let Some(link_property) = self.links.get(&link) {
            require!(
                !link_property.relays().contains(&relay),
                format!(
                    "{}",
                    BmcError::RelayExist {
                        link: link.to_string()
                    }
                )
            );
        }
    }
    /// Checks for the existing relay
    /// # Arguments
    /// * `link` - It should be in the form of btp://0x1.near/account.testnet
    /// * `relay` - AccountId should be given
    /// 
    pub fn assert_relay_exists(&self, link: &BTPAddress, relay: &AccountId) {
        if let Some(link_property) = self.links.get(&link).as_mut() {
            require!(
                link_property.relays().contains(&relay),
                format!(
                    "{}",
                    BmcError::RelayNotExist {
                        link: link.to_string()
                    }
                )
            );
        }
        }
    /// Checks for the existing service
    /// # Arguments 
    /// * `name` - service name should be given which is in string format
    /// 
    pub fn ensure_service_exists(&self, name: &str) -> Result<(), BmcError> {
        if !self.services.contains(name) {
            return Err(BmcError::ServiceNotExist);
        }
        Ok(())
    }
    /// Checks for the valid sequence
    /// # Arguments
    /// * `link`
    /// * `event`
    /// 
    pub fn ensure_valid_sequence(&self, link: &Link, event: &Event) -> Result<(), BmcError> {
        if link.rx_seq() != event.sequence() {
            return Err(BmcError::InvalidSequence);
        }
        Ok(())
    }
}
