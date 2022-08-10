use super::{ Nullable, Proofs};
use libraries::rlp::{self, Decodable};

#[derive(Clone, PartialEq, Eq, Debug)]
pub struct ReceiptProof {
    index: u64,
    events: Proofs,
    height: u64,
}

impl ReceiptProof {
    pub fn index(&self) -> u64 {
        self.index
    }

    pub fn events(&self) -> &Proofs {
        &self.events
    }

    pub fn height(&self) -> u64 {
        self.height
    }

    // pub fn event_len(&self) -> Result<usize,String>{
    //     match  self.events.get() {
    //         Ok(proofs) => {
    //             Ok(proofs.len())
    //         },
    //         Err(err) => {
    //             Err(err.to_string())
    //         },
    //     }
    // }
}

impl Decodable for ReceiptProof {
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
