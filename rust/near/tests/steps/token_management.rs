use super::*;
use super::{BMC_CONTRACT, BTS_CONTRACT};
use serde_json::{from_value, json};
use test_helper::types::Context;

pub static BOB_IS_BTS_CONTRACT_OWNER: fn(Context) -> Context = |mut context: Context| {
    let bsh_signer = context.contracts().get("bts").as_account().clone();
    context.accounts_mut().add("bob", bsh_signer);
    context
};

pub static BOB_INVOKES_REGISTER_NEW_COIN_FORM_BSH: fn(Context) -> Context =
    |mut context: Context| {
        let signer = context.accounts().get("bob").to_owned();
        context.set_signer(&signer);
        BTS_CONTRACT.register(context)
    };

pub static BSH_OWNER_REGISTERS_BNUSD: fn(Context) -> Context = |mut context: Context| {
    BOBS_ACCOUNT_IS_CREATED(context)
        .pipe(BOB_IS_BTS_CONTRACT_OWNER)
        .pipe(BNUSD_COIN_NAME_IS_PROVIDED_AS_REGISTER_PARAM)
        .pipe(BOB_INVOKES_REGISTER_NEW_COIN_FORM_BSH)
};

pub static BNUSD_COIN_NAME_IS_PROVIDED_AS_REGISTER_PARAM: fn(Context) -> Context =
    |mut context: Context| {
        let uri = format!("bnusd.{}", context.contracts().get("bts").id());
        let mut context = register_token_account("bnusd", &uri, context);
        context.add_method_params(
        "register",
        json!({
            "coin": {
                "metadata": {
                    "name": format!("btp-{}-bnUSD", ICON_NETWORK),
                    "label": "Wrapped bnUSD From ICON",
                    "symbol": "bnUSD",
                    "uri": uri,
                    "network": ICON_NETWORK,
                    "fee_numerator": "100",
                    "fixed_fee": "1500000000000000000",
                    "extras": {
                        "spec": "ft-1.0.0",
                        "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAMAAABHPGVmAAABTVBMVEUMKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0wwqAyyaQYYGcuu5saa20VVWEkkoMSSFoPOlQiiX0ddnIoo40mmogff3gss5Yqq5L26tYlAAAAXnRSTlMA8+/uphQB4/v98RVrJ1wg4S0428iVbpDLJnMoZde/nD6y8PckpxPKmvaZIwfcM2xDyfJEpTSRiRybRbEahBIG9N0s4MCLsGTfGSJwCBY3XRhjlu0plNhypOQXiuK+gdXRMQAAA25JREFUeF692ldzGkkUBeDDgAAJEBJCgEDBcsBJsmKwvJLTOqxz3rxnhhwk+/8/bgF2F1BjcXvS9zZPp6Z77q2e7oZUYjuXf5bdDIfIUHgz+yyf207AS1eTkTnamIskr8ITpdMFXmDhtASX7m9lOVF26z6cS926QpErt1JwZunPaYpN7y9B39RinFrii1PQE/1gUJvxPAoN6VU6spqG2OwMHZqZhczGMV043oDA++t05foLTPR4ji7NPcYEmTBdC2dwoUsheiB0CRd4G6Indt/ip+Yv0yOX52EPJzF6JnYXth5k6aHsA9g5pqfWYGOWHrPpMOkZemwmjTHR/+i51ShGPacPljGi/Ik+CJcx7Atl2lZfmzKLGLIU9yckvgQF+/QnhBEoqWm/QqZLKuQG/QrhDbXEuulfyM0pDPxB/0JYxMAKhepnVavPPKtTaAV9dyjTaliKxsvcQc+RMOOrNaJBmSP0LFCka42pUWQBAG5TpGKNq1LmNoCk7osoHYokAVyjREtNeKXyI69LkWsA/tUarfOhtzIpEgMSFDm3BtjTVA8iCWxrlXr1e2S7jzL3kNMKsehADnmK1NSc6MvjF4rULecp6yhQxlQ12KGmAn4TN3nFbNSp41cYlGmZ1pCvHcoZ2CV1ZkWpVigVAsUqpmavV0C55pk1ok0h7FJDszHyNnXpcBnUU+tayhlFDPUJOxk1U/oJF6ivpjdeBbu2Iq/MirCt5OlASyvkUNjqvw2oOtcKyWGPEmP119QKuYcEJX4sT1tqTjQm/m8gRoHuSGV0tD7hmHRJ1FF98bxSU3XSFS+JkvLxUrSmZEu8TK3ZZHyTL1OxSYmGNc6sU2ITPf/QUYpZo8iR1k9QbaTNV+uUOdD7nWu1qyqi3aLM7xgoUqyhKl+qGMgvdpCbBSj5t+3xGUrEr5BIEFtRKQxZ9CfkI4aVwxRpVvqaFDHKGLEcwEYnoqt+bNkGsfkcwDb6O9hYC+BAABs79FBhA7bu+nRIE/xxE5AJ0ROhTJBHgPYyBl0zMphg3v2x7DwmeuH2gPkEAq9dVeXaa/8P/d9B7OUTOvLkJTRElw1qM5ajQV4pkUvt61yOiaT8v+ZTgnOJ4gonWikm4NLB04uvXj098OwSWYw2/pJdIpN7s/fqcH3nUfgh+TD8aGf98NXeGwj9D4eiMfXOYUaZAAAAAElFTkSuQmCC",
                        "decimals": 18
                    }
                }
            }
        }),
    );
        context
    };

