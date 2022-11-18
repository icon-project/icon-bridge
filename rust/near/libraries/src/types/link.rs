use super::{relay::BmrStatus, *};

#[derive(Debug, Default, BorshDeserialize, BorshSerialize, Eq, PartialEq)]
pub struct Link {
    rx_seq: u128,
    tx_seq: u128,
    relays: Relays,
    reachable: HashSet<BTPAddress>,
    relay_index: u64,
    rotate_height: u64,
    delay_limit: u64,
    max_aggregation: u64,
    rx_height_src: u64,
    rx_height: u64,
    block_interval_dst: u64,
    block_interval_src: u64,
}

impl Link {
    pub fn tx_seq(&self) -> u128 {
        self.tx_seq
    }

    pub fn tx_seq_mut(&mut self) -> &mut u128 {
        &mut self.tx_seq
    }

    pub fn rx_seq(&self) -> u128 {
        self.rx_seq
    }

    pub fn rx_seq_mut(&mut self) -> &mut u128 {
        &mut self.rx_seq
    }

    pub fn reachable(&self) -> &HashSet<BTPAddress> {
        &self.reachable
    }

    pub fn reachable_mut(&mut self) -> &mut HashSet<BTPAddress> {
        &mut self.reachable
    }

    pub fn block_interval_dst(&self) -> u64 {
        self.block_interval_dst
    }

    pub fn block_interval_dst_mut(&mut self) -> &mut u64 {
        &mut self.block_interval_dst
    }

    pub fn max_aggregation(&self) -> u64 {
        self.max_aggregation
    }

    pub fn max_aggregation_mut(&mut self) -> &mut u64 {
        &mut self.max_aggregation
    }

    pub fn delay_limit(&mut self) -> u64 {
        self.delay_limit
    }

    pub fn delay_limit_mut(&mut self) -> &mut u64 {
        &mut self.delay_limit
    }

    pub fn relays_mut(&mut self) -> &mut Relays {
        self.relays.as_mut()
    }

    pub fn relays(&self) -> &Relays {
        &self.relays
    }

    fn relay_index(&self) -> u64 {
        self.relay_index
    }

    fn relay_index_mut(&mut self) -> &mut u64 {
        &mut self.relay_index
    }

    pub fn rotate_height(&self) -> u64 {
        self.rotate_height
    }

    pub fn rotate_height_mut(&mut self) -> &mut u64 {
        &mut self.rotate_height
    }

    pub fn rx_height(&self) -> u64 {
        self.rx_height
    }

    pub fn rx_height_mut(&mut self) -> &mut u64 {
        &mut self.rx_height
    }

    pub fn rx_height_src(&self) -> u64 {
        self.rx_height_src
    }

    pub fn rx_height_src_mut(&mut self) -> &mut u64 {
        &mut self.rx_height_src
    }

    fn scale(&self) -> u64 {
        if self.block_interval_src < 1 || self.block_interval_dst < 1 {
            0
        } else {
            *self
                .block_interval_src
                .clone()
                .div_ceil(self.block_interval_dst)
        }
    }

    pub fn rotate_term(&self) -> u64 {
        let scale = self.scale();
        if scale > 0 {
            *self.max_aggregation.clone().div_ceil(scale)
        } else {
            0
        }
    }

