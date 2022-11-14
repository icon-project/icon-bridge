use super::*;

#[near_bindgen]
impl BtpMessageCenter {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * Link Management * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    pub fn add_link(&mut self, link: BTPAddress) {
        self.assert_have_permission();
        self.assert_link_does_not_exists(&link);

        self.propogate_internal(BmcServiceMessage::new(BmcServiceType::Link {
            link: link.clone(),
        }));

        self.links.add(&link, self.block_interval);
        self.connections
            .add(&Connection::Link(link.network_address().unwrap()), &link);

        self.send_internal_service_message(
            &link,
            &BmcServiceMessage::new(BmcServiceType::Init {
                links: self.links.to_vec(),
            }),
        );
    }

    pub fn remove_link(&mut self, link: BTPAddress) {
        self.assert_have_permission();
        self.assert_link_exists(&link);
        self.assert_link_does_not_have_route_connection(&link);

        self.propogate_internal(BmcServiceMessage::new(BmcServiceType::Unlink {
            link: link.clone(),
        }));

        self.links.remove(&link);
        self.connections
            .remove(&Connection::Link(link.network_address().unwrap()), &link);
    }

    #[cfg(feature = "testable")]
    pub fn get_reachable_link(&self, link: BTPAddress) -> HashedCollection<BTPAddress> {
        if let Some(link) = self.links.get(&link) {
            return link
                .reachable()
                .to_owned()
                .into_iter()
                .collect::<HashedCollection<BTPAddress>>();
        }
        HashedCollection::new()
    }

    pub fn get_links(&self) -> serde_json::Value {
        self.links.to_vec().into()
    }

    pub fn set_link(
        &mut self,
        link: BTPAddress,
        block_interval: u64,
        max_aggregation: u64,
        delay_limit: u64,
    ) {
        self.assert_have_permission();
        self.assert_link_exists(&link);
        self.assert_valid_set_link_param(max_aggregation, delay_limit);
        if let Some(link_property) = self.links.get(&link).as_mut() {
            let previous_rotate_term = link_property.rotate_term();

            link_property
                .block_interval_dst_mut()
                .clone_from(&block_interval);
            link_property
                .max_aggregation_mut()
                .clone_from(&max_aggregation);
            link_property.delay_limit_mut().clone_from(&delay_limit);

            let current_rotate_term = link_property.rotate_term();
            if previous_rotate_term == 0 && current_rotate_term > 0 {
                link_property
                    .rotate_height_mut()
                    .clone_from(&(env::block_height() + current_rotate_term));
            }
            self.links.set(&link, link_property);
        }
    }

    pub fn get_status(&self, link: BTPAddress) -> LinkStatus {
        self.assert_link_exists(&link);
        self.links.get(&link).unwrap().status()
    }

    pub fn set_link_rx_height(&mut self, link: BTPAddress, height: u64) {
        self.assert_have_permission();
        self.assert_link_exists(&link);

        if let Some(link_property) = self.links.get(&link).as_mut() {
            link_property.rx_height_mut().clone_from(&height);
            self.links.set(&link, &link_property);
        }
    }
}

impl BtpMessageCenter {
    pub fn increment_link_rx_seq(&mut self, link: &BTPAddress) {
        if let Some(link_property) = self.links.get(link).as_mut() {
            link_property.rx_seq_mut().add(1).unwrap();
            self.links.set(link, link_property);
        }
    }

    pub fn get_link(&self, link: BTPAddress) -> Link {
        self.links.get(&link).unwrap()
    }
}
