mod steps;

#[cfg(test)]
mod manage_tokens {
    use super::*;
    use kitten::*;
    use steps::*;

    #[tokio::test(flavor = "multi_thread")]
    async fn bts_owner_can_register_a_cross_chain_token() {
        Kitten::given(NEW_CONTEXT)
            .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
            .and(BTS_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
            .when(BSH_OWNER_REGISTERS_BNUSD)
            .then(BNUSD_CONTRACT_SHOULD_BE_REGISTERED_WITH_PROVIDED_METADATA);
    }
}
