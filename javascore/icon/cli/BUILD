load("@icon//cli:defs.bzl", "configure_bmr", "configure_link", "create_keystore", "get_score_address")

package(default_visibility = ["//visibility:public"])

create_keystore(
    name = "bmc",
)

get_score_address(
    name = "get_score_address_bmc",
    contract = "bmc",
)

get_score_address(
    name = "get_score_address_bsr",
    contract = "bsr",
)

get_score_address(
    name = "get_score_address_irc",
    contract = "irc",
)

get_score_address(
    name = "get_score_address_bts",
    contract = "bts",
)

configure_bmr(
    name = "icon",
)

genrule(
    name = "deploy_bmc",
    srcs = ["@icon//:goloop_config_dir"],
    outs = ["deploy_bmc.out"],
    cmd = "$(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx deploy $(location @javascore//bmc:javascore) --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --content_type application/java --step_limit 13610920001 --param _net=\"$$(cat $(location @icon//:network_address))\" | echo $$(tr -d '\"') > $@",
    executable = True,
    local = True,
    tools = [
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:network_address",
        "@icon//:network_id",
        "@icon//:node_url",
        "@javascore//bmc:javascore",
    ],
)

genrule(
    name = "deploy_bts",
    srcs = ["@icon//:goloop_config_dir", "@javascore//bts:javascore"],
    outs = ["deploy_bts.out"],
    cmd = "$(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url))  sendtx deploy $(location @javascore//bts:javascore) --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --content_type application/java --step_limit 13610920001 --param _bmc=$$(cat $(location @icon//cli:get_score_address_bmc)) --param  _serializedIrc2=$$(cat $(location @javascore//irc2Tradeable:javascore_serialized)) --param _name=\"btp-$$(cat $(location @icon//:network_id))-ICX\" --param _decimals=\"0x$$(printf '%x' 18)\" --param _feeNumerator=\"0x$$(printf '%x' 100)\" --param _fixedFee=\"0x$$(printf '%x' 4300000000000000000)\" | echo $$(tr -d '\"') > $@",
    executable = True,
    local = True,
    tools = [
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:network_id",
        "@icon//:node_url",
        "@icon//cli:deploy_bmc",
        "@javascore//irc2Tradeable:javascore_serialized",
        "@icon//cli:get_score_address_bmc"
    ],
)

genrule(
    name = "deploy_bsr",
    srcs = ["@icon//:goloop_config_dir", "@javascore//bsr:javascore"],
    outs = ["deploy_bsr.out"],
    cmd = "$(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx deploy $(location @javascore//bsr:javascore) --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --content_type application/java --step_limit 13610920001 | echo $$(tr -d '\"') > $@",
    executable = True,
    local = True,
    tools = [
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:network_id",
        "@icon//:node_url",
    ],
)

genrule(
    name = "deploy_irc",
    srcs = ["@icon//:goloop_config_dir", "@javascore//lib:irc2-token-0.9.1-optimized.jar"],
    outs = ["deploy_irc.out"],
    cmd = "$(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx deploy $(location @javascore//lib:irc2-token-0.9.1-optimized.jar) --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --content_type application/java --step_limit 13610920001 --param _name=\"TICX\" --param _symbol=\"TICX\" --param _decimals=\"0x12\" --param _initialSupply=\"0x186a0\" | echo $$(tr -d '\"') > $@",
    executable = True,
    local = True,
    tools = [
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:network_id",
        "@icon//:node_url",
    ],
)

genrule(
    name = "set_near_link_rotate_term",
    srcs = [
        "@near//:btp_address",
        "@icon//:goloop_config_dir",
        "@icon//cli:add_near_link",
    ],
    outs = ["set_near_link_rotate_term.out"],
    cmd = """$(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx  call --to $$(cat $(location @icon//cli:get_score_address_bmc))  --method setLinkRotateTerm --param _link=$$(cat $(location @near//:btp_address)) --param _block_interval=0x1770 --param _max_agg=0x08 --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --step_limit 13610920001  | jq -r . >$@ """,
    executable = True,
    local = True,
    tools = [
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:network_id",
        "@icon//:node_url",
        "@icon//cli:get_score_address_bmc",
        "@javascore//bmc:javascore",
    ],
)

genrule(
    name = "set_near_link_height",
    srcs = [
        "@near//:btp_address",
        "@icon//:goloop_config_dir",
        "@icon//cli:add_near_link",
        "@near//:latest_block_height",
    ],
    outs = ["set_near_link_height.out"],
    cmd = """$(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx  call --to $$(cat $(location @icon//cli:get_score_address_bmc))  --method setLinkRxHeight --param _link=$$(cat $(location @near//:btp_address)) --param _height=$$(($$(cat $(location  @near//:latest_block_height))+1)) --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --step_limit 13610920001  | jq -r . >$@ """,
    executable = True,
    local = True,
    tools = [
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:network_id",
        "@icon//:node_url",
        "@icon//cli:get_score_address_bmc",
        "@javascore//bmc:javascore",
    ],
)

genrule(
    name = "set_near_link_delay_limit",
    srcs = [
        "@near//:btp_address",
        "@icon//:goloop_config_dir",
        "@icon//cli:add_near_link",
    ],
    outs = ["set_near_link_delay_limit.out"],
    cmd = """$(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx  call --to $$(cat $(location @icon//cli:get_score_address_bmc)) --method setLinkDelayLimit --param _link=$$(cat $(location @near//:btp_address)) --param _value=4 --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --step_limit 13610920001  | jq -r . >$@ """,
    executable = True,
    local = True,
    tools = [
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:network_id",
        "@icon//:node_url",
        "@icon//cli:get_score_address_bmc",
        "@javascore//bmc:javascore",
    ],
)

