use super::*;

impl BtpMessageCenter {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Interval Services * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    ///Handling the init process
    /// # Arguments
    /// * source - should be in the form of btp://0x1.near/account.testnet
    /// * links - Given by the vectors
    /// 
    pub fn handle_init(
        &mut self,
        source: &BTPAddress,
        links: &Vec<BTPAddress>,
    ) -> Result<(), BmcError> {
        if let Some(mut link) = self.links.get(source) {
            for source_link in links.iter() {
                // Add to Reachable list of the link
                link.reachable_mut().insert(source_link.to_owned());

                // Add to the connections for quickily quering for routing
                self.connections.add(
                    &Connection::LinkReachable(
                        source_link
                            .network_address()
                            .map_err(|error| BmcError::InvalidAddress { description: error })?,
                    ),
                    source,
                )
            }
            self.links.set(source, &link);
            Ok(())
        } else {
            Err(BmcError::LinkNotExist)
        }
    }
    /// Handling the links in bmc
    /// # Arguments
    /// * source - should be in the form of btp://0x1.near/account.testnet
    /// * souce_link - should be in the form of btp://0x1.near/account.testnet
    /// 
    pub fn handle_link(
        &mut self,
        source: &BTPAddress,
        source_link: &BTPAddress,
    ) -> Result<(), BmcError> {
        if let Some(mut link) = self.links.get(source) {
            if !link.reachable().contains(source_link) {
                link.reachable_mut().insert(source_link.to_owned());

                // Add to the connections for quickily quering for routing
                self.connections.add(
                    &Connection::LinkReachable(
                        source_link
                            .network_address()
                            .map_err(|error| BmcError::InvalidAddress { description: error })?,
                    ),
                    source,
                );
            }

            self.links.set(source, &link);
            Ok(())
        } else {
            Err(BmcError::LinkNotExist)
        }
    }

     /// Handling the unlinks in bmc
    /// # Arguments
    /// * source - should be in the form of btp://0x1.near/account.testnet
    /// * souce_link - should be in the form of btp://0x1.near/account.testnet
    ///  
    
    pub fn handle_unlink(
        &mut self,
        source: &BTPAddress,
        source_link: &BTPAddress,
    ) -> Result<(), BmcError> {
        if let Some(mut link) = self.links.get(source) {
            if link.reachable().contains(source_link) {
                link.reachable_mut().remove(source_link);

                // Remove from the connections for quickily quering for routing
                self.connections.remove(
                    &Connection::LinkReachable(
                        source_link
                            .network_address()
                            .map_err(|error| BmcError::InvalidAddress { description: error })?,
                    ),
                    source,
                );
            }

            self.links.set(source, &link);
            Ok(())
        } else {
            Err(BmcError::LinkNotExist)
        }
    }
    
    /// Handle fee gathering in bmc
    /// # Arguments
    /// * source - should be in the form of btp://0x1.near/account.testnet
    /// * fee_aggregator - should be in the form of btp://0x1.near/account.testnet
    /// * services - should be given in vector with the format of string.
    /// 
    pub fn handle_fee_gathering(
        &self,
        source: &BTPAddress,
        fee_aggregator: &BTPAddress,
        services: &Vec<String>,
    ) -> Result<(), BmcError> {
        if source.network_address() != fee_aggregator.network_address() {
            return Err(BmcError::FeeAggregatorNotAllowed {
                source: source.to_string(),
            });
        }

        services.iter().for_each(|service| {
            //TODO: Handle Services that are not available
            #[allow(unused_variables)]
            if let Some(account_id) = self.services.get(service) {
                #[cfg(not(feature = "testable"))]
                bsh_contract::ext(account_id.clone()).handle_fee_gathering(
                    fee_aggregator.clone(),
                    service.clone(),
                );
            }
        });
        Ok(())
    }
}
