use super::*;

#[derive(Debug, Eq, PartialEq, BorshDeserialize, BorshSerialize)]
pub struct Request {
    sender: String,
    receiver: String,
    assets: Vec<TransferableAsset>,
}

impl Request {
    pub fn new(sender: String, receiver: String, assets: Vec<TransferableAsset>) -> Self {
        Self {
            sender,
            receiver,
            assets,
        }
    }

    pub fn sender(&self) -> &String {
        &self.sender
    }

    pub fn receiver(&self) -> &String {
        &self.receiver
    }

    pub fn assets(&self) -> &Vec<TransferableAsset> {
        &self.assets
    }
}

#[derive(BorshDeserialize, BorshSerialize)]
pub struct Requests(UnorderedMap<i128, Request>);

impl Requests {
    pub fn new() -> Self {
        Self(UnorderedMap::new(StorageKey::Requests))
    }

    pub fn add(&mut self, serial_no: i128, request: &Request) {
        self.0.insert(&serial_no, request);
    }

    pub fn remove(&mut self, serial_no: i128) {
        self.0.remove(&serial_no);
    }

    pub fn get(&self, serial_no: i128) -> Option<Request> {
        if let Some(request) = self.0.get(&serial_no) {
            return Some(request);
        }
        None
    }

    pub fn contains(&self, serial_no: i128) -> bool {
        self.0.get(&serial_no).is_some()
    }
}

impl Default for Requests {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::types::TransferableAsset;

    #[test]
    fn add_request() {
        let mut requests = Requests::new();
        let request = Request::new(
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4".to_string(),
            "78bd0675686be0a5df7da33b6f1089eghea3769b19dbb2477fe0cd6e0f12667".to_string(),
            vec![TransferableAsset::new("ABC".to_string(), 100, 1)],
        );
        requests.add(1, &request);
        let result = requests.get(1).unwrap();
        assert_eq!(result, request);
    }

    #[test]
    fn add_request_existing() {
        let mut requests = Requests::new();
        let request = Request::new(
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4".to_string(),
            "78bd0675686be0a5df7da33b6f1089eghea3769b19dbb2477fe0cd6e0f12667".to_string(),
            vec![TransferableAsset::new("ABC".to_string(), 100, 1)],
        );
        requests.add(1, &request);
        requests.add(1, &request);
        let result = requests.get(1).unwrap();
        assert_eq!(result, request);
    }

    #[test]
    fn remove_request() {
        let mut requests = Requests::new();
        let request = Request::new(
            "88bd05442686be0a5df7da33b6f1089ebfea3769b19dbb2477fe0cd6e0f126e4".to_string(),
            "78bd0675686be0a5df7da33b6f1089eghea3769b19dbb2477fe0cd6e0f12667".to_string(),
            vec![TransferableAsset::new("ABC".to_string(), 100, 1)],
        );
        requests.add(1, &request);
        requests.remove(1);
        let result = requests.get(1);
        assert_eq!(result, None);
    }

    #[test]
    fn remove_request_non_existing() {
        let mut requests = Requests::new();
        requests.remove(1);
        let result = requests.get(1);
        assert_eq!(result, None);
    }
}
