if [ $# -eq 2 ]
then
    SOURCE=$1
    DESTINATION=$2
else
    echo "source and destination chain is mandatory"
    exit
fi

if !(command -v bazel &> /dev/null)
then
    echo "bazel is not found in path. Please follow https://bazel.build"
    exit
fi

bazel run @btp//cmd/iconbridge:iconbridge_image --incompatible_enable_cc_toolchain_resolution --extra_toolchains @zig_sdk//:linux_amd64_gnu.2.19_toolchain

bazel run "@btp//bridge:${SOURCE}_to_${DESTINATION}" --define near_network=testnet --define icon_network=lisbon
bazel run "@btp//bridge:${DESTINATION}_to_${SOURCE}" --define near_network=testnet --define icon_network=lisbon
