use super::*;

#[near_bindgen]
impl BtpMessageCenter {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Route Management  * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *


    /// Adding the route to bmc
    /// # Arguments
    /// * destination - Should give the btp address that should be in the form of btp://0x1.near/account.testnet
    /// * link - Link should be given which is to be added for the route.
    pub fn add_route(&mut self, destination: BTPAddress, link: BTPAddress) {
        self.assert_have_permission();
        self.assert_route_does_not_exists(&destination);
        self.assert_link_exists(&link);
        self.routes.add(&destination, &link);
        self.connections.add(
            &Connection::Route(destination.network_address().unwrap()),
            &link,
        );
    }

    /// Removing the route from bmc
    /// 
    /// # Arguments
    /// * destination - Should give the btp address that should be in the form of btp://0x1.near/account.testnet

    pub fn remove_route(&mut self, destination: BTPAddress) {
        self.assert_have_permission();
        self.assert_route_exists(&destination);
        let link = self.routes.get(&destination).unwrap_or_default();
        self.routes.remove(&destination);
        self.connections.remove(
            &Connection::Route(destination.network_address().unwrap()),
            &link,
        )
    }


    /// Returns the list of routes present in the bmc
    pub fn get_routes(&self) -> Value {
        self.routes.to_vec().into()
    }


    /// ```
    /// self.redolve_route(&destination)
    /// ```
    #[cfg(feature = "testable")]
    pub fn resolve_route_pub(&self, destination: BTPAddress) -> Option<BTPAddress> {
        self.resolve_route(&destination)
    }
}

impl BtpMessageCenter {

    /// Resolving the route in bmc
    /// # Arguments
    /// * destination - Should be in the form of btp://0x1.near/account.testnet
    /// return Some of destination address else None
    pub fn resolve_route(&self, destination: &BTPAddress) -> Option<BTPAddress> {
        if self.links.contains(destination) {
            Some(destination.clone())
        } else if self
            .connections
            .contains(&Connection::Link(destination.network_address().unwrap()))
        {
            self.connections
                .get(&Connection::Link(destination.network_address().unwrap()))
        } else if self
            .connections
            .contains(&Connection::Route(destination.network_address().unwrap()))
        {
            self.connections
                .get(&Connection::Route(destination.network_address().unwrap()))
        } else if self.connections.contains(&Connection::LinkReachable(
            destination.network_address().unwrap(),
        )) {
            self.connections.get(&Connection::LinkReachable(
                destination.network_address().unwrap(),
            ))
        } else {
            None
        }
    }
}
