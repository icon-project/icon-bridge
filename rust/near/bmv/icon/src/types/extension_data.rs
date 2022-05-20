use libraries::rlp::{self, Decodable, Encodable};
use super::Nullable;

#[derive(Default, PartialEq, Eq, Debug, Clone)]
pub struct ExtensionData {
    data: Vec<Nullable<Vec<u8>>>,
}

impl Decodable for ExtensionData {
    fn decode(rlp: &rlp::Rlp) -> Result<Self, rlp::DecoderError> {
        let data = rlp.as_val::<Vec<u8>>()?;
        let rlp = rlp::Rlp::new(&data);
        Ok(Self {
            data: rlp.as_list()?,
        })
    }
}

impl Encodable for ExtensionData {
    fn rlp_append(&self, stream: &mut rlp::RlpStream) {
        stream.append_list(&self.data);
    }
}