    pub fn rotate_relay(
        &mut self,
        last_height: BlockHeight,
        has_message: bool,
    ) -> Option<&AccountId> {
        let rotate_term = self.rotate_term();
        let current_height = env::block_height();
        if rotate_term > 0 {
            let (rotate_count, base_height) = match has_message {
                true => {
                    let mut guess_height = (self.rx_height
                        + ((last_height.checked_sub(self.rx_height_src()).unwrap())
                            .div_ceil(self.scale()))
                        .deref())
                    .checked_sub(1)
                    .unwrap();
                    if guess_height > current_height {
                        guess_height = current_height;
                    };

                    let mut rotate_count = {
                        let mut count = if guess_height < self.rotate_height {
                            0
                        } else {
                            guess_height - self.rotate_height
                        };
                        let rotate_count = count.div_ceil(rotate_term);
                        *rotate_count.deref()
                    };
                    #[allow(unused)]
                    let mut base_height: u64 = 0;
                    if rotate_count > 0_u64 {
                        base_height = self.rotate_height + ((rotate_count - 1) * rotate_term);
                    } else {
                        base_height = self.rotate_height - rotate_term;
                    }

                    let mut skip_count = (current_height - guess_height)
                        .div_ceil(self.delay_limit)
                        .deref()
                        .to_owned();
                    if skip_count > 0 {
                        skip_count -= 1;
                        rotate_count.add(skip_count).unwrap();
                        base_height = current_height;
                    }
                    self.rx_height.clone_from(&current_height);
                    self.rx_height_src_mut().clone_from(&last_height);
                    (rotate_count, base_height)
                }
                false => {
                    let mut count = if current_height < self.rotate_height {
                        0
                    } else {
                        current_height - self.rotate_height
                    };
                    let rotate_count = count.div_ceil(rotate_term);
                    #[allow(unused)]
                    let mut base_height: u64 = 0;
                    if *rotate_count > 0_u64 {
                        base_height = self.rotate_height + ((*rotate_count - 1) * rotate_term);
                    } else {
                        base_height = self.rotate_height - rotate_term;
                    }
                    (*rotate_count, base_height)
                }
            };

            if rotate_count > 0 {
                self.rotate_height_mut()
                    .add(base_height + rotate_term)
                    .unwrap();
                self.relay_index_mut().add(rotate_count).unwrap();
                if self.relay_index() >= self.relays().len().try_into().unwrap() {
                    let relays_count: u64 = self.relays.len().try_into().unwrap();
                    let relay_index = self.relay_index();
                    self.relay_index_mut()
                        .clone_from(&(relay_index % relays_count));
                }
            }
            self.relays().get(self.relay_index())
        } else {
            None
        }
    }

    //TODO: Confirm if relay status is linked with link
    pub fn status(&self) -> LinkStatus {
        LinkStatus {
            rx_seq: self.rx_seq - 1,
            tx_seq: self.tx_seq,
            relays: self.relays.bmr_status(),
            relay_index: self.relay_index,
            rotate_height: self.rotate_height,
            rotate_term: self.rotate_term(),
            delay_limit: self.delay_limit,
            rx_height_src: self.rx_height_src,
            rx_height: self.rx_height,
            block_interval_dst: self.block_interval_dst,
            block_interval_src: self.block_interval_src,
            current_height: env::block_height(),
            max_aggregation: self.max_aggregation(),
        }
    }
}

#[derive(Clone, Deserialize, Serialize, Debug, PartialEq, Eq)]
#[serde(crate = "near_sdk::serde")]
pub struct LinkStatus {
    rx_seq: u128,
    tx_seq: u128,
    relays: Vec<BmrStatus>,
    relay_index: u64,
    rotate_height: u64,
    rotate_term: u64,
    delay_limit: u64,
    rx_height_src: u64,
    rx_height: u64,
    block_interval_dst: u64,
    block_interval_src: u64,
    current_height: BlockHeight,
    max_aggregation: u64,
}

impl LinkStatus {
    pub fn delay_limit(&self) -> u64 {
        self.delay_limit
    }

    pub fn max_aggregation(&self) -> u64 {
        self.max_aggregation
    }

    pub fn block_interval_dst(&self) -> u64 {
        self.block_interval_dst
    }

    pub fn rx_seq(&self) -> u128 {
        self.rx_seq
    }

    pub fn rx_height(&self) -> u64 {
        self.rx_height
    }
}

#[derive(BorshDeserialize, BorshSerialize)]

pub struct Links {
    keys: UnorderedSet<BTPAddress>,
    values: LookupMap<BTPAddress, Link>,
}

impl Deref for Links {
    type Target = LookupMap<BTPAddress, Link>;

    fn deref(&self) -> &Self::Target {
        &self.values
    }
}

impl DerefMut for Links {
    fn deref_mut(&mut self) -> &mut Self::Target {
        &mut self.values
    }
}

impl Links {
    pub fn new() -> Self {
        Self {
            keys: UnorderedSet::new(StorageKey::Links(KeyType::Key)),
            values: LookupMap::new(StorageKey::Links(KeyType::Value)),
        }
    }

    pub fn add(&mut self, link: &BTPAddress, block_interval_src: u64) {
        self.keys.insert(link);
        self.values.insert(
            link,
            &Link {
                relays: Relays::new(link),
                block_interval_src,
                rx_seq: 1,
                ..Default::default()
            },
        );
    }

    pub fn set(&mut self, link: &BTPAddress, property: &Link) {
        self.values.insert(link, property);
    }

    pub fn remove(&mut self, link: &BTPAddress) {
        self.keys.remove(link);
        self.values.remove(link);
    }

    pub fn to_vec(&self) -> Vec<BTPAddress> {
        self.keys.to_vec()
    }

