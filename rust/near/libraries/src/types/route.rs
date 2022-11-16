use super::{Address, BTPAddress};
use near_sdk::borsh::{self, BorshDeserialize, BorshSerialize};
use near_sdk::collections::{LookupMap, UnorderedMap, UnorderedSet};
use near_sdk::serde::{Deserialize, Serialize};
use near_sdk::serde_json::{json, Value};
use std::collections::HashMap;
use std::collections::HashSet;

#[derive(Serialize, Deserialize, Debug, Eq, PartialEq, Hash)]
pub struct Route {
    destination: BTPAddress,
    next: BTPAddress,
}

impl From<Route> for Value {
    fn from(route: Route) -> Self {
        json!({
            "dst": route.destination,
            "next": route.next
        })
    }
}

#[derive(BorshDeserialize, BorshSerialize)]
pub struct Routes {
    keys: UnorderedSet<String>,
    values: LookupMap<String, HashMap<BTPAddress, BTPAddress>>,
}

impl Routes {
    pub fn new() -> Self {
        Self {
            keys: UnorderedSet::new(b"route_keys".to_vec()),
            values: LookupMap::new(b"route_values".to_vec()),
        }
    }

    pub fn add(&mut self, destination: &BTPAddress, link: &BTPAddress) {
        self.keys.insert(&destination.network_address().unwrap());
        let mut list = self
            .values
            .get(&destination.network_address().unwrap())
            .unwrap_or_default();
        list.insert(destination.to_owned(), link.to_owned());
        self.values
            .insert(&destination.network_address().unwrap(), &list);
    }

    pub fn remove(&mut self, destination: &BTPAddress) {
        self.keys.remove(&destination.network_address().unwrap());
        let mut list = self
            .values
            .get(&destination.network_address().unwrap())
            .unwrap_or_default();
        list.remove(&destination);

        if list.is_empty() {
            self.values.remove(&destination.network_address().unwrap());
        } else {
            self.values
                .insert(&destination.network_address().unwrap(), &list);
        }
    }

    pub fn get(&self, destination: &BTPAddress) -> Option<BTPAddress> {
        let list = self
            .values
            .get(&destination.network_address().unwrap())
            .unwrap_or_default();
        list.get(destination).map(|link| link.to_owned())
    }

    pub fn contains_network(&self, network: &str) -> bool {
        self.keys.contains(&network.to_string())
    }

    pub fn contains(&self, destination: &BTPAddress) -> bool {
        let list = self
            .values
            .get(&destination.network_address().unwrap())
            .unwrap_or_default();
        list.contains_key(destination)
    }

    pub fn to_vec(&self) -> Vec<Route> {
        let mut routes: HashSet<Route> = HashSet::new();
        if !self.keys.is_empty() {
            self.keys.iter().for_each(|network| {
                let values = self.values.get(&network).unwrap();
                values.into_iter().for_each(|(destination, next)| {
                    routes.insert(Route { destination, next });
                });
            });
        }
        routes.into_iter().collect()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashSet;

    #[test]
    fn add_route() {
        let destination = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let link = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        );
        let mut routes = Routes::new();
        routes.add(&destination, &link);
        let link = routes.get(&destination);
        assert_eq!(
            link,
            Some(BTPAddress::new(
                "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .to_string(),
            ))
        );
    }

    #[test]
    fn get_route() {
        let destination_1 = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let next_1 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        );
        let destination_2 =
            BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
        let next_2 = BTPAddress::new(
            "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let destination_3 =
            BTPAddress::new("btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
        let next_3 = BTPAddress::new(
            "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut routes = Routes::new();
        routes.add(&destination_1, &next_1);
        routes.add(&destination_2, &next_2);
        routes.add(&destination_3, &next_3);
        let result = routes.get(&destination_2);
        assert_eq!(
            result,
            Some(BTPAddress::new(
                "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
            ))
        );
    }

    #[test]
    fn remove_route() {
        let destination_1 = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let next_1 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        );
        let destination_2 =
            BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
        let next_2 = BTPAddress::new(
            "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let destination_3 =
            BTPAddress::new("btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
        let next_3 = BTPAddress::new(
            "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut routes = Routes::new();
        routes.add(&destination_1, &next_1);
        routes.add(&destination_2, &next_2);
        routes.add(&destination_3, &next_3);
        routes.remove(&destination_2);
        let links = routes.get(&destination_2);
        assert_eq!(links, None);
    }

    #[test]
    fn contains_route() {
        let destination_1 = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let next_1 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        );
        let destination_2 =
            BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
        let next_2 = BTPAddress::new(
            "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let destination_3 =
            BTPAddress::new("btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
        let next_3 = BTPAddress::new(
            "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut routes = Routes::new();
        routes.add(&destination_1, &next_1);
        routes.add(&destination_2, &next_2);
        routes.add(&destination_3, &next_3);
        let result = routes.contains_network(&destination_1.network_address().unwrap());
        assert_eq!(result, true);
        routes.remove(&destination_1);
        let result = routes.contains_network(&destination_1.network_address().unwrap());
        assert_eq!(result, false);
        routes.remove(&destination_2);
        let result = routes.contains_network(&destination_3.network_address().unwrap());
        assert_eq!(result, true);
    }

    #[test]
    fn contains_network() {
        let dst = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let next = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        );
        let mut routes = Routes::new();
        routes.add(&dst, &next);
        let result = routes.contains_network(&dst.network_address().unwrap());
        assert_eq!(result, true);
    }

    #[test]
    fn to_vec_route() {
        let destination_1 = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let next_1 = BTPAddress::new(
            "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .to_string(),
        );
        let destination_2 =
            BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
        let next_2 = BTPAddress::new(
            "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let destination_3 =
            BTPAddress::new("btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string());
        let next_3 = BTPAddress::new(
            "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut routes = Routes::new();
        routes.add(&destination_1, &next_1);
        routes.add(&destination_2, &next_2);
        routes.add(&destination_3, &next_3);
        let routes = routes.to_vec();
        let expected_routes = vec![
            Route {
                destination: BTPAddress::new(
                    "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
                ),
                next:
                BTPAddress::new(
                    "btp://0x1.bsc/88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                        .to_string(),
                ),
            },
            Route {
                destination: BTPAddress::new("btp://0x1.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
                next: BTPAddress::new(
                    "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
                ),
            },
            Route {
                destination: BTPAddress::new("btp://0x5.pra/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string()),
                next: BTPAddress::new(
                    "btp://0x3.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
                ),
            },
        ];
        let result: HashSet<_> = routes.iter().collect();
        let expected: HashSet<_> = expected_routes.iter().collect();
        assert_eq!(result, expected);
    }
}
