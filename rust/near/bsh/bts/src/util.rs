use super::*;

impl BtpTokenService {
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * Utils * * * * * * *
    // * * * * * * * * * * * * * * * * *
    // * * * * * * * * * * * * * * * * *
/// Hash coin id funcion got created with the coin name as arguments
/// Returns the coinId
    pub fn hash_coin_id(coin_name: &String) -> CoinId {
        env::sha256(coin_name.as_bytes()).try_into().unwrap()
    }
}
