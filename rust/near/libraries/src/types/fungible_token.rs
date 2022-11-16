use super::*;

#[derive(BorshDeserialize, BorshSerialize, Clone, Debug, Deserialize, Serialize, Eq, PartialEq)]
#[serde(crate = "near_sdk::serde")]
pub struct AssetMetadataExtras {
    pub spec: String,
    pub icon: Option<String>,
    pub reference: Option<String>,
    pub reference_hash: Option<Base64VecU8>,
    pub decimals: u8,
}

#[derive(BorshDeserialize, BorshSerialize, Clone, Debug, Deserialize, Serialize, Eq, PartialEq)]
#[serde(crate = "near_sdk::serde")]
pub struct FungibleToken {
    name: String,
    label: String,
    symbol: String,
    uri: Option<AccountId>,
    network: Network,
    #[serde(
        deserialize_with = "deserialize_u128",
        serialize_with = "serialize_u128"
    )]
    fee_numerator: u128,
    #[serde(
        deserialize_with = "deserialize_u128",
        serialize_with = "serialize_u128"
    )]
    fixed_fee: u128,
    extras: Option<AssetMetadataExtras>,
}

fn deserialize_u128<'de, D>(deserializer: D) -> Result<u128, D::Error>
where
    D: Deserializer<'de>,
{
    <U128 as Deserialize>::deserialize(deserializer).map(|s| s.into())
}

fn serialize_u128<S>(x: &u128, s: S) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    <U128 as Serialize>::serialize(&U128::from(*x), s)
}

impl FungibleToken {
    #[allow(clippy::too_many_arguments)]
    pub fn new(
        name: String,
        label: String,
        symbol: String,
        uri: Option<AccountId>,
        network: Network,
        fee_numerator: u128,
        fixed_fee: u128,
        extras: Option<AssetMetadataExtras>,
    ) -> FungibleToken {
        Self {
            name,
            label,
            symbol,
            uri,
            network,
            fee_numerator,
            fixed_fee,
            extras,
        }
    }
}

impl FungibleToken {
    pub fn uri(&self) -> &Option<AccountId> {
        &self.uri
    }

    pub fn uri_deref(&self) -> Option<AccountId> {
        self.uri.clone()
    }
}

impl AssetMetadata for FungibleToken {
    fn name(&self) -> &String {
        &self.name
    }

    fn label(&self) -> &String {
        &self.label
    }

    fn network(&self) -> &Network {
        &self.network
    }

    fn symbol(&self) -> &String {
        &self.symbol
    }

    fn extras(&self) -> &Option<AssetMetadataExtras> {
        &self.extras
    }

    fn metadata(&self) -> &Self {
        self
    }

    fn fee_numerator(&self) -> u128 {
        self.fee_numerator
    }

    fn fee_numerator_mut(&mut self) -> &mut u128 {
        &mut self.fee_numerator
    }

    fn fixed_fee(&self) -> u128 {
        self.fixed_fee
    }

    fn fixed_fee_mut(&mut self) -> &mut u128 {
        &mut self.fixed_fee
    }
}
