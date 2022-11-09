mod steps;

#[cfg(test)]
mod manage_tokens {
    use super::*;
    use kitten::*;
    use steps::*;

    #[tokio::test(flavor = "multi_thread")]
    async fn withdraw_icx_success() {
        Kitten::given(NEW_CONTEXT)
            .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
            .and(BTS_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
            .and(BSH_OWNER_REGISTERS_ICX)
            .and(ICON_LINK_IS_PRESENT_IN_BMC)
            .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
            .and(NATIVE_COIN_BSH_IS_REGISTERED)
            .and(TOKEN_TRANSFER_MESSAGE_IS_PROVIDED_AS_HANDLE_BTP_MESSAGE_PARAM)
            .and(BMC_OWNER_INVOKES_HANDLE_BTP_MESSAGE_IN_BTS)
            .and(WITHDRAW_IS_INVOKED_BY_CHUCK)
            .when(STROAGE_BALANCE_IS_INVOKED)
            .then(STROAGE_BALANCE_FOR_ICX_WITH_CHUCK_ACCOUNT_SHOULD_BE_ZERO);
    }
}
