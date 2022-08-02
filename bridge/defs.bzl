def bridge(name, SRC, DST):
    native.genrule(
        name = "set_link_%s" % DST,
        outs = ["set_link_%s.out" % DST],
        srcs = ["@%s//cli:set_%s_link" % (DST, SRC)],
        cmd = "echo 'done' > $@",
        local = True,
        executable = True,
    )

    native.genrule(
        name = "%s_to_%s" % (SRC, DST),
        outs = ["%s_to_%s.out" % (SRC, DST)],
        srcs = [
            ":set_link_%s" % DST,
            ":deploy_%s_bmr" % SRC,
        ],
        cmd = "echo 'done' > $@",
        local = True,
        executable = True,
    )

    configuration = "".join([
        """ .name = "%s2%s" | """ % (SRC, DST),
        " .src.address = $$src_address |",
        " .src.endpoint = [$$src_endpoint] |",
        " .src.offset = $$offset |",
        " .src.options = {} |",
        " .dst.address = $$dst_address |",
        " .dst.endpoint = [$$dst_endpoint] |",
        " .dst.options = {} |",
        " .dst.key_store = $$keystore |",
        " .dst.key_password = $$keypassword"
    ])

    native.genrule(
        name = "iconbridge_%s_configuration" % (DST),
        outs = ["iconbridge_%s_configuration.json" % (DST)],
        local = True,
        executable = True,
        srcs = [
            "@%s//:btp_address" % SRC,
            "@%s//:endpoint_docker" % SRC,
            "@%s//:latest_block_height" % SRC,
            "@%s//:btp_address" % DST,
            "@%s//:endpoint_docker" % DST,
            "@%s//cli:keysecret" % DST,
            "@btp//cmd/iconbridge:iconbridge",
            "@%s//:bmr_config_dir" % (DST),
            "@%s//cli:get_%s_bmr_keystore" % (DST, DST),
        ],
        cmd = "".join([
            "echo $$(cat $(location @%s//cli:get_%s_bmr_keystore)) > $$(cat $(location @%s//:bmr_config_dir))/keystore.json" % (DST, DST, DST),
            """
            jq <<<{} '
                .base_dir = "/iconbridge/%s" |
                .log_level = $$log_level |
                .console_level = "trace" |
                .log_writer.filename = "/iconbridge/%s/"+$$log_level+".log"|
                .relays = [ $$chain ]' --arg log_level "debug" \\
            """ % (SRC, SRC),
            """  --argjson chain "$$(jq <<<{} ' %s ' --arg src_address $$(cat $(location @%s//:btp_address)) --arg src_endpoint $$(cat $(location @%s//:endpoint_docker))""" % (configuration, SRC, SRC),
            """ --argjson offset $$(($$(cat $(location  @%s//:latest_block_height))+1)) """ % SRC,
            """ --arg dst_address  $$(cat $(location @%s//:btp_address)) """ % DST,
            """ --arg dst_endpoint  $$(cat $(location @%s//:endpoint_docker)) """ % DST,
            """ --arg keypassword $$(cat $(location @%s//cli:keysecret)) """ % DST,
            """ --argjson keystore \"$$(cat $$(cat $(location @%s//:bmr_config_dir))/keystore.json)\" """ % DST,
            """)" > $@"""
        ]),
    )

    native.genrule(
        name = "deploy_%s_bmr" % (SRC),
        outs = ["deploy_%s_bmr.out" % (SRC)],
        cmd = """
            docker run -d -v $$(cat $(location @%s//:bmr_config_dir)):/config bazel/cmd/iconbridge:iconbridge_image start --config "/config/iconbridge_%s_configuration.json" --key_password $$(cat $(location @%s//cli:keysecret))
            echo 'done'> \"$@\" """ % (SRC, SRC, DST),
        executable = True,
        output_to_bindir = True,
        srcs = [
            ":iconbridge_%s_configuration" % (DST),
            "@%s//:bmr_config_dir" % (SRC),
            "@%s//cli:keysecret" % DST,
        ],
    )