genrule(
    name = "set_near_link",
    srcs = [
        "@icon//cli:set_near_link_rotate_term",
        "@icon//cli:set_near_link_delay_limit",
        "@icon//cli:set_near_link_height",
    ],
    outs = ["set_near_link.out"],
    cmd = "echo 'set_link_completed' >$@",
    executable = True,
    local = True,
)

configure_link(
    name = "near",
)

genrule(
    name = "keysecret",
    outs = [
        "keysecret.out",
    ],
    cmd = "echo 'iconbridge' > $@",
    executable = True,
    local = True,
)

genrule(
    name = "add_bts_icon",
    srcs = ["@icon//cli:deploy_bts"],
    outs = ["add_bts_icon.out"],
    cmd = """$(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx  call --to $$(cat $(location @icon//cli:get_score_address_bmc)) --method addService --param _addr=$$(cat $(location @icon//cli:get_score_address_bts)) --param _svc=bts --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --step_limit 13610920001  | jq -r . >$@ """,
    executable = True,
    local = True,
    tools = [
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:goloop_config_dir",
        "@icon//:network_id",
        "@icon//:node_url",
        "@icon//cli:get_score_address_bmc",
        "@icon//cli:get_score_address_bts",
    ],
)

genrule(
    name = "get_icon_validator_hash",
    srcs = ["@btp//cmd/iconvalidators"],
    outs = ["get_icon_validator_hash.out"],
    cmd = """ URI=\"$$(cat $(location @icon//:node_url))\" HEIGHT=\"0x$$(printf '%x' $$(cat $(location @icon//:latest_block_height)))\" $(execpath @btp//cmd/iconvalidators:iconvalidators) | jq -r '.hash' |  echo $$(tr -d '\"') > $@""",
    executable = True,
    local = True,
    tools = [
        "@icon//:latest_block_height",
        "@icon//:node_url",
    ],
)

genrule(
    name = "add_near_bmr",
    srcs = [
        ":set_near_link",
        "@icon//cli:get_wallet_icon_keystore",
        "@near//:btp_address",
    ],
    outs = ["add_near_bmr.out"],
    cmd = """
    $(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx  call --to $$(cat $(location @icon//cli:get_score_address_bmc)) --method addRelay  --param _link=$$(cat $(location @near//:btp_address)) --param _addr=$$(cat $(location @icon//cli:get_wallet_icon_keystore))  --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --step_limit 13610920001  | jq -r . >$@
    """,
    executable = True,
    local = True,
    tools = [
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:goloop_config_dir",
        "@icon//:network_id",
        "@icon//:node_url",
        "@icon//cli:get_score_address_bmc",
    ],
)

# genrule(
#     name = "add_bsr",
#     srcs= [
#         ":deploy_bsr",
#     ],
#     outs = ["add_bsr.out"],
#     cmd = """
#     $(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx  call --to $$(cat $(location @icon//cli:get_score_address_bts)) --method addRestrictor  --param _addr=$$(cat $(location @icon//cli:get_score_address_bsr))  --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --step_limit 13610920001  | jq -r . >$@
#     """,
#     executable = True,
#     local = True,
#     tools = [
#         "@com_github_icon_project_goloop//cmd/goloop",
#         "@icon//:goloop_config_dir",
#         "@icon//:network_id",
#         "@icon//:node_url",
#         "@icon//cli:get_score_address_bsr",
#         "@icon//cli:get_score_address_bts",
#     ]
# )

genrule(
    name = "send_icx_to_near",
    srcs = [
        ":add_bts_icon",
        ":set_near_link",
        "@com_github_icon_project_goloop//cmd/goloop",
        "@icon//:goloop_config_dir",
        "@icon//cli:get_score_address_bts",
        "@near//cli:encode_public_key_alice",
        "@near//cli:register_wrapped_coin_bts",
    ],
    outs = ["send_icx_to_near.out"],
    cmd = """
     $(execpath @com_github_icon_project_goloop//cmd/goloop:goloop) rpc --uri $$(cat $(location @icon//:node_url)) sendtx  call --to $$(cat $(location get_score_address_bts)) --method transferNativeCoin  --param _to=\"btp://0x1.near/$$(cat $(location @near//cli:encode_public_key_alice))\"  --value 0x$$(printf '%x' 1000000000000000000000000) --key_store $$(cat $(location @icon//:goloop_config_dir))/keystore.json --key_secret $$(cat $(location @icon//:goloop_config_dir))/keysecret --nid \"$$(cat $(location @icon//:network_id))\" --step_limit 13610920001 >$@
    """,
    executable = True,
    local = True,
    tools = [
        "@icon//:network_id",
        "@icon//:node_url",
    ],
)

select({
    "@icon//:localnet": exports_files(glob(["configs/localnet/*"])),
    "@icon//:docker_localnet": exports_files(glob(["configs/localnet/*"])),
    "@icon//:lisbon": exports_files(glob(["configs/lisbon/*"])),
    "@icon//:berlin": exports_files(glob(["configs/berlin/*"])),
    "@icon//:mainnet": exports_files(glob(["configs/mainnet/*"])),
    "//conditions:default": []
})
