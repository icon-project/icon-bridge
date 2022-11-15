use super::*;

impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * Utils * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *

    pub fn hash_token_id(token_name: &String) -> TokenId {
        env::sha256(token_name.as_bytes()).try_into().unwrap()
    }
}
