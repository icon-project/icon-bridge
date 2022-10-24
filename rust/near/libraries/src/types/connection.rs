use super::BTPAddress;
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::collections::LookupMap;
use std::collections::HashSet;

#[derive(BorshDeserialize, BorshSerialize, Eq, PartialEq, PartialOrd, Hash, Clone)]
pub enum Connection {
    Link(String),
    Route(String),
    LinkReachable(String),
}

#[derive(BorshDeserialize, BorshSerialize)]
pub struct Connections(LookupMap<Connection, HashSet<BTPAddress>>);

impl Connections {
    pub fn new() -> Self {
        Self(LookupMap::new(b"connections".to_vec()))
    }

    pub fn add(&mut self, connection: &Connection, link: &BTPAddress) {
        let mut list = self.0.get(connection).unwrap_or_default();
        list.insert(link.to_owned());
        self.0.insert(&connection, &list);
    }

    pub fn remove(&mut self, connection: &Connection, link: &BTPAddress) {
        let mut list = self.0.get(connection).unwrap_or_default();
        list.remove(link);

        if list.is_empty() {
            self.0.remove(connection);
        } else {
            self.0.insert(connection, &list);
        }
    }

    pub fn get(&self, connection: &Connection) -> Option<BTPAddress> {
        self.0
            .get(connection)
            .unwrap_or_default()
            .into_iter()
            .last()
    }

    pub fn contains(&self, connection: &Connection) -> bool {
        self.0.contains_key(connection)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::types::Address;

    #[test]
    fn add_connection() {

        let destination = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let link = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126ea"
                .to_string(),
        );
        let mut connections = Connections::new();
        connections.add(
            &Connection::Route(destination.network_address().unwrap()),
            &link,
        );
        let link = connections.get(&Connection::Route(destination.network_address().unwrap()));
        assert_eq!(
            link,
            Some(BTPAddress::new(
                "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126ea"
                    .to_string(),
            ))
        );
    }

    #[test]
    fn remove_connection() {

        let destination = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let link_1 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126ea"
                .to_string(),
        );
        let link_2 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126eb"
                .to_string(),
        );
        let mut connections = Connections::new();
        connections.add(
            &Connection::Route(destination.network_address().unwrap()),
            &link_1,
        );
        connections.add(
            &Connection::LinkReachable(destination.network_address().unwrap()),
            &link_1,
        );
        connections.add(
            &Connection::LinkReachable(destination.network_address().unwrap()),
            &link_2,
        );
        connections.remove(&Connection::Route(destination.network_address().unwrap()), &link_1);
        let link = connections.get(&Connection::Route(destination.network_address().unwrap()));
        assert_eq!(
            link,
            None
        );
    }

    #[test]
    fn contains_connection() {

        let destination = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let link_1 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126ea"
                .to_string(),
        );
        let link_2 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126eb"
                .to_string(),
        );
        let mut connections = Connections::new();
        connections.add(
            &Connection::Route(destination.network_address().unwrap()),
            &link_1,
        );
        connections.add(
            &Connection::LinkReachable(destination.network_address().unwrap()),
            &link_1,
        );
        connections.add(
            &Connection::LinkReachable(destination.network_address().unwrap()),
            &link_2,
        );
        connections.remove(&Connection::Route(destination.network_address().unwrap()), &link_1);
        let result = connections.contains(&Connection::Route(destination.network_address().unwrap()));
        assert_eq!(result, false);

        connections.remove(&Connection::LinkReachable(destination.network_address().unwrap()), &link_1);
        let result = connections.contains(&Connection::LinkReachable(destination.network_address().unwrap()));
        assert_eq!(result, true);
    }

    #[test]
    fn get_connection() {

        let destination = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let link_1 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126ea"
                .to_string(),
        );
        let link_2 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126eb"
                .to_string(),
        );
        let mut connections = Connections::new();
        connections.add(
            &Connection::Route(destination.network_address().unwrap()),
            &link_1,
        );
        connections.add(
            &Connection::LinkReachable(destination.network_address().unwrap()),
            &link_2,
        );
        let link = connections.get(&Connection::Route(destination.network_address().unwrap()));
        assert_eq!(
            link,
            Some(BTPAddress::new(
                "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126ea"
                    .to_string(),
            ))
        );

        let link = connections.get(&Connection::LinkReachable(destination.network_address().unwrap()));
        assert_eq!(
            link,
            Some(BTPAddress::new(
                "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126eb"
                    .to_string(),
            ))
        );
    }

}