pub static BNUSD_CONTRACT_SHOULD_BE_REGISTERED_WITH_PROVIDED_METADATA: fn(Context) =
    |context: Context| {
        let expected = json!({
            "spec": "ft-1.0.0",
            "name": "Wrapped bnUSD From ICON",
            "symbol": "bnUSD",
            "icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAMAAABHPGVmAAABTVBMVEUMKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0MKk0wwqAyyaQYYGcuu5saa20VVWEkkoMSSFoPOlQiiX0ddnIoo40mmogff3gss5Yqq5L26tYlAAAAXnRSTlMA8+/uphQB4/v98RVrJ1wg4S0428iVbpDLJnMoZde/nD6y8PckpxPKmvaZIwfcM2xDyfJEpTSRiRybRbEahBIG9N0s4MCLsGTfGSJwCBY3XRhjlu0plNhypOQXiuK+gdXRMQAAA25JREFUeF692ldzGkkUBeDDgAAJEBJCgEDBcsBJsmKwvJLTOqxz3rxnhhwk+/8/bgF2F1BjcXvS9zZPp6Z77q2e7oZUYjuXf5bdDIfIUHgz+yyf207AS1eTkTnamIskr8ITpdMFXmDhtASX7m9lOVF26z6cS926QpErt1JwZunPaYpN7y9B39RinFrii1PQE/1gUJvxPAoN6VU6spqG2OwMHZqZhczGMV043oDA++t05foLTPR4ji7NPcYEmTBdC2dwoUsheiB0CRd4G6Indt/ip+Yv0yOX52EPJzF6JnYXth5k6aHsA9g5pqfWYGOWHrPpMOkZemwmjTHR/+i51ShGPacPljGi/Ik+CJcx7Atl2lZfmzKLGLIU9yckvgQF+/QnhBEoqWm/QqZLKuQG/QrhDbXEuulfyM0pDPxB/0JYxMAKhepnVavPPKtTaAV9dyjTaliKxsvcQc+RMOOrNaJBmSP0LFCka42pUWQBAG5TpGKNq1LmNoCk7osoHYokAVyjREtNeKXyI69LkWsA/tUarfOhtzIpEgMSFDm3BtjTVA8iCWxrlXr1e2S7jzL3kNMKsehADnmK1NSc6MvjF4rULecp6yhQxlQ12KGmAn4TN3nFbNSp41cYlGmZ1pCvHcoZ2CV1ZkWpVigVAsUqpmavV0C55pk1ok0h7FJDszHyNnXpcBnUU+tayhlFDPUJOxk1U/oJF6ivpjdeBbu2Iq/MirCt5OlASyvkUNjqvw2oOtcKyWGPEmP119QKuYcEJX4sT1tqTjQm/m8gRoHuSGV0tD7hmHRJ1FF98bxSU3XSFS+JkvLxUrSmZEu8TK3ZZHyTL1OxSYmGNc6sU2ITPf/QUYpZo8iR1k9QbaTNV+uUOdD7nWu1qyqi3aLM7xgoUqyhKl+qGMgvdpCbBSj5t+3xGUrEr5BIEFtRKQxZ9CfkI4aVwxRpVvqaFDHKGLEcwEYnoqt+bNkGsfkcwDb6O9hYC+BAABs79FBhA7bu+nRIE/xxE5AJ0ROhTJBHgPYyBl0zMphg3v2x7DwmeuH2gPkEAq9dVeXaa/8P/d9B7OUTOvLkJTRElw1qM5ajQV4pkUvt61yOiaT8v+ZTgnOJ4gonWikm4NLB04uvXj098OwSWYw2/pJdIpN7s/fqcH3nUfgh+TD8aGf98NXeGwj9D4eiMfXOYUaZAAAAAElFTkSuQmCC",
            "decimals": 18,
            "reference": null,
            "reference_hash": null
        });

        let result = nep141_contract("bnusd")
            .get_metadata(context)
            .method_responses("ft_metadata");

        assert_eq!(result, expected);
    };
