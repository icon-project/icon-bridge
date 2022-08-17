use super:: Events;
use libraries::rlp::{self, Decodable};

#[derive(Clone, PartialEq, Eq, Debug)]
pub struct Receipt {
    index: u64,
    events: Events,
    height: u64,
}

impl Receipt {
    pub fn index(&self) -> u64 {
        self.index
    }

    pub fn events(&self) -> &Events {
        &self.events
    }

    pub fn height(&self) -> u64 {
        self.height
    }
}

impl Decodable for Receipt {
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