    pub fn get(&self, link: &BTPAddress) -> Option<Link> {
        if let Some(value) = self.values.get(link) {
            return Some(value);
        }
        None
    }

    pub fn contains(&self, link: &BTPAddress) -> bool {
        self.keys.contains(link)
    }
}

impl Default for Links {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use std::vec;

    use super::*;
    use near_sdk::AccountId;

    #[test]
    fn add_link() {
        let link = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut links = Links::new();
        links.add(&link, 0);
        let expected = Link {
            rx_seq: 1,
            ..Default::default()
        };
        assert_eq!(links.get(&link).unwrap(), expected);
    }

    #[test]
    fn add_link_relays_pass() {
        let link_1 = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut links = Links::new();
        links.add(&link_1, 0);

        if let Some(link) = links.get(&link_1).as_mut() {
            link.relays.add(
                &"88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .parse::<AccountId>()
                    .unwrap(),
            );
            links.set(&link_1, &link);
        }
        let mut expected = Link {
            rx_seq: 1,
            ..Default::default()
        };
        expected.relays.add(
            &"88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        );
        assert_eq!(links.get(&link_1).unwrap(), expected);
    }

    #[test]
    fn set_link_block_interval_dst() {
        let link_1 = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut links = Links::new();
        links.add(&link_1, 0);
        if let Some(link) = links.get(&link_1).as_mut() {
            link.block_interval_dst_mut().clone_from(&1000);
            links.set(&link_1, &link);
        }
        let expected = Link {
            rx_seq: 1,
            block_interval_dst: 1000,
            ..Default::default()
        };
        assert_eq!(links.get(&link_1).unwrap(), expected);
    }

    #[test]
    fn set_link_max_aggregation_src() {
        let link_1 = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut links = Links::new();
        links.add(&link_1, 0);
        if let Some(link) = links.get(&link_1).as_mut() {
            link.max_aggregation_mut().clone_from(&10);
            links.set(&link_1, &link);
        }
        let expected = Link {
            rx_seq: 1,
            max_aggregation: 10,
            ..Default::default()
        };
        assert_eq!(links.get(&link_1).unwrap(), expected);
    }

    #[test]
    fn set_link_delay_limit_src() {
        let link_1 = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut links = Links::new();
        links.add(&link_1, 0);
        if let Some(link) = links.get(&link_1).as_mut() {
            link.delay_limit_mut().clone_from(&100);
            links.set(&link_1, &link);
        }
        let expected = Link {
            rx_seq: 1,
            delay_limit: 100,
            ..Default::default()
        };
        assert_eq!(links.get(&link_1).unwrap(), expected);
    }

    #[test]
    fn set_link_relays_pass() {
        let link_1 = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut links = Links::new();
        links.add(&link_1, 0);

        if let Some(link) = links.get(&link_1).as_mut() {
            link.relays.set(&vec![
                "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e2"
                    .parse::<AccountId>()
                    .unwrap(),
                "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e3"
                    .parse::<AccountId>()
                    .unwrap(),
                "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .parse::<AccountId>()
                    .unwrap(),
            ]);
            links.set(&link_1, &link);
        }
        let mut expected = Link {
            rx_seq: 1,
            ..Default::default()
        };
        expected.relays.set(&vec![
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e2"
                .parse::<AccountId>()
                .unwrap(),
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e3"
                .parse::<AccountId>()
                .unwrap(),
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        ]);
        assert_eq!(links.get(&link_1).unwrap(), expected);
    }

    #[test]
    fn remove_link() {
        let mut links = Links::new();
        let link = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        links.add(&link, 0);
        links.remove(&link);
        let result = links.contains(&link);
        assert_eq!(result, false);
    }

    #[test]
    fn remove_link_non_existing() {
        let mut links = Links::new();
        let link_1 = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let link_2 = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        links.add(&link_1, 0);
        links.remove(&link_2);
        let result = links.contains(&link_2);
        assert_eq!(result, false);
    }

    #[test]
    fn to_vec_links() {
        let mut links = Links::new();
        let link_1 = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let link_2 = BTPAddress::new(
            "btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let link_3 = BTPAddress::new(
            "btp://0x1.iconee/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        links.add(&link_1, 0);
        links.add(&link_2, 0);
        links.add(&link_3, 0);
        let result = links.to_vec();
        let expected: Vec<BTPAddress> = vec![link_1, link_2, link_3];
        assert_eq!(result, expected);
    }
}
