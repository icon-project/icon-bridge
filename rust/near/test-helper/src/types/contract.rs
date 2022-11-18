use duplicate::duplicate;
use std::collections::HashMap;
use std::marker::PhantomData;
use workspaces::Contract as WorkspaceContract;

#[derive(Default)]
pub struct Contracts(HashMap<String, WorkspaceContract>);

impl Contracts {
    pub fn add(&mut self, name: &str, contract: WorkspaceContract) {
        self.0.insert(name.to_owned(), contract);
    }

    pub fn get(&self, name: &str) -> &WorkspaceContract {
        self.0.get(name).unwrap()
    }

    pub fn as_mut(&mut self) -> &mut Self {
        self
    }
}

pub struct Contract<'a, T> {
    name: &'a str,
    source: &'a str,
    phantom: PhantomData<&'a T>,
}

#[duplicate(
    contract_type;
    [ Bmc ];
    [ Bts ];
    [ Nep141 ];
  )]
impl Contract<'_, contract_type> {
    fn new(name: &'static str, source: &'static str) -> Contract<'static, contract_type> {
        Contract {
            name,
            source,
            phantom: PhantomData::<&contract_type>,
        }
    }

    pub fn name(&self) -> &str {
        self.name
    }

    pub fn source(&self) -> &str {
        self.source
    }
}

pub struct Bmc {}

pub struct Bts {}

pub struct Nep141 {}

pub struct BmcContract {}

impl BmcContract {
    pub fn new(name: &'static str, source: &'static str) -> Contract<'static, Bmc> {
        Contract::<Bmc>::new(name, source)
    }
}

pub struct BtsContract {}

impl BtsContract {
    pub fn new(name: &'static str, source: &'static str) -> Contract<'static, Bts> {
        Contract::<Bts>::new(name, source)
    }
}

pub struct Nep141Contract {}

impl Nep141Contract {
    pub fn new(name: &'static str, source: &'static str) -> Contract<'static, Nep141> {
        Contract::<Nep141>::new(name, source)
    }
}

pub struct Nep141{}

pub struct Nep141Contract{}

impl Nep141Contract{
    pub fn new(name:&'static str, source: &'static str) -> Contract<'static,Nep141>{
        Contract::<Nep141>::new(name,source)
    }
}

pub struct WNear{}

pub struct WNearContract{}

impl WNearContract{
    pub fn new(name:&'static str, source: &'static str) -> Contract<'static,WNear>{
        Contract::<WNear>::new(name,source)
    }
}

pub struct Nep141Testable{}

pub struct Nep141TestableContract{}

impl Nep141TestableContract{
    pub fn new(name:&'static str, source: &'static str) -> Contract<'static,Nep141Testable>{
        Contract::<Nep141Testable>::new(name,source)
    }
}