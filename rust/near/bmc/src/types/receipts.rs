use super:: {Events};
use libraries::rlp::{self, Decodable};
/// Representation of Receipt
#[derive(Clone, PartialEq, Eq, Debug)]
pub struct Receipt {
    /// A receipt must have index
    index: u64,
    /// A receipt must have evens
    events: Events,
    ///A receipt contains height
    height: u64,
}

impl Receipt {
    /// Returns the index
    pub fn index(&self) -> u64 {
        self.index
    }
    /// returns the events
    pub fn events(&self) -> &Events {
        &self.events
    }
    /// Returns the height
    pub fn height(&self) -> u64 {
        self.height
    }
}

impl Decodable for Receipt {
    /// Decodes the obtained result
    fn decode(rlp: &rlp::Rlp) -> Result<Self, rlp::DecoderError> {
        let data = rlp.as_val::<Vec<u8>>()?;
        let rlp = rlp::Rlp::new(&data);
        Ok(Self {
            index: rlp.val_at(0)?,
            events: rlp.val_at(1)?,
            height: rlp.val_at(2)?,
        })
    }
}
