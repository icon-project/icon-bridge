mod steps;

#[cfg(test)]
mod manage_owner_accounts {
    use super::*;
    use kitten::*;
    use steps::*;

    mod bmc {
        use super::*;

        #[tokio::test(flavor = "multi_thread")]
        async fn non_bmc_owner_cannot_remove_owners() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIES_ACCOUNT_IS_CREATED)
                .and(CHARLIES_ACCOUNT_ID_IS_PROVIDED_AS_REMOVE_OWNER_PARAM)
                .when(ALICE_INVOKES_REMOVE_OWNER_IN_BMC)
                .then(BMC_SHOULD_THROW_OWNER_DOES_NOT_EXIST_ON_REMOVING_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn bmc_contract_account_owner_can_add_new_bmc_owner() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIES_ACCOUNT_IS_CREATED)
                .and(CHARLIES_ACCOUNT_ID_IS_PROVIDED_AS_ADD_OWNER_PARAM)
                .when(ALICE_INVOKES_ADD_OWNER_IN_BMC)
                .then(CHARLIES_ACCOUNT_ID_SHOULD_BE_IN_THE_LIST_OF_BMC_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn recently_added_bmc_owner_can_add_new_bmc_owner() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIE_IS_AN_EXISITNG_OWNER_IN_BMC)
                .and(BOBS_ACCOUNT_IS_CREATED)
                .and(BOBS_ACCOUNT_ID_IS_PROVIDED_AS_ADD_OWNER_PARAM)
                .when(CHARLIE_INVOKES_ADD_OWNER_IN_BMC)
                .then(BOBS_ACCOUNT_ID_SHOULD_BE_IN_THE_LIST_OF_BMC_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn bmc_contract_account_owner_can_remove_other_owners() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIE_IS_AN_EXISITNG_OWNER_IN_BMC)
                .and(CHARLIES_ACCOUNT_ID_IS_PROVIDED_AS_REMOVE_OWNER_PARAM)
                .when(ALICE_INVOKES_REMOVE_OWNER_IN_BMC)
                .then(CHARLIES_ACCOUNT_ID_SHOULD_NOT_BE_IN_THE_LIST_OF_BMC_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn add_new_bmc_owner_by_existing_bmc_owner_success() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIES_ACCOUNT_IS_CREATED)
                .and(CHARLIES_ACCOUNT_ID_IS_PROVIDED_AS_ADD_OWNER_PARAM)
                .when(ALICE_INVOKES_ADD_OWNER_IN_BMC)
                .then(CHARLIES_ACCOUNT_ID_SHOULD_BE_IN_THE_LIST_OF_BMC_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn non_bmc_owner_cannot_add_new_owners() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHUCKS_ACCOUNT_IS_CREATED)
                .and(CHARLIES_ACCOUNT_IS_CREATED)
                .and(CHARLIES_ACCOUNT_ID_IS_PROVIDED_AS_ADD_OWNER_PARAM)
                .when(CHUCK_INVOKES_ADD_OWNER_IN_BMC)
                .then(BMC_SHOULD_THROW_UNAUTHORISED_ERROR_ON_ADDING_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn bmc_owner_cannot_remove_self_if_self_is_the_last_owner() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(ALICES_ACCOUNT_ID_IS_PROVIDED_AS_REMOVE_OWNER_PARAM)
                .when(ALICE_INVOKES_REMOVE_OWNER_IN_BMC)
                .then(BMC_SHOULD_THROW_LAST_OWNER_ERROR_ON_REMOVING_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn new_bmc_owner_cannot_add_already_existing_owner() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIE_IS_AN_EXISITNG_OWNER_IN_BMC)
                .and(CHARLIES_ACCOUNT_ID_IS_PROVIDED_AS_ADD_OWNER_PARAM)
                .when(ALICE_INVOKES_ADD_OWNER_IN_BMC)
                .then(BMC_SHOULD_THROW_OWNER_ALREADY_EXISTS_ERROR_ON_ADDING_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn newly_added_bmc_owner_can_remove_bmc_owner() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIE_IS_AN_EXISITNG_OWNER_IN_BMC)
                .and(ALICES_ACCOUNT_ID_IS_PROVIDED_AS_REMOVE_OWNER_PARAM)
                .when(CHARLIE_INVOKES_REMOVE_OWNER_IN_BMC)
                .then(ALICES_ACCOUNT_ID_SHOULD_NOT_BE_IN_THE_LIST_OF_BMC_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn bmc_owner_can_remove_another_bmc_owner_success() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIE_IS_AN_EXISITNG_OWNER_IN_BMC)
                .and(CHARLIES_ACCOUNT_ID_IS_PROVIDED_AS_REMOVE_OWNER_PARAM)
                .when(ALICE_INVOKES_REMOVE_OWNER_IN_BMC)
                .then(CHARLIES_ACCOUNT_ID_SHOULD_NOT_BE_IN_THE_LIST_OF_BMC_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn bmc_owner_cannot_remove_non_existing_owner() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIES_ACCOUNT_IS_CREATED)
                .and(CHARLIES_ACCOUNT_ID_IS_PROVIDED_AS_REMOVE_OWNER_PARAM)
                .when(ALICE_INVOKES_REMOVE_OWNER_IN_BMC)
                .then(BMC_SHOULD_THROW_OWNER_DOES_NOT_EXIST_ON_REMOVING_OWNERS);
        }

        #[tokio::test(flavor = "multi_thread")]
        async fn bmc_owner_can_remove_self() {
            Kitten::given(NEW_CONTEXT)
                .and(BMC_CONTRACT_IS_DEPLOYED_AND_INITIALIZED)
                .and(BMC_CONTRACT_IS_OWNED_BY_ALICE)
                .and(CHARLIE_IS_AN_EXISITNG_OWNER_IN_BMC)
                .and(ALICES_ACCOUNT_ID_IS_PROVIDED_AS_REMOVE_OWNER_PARAM)
                .when(ALICE_INVOKES_REMOVE_OWNER_IN_BMC)
                .then(ALICES_ACCOUNT_ID_SHOULD_NOT_BE_IN_THE_LIST_OF_BMC_OWNERS);
        }
    }
}

