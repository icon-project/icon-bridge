load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary", "npm_package_bin")
load("@near//cli:defs.bzl", "configure_bmr", "configure_link", "create_account", "create_sub_account")
load("@rules_python//python:defs.bzl", "py_binary")
load("@near_cli_util_python_deps//:requirements.bzl", "requirement")
load("@bazel_skylib//rules:common_settings.bzl", "string_flag")

package(default_visibility = ["//visibility:public"])

py_binary(
    name = "encode_base58",
    srcs = ["encode_base58.py"],
    deps = [
        requirement("base58"),
    ],
)

nodejs_binary(
    name = "near_binary",
    data = [
        "@near_npm//near-cli",
    ],
    entry_point = "@near_npm//:node_modules/near-cli/bin/near",
    env = select({
        "@near//:localnet": {"NEAR_ENV": "local"},
        "@near//:docker_localnet": {"NEAR_ENV": "local"},
        "@near//:testnet": {},
    }),
)

create_account(
    name = "bmc",
)

create_account(
    name = "alice",
)

create_sub_account(
    name = "bts",
)

create_sub_account(
    name = "tokenbsh",
)

genrule(
    name = "deploy_bmc",
    srcs = [
        "@near//bmc:bmc_wasm",
    ],
    outs = ["deploy_bmc.out"],
    cmd = "$(execpath :near_binary) deploy $$(cat $(location @near//cli:encode_public_key_bmc)) --nodeUrl $$(cat $(locations @near//:node_url)) --wasmFile $(locations @near//bmc:bmc_wasm) --accountId $$(cat $(location @near//cli:encode_public_key_bmc)) --initFunction \"new\" --initArgs '{\"network\": \"0x1.near\", \"block_interval\": 3000}'> $@",
    executable = True,
    local = True,
    tools = [
        "@near//:node_url",
        "@near//cli:create_account_bmc",
        "@near//cli:encode_public_key_bmc",
        "@near//cli:near_binary",
    ],
)

genrule(
    name = "deploy_bts",
    srcs = [
        "@near//bsh/bts:bts_wasm",
    ],
    outs = ["deploy_natvie_coin_bsh.out"],
    cmd = """$(execpath :near_binary) deploy $$(cat $(location @near//cli:create_sub_account_bts)) --nodeUrl $$(cat $(locations @near//:node_url)) --wasmFile $(locations @near//bsh/bts:bts_wasm) --accountId $$(cat $(location @near//cli:create_sub_account_bts)) --initFunction \"new\" --initArgs  \\'\\{\\"service_name\\":\\"bts\\"\\,\\"bmc\\":\\"$$(cat $(location @near//cli:encode_public_key_bmc))\\"\\,\\"network\\":\\"0x1.near\\"\\,\\"native_coin\\":\\{\\"metadata\\":\\{\\"name\\":\\"NEAR\\"\\,\\"symbol\\":\\"NEAR\\"\\,\\"uri\\":null\\,\\"network\\":\\"0x1.near\\"\\,\\"fee_numerator\\":\\"0\\"\\,\\"fixed_fee\\":\\"0\\"\\,\\"extras\\":null\\}\\}\\}\\' > $@""",
    executable = True,
    local = True,
    tools = [
        "@near//:node_url",
        "@near//cli:create_account_bmc",
        "@near//cli:create_sub_account_bts",
        "@near//cli:encode_public_key_bmc",
        "@near//cli:near_binary",
    ],
)

genrule(
    name = "deploy_token_bsh",
    srcs = [
        "@near//bsh/token:token_bsh_wasm",
    ],
    outs = ["deploy_token_bsh.out"],
    cmd = """$(execpath :near_binary) deploy $$(cat $(location @near//cli:encode_public_key_tokenbsh)) --nodeUrl $$(cat $(locations @near//:node_url)) --wasmFile $(locations @near//bsh/token:token_bsh_wasm) --accountId $$(cat $(location @near//cli:create_sub_account_tokenbsh)) --initFunction \"new\" --initArgs \\'\\{\\"service_name\\":\\"tokenbsh\\"\\,\\"bmc\\":\\"$$(cat $(location @near//cli:encode_public_key_bmc))\\"\\,\\"network\\":\\"0x1.near\\"\\}\\' > $@""",
    executable = True,
    local = True,
    tools = [
        "@near//:node_url",
        "@near//cli:create_account_tokenbsh",
        "@near//cli:create_sub_account_tokenbsh",
        "@near//cli:encode_public_key_bmc",
        "@near//cli:near_binary",
    ],
)

genrule(
    name = "set_icon_link",
    srcs = [
        "@near//cli:deploy_bmc",
        "@near//cli:add_icon_link",
    ],
    outs = ["set_icon_link.out"],
    cmd = """$(execpath :near_binary) call $$(cat $(location @near//cli:encode_public_key_bmc)) set_link \\'\\{\\"link\\":\\"$$(cat $(location @icon//:btp_address))\\"\\,\\"block_interval\\":2000\\,\\"max_aggregation\\":100\\,\\"delay_limit\\":2\\}\\' --nodeUrl $$(cat $(locations @near//:node_url)) --accountId $$(cat $(location @near//cli:encode_public_key_bmc)) > $@""",
    executable = True,
    local = True,
    tools = [
        "@icon//:btp_address",
        "@near//:node_url",
        "@near//cli:encode_public_key_bmc",
        "@near//cli:near_binary",
    ],
)

