use crate::types::{Bmc, Bts, Context, Contract, Nep141};
use duplicate::duplicate;
use std::path::Path;
use std::str::FromStr;
use tokio::runtime::Handle;
use workspaces::{prelude::*, AccountId};
use workspaces::{
    types::KeyType, types::SecretKey, Contract as WorkspaceContract, DevNetwork, Worker,
};

pub async fn deploy(
    path: &str,
    worker: Worker<impl DevNetwork>,
) -> Result<WorkspaceContract, workspaces::error::Error> {
    worker
        .dev_deploy(&std::fs::read(Path::new(path)).unwrap())
        .await
}

pub async fn create_empty_contract(worker: Worker<impl DevNetwork>) -> anyhow::Result<WorkspaceContract> {
    worker.create_empty_contract().await
}

#[duplicate(
    contract_type;
    [ Bmc ];
    [ Bts ];
)]
impl Contract<'_, contract_type> {
    pub fn deploy(&self, mut context: Context) -> Context {
        let worker = context.worker().clone();
        let handle = Handle::current();
        let contract = tokio::task::block_in_place(move || {
            handle.block_on(async { deploy(self.source(), worker).await.unwrap() })
        });
        context.set_signer(contract.as_account());
        context.contracts_mut().add(self.name(), contract);
        context
    }
}

#[duplicate(
    contract_type;
    [ Nep141 ];
)]
impl Contract<'_, contract_type> {
    pub fn setup(&self, mut context: Context, account_id: &str) -> Context {
        let contract = WorkspaceContract::from_secret_key(
            AccountId::from_str(account_id).unwrap(),
            SecretKey::from_random(KeyType::ED25519),
            context.worker(),
        );
        context.set_signer(contract.as_account());
        context.contracts_mut().add(self.name(), contract);
        context
    }
}
