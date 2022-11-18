use super::*;

#[derive(Serialize, Debug, Eq, PartialEq, Hash, Deserialize)]
pub struct Service {
    name: String,
    service: AccountId,
}

#[derive(BorshDeserialize, BorshSerialize)]
pub struct Services(TreeMap<String, AccountId>);

impl Services {
    pub fn new() -> Self {
        Self(TreeMap::new(StorageKey::Services))
    }

    pub fn add(&mut self, name: &str, service: &AccountId) {
        self.0.insert(&name.to_string(), &service.to_owned());
    }

    pub fn remove(&mut self, name: &str) {
        if self.contains(name) {
            self.0.remove(&name.to_string());
        }
    }

    pub fn contains(&self, name: &str) -> bool {
        self.0.contains_key(&name.to_string())
    }

    pub fn get(&self, name: &str) -> Option<AccountId> {
        self.0.get(&name.to_string())
    }

    pub fn to_vec(&self) -> Vec<Service> {
        if !self.0.is_empty() {
            return self
                .0
                .to_vec()
                .into_iter()
                .map(|v| Service {
                    name: v.0,
                    service: v.1,
                })
                .collect();
        }
        vec![]
    }
}

impl Default for Services {
    fn default() -> Self {
        Self::new()        
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashSet;

    #[test]
    fn add_service() {
        let name = String::from("service_1");
        let service = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let mut services = Services::new();
        services.add(&name, &service);
        let result = services.get(&name);
        assert_eq!(
            result,
            Some(
                "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .parse::<AccountId>()
                    .unwrap()
            )
        );
    }

    #[test]
    fn get_service() {
        let name_1 = String::from("service_1");
        let service_1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let name_2 = String::from("service_2");
        let service_2 = "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let name_3 = String::from("service_3");
        let service_3 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let mut services = Services::new();
        services.add(&name_1, &service_1);
        services.add(&name_2, &service_2);
        services.add(&name_3, &service_3);
        let result = services.get(&name_2);
        assert_eq!(
            result,
            Some(
                "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .parse::<AccountId>()
                    .unwrap()
            )
        );
    }

    #[test]
    fn remove_service() {
        let name_1 = String::from("service_1");
        let service_1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let name_2 = String::from("service_2");
        let service_2 = "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let name_3 = String::from("service_3");
        let service_3 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let mut services = Services::new();
        services.add(&name_1, &service_1);
        services.add(&name_2, &service_2);
        services.add(&name_3, &service_3);
        services.remove(&name_2);
        let result = services.get(&name_2);
        assert_eq!(result, None);
    }

    #[test]
    fn contains_service() {
        let name = String::from("service_1");
        let service = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let mut services = Services::new();
        services.add(&name, &service);
        let result = services.contains(&name);
        assert_eq!(result, true);
    }

    #[test]
    fn to_vec_service() {
        let name_1 = String::from("service_1");
        let service_1 = "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let name_2 = String::from("service_2");
        let service_2 = "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let name_3 = String::from("service_3");
        let service_3 = "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
            .parse::<AccountId>()
            .unwrap();
        let mut services = Services::new();
        services.add(&name_1, &service_1);
        services.add(&name_2, &service_2);
        services.add(&name_3, &service_3);
        let services = services.to_vec();
        let expected_services = vec![
            Service {
                name: "service_1".to_string(),
                service: "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .parse::<AccountId>()
                    .unwrap(),
            },
            Service {
                name: "service_2".to_string(),
                service: "68bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .parse::<AccountId>()
                    .unwrap(),
            },
            Service {
                name: "service_3".to_string(),
                service: "78bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4"
                    .parse::<AccountId>()
                    .unwrap(),
            },
        ];
        let result: HashSet<_> = services.iter().collect();
        let expected: HashSet<_> = expected_services.iter().collect();

        assert_eq!(result, expected);
    }
}