genrule(
    name = "add_bts_service",
    outs = ["add_bts_service.out"],
    srcs = ["@near//cli:deploy_bts"],
    cmd = """$(execpath :near_binary) call $$(cat $(location @near//cli:encode_public_key_bmc)) add_service \\'\\{\\"name\\":\\"bts\\"\\,\\"service\\":\\"$$(cat $(location @near//cli:create_sub_account_bts))\\"\\}\\' --nodeUrl $$(cat $(locations @near//:node_url)) --accountId $$(cat $(location @near//cli:encode_public_key_bmc)) > $@ """,
    executable = True,
    local = True,
    tools = [
        "@near//:node_url",
        "@near//cli:create_sub_account_bts",
        "@near//cli:encode_public_key_bmc",
        "@near//cli:near_binary",
    ],
)

genrule(
    name = "register_wrapped_coin_bts",
    srcs = [
        "@near//cli:add_bts_service",
        "@icon//:network_address",
    ],
    outs = ["regsiter_new_wrapped_coin.out"],
    cmd = """$(execpath :near_binary) call $$(cat $(location @near//cli:create_sub_account_bts)) register \\'\\{\\"coin\\":\\{\\"metadata\\":\\{\\"name\\":\\"btp-$$(cat $(location @icon//:network_id))-ICX\\"\\,\\"symbol\\":\\"ICX\\",\\"uri\\":\\"btp-$$(cat $(location @icon//:network_id)).$$(cat $(location @near//cli:create_sub_account_bts))\\"\\,\\"network\\":\\"$$(cat $(location @icon//:network_address))\\"\\,\\"fee_numerator\\":\\"100\\"\\,\\"fixed_fee\\":\\"4300000000000000000\\"\\,\\"extras\\":\\{\\"spec\\":\\"ft-1.0.0\\",\\"icon\\":null,\\"reference\\":null,\\"reference_hash\\":null,\\"decimals\\":18\\}\\}\\}\\}\\' --nodeUrl $$(cat $(locations @near//:node_url)) --accountId $$(cat $(location @near//cli:create_sub_account_bts)) --amount 10 --gas 300000000000000 >$@""",
    executable = True,
    local = True,
    tools = [
        "@near//:node_url",
        "@near//cli:create_sub_account_bts",
        "@near//cli:near_binary",
        "@icon//:network_id"
    ],
)

genrule(
    name = "add_icon_bmr",
    srcs = [
        "@near//cli:set_icon_height",
        "@near//cli:get_wallet_near_keystore",
    ],
    outs = ["add_icon_bmr.out"],
    cmd = """
    $(execpath :near_binary) call $$(cat $(location @near//cli:encode_public_key_bmc)) add_relay \\'\\{\\"link\\":\\"$$(cat $(location @icon//:btp_address))\\"\\,\\"relay\\":\\"$$(cat $(location @near//cli:get_wallet_near_keystore))\\"\\}\\' --nodeUrl $$(cat $(locations @near//:node_url)) --accountId $$(cat $(location @near//cli:encode_public_key_bmc)) > $@ 
    """,
    executable = True,
    local = True,
    tools = [
        "@icon//:btp_address",
        "@near//:node_url",
        "@near//cli:encode_public_key_bmc",
        "@near//cli:near_binary",
    ],
)

genrule(
    name = "generate_keystore",
    srcs = [
        "@com_github_hugobyte_keygen//:keygen",
        ":keysecret",
    ],
    outs = ["keystore.json"],
    cmd = "$(execpath @com_github_hugobyte_keygen//:keygen) generate -p $$(cat $(location :keysecret)) -o $@",
    executable = True,
    local = True,
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

configure_link(
    name = "icon",
)

configure_bmr(
    name = "near",
)

genrule(
    name = "get_master_account",
    srcs = select({
        "@near//:localnet": ["@near//cli:configs/localnet/master_key.json"],
        "@near//:docker_localnet": ["@near//cli:configs/localnet/master_key.json"],
        "@near//:testnet": ["@near//cli:configs/testnet/master_key.json"],
        "@near//:mainnet": ["@near//cli:configs/mainnet/master_key.json"],
    }),
    outs = [
        "master_account.out",
    ],
    cmd = "echo $$(cat $(SRCS)) | jq -r .account_id > $@",
    output_to_bindir = True,
    executable = True,
    local = True,
)

genrule(
    name = "set_icon_height",
   srcs = [
        "@near//cli:set_icon_link",
        "@near//cli:get_wallet_near_keystore",
        "@icon//:latest_block_height"
    ],
    outs = ["set_icon_height.out"],
    cmd = """$(execpath :near_binary) call $$(cat $(location @near//cli:encode_public_key_bmc)) set_link_rx_height \\'\\{\\"link\\":\\"$$(cat $(location @icon//:btp_address))\\"\\,\\"height\\":$$(($$(cat $(location  @icon//:latest_block_height))+1))\\}\\' --nodeUrl $$(cat $(locations @near//:node_url)) --accountId $$(cat $(location @near//cli:encode_public_key_bmc)) > $@""",
    executable = True,
    local = True,
    tools = [
        "@icon//:btp_address",
        "@near//:node_url",
        "@near//cli:encode_public_key_bmc",
        "@near//cli:near_binary",
    ],

)

select({
    "@near//:localnet": exports_files(glob(["configs/localnet/*"])),
    "@near//:docker_localnet": exports_files(glob(["configs/localnet/*"])),
    "@near//:testnet": exports_files(glob(["configs/testnet/*"])),
    "@near//:mainnet": exports_files(glob(["configs/mainnet/*"])),
    "//conditions:default": []
})