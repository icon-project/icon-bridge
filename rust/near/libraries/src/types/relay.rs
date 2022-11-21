use super::*;

#[derive(BorshDeserialize, BorshSerialize)]
pub struct Relays(Vec<AccountId>, UnorderedMap<AccountId, RelayStatus>);

impl Default for Relays {
    fn default() -> Relays {
        let relays: Vec<AccountId> = Vec::new();
        let relay_status: UnorderedMap<AccountId, RelayStatus> =
            UnorderedMap::new(keccak256("relay_status".to_string().as_bytes()));
        Self(relays, relay_status)
    }
}

impl Eq for Relays {}
impl PartialEq for Relays {
    fn eq(&self, other: &Self) -> bool {
        self.0.to_vec() == other.0.to_vec()
    }
}

impl std::fmt::Debug for Relays {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_list().entries(self.0.iter()).finish()
    }
}

impl AsMut<Relays> for Relays {
    fn as_mut(&mut self) -> &mut Self {
        self
    }
}

impl Relays {
    pub fn new(link: &BTPAddress) -> Self {
        Self(
            Vec::new(),
            UnorderedMap::new(keccak256(format!("{}_relay_status", link).as_bytes())),
        )
    }

    pub fn add(&mut self, account_id: &AccountId) {
        if !self.contains(account_id) {
            self.0.push(account_id.to_owned())
        }
    }

    pub fn get(&self, index: u64) -> Option<&AccountId> {
        self.0.get(index as usize)
    }

    pub fn set(&mut self, account_ids: &[AccountId]) {
        self.clear();

        for account_id in account_ids {
            self.add(account_id);
        }
    }

    pub fn status(&self, account_id: &AccountId) -> Option<RelayStatus> {
        self.1.get(account_id)
    }

    pub fn set_status(&mut self, account_id: &AccountId, relay_status: &RelayStatus) {
        self.1.insert(account_id, relay_status);
    }

    pub fn clear(&mut self) {
        self.0.clear();
    }

    pub fn contains(&self, account_id: &AccountId) -> bool {
        self.0.contains(account_id)
    }

    pub fn remove(&mut self, account_id: &AccountId) {
        let index = self.0.iter().position(|item| item == account_id);
        if let Some(index) = index {
            self.0.swap_remove(index);
        }
    }

    pub fn len(&self) -> usize {
        self.0.len()
    }

    pub fn to_vec(&self) -> Vec<AccountId> {
        self.0.to_vec()
    }

    pub fn bmr_status(&self) -> Vec<BmrStatus> {
        self.to_vec()
            .iter()
            .map(|relay| {
                let relay_status = self.status(relay).unwrap_or_default();
                {
                    BmrStatus {
                        account_id: relay.to_owned(),
                        block_count: relay_status.block_count,
                        message_count: relay_status.message_count,
                    }
                }
            })
            .collect()
    }
    pub fn is_empty(&self) -> bool {
        self.0.is_empty()
    }
}
#[derive(Clone, Deserialize, Serialize, Debug, PartialEq, Eq)]
#[serde(crate = "near_sdk::serde")]
pub struct BmrStatus {
    account_id: AccountId,
    block_count: u64,
    message_count: u64,
}

#[derive(Default, BorshDeserialize, BorshSerialize)]
pub struct RelayStatus {
    block_count: u64,
    message_count: u64,
}

impl RelayStatus {
    pub fn new(block_count: u64, message_count: u64) -> RelayStatus {
        RelayStatus {
            block_count,
            message_count,
        }
    }

    pub fn block_count(&self) -> u64 {
        self.block_count
    }

    pub fn block_count_mut(&mut self) -> &mut u64 {
        &mut self.block_count
    }

    pub fn message_count(&self) -> u64 {
        self.message_count
    }

    pub fn message_count_mut(&mut self) -> &mut u64 {
        &mut self.message_count
    }
}
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn add_relay() {
        let link = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut relays = Relays::new(&link);
        let relay = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        relays.add(&relay);
        let result = relays.contains(&relay);
        assert_eq!(result, true);
    }

    #[test]
    fn add_existing_relay() {
        let link = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut relays = Relays::new(&link);
        let relay_1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let relay_1_duplicate = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        relays.add(&relay_1);
        relays.add(&relay_1_duplicate);
        let result = relays.to_vec();
        let expected: Vec<AccountId> = vec![
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        ];
        assert_eq!(result, expected);
    }

    #[test]
    fn remove_relay() {
        let link = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut relays = Relays::new(&link);
        let relay = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        relays.add(&relay);
        relays.remove(&relay);
        let result = relays.contains(&relay);
        assert_eq!(result, false);
    }

    #[test]
    fn remove_relay_non_existing() {
        let link = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut relays = Relays::new(&link);
        let relay_1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e8"
            .parse::<AccountId>()
            .unwrap();
        let relay_2 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        relays.add(&relay_1);
        relays.remove(&relay_2);
        let result = relays.contains(&relay_2);
        assert_eq!(result, false);
    }

    #[test]
    fn clear_relays() {
        let link = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut relays = Relays::new(&link);
        let relay = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        relays.add(&relay);
        relays.clear();
        let result = relays.to_vec();
        let expected: Vec<AccountId> = vec![];
        assert_eq!(result, expected);
    }

    #[test]
    fn to_vec_relays() {
        let link = BTPAddress::new(
            "btp://0x1.near/cx87ed9048b594b95199f326fc76e76a9d33dd665b".to_string(),
        );
        let mut relays = Relays::new(&link);
        let relay_1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let relay_2 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let relay_3 = "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        relays.add(&relay_1);
        relays.add(&relay_2);
        relays.add(&relay_3);
        let result = relays.to_vec();
        let expected: Vec<AccountId> = vec![
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
            "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
            "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                .parse::<AccountId>()
                .unwrap(),
        ];
        assert_eq!(result, expected);
    }
}
