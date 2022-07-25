def bridge(name, chains):
    native.genrule(
        name = "set_link_%s" % chains[0],
        outs = ["set_link_%s.out" % chains[0]],
        srcs = ["@%s//cli:set_%s_link" % (chains[0], chains[1])],
        cmd = "echo 'done' > $@",
        local = True,
        executable = True,
    )
    native.genrule(
        name = "set_link_%s" % chains[1],
        outs = ["set_link_%s.out" % chains[1]],
        srcs = ["@%s//cli:set_%s_link" % (chains[1], chains[0])],
        cmd = "echo 'done' > $@",
        local = True,
        executable = True,
    )
    native.genrule(
        name = "%s_to_%s" % (chains[0], chains[1]),
        outs = ["%s_to_%s.out" % (chains[0], chains[1])],
        srcs = [
            ":set_link_%s" % chains[1],
            ":deploy_%s_bmr" % chains[0],
        ],
        cmd = "echo 'done' > $@",
        local = True,
        executable = True,
    )

    native.genrule(
        name = "%s_to_%s" % (chains[1], chains[0]),
        outs = ["%s_to_%s.out" % (chains[1], chains[0])],
        srcs = [
            ":set_link_%s" % chains[0],
            ":deploy_%s_bmr" % chains[1],
        ],
        cmd = "echo 'done' > $@",
        local = True,
        executable = True,
    )

    native.genrule(
        name = "iconbridge_%s_configuration" % (chains[0]),
        outs = ["iconbridge_%s_configuration.json" % (chains[0])],
        srcs = [
            "@btp//cmd/iconbridge:iconbridge",
            "@%s//:bmr_config_dir" % (chains[0]),
            "@%s//cli:add_%s_bmr" % (chains[1], chains[0]),
            "@%s//cli:transfer_amount_%s_address" % (chains[1], chains[1]),
            "@%s//cli:keysecret" % chains[1],
            "@%s//:endpoint_docker" % chains[0],
            "@%s//:endpoint_docker" % chains[1],
            "@%s//:btp_address" % chains[0],
            "@%s//:btp_address" % chains[1],
            "@%s//:latest_block_height" % chains[0],
        ],
        cmd = """echo $$(cat $(location @%s//cli:transfer_amount_%s_address)) > $$(cat $(location @%s//:bmr_config_dir))/keystore.json
        export ICONBRIDGE_SRC_OPTIONS=[mtaRootSize=8]
        export ICONBRIDGE_DST_OPTIONS=[mtaRootSize=8]
        export ICONBRIDGE_LOG_WRITER_MAXSIZE=1024
        export ICONBRIDGE_BASE_DIR=\"./config/data\"
        export ICONBRIDGE_OFFSET=$$(($$(cat $(location  @%s//:latest_block_height))+1))
        export ICONBRIDGE_LOG_WRITER_FILENAME=\"./config/log/iconbridge.log\"
        $(execpath @btp//cmd/iconbridge:iconbridge) --key_password $$(cat $(location @%s//cli:keysecret)) --key_store $$(cat $(location @%s//:bmr_config_dir))/keystore.json \
            --src.address $$(cat $(location @%s//:btp_address)) \
            --src.endpoint $$(cat $(location @%s//:endpoint_docker)) \
            --dst.address $$(cat $(location @%s//:btp_address)) \
            --dst.endpoint $$(cat $(location @%s//:endpoint_docker)) \
        save $@; cp $@ $$(cat $(location @%s//:bmr_config_dir))""" % (chains[1], chains[1], chains[0], chains[0], chains[1], chains[0], chains[0], chains[0], chains[1], chains[1], chains[0]),
        local = True,
        executable = True,
    )

    native.genrule(
        name = "iconbridge_%s_configuration" % (chains[1]),
        outs = ["iconbridge_%s_configuration.json" % (chains[1])],
        srcs = [
            "@btp//cmd/iconbridge:iconbridge",
            "@%s//:bmr_config_dir" % (chains[1]),
            "@%s//cli:add_%s_bmr" % (chains[0], chains[1]),
            "@%s//cli:transfer_amount_%s_address" % (chains[0], chains[0]),
            "@%s//cli:keysecret" % chains[0],
            "@%s//:endpoint_docker" % chains[1],
            "@%s//:endpoint_docker" % chains[0],
            "@%s//:btp_address" % chains[0],
            "@%s//:btp_address" % chains[1],
            "@%s//:latest_block_height" % chains[1],
        ],
        cmd = """echo $$(cat $(location @%s//cli:transfer_amount_%s_address)) > $$(cat $(location @%s//:bmr_config_dir))/keystore.json
        export ICONBRIDGE_SRC_OPTIONS=[mtaRootSize=8]
        export ICONBRIDGE_DST_OPTIONS=[mtaRootSize=8]
        export ICONBRIDGE_LOG_WRITER_MAXSIZE=1024
        export ICONBRIDGE_BASE_DIR=\"./config/data\"
        export ICONBRIDGE_OFFSET=$$(($$(cat $(location  @%s//:latest_block_height))+1))
        export ICONBRIDGE_LOG_WRITER_FILENAME=\"./config/log/iconbridge.log\"
        $(execpath @btp//cmd/iconbridge:iconbridge) --key_password $$(cat $(location @%s//cli:keysecret)) --key_store $$(cat $(location @%s//:bmr_config_dir))/keystore.json \
            --src.address $$(cat $(location @%s//:btp_address)) \
            --src.endpoint $$(cat $(location @%s//:endpoint_docker)) \
            --dst.address $$(cat $(location @%s//:btp_address)) \
            --dst.endpoint $$(cat $(location @%s//:endpoint_docker)) \
        save $@; cp $@ $$(cat $(location @%s//:bmr_config_dir))""" % (chains[0], chains[0], chains[1], chains[1], chains[0], chains[1], chains[1], chains[1], chains[0], chains[0], chains[1]),
        local = True,
        executable = True,
    )

    native.genrule(
        name = "deploy_%s_bmr" % (chains[0]),
        outs = ["deploy_%s_bmr.out" % (chains[0])],
        cmd = """
            docker run -d -v $$(cat $(location @%s//:bmr_config_dir)):/config bazel/cmd/iconbridge:iconbridge_image start --config "/config/iconbridge_%s_configuration.json" --key_password $$(cat $(location @%s//cli:keysecret))
            echo 'done'> \"$@\" """ % (chains[0], chains[0], chains[1]),
        executable = True,
        output_to_bindir = True,
        srcs = [
            ":iconbridge_%s_configuration" % (chains[0]),
            "@%s//:bmr_config_dir" % (chains[0]),
            "@%s//cli:keysecret" % chains[1],
        ],
    )

    native.genrule(
        name = "deploy_%s_bmr" % (chains[1]),
        outs = ["deploy_%s_bmr.out" % (chains[1])],
        cmd = """
            docker run -d -v $$(cat $(location @%s//:bmr_config_dir)):/config bazel/cmd/iconbridge:iconbridge_image start --config "/config/iconbridge_%s_configuration.json" --key_password $$(cat $(location @%s//cli:keysecret))
            echo 'done'> \"$@\" """ % (chains[1], chains[1], chains[0]),
        executable = True,
        output_to_bindir = True,
        srcs = [
            ":iconbridge_%s_configuration" % (chains[1]),
            "@%s//:bmr_config_dir" % (chains[1]),
            "@%s//cli:keysecret" % chains[0],
        ],
    )
