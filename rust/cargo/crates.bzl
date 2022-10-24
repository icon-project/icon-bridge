"""
@generated
cargo-raze generated Bazel file.

DO NOT EDIT! Replaced on runs of cargo-raze
"""

load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")  # buildifier: disable=load
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")  # buildifier: disable=load
load("@bazel_tools//tools/build_defs/repo:utils.bzl", "maybe")  # buildifier: disable=load

def raze_fetch_remote_crates():
    """This function defines a collection of repos and should be called in a WORKSPACE file"""
    maybe(
        http_archive,
        name = "raze__Inflector__0_11_4",
        url = "https://crates.io/api/v1/crates/Inflector/0.11.4/download",
        type = "tar.gz",
        sha256 = "fe438c63458706e03479442743baae6c88256498e6431708f6dfc520a26515d3",
        strip_prefix = "Inflector-0.11.4",
        build_file = Label("//cargo/remote:BUILD.Inflector-0.11.4.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__ahash__0_7_6",
        url = "https://crates.io/api/v1/crates/ahash/0.7.6/download",
        type = "tar.gz",
        sha256 = "fcb51a0695d8f838b1ee009b3fbf66bda078cd64590202a864a8f3e8c4315c47",
        strip_prefix = "ahash-0.7.6",
        build_file = Label("//cargo/remote:BUILD.ahash-0.7.6.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__android_system_properties__0_1_5",
        url = "https://crates.io/api/v1/crates/android_system_properties/0.1.5/download",
        type = "tar.gz",
        sha256 = "819e7219dbd41043ac279b19830f2efc897156490d7fd6ea916720117ee66311",
        strip_prefix = "android_system_properties-0.1.5",
        build_file = Label("//cargo/remote:BUILD.android_system_properties-0.1.5.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__arrayref__0_3_6",
        url = "https://crates.io/api/v1/crates/arrayref/0.3.6/download",
        type = "tar.gz",
        sha256 = "a4c527152e37cf757a3f78aae5a06fbeefdb07ccc535c980a3208ee3060dd544",
        strip_prefix = "arrayref-0.3.6",
        build_file = Label("//cargo/remote:BUILD.arrayref-0.3.6.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__arrayvec__0_5_2",
        url = "https://crates.io/api/v1/crates/arrayvec/0.5.2/download",
        type = "tar.gz",
        sha256 = "23b62fc65de8e4e7f52534fb52b0f3ed04746ae267519eef2a83941e8085068b",
        strip_prefix = "arrayvec-0.5.2",
        build_file = Label("//cargo/remote:BUILD.arrayvec-0.5.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__arrayvec__0_7_2",
        url = "https://crates.io/api/v1/crates/arrayvec/0.7.2/download",
        type = "tar.gz",
        sha256 = "8da52d66c7071e2e3fa2a1e5c6d088fec47b593032b254f5e980de8ea54454d6",
        strip_prefix = "arrayvec-0.7.2",
        build_file = Label("//cargo/remote:BUILD.arrayvec-0.7.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__autocfg__1_1_0",
        url = "https://crates.io/api/v1/crates/autocfg/1.1.0/download",
        type = "tar.gz",
        sha256 = "d468802bab17cbc0cc575e9b053f41e72aa36bfa6b7f55e3529ffa43161b97fa",
        strip_prefix = "autocfg-1.1.0",
        build_file = Label("//cargo/remote:BUILD.autocfg-1.1.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__base64__0_11_0",
        url = "https://crates.io/api/v1/crates/base64/0.11.0/download",
        type = "tar.gz",
        sha256 = "b41b7ea54a0c9d92199de89e20e58d49f02f8e699814ef3fdf266f6f748d15c7",
        strip_prefix = "base64-0.11.0",
        build_file = Label("//cargo/remote:BUILD.base64-0.11.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__base64__0_13_1",
        url = "https://crates.io/api/v1/crates/base64/0.13.1/download",
        type = "tar.gz",
        sha256 = "9e1b586273c5702936fe7b7d6896644d8be71e6314cfe09d3167c95f712589e8",
        strip_prefix = "base64-0.13.1",
        build_file = Label("//cargo/remote:BUILD.base64-0.13.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__bitvec__0_20_4",
        url = "https://crates.io/api/v1/crates/bitvec/0.20.4/download",
        type = "tar.gz",
        sha256 = "7774144344a4faa177370406a7ff5f1da24303817368584c6206c8303eb07848",
        strip_prefix = "bitvec-0.20.4",
        build_file = Label("//cargo/remote:BUILD.bitvec-0.20.4.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__blake2__0_9_2",
        url = "https://crates.io/api/v1/crates/blake2/0.9.2/download",
        type = "tar.gz",
        sha256 = "0a4e37d16930f5459780f5621038b6382b9bb37c19016f39fb6b5808d831f174",
        strip_prefix = "blake2-0.9.2",
        build_file = Label("//cargo/remote:BUILD.blake2-0.9.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__block_buffer__0_10_3",
        url = "https://crates.io/api/v1/crates/block-buffer/0.10.3/download",
        type = "tar.gz",
        sha256 = "69cce20737498f97b993470a6e536b8523f0af7892a4f928cceb1ac5e52ebe7e",
        strip_prefix = "block-buffer-0.10.3",
        build_file = Label("//cargo/remote:BUILD.block-buffer-0.10.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__block_buffer__0_9_0",
        url = "https://crates.io/api/v1/crates/block-buffer/0.9.0/download",
        type = "tar.gz",
        sha256 = "4152116fd6e9dadb291ae18fc1ec3575ed6d84c29642d97890f4b4a3417297e4",
        strip_prefix = "block-buffer-0.9.0",
        build_file = Label("//cargo/remote:BUILD.block-buffer-0.9.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__borsh__0_9_3",
        url = "https://crates.io/api/v1/crates/borsh/0.9.3/download",
        type = "tar.gz",
        sha256 = "15bf3650200d8bffa99015595e10f1fbd17de07abbc25bb067da79e769939bfa",
        strip_prefix = "borsh-0.9.3",
        build_file = Label("//cargo/remote:BUILD.borsh-0.9.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__borsh_derive__0_9_3",
        url = "https://crates.io/api/v1/crates/borsh-derive/0.9.3/download",
        type = "tar.gz",
        sha256 = "6441c552f230375d18e3cc377677914d2ca2b0d36e52129fe15450a2dce46775",
        strip_prefix = "borsh-derive-0.9.3",
        build_file = Label("//cargo/remote:BUILD.borsh-derive-0.9.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__borsh_derive_internal__0_9_3",
        url = "https://crates.io/api/v1/crates/borsh-derive-internal/0.9.3/download",
        type = "tar.gz",
        sha256 = "5449c28a7b352f2d1e592a8a28bf139bc71afb0764a14f3c02500935d8c44065",
        strip_prefix = "borsh-derive-internal-0.9.3",
        build_file = Label("//cargo/remote:BUILD.borsh-derive-internal-0.9.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__borsh_schema_derive_internal__0_9_3",
        url = "https://crates.io/api/v1/crates/borsh-schema-derive-internal/0.9.3/download",
        type = "tar.gz",
        sha256 = "cdbd5696d8bfa21d53d9fe39a714a18538bad11492a42d066dbbc395fb1951c0",
        strip_prefix = "borsh-schema-derive-internal-0.9.3",
        build_file = Label("//cargo/remote:BUILD.borsh-schema-derive-internal-0.9.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__bs58__0_4_0",
        url = "https://crates.io/api/v1/crates/bs58/0.4.0/download",
        type = "tar.gz",
        sha256 = "771fe0050b883fcc3ea2359b1a96bcfbc090b7116eae7c3c512c7a083fdf23d3",
        strip_prefix = "bs58-0.4.0",
        build_file = Label("//cargo/remote:BUILD.bs58-0.4.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__bumpalo__3_11_1",
        url = "https://crates.io/api/v1/crates/bumpalo/3.11.1/download",
        type = "tar.gz",
        sha256 = "572f695136211188308f16ad2ca5c851a712c464060ae6974944458eb83880ba",
        strip_prefix = "bumpalo-3.11.1",
        build_file = Label("//cargo/remote:BUILD.bumpalo-3.11.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__byte_slice_cast__1_2_1",
        url = "https://crates.io/api/v1/crates/byte-slice-cast/1.2.1/download",
        type = "tar.gz",
        sha256 = "87c5fdd0166095e1d463fc6cc01aa8ce547ad77a4e84d42eb6762b084e28067e",
        strip_prefix = "byte-slice-cast-1.2.1",
        build_file = Label("//cargo/remote:BUILD.byte-slice-cast-1.2.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__byteorder__1_4_3",
        url = "https://crates.io/api/v1/crates/byteorder/1.4.3/download",
        type = "tar.gz",
        sha256 = "14c189c53d098945499cdfa7ecc63567cf3886b3332b312a5b4585d8d3a6a610",
        strip_prefix = "byteorder-1.4.3",
        build_file = Label("//cargo/remote:BUILD.byteorder-1.4.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__bytes__1_2_1",
        url = "https://crates.io/api/v1/crates/bytes/1.2.1/download",
        type = "tar.gz",
        sha256 = "ec8a7b6a70fde80372154c65702f00a0f56f3e1c36abbc6c440484be248856db",
        strip_prefix = "bytes-1.2.1",
        build_file = Label("//cargo/remote:BUILD.bytes-1.2.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__bytesize__1_1_0",
        url = "https://crates.io/api/v1/crates/bytesize/1.1.0/download",
        type = "tar.gz",
        sha256 = "6c58ec36aac5066d5ca17df51b3e70279f5670a72102f5752cb7e7c856adfc70",
        strip_prefix = "bytesize-1.1.0",
        build_file = Label("//cargo/remote:BUILD.bytesize-1.1.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__c2_chacha__0_3_3",
        url = "https://crates.io/api/v1/crates/c2-chacha/0.3.3/download",
        type = "tar.gz",
        sha256 = "d27dae93fe7b1e0424dc57179ac396908c26b035a87234809f5c4dfd1b47dc80",
        strip_prefix = "c2-chacha-0.3.3",
        build_file = Label("//cargo/remote:BUILD.c2-chacha-0.3.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__cc__1_0_73",
        url = "https://crates.io/api/v1/crates/cc/1.0.73/download",
        type = "tar.gz",
        sha256 = "2fff2a6927b3bb87f9595d67196a70493f627687a71d87a0d692242c33f58c11",
        strip_prefix = "cc-1.0.73",
        build_file = Label("//cargo/remote:BUILD.cc-1.0.73.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__cfg_if__0_1_10",
        url = "https://crates.io/api/v1/crates/cfg-if/0.1.10/download",
        type = "tar.gz",
        sha256 = "4785bdd1c96b2a846b2bd7cc02e86b6b3dbf14e7e53446c4f54c92a361040822",
        strip_prefix = "cfg-if-0.1.10",
        build_file = Label("//cargo/remote:BUILD.cfg-if-0.1.10.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__cfg_if__1_0_0",
        url = "https://crates.io/api/v1/crates/cfg-if/1.0.0/download",
        type = "tar.gz",
        sha256 = "baf1de4339761588bc0619e3cbc0120ee582ebb74b53b4efbf79117bd2da40fd",
        strip_prefix = "cfg-if-1.0.0",
        build_file = Label("//cargo/remote:BUILD.cfg-if-1.0.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__chrono__0_4_22",
        url = "https://crates.io/api/v1/crates/chrono/0.4.22/download",
        type = "tar.gz",
        sha256 = "bfd4d1b31faaa3a89d7934dbded3111da0d2ef28e3ebccdb4f0179f5929d1ef1",
        strip_prefix = "chrono-0.4.22",
        build_file = Label("//cargo/remote:BUILD.chrono-0.4.22.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__cipher__0_2_5",
        url = "https://crates.io/api/v1/crates/cipher/0.2.5/download",
        type = "tar.gz",
        sha256 = "12f8e7987cbd042a63249497f41aed09f8e65add917ea6566effbc56578d6801",
        strip_prefix = "cipher-0.2.5",
        build_file = Label("//cargo/remote:BUILD.cipher-0.2.5.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__codespan_reporting__0_11_1",
        url = "https://crates.io/api/v1/crates/codespan-reporting/0.11.1/download",
        type = "tar.gz",
        sha256 = "3538270d33cc669650c4b093848450d380def10c331d38c768e34cac80576e6e",
        strip_prefix = "codespan-reporting-0.11.1",
        build_file = Label("//cargo/remote:BUILD.codespan-reporting-0.11.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__convert_case__0_4_0",
        url = "https://crates.io/api/v1/crates/convert_case/0.4.0/download",
        type = "tar.gz",
        sha256 = "6245d59a3e82a7fc217c5828a6692dbc6dfb63a0c8c90495621f7b9d79704a0e",
        strip_prefix = "convert_case-0.4.0",
        build_file = Label("//cargo/remote:BUILD.convert_case-0.4.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__core_foundation_sys__0_8_3",
        url = "https://crates.io/api/v1/crates/core-foundation-sys/0.8.3/download",
        type = "tar.gz",
        sha256 = "5827cebf4670468b8772dd191856768aedcb1b0278a04f989f7766351917b9dc",
        strip_prefix = "core-foundation-sys-0.8.3",
        build_file = Label("//cargo/remote:BUILD.core-foundation-sys-0.8.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__cpufeatures__0_2_5",
        url = "https://crates.io/api/v1/crates/cpufeatures/0.2.5/download",
        type = "tar.gz",
        sha256 = "28d997bd5e24a5928dd43e46dc529867e207907fe0b239c3477d924f7f2ca320",
        strip_prefix = "cpufeatures-0.2.5",
        build_file = Label("//cargo/remote:BUILD.cpufeatures-0.2.5.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__crunchy__0_2_2",
        url = "https://crates.io/api/v1/crates/crunchy/0.2.2/download",
        type = "tar.gz",
        sha256 = "7a81dae078cea95a014a339291cec439d2f232ebe854a9d672b796c6afafa9b7",
        strip_prefix = "crunchy-0.2.2",
        build_file = Label("//cargo/remote:BUILD.crunchy-0.2.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__crypto_common__0_1_6",
        url = "https://crates.io/api/v1/crates/crypto-common/0.1.6/download",
        type = "tar.gz",
        sha256 = "1bfb12502f3fc46cca1bb51ac28df9d618d813cdc3d2f25b9fe775a34af26bb3",
        strip_prefix = "crypto-common-0.1.6",
        build_file = Label("//cargo/remote:BUILD.crypto-common-0.1.6.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__crypto_mac__0_8_0",
        url = "https://crates.io/api/v1/crates/crypto-mac/0.8.0/download",
        type = "tar.gz",
        sha256 = "b584a330336237c1eecd3e94266efb216c56ed91225d634cb2991c5f3fd1aeab",
        strip_prefix = "crypto-mac-0.8.0",
        build_file = Label("//cargo/remote:BUILD.crypto-mac-0.8.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__curve25519_dalek__3_2_1",
        url = "https://crates.io/api/v1/crates/curve25519-dalek/3.2.1/download",
        type = "tar.gz",
        sha256 = "90f9d052967f590a76e62eb387bd0bbb1b000182c3cefe5364db6b7211651bc0",
        strip_prefix = "curve25519-dalek-3.2.1",
        build_file = Label("//cargo/remote:BUILD.curve25519-dalek-3.2.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__cxx__1_0_80",
        url = "https://crates.io/api/v1/crates/cxx/1.0.80/download",
        type = "tar.gz",
        sha256 = "6b7d4e43b25d3c994662706a1d4fcfc32aaa6afd287502c111b237093bb23f3a",
        strip_prefix = "cxx-1.0.80",
        build_file = Label("//cargo/remote:BUILD.cxx-1.0.80.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__cxx_build__1_0_80",
        url = "https://crates.io/api/v1/crates/cxx-build/1.0.80/download",
        type = "tar.gz",
        sha256 = "84f8829ddc213e2c1368e51a2564c552b65a8cb6a28f31e576270ac81d5e5827",
        strip_prefix = "cxx-build-1.0.80",
        build_file = Label("//cargo/remote:BUILD.cxx-build-1.0.80.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__cxxbridge_flags__1_0_80",
        url = "https://crates.io/api/v1/crates/cxxbridge-flags/1.0.80/download",
        type = "tar.gz",
        sha256 = "e72537424b474af1460806647c41d4b6d35d09ef7fe031c5c2fa5766047cc56a",
        strip_prefix = "cxxbridge-flags-1.0.80",
        build_file = Label("//cargo/remote:BUILD.cxxbridge-flags-1.0.80.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__cxxbridge_macro__1_0_80",
        url = "https://crates.io/api/v1/crates/cxxbridge-macro/1.0.80/download",
        type = "tar.gz",
        sha256 = "309e4fb93eed90e1e14bea0da16b209f81813ba9fc7830c20ed151dd7bc0a4d7",
        strip_prefix = "cxxbridge-macro-1.0.80",
        build_file = Label("//cargo/remote:BUILD.cxxbridge-macro-1.0.80.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__derive_more__0_99_17",
        url = "https://crates.io/api/v1/crates/derive_more/0.99.17/download",
        type = "tar.gz",
        sha256 = "4fb810d30a7c1953f91334de7244731fc3f3c10d7fe163338a35b9f640960321",
        strip_prefix = "derive_more-0.99.17",
        build_file = Label("//cargo/remote:BUILD.derive_more-0.99.17.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__digest__0_10_5",
        url = "https://crates.io/api/v1/crates/digest/0.10.5/download",
        type = "tar.gz",
        sha256 = "adfbc57365a37acbd2ebf2b64d7e69bb766e2fea813521ed536f5d0520dcf86c",
        strip_prefix = "digest-0.10.5",
        build_file = Label("//cargo/remote:BUILD.digest-0.10.5.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__digest__0_9_0",
        url = "https://crates.io/api/v1/crates/digest/0.9.0/download",
        type = "tar.gz",
        sha256 = "d3dd60d1080a57a05ab032377049e0591415d2b31afd7028356dbf3cc6dcb066",
        strip_prefix = "digest-0.9.0",
        build_file = Label("//cargo/remote:BUILD.digest-0.9.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__dyn_clone__1_0_9",
        url = "https://crates.io/api/v1/crates/dyn-clone/1.0.9/download",
        type = "tar.gz",
        sha256 = "4f94fa09c2aeea5b8839e414b7b841bf429fd25b9c522116ac97ee87856d88b2",
        strip_prefix = "dyn-clone-1.0.9",
        build_file = Label("//cargo/remote:BUILD.dyn-clone-1.0.9.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__easy_ext__0_2_9",
        url = "https://crates.io/api/v1/crates/easy-ext/0.2.9/download",
        type = "tar.gz",
        sha256 = "53aff6fdc1b181225acdcb5b14c47106726fd8e486707315b1b138baed68ee31",
        strip_prefix = "easy-ext-0.2.9",
        build_file = Label("//cargo/remote:BUILD.easy-ext-0.2.9.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__ed25519__1_5_2",
        url = "https://crates.io/api/v1/crates/ed25519/1.5.2/download",
        type = "tar.gz",
        sha256 = "1e9c280362032ea4203659fc489832d0204ef09f247a0506f170dafcac08c369",
        strip_prefix = "ed25519-1.5.2",
        build_file = Label("//cargo/remote:BUILD.ed25519-1.5.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__ed25519_dalek__1_0_1",
        url = "https://crates.io/api/v1/crates/ed25519-dalek/1.0.1/download",
        type = "tar.gz",
        sha256 = "c762bae6dcaf24c4c84667b8579785430908723d5c889f469d76a41d59cc7a9d",
        strip_prefix = "ed25519-dalek-1.0.1",
        build_file = Label("//cargo/remote:BUILD.ed25519-dalek-1.0.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__fixed_hash__0_7_0",
        url = "https://crates.io/api/v1/crates/fixed-hash/0.7.0/download",
        type = "tar.gz",
        sha256 = "cfcf0ed7fe52a17a03854ec54a9f76d6d84508d1c0e66bc1793301c73fc8493c",
        strip_prefix = "fixed-hash-0.7.0",
        build_file = Label("//cargo/remote:BUILD.fixed-hash-0.7.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__funty__1_1_0",
        url = "https://crates.io/api/v1/crates/funty/1.1.0/download",
        type = "tar.gz",
        sha256 = "fed34cd105917e91daa4da6b3728c47b068749d6a62c59811f06ed2ac71d9da7",
        strip_prefix = "funty-1.1.0",
        build_file = Label("//cargo/remote:BUILD.funty-1.1.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__generic_array__0_14_6",
        url = "https://crates.io/api/v1/crates/generic-array/0.14.6/download",
        type = "tar.gz",
        sha256 = "bff49e947297f3312447abdca79f45f4738097cc82b06e72054d2223f601f1b9",
        strip_prefix = "generic-array-0.14.6",
        build_file = Label("//cargo/remote:BUILD.generic-array-0.14.6.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__getrandom__0_1_16",
        url = "https://crates.io/api/v1/crates/getrandom/0.1.16/download",
        type = "tar.gz",
        sha256 = "8fc3cb4d91f53b50155bdcfd23f6a4c39ae1969c2ae85982b135750cccaf5fce",
        strip_prefix = "getrandom-0.1.16",
        build_file = Label("//cargo/remote:BUILD.getrandom-0.1.16.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__getrandom__0_2_8",
        url = "https://crates.io/api/v1/crates/getrandom/0.2.8/download",
        type = "tar.gz",
        sha256 = "c05aeb6a22b8f62540c194aac980f2115af067bfe15a0734d7277a768d396b31",
        strip_prefix = "getrandom-0.2.8",
        build_file = Label("//cargo/remote:BUILD.getrandom-0.2.8.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__hashbrown__0_11_2",
        url = "https://crates.io/api/v1/crates/hashbrown/0.11.2/download",
        type = "tar.gz",
        sha256 = "ab5ef0d4909ef3724cc8cce6ccc8572c5c817592e9285f5464f8e86f8bd3726e",
        strip_prefix = "hashbrown-0.11.2",
        build_file = Label("//cargo/remote:BUILD.hashbrown-0.11.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__heck__0_4_0",
        url = "https://crates.io/api/v1/crates/heck/0.4.0/download",
        type = "tar.gz",
        sha256 = "2540771e65fc8cb83cd6e8a237f70c319bd5c29f78ed1084ba5d50eeac86f7f9",
        strip_prefix = "heck-0.4.0",
        build_file = Label("//cargo/remote:BUILD.heck-0.4.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__hex__0_4_3",
        url = "https://crates.io/api/v1/crates/hex/0.4.3/download",
        type = "tar.gz",
        sha256 = "7f24254aa9a54b5c858eaee2f5bccdb46aaf0e486a595ed5fd8f86ba55232a70",
        strip_prefix = "hex-0.4.3",
        build_file = Label("//cargo/remote:BUILD.hex-0.4.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__iana_time_zone__0_1_51",
        url = "https://crates.io/api/v1/crates/iana-time-zone/0.1.51/download",
        type = "tar.gz",
        sha256 = "f5a6ef98976b22b3b7f2f3a806f858cb862044cfa66805aa3ad84cb3d3b785ed",
        strip_prefix = "iana-time-zone-0.1.51",
        build_file = Label("//cargo/remote:BUILD.iana-time-zone-0.1.51.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__iana_time_zone_haiku__0_1_1",
        url = "https://crates.io/api/v1/crates/iana-time-zone-haiku/0.1.1/download",
        type = "tar.gz",
        sha256 = "0703ae284fc167426161c2e3f1da3ea71d94b21bedbcc9494e92b28e334e3dca",
        strip_prefix = "iana-time-zone-haiku-0.1.1",
        build_file = Label("//cargo/remote:BUILD.iana-time-zone-haiku-0.1.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__impl_codec__0_5_1",
        url = "https://crates.io/api/v1/crates/impl-codec/0.5.1/download",
        type = "tar.gz",
        sha256 = "161ebdfec3c8e3b52bf61c4f3550a1eea4f9579d10dc1b936f3171ebdcd6c443",
        strip_prefix = "impl-codec-0.5.1",
        build_file = Label("//cargo/remote:BUILD.impl-codec-0.5.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__impl_trait_for_tuples__0_2_2",
        url = "https://crates.io/api/v1/crates/impl-trait-for-tuples/0.2.2/download",
        type = "tar.gz",
        sha256 = "11d7a9f6330b71fea57921c9b61c47ee6e84f72d394754eff6163ae67e7395eb",
        strip_prefix = "impl-trait-for-tuples-0.2.2",
        build_file = Label("//cargo/remote:BUILD.impl-trait-for-tuples-0.2.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__itoa__1_0_4",
        url = "https://crates.io/api/v1/crates/itoa/1.0.4/download",
        type = "tar.gz",
        sha256 = "4217ad341ebadf8d8e724e264f13e593e0648f5b3e94b3896a5df283be015ecc",
        strip_prefix = "itoa-1.0.4",
        build_file = Label("//cargo/remote:BUILD.itoa-1.0.4.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__js_sys__0_3_60",
        url = "https://crates.io/api/v1/crates/js-sys/0.3.60/download",
        type = "tar.gz",
        sha256 = "49409df3e3bf0856b916e2ceaca09ee28e6871cf7d9ce97a692cacfdb2a25a47",
        strip_prefix = "js-sys-0.3.60",
        build_file = Label("//cargo/remote:BUILD.js-sys-0.3.60.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__keccak__0_1_2",
        url = "https://crates.io/api/v1/crates/keccak/0.1.2/download",
        type = "tar.gz",
        sha256 = "f9b7d56ba4a8344d6be9729995e6b06f928af29998cdf79fe390cbf6b1fee838",
        strip_prefix = "keccak-0.1.2",
        build_file = Label("//cargo/remote:BUILD.keccak-0.1.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__lazy_static__1_4_0",
        url = "https://crates.io/api/v1/crates/lazy_static/1.4.0/download",
        type = "tar.gz",
        sha256 = "e2abad23fbc42b3700f2f279844dc832adb2b2eb069b2df918f455c4e18cc646",
        strip_prefix = "lazy_static-1.4.0",
        build_file = Label("//cargo/remote:BUILD.lazy_static-1.4.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__libc__0_2_135",
        url = "https://crates.io/api/v1/crates/libc/0.2.135/download",
        type = "tar.gz",
        sha256 = "68783febc7782c6c5cb401fbda4de5a9898be1762314da0bb2c10ced61f18b0c",
        strip_prefix = "libc-0.2.135",
        build_file = Label("//cargo/remote:BUILD.libc-0.2.135.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__link_cplusplus__1_0_7",
        url = "https://crates.io/api/v1/crates/link-cplusplus/1.0.7/download",
        type = "tar.gz",
        sha256 = "9272ab7b96c9046fbc5bc56c06c117cb639fe2d509df0c421cad82d2915cf369",
        strip_prefix = "link-cplusplus-1.0.7",
        build_file = Label("//cargo/remote:BUILD.link-cplusplus-1.0.7.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__log__0_4_17",
        url = "https://crates.io/api/v1/crates/log/0.4.17/download",
        type = "tar.gz",
        sha256 = "abb12e687cfb44aa40f41fc3978ef76448f9b6038cad6aef4259d3c095a2382e",
        strip_prefix = "log-0.4.17",
        build_file = Label("//cargo/remote:BUILD.log-0.4.17.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__memory_units__0_4_0",
        url = "https://crates.io/api/v1/crates/memory_units/0.4.0/download",
        type = "tar.gz",
        sha256 = "8452105ba047068f40ff7093dd1d9da90898e63dd61736462e9cdda6a90ad3c3",
        strip_prefix = "memory_units-0.4.0",
        build_file = Label("//cargo/remote:BUILD.memory_units-0.4.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_account_id__0_14_0",
        url = "https://crates.io/api/v1/crates/near-account-id/0.14.0/download",
        type = "tar.gz",
        sha256 = "71d258582a1878e6db67400b0504a5099db85718d22c2e07f747fe1706ae7150",
        strip_prefix = "near-account-id-0.14.0",
        build_file = Label("//cargo/remote:BUILD.near-account-id-0.14.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_contract_standards__4_1_0_pre_3",
        url = "https://crates.io/api/v1/crates/near-contract-standards/4.1.0-pre.3/download",
        type = "tar.gz",
        sha256 = "719edff4fe1558fe68fc1325a7d668288d1e2dd47177c37d5d43028ca59f24a5",
        strip_prefix = "near-contract-standards-4.1.0-pre.3",
        build_file = Label("//cargo/remote:BUILD.near-contract-standards-4.1.0-pre.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_crypto__0_14_0",
        url = "https://crates.io/api/v1/crates/near-crypto/0.14.0/download",
        type = "tar.gz",
        sha256 = "1e75673d69fd7365508f3d32483669fe45b03bfb34e4d9363e90adae9dfb416c",
        strip_prefix = "near-crypto-0.14.0",
        build_file = Label("//cargo/remote:BUILD.near-crypto-0.14.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_primitives__0_14_0",
        url = "https://crates.io/api/v1/crates/near-primitives/0.14.0/download",
        type = "tar.gz",
        sha256 = "8ad1a9a1640539c81f065425c31bffcfbf6b31ef1aeaade59ce905f5df6ac860",
        strip_prefix = "near-primitives-0.14.0",
        build_file = Label("//cargo/remote:BUILD.near-primitives-0.14.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_primitives_core__0_14_0",
        url = "https://crates.io/api/v1/crates/near-primitives-core/0.14.0/download",
        type = "tar.gz",
        sha256 = "91d508f0fc340f6461e4e256417685720d3c4c00bb5a939b105160e49137caba",
        strip_prefix = "near-primitives-core-0.14.0",
        build_file = Label("//cargo/remote:BUILD.near-primitives-core-0.14.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_rpc_error_core__0_14_0",
        url = "https://crates.io/api/v1/crates/near-rpc-error-core/0.14.0/download",
        type = "tar.gz",
        sha256 = "93ee0b41c75ef859c193a8ff1dadfa0c8207bc0ac447cc22259721ad769a1408",
        strip_prefix = "near-rpc-error-core-0.14.0",
        build_file = Label("//cargo/remote:BUILD.near-rpc-error-core-0.14.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_rpc_error_macro__0_14_0",
        url = "https://crates.io/api/v1/crates/near-rpc-error-macro/0.14.0/download",
        type = "tar.gz",
        sha256 = "8e837bd4bacd807073ec5ceb85708da7f721b46a4c2a978de86027fb0034ce31",
        strip_prefix = "near-rpc-error-macro-0.14.0",
        build_file = Label("//cargo/remote:BUILD.near-rpc-error-macro-0.14.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_sdk__4_1_0_pre_3",
        url = "https://crates.io/api/v1/crates/near-sdk/4.1.0-pre.3/download",
        type = "tar.gz",
        sha256 = "c5950d57f8d412a6603eac3097d3f6d528d3067ee09bd347cd8e7448bb7215b6",
        strip_prefix = "near-sdk-4.1.0-pre.3",
        build_file = Label("//cargo/remote:BUILD.near-sdk-4.1.0-pre.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_sdk_macros__4_1_0_pre_3",
        url = "https://crates.io/api/v1/crates/near-sdk-macros/4.1.0-pre.3/download",
        type = "tar.gz",
        sha256 = "a31526ef660a7442f216a1da5ce19e844185351013e29d4f1f86484ddab3b7e6",
        strip_prefix = "near-sdk-macros-4.1.0-pre.3",
        build_file = Label("//cargo/remote:BUILD.near-sdk-macros-4.1.0-pre.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_sys__0_2_0",
        url = "https://crates.io/api/v1/crates/near-sys/0.2.0/download",
        type = "tar.gz",
        sha256 = "e307313276eaeced2ca95740b5639e1f3125b7c97f0a1151809d105f1aa8c6d3",
        strip_prefix = "near-sys-0.2.0",
        build_file = Label("//cargo/remote:BUILD.near-sys-0.2.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_vm_errors__0_14_0",
        url = "https://crates.io/api/v1/crates/near-vm-errors/0.14.0/download",
        type = "tar.gz",
        sha256 = "d0da466a30f0446639cbd788c30865086fac3e8dcb07a79e51d2b0775ed4261e",
        strip_prefix = "near-vm-errors-0.14.0",
        build_file = Label("//cargo/remote:BUILD.near-vm-errors-0.14.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__near_vm_logic__0_14_0",
        url = "https://crates.io/api/v1/crates/near-vm-logic/0.14.0/download",
        type = "tar.gz",
        sha256 = "81b534828419bacbf1f7b11ef7b00420f248c548c485d3f0cfda8bb6931152f2",
        strip_prefix = "near-vm-logic-0.14.0",
        build_file = Label("//cargo/remote:BUILD.near-vm-logic-0.14.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__num_bigint__0_3_3",
        url = "https://crates.io/api/v1/crates/num-bigint/0.3.3/download",
        type = "tar.gz",
        sha256 = "5f6f7833f2cbf2360a6cfd58cd41a53aa7a90bd4c202f5b1c7dd2ed73c57b2c3",
        strip_prefix = "num-bigint-0.3.3",
        build_file = Label("//cargo/remote:BUILD.num-bigint-0.3.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__num_integer__0_1_45",
        url = "https://crates.io/api/v1/crates/num-integer/0.1.45/download",
        type = "tar.gz",
        sha256 = "225d3389fb3509a24c93f5c29eb6bde2586b98d9f016636dff58d7c6f7569cd9",
        strip_prefix = "num-integer-0.1.45",
        build_file = Label("//cargo/remote:BUILD.num-integer-0.1.45.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__num_rational__0_3_2",
        url = "https://crates.io/api/v1/crates/num-rational/0.3.2/download",
        type = "tar.gz",
        sha256 = "12ac428b1cb17fce6f731001d307d351ec70a6d202fc2e60f7d4c5e42d8f4f07",
        strip_prefix = "num-rational-0.3.2",
        build_file = Label("//cargo/remote:BUILD.num-rational-0.3.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__num_traits__0_2_15",
        url = "https://crates.io/api/v1/crates/num-traits/0.2.15/download",
        type = "tar.gz",
        sha256 = "578ede34cf02f8924ab9447f50c28075b4d3e5b269972345e7e0372b38c6cdcd",
        strip_prefix = "num-traits-0.2.15",
        build_file = Label("//cargo/remote:BUILD.num-traits-0.2.15.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__once_cell__1_15_0",
        url = "https://crates.io/api/v1/crates/once_cell/1.15.0/download",
        type = "tar.gz",
        sha256 = "e82dad04139b71a90c080c8463fe0dc7902db5192d939bd0950f074d014339e1",
        strip_prefix = "once_cell-1.15.0",
        build_file = Label("//cargo/remote:BUILD.once_cell-1.15.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__opaque_debug__0_3_0",
        url = "https://crates.io/api/v1/crates/opaque-debug/0.3.0/download",
        type = "tar.gz",
        sha256 = "624a8340c38c1b80fd549087862da4ba43e08858af025b236e509b6649fc13d5",
        strip_prefix = "opaque-debug-0.3.0",
        build_file = Label("//cargo/remote:BUILD.opaque-debug-0.3.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__parity_scale_codec__2_3_1",
        url = "https://crates.io/api/v1/crates/parity-scale-codec/2.3.1/download",
        type = "tar.gz",
        sha256 = "373b1a4c1338d9cd3d1fa53b3a11bdab5ab6bd80a20f7f7becd76953ae2be909",
        strip_prefix = "parity-scale-codec-2.3.1",
        build_file = Label("//cargo/remote:BUILD.parity-scale-codec-2.3.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__parity_scale_codec_derive__2_3_1",
        url = "https://crates.io/api/v1/crates/parity-scale-codec-derive/2.3.1/download",
        type = "tar.gz",
        sha256 = "1557010476e0595c9b568d16dcfb81b93cdeb157612726f5170d31aa707bed27",
        strip_prefix = "parity-scale-codec-derive-2.3.1",
        build_file = Label("//cargo/remote:BUILD.parity-scale-codec-derive-2.3.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__parity_secp256k1__0_7_0",
        url = "https://crates.io/api/v1/crates/parity-secp256k1/0.7.0/download",
        type = "tar.gz",
        sha256 = "4fca4f82fccae37e8bbdaeb949a4a218a1bbc485d11598f193d2a908042e5fc1",
        strip_prefix = "parity-secp256k1-0.7.0",
        build_file = Label("//cargo/remote:BUILD.parity-secp256k1-0.7.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__ppv_lite86__0_2_16",
        url = "https://crates.io/api/v1/crates/ppv-lite86/0.2.16/download",
        type = "tar.gz",
        sha256 = "eb9f9e6e233e5c4a35559a617bf40a4ec447db2e84c20b55a6f83167b7e57872",
        strip_prefix = "ppv-lite86-0.2.16",
        build_file = Label("//cargo/remote:BUILD.ppv-lite86-0.2.16.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__primitive_types__0_10_1",
        url = "https://crates.io/api/v1/crates/primitive-types/0.10.1/download",
        type = "tar.gz",
        sha256 = "05e4722c697a58a99d5d06a08c30821d7c082a4632198de1eaa5a6c22ef42373",
        strip_prefix = "primitive-types-0.10.1",
        build_file = Label("//cargo/remote:BUILD.primitive-types-0.10.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__proc_macro_crate__0_1_5",
        url = "https://crates.io/api/v1/crates/proc-macro-crate/0.1.5/download",
        type = "tar.gz",
        sha256 = "1d6ea3c4595b96363c13943497db34af4460fb474a95c43f4446ad341b8c9785",
        strip_prefix = "proc-macro-crate-0.1.5",
        build_file = Label("//cargo/remote:BUILD.proc-macro-crate-0.1.5.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__proc_macro_crate__1_2_1",
        url = "https://crates.io/api/v1/crates/proc-macro-crate/1.2.1/download",
        type = "tar.gz",
        sha256 = "eda0fc3b0fb7c975631757e14d9049da17374063edb6ebbcbc54d880d4fe94e9",
        strip_prefix = "proc-macro-crate-1.2.1",
        build_file = Label("//cargo/remote:BUILD.proc-macro-crate-1.2.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__proc_macro2__1_0_47",
        url = "https://crates.io/api/v1/crates/proc-macro2/1.0.47/download",
        type = "tar.gz",
        sha256 = "5ea3d908b0e36316caf9e9e2c4625cdde190a7e6f440d794667ed17a1855e725",
        strip_prefix = "proc-macro2-1.0.47",
        build_file = Label("//cargo/remote:BUILD.proc-macro2-1.0.47.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__quote__1_0_21",
        url = "https://crates.io/api/v1/crates/quote/1.0.21/download",
        type = "tar.gz",
        sha256 = "bbe448f377a7d6961e30f5955f9b8d106c3f5e449d493ee1b125c1d43c2b5179",
        strip_prefix = "quote-1.0.21",
        build_file = Label("//cargo/remote:BUILD.quote-1.0.21.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__radium__0_6_2",
        url = "https://crates.io/api/v1/crates/radium/0.6.2/download",
        type = "tar.gz",
        sha256 = "643f8f41a8ebc4c5dc4515c82bb8abd397b527fc20fd681b7c011c2aee5d44fb",
        strip_prefix = "radium-0.6.2",
        build_file = Label("//cargo/remote:BUILD.radium-0.6.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rand__0_7_3",
        url = "https://crates.io/api/v1/crates/rand/0.7.3/download",
        type = "tar.gz",
        sha256 = "6a6b1679d49b24bbfe0c803429aa1874472f50d9b363131f0e89fc356b544d03",
        strip_prefix = "rand-0.7.3",
        build_file = Label("//cargo/remote:BUILD.rand-0.7.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rand__0_8_5",
        url = "https://crates.io/api/v1/crates/rand/0.8.5/download",
        type = "tar.gz",
        sha256 = "34af8d1a0e25924bc5b7c43c079c942339d8f0a8b57c39049bef581b46327404",
        strip_prefix = "rand-0.8.5",
        build_file = Label("//cargo/remote:BUILD.rand-0.8.5.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rand_chacha__0_2_2",
        url = "https://crates.io/api/v1/crates/rand_chacha/0.2.2/download",
        type = "tar.gz",
        sha256 = "f4c8ed856279c9737206bf725bf36935d8666ead7aa69b52be55af369d193402",
        strip_prefix = "rand_chacha-0.2.2",
        build_file = Label("//cargo/remote:BUILD.rand_chacha-0.2.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rand_chacha__0_3_1",
        url = "https://crates.io/api/v1/crates/rand_chacha/0.3.1/download",
        type = "tar.gz",
        sha256 = "e6c10a63a0fa32252be49d21e7709d4d4baf8d231c2dbce1eaa8141b9b127d88",
        strip_prefix = "rand_chacha-0.3.1",
        build_file = Label("//cargo/remote:BUILD.rand_chacha-0.3.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rand_core__0_5_1",
        url = "https://crates.io/api/v1/crates/rand_core/0.5.1/download",
        type = "tar.gz",
        sha256 = "90bde5296fc891b0cef12a6d03ddccc162ce7b2aff54160af9338f8d40df6d19",
        strip_prefix = "rand_core-0.5.1",
        build_file = Label("//cargo/remote:BUILD.rand_core-0.5.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rand_core__0_6_4",
        url = "https://crates.io/api/v1/crates/rand_core/0.6.4/download",
        type = "tar.gz",
        sha256 = "ec0be4795e2f6a28069bec0b5ff3e2ac9bafc99e6a9a7dc3547996c5c816922c",
        strip_prefix = "rand_core-0.6.4",
        build_file = Label("//cargo/remote:BUILD.rand_core-0.6.4.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rand_hc__0_2_0",
        url = "https://crates.io/api/v1/crates/rand_hc/0.2.0/download",
        type = "tar.gz",
        sha256 = "ca3129af7b92a17112d59ad498c6f81eaf463253766b90396d39ea7a39d6613c",
        strip_prefix = "rand_hc-0.2.0",
        build_file = Label("//cargo/remote:BUILD.rand_hc-0.2.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__reed_solomon_erasure__4_0_2",
        url = "https://crates.io/api/v1/crates/reed-solomon-erasure/4.0.2/download",
        type = "tar.gz",
        sha256 = "a415a013dd7c5d4221382329a5a3482566da675737494935cbbbcdec04662f9d",
        strip_prefix = "reed-solomon-erasure-4.0.2",
        build_file = Label("//cargo/remote:BUILD.reed-solomon-erasure-4.0.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__ripemd__0_1_3",
        url = "https://crates.io/api/v1/crates/ripemd/0.1.3/download",
        type = "tar.gz",
        sha256 = "bd124222d17ad93a644ed9d011a40f4fb64aa54275c08cc216524a9ea82fb09f",
        strip_prefix = "ripemd-0.1.3",
        build_file = Label("//cargo/remote:BUILD.ripemd-0.1.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rlp_derive__0_1_0",
        url = "https://crates.io/api/v1/crates/rlp-derive/0.1.0/download",
        type = "tar.gz",
        sha256 = "e33d7b2abe0c340d8797fe2907d3f20d3b5ea5908683618bfe80df7f621f672a",
        strip_prefix = "rlp-derive-0.1.0",
        build_file = Label("//cargo/remote:BUILD.rlp-derive-0.1.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rustc_hex__2_1_0",
        url = "https://crates.io/api/v1/crates/rustc-hex/2.1.0/download",
        type = "tar.gz",
        sha256 = "3e75f6a532d0fd9f7f13144f392b6ad56a32696bfcd9c78f797f16bbb6f072d6",
        strip_prefix = "rustc-hex-2.1.0",
        build_file = Label("//cargo/remote:BUILD.rustc-hex-2.1.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rustc_version__0_4_0",
        url = "https://crates.io/api/v1/crates/rustc_version/0.4.0/download",
        type = "tar.gz",
        sha256 = "bfa0f585226d2e68097d4f95d113b15b83a82e819ab25717ec0590d9584ef366",
        strip_prefix = "rustc_version-0.4.0",
        build_file = Label("//cargo/remote:BUILD.rustc_version-0.4.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__rustversion__1_0_9",
        url = "https://crates.io/api/v1/crates/rustversion/1.0.9/download",
        type = "tar.gz",
        sha256 = "97477e48b4cf8603ad5f7aaf897467cf42ab4218a38ef76fb14c2d6773a6d6a8",
        strip_prefix = "rustversion-1.0.9",
        build_file = Label("//cargo/remote:BUILD.rustversion-1.0.9.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__ryu__1_0_11",
        url = "https://crates.io/api/v1/crates/ryu/1.0.11/download",
        type = "tar.gz",
        sha256 = "4501abdff3ae82a1c1b477a17252eb69cee9e66eb915c1abaa4f44d873df9f09",
        strip_prefix = "ryu-1.0.11",
        build_file = Label("//cargo/remote:BUILD.ryu-1.0.11.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__schemars__0_8_11",
        url = "https://crates.io/api/v1/crates/schemars/0.8.11/download",
        type = "tar.gz",
        sha256 = "2a5fb6c61f29e723026dc8e923d94c694313212abbecbbe5f55a7748eec5b307",
        strip_prefix = "schemars-0.8.11",
        build_file = Label("//cargo/remote:BUILD.schemars-0.8.11.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__schemars_derive__0_8_11",
        url = "https://crates.io/api/v1/crates/schemars_derive/0.8.11/download",
        type = "tar.gz",
        sha256 = "f188d036977451159430f3b8dc82ec76364a42b7e289c2b18a9a18f4470058e9",
        strip_prefix = "schemars_derive-0.8.11",
        build_file = Label("//cargo/remote:BUILD.schemars_derive-0.8.11.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__scratch__1_0_2",
        url = "https://crates.io/api/v1/crates/scratch/1.0.2/download",
        type = "tar.gz",
        sha256 = "9c8132065adcfd6e02db789d9285a0deb2f3fcb04002865ab67d5fb103533898",
        strip_prefix = "scratch-1.0.2",
        build_file = Label("//cargo/remote:BUILD.scratch-1.0.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__semver__1_0_14",
        url = "https://crates.io/api/v1/crates/semver/1.0.14/download",
        type = "tar.gz",
        sha256 = "e25dfac463d778e353db5be2449d1cce89bd6fd23c9f1ea21310ce6e5a1b29c4",
        strip_prefix = "semver-1.0.14",
        build_file = Label("//cargo/remote:BUILD.semver-1.0.14.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__serde__1_0_147",
        url = "https://crates.io/api/v1/crates/serde/1.0.147/download",
        type = "tar.gz",
        sha256 = "d193d69bae983fc11a79df82342761dfbf28a99fc8d203dca4c3c1b590948965",
        strip_prefix = "serde-1.0.147",
        build_file = Label("//cargo/remote:BUILD.serde-1.0.147.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__serde_derive__1_0_147",
        url = "https://crates.io/api/v1/crates/serde_derive/1.0.147/download",
        type = "tar.gz",
        sha256 = "4f1d362ca8fc9c3e3a7484440752472d68a6caa98f1ab81d99b5dfe517cec852",
        strip_prefix = "serde_derive-1.0.147",
        build_file = Label("//cargo/remote:BUILD.serde_derive-1.0.147.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__serde_derive_internals__0_26_0",
        url = "https://crates.io/api/v1/crates/serde_derive_internals/0.26.0/download",
        type = "tar.gz",
        sha256 = "85bf8229e7920a9f636479437026331ce11aa132b4dde37d121944a44d6e5f3c",
        strip_prefix = "serde_derive_internals-0.26.0",
        build_file = Label("//cargo/remote:BUILD.serde_derive_internals-0.26.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__serde_json__1_0_87",
        url = "https://crates.io/api/v1/crates/serde_json/1.0.87/download",
        type = "tar.gz",
        sha256 = "6ce777b7b150d76b9cf60d28b55f5847135a003f7d7350c6be7a773508ce7d45",
        strip_prefix = "serde_json-1.0.87",
        build_file = Label("//cargo/remote:BUILD.serde_json-1.0.87.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__sha2__0_10_6",
        url = "https://crates.io/api/v1/crates/sha2/0.10.6/download",
        type = "tar.gz",
        sha256 = "82e6b795fe2e3b1e845bafcb27aa35405c4d47cdfc92af5fc8d3002f76cebdc0",
        strip_prefix = "sha2-0.10.6",
        build_file = Label("//cargo/remote:BUILD.sha2-0.10.6.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__sha2__0_9_9",
        url = "https://crates.io/api/v1/crates/sha2/0.9.9/download",
        type = "tar.gz",
        sha256 = "4d58a1e1bf39749807d89cf2d98ac2dfa0ff1cb3faa38fbb64dd88ac8013d800",
        strip_prefix = "sha2-0.9.9",
        build_file = Label("//cargo/remote:BUILD.sha2-0.9.9.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__sha3__0_10_6",
        url = "https://crates.io/api/v1/crates/sha3/0.10.6/download",
        type = "tar.gz",
        sha256 = "bdf0c33fae925bdc080598b84bc15c55e7b9a4a43b3c704da051f977469691c9",
        strip_prefix = "sha3-0.10.6",
        build_file = Label("//cargo/remote:BUILD.sha3-0.10.6.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__signature__1_6_4",
        url = "https://crates.io/api/v1/crates/signature/1.6.4/download",
        type = "tar.gz",
        sha256 = "74233d3b3b2f6d4b006dc19dee745e73e2a6bfb6f93607cd3b02bd5b00797d7c",
        strip_prefix = "signature-1.6.4",
        build_file = Label("//cargo/remote:BUILD.signature-1.6.4.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__smallvec__1_10_0",
        url = "https://crates.io/api/v1/crates/smallvec/1.10.0/download",
        type = "tar.gz",
        sha256 = "a507befe795404456341dfab10cef66ead4c041f62b8b11bbb92bffe5d0953e0",
        strip_prefix = "smallvec-1.10.0",
        build_file = Label("//cargo/remote:BUILD.smallvec-1.10.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__smart_default__0_6_0",
        url = "https://crates.io/api/v1/crates/smart-default/0.6.0/download",
        type = "tar.gz",
        sha256 = "133659a15339456eeeb07572eb02a91c91e9815e9cbc89566944d2c8d3efdbf6",
        strip_prefix = "smart-default-0.6.0",
        build_file = Label("//cargo/remote:BUILD.smart-default-0.6.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__spin__0_5_2",
        url = "https://crates.io/api/v1/crates/spin/0.5.2/download",
        type = "tar.gz",
        sha256 = "6e63cff320ae2c57904679ba7cb63280a3dc4613885beafb148ee7bf9aa9042d",
        strip_prefix = "spin-0.5.2",
        build_file = Label("//cargo/remote:BUILD.spin-0.5.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__static_assertions__1_1_0",
        url = "https://crates.io/api/v1/crates/static_assertions/1.1.0/download",
        type = "tar.gz",
        sha256 = "a2eb9349b6444b326872e140eb1cf5e7c522154d69e7a0ffb0fb81c06b37543f",
        strip_prefix = "static_assertions-1.1.0",
        build_file = Label("//cargo/remote:BUILD.static_assertions-1.1.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__strum__0_24_1",
        url = "https://crates.io/api/v1/crates/strum/0.24.1/download",
        type = "tar.gz",
        sha256 = "063e6045c0e62079840579a7e47a355ae92f60eb74daaf156fb1e84ba164e63f",
        strip_prefix = "strum-0.24.1",
        build_file = Label("//cargo/remote:BUILD.strum-0.24.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__strum_macros__0_24_3",
        url = "https://crates.io/api/v1/crates/strum_macros/0.24.3/download",
        type = "tar.gz",
        sha256 = "1e385be0d24f186b4ce2f9982191e7101bb737312ad61c1f2f984f34bcf85d59",
        strip_prefix = "strum_macros-0.24.3",
        build_file = Label("//cargo/remote:BUILD.strum_macros-0.24.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__subtle__2_4_1",
        url = "https://crates.io/api/v1/crates/subtle/2.4.1/download",
        type = "tar.gz",
        sha256 = "6bdef32e8150c2a081110b42772ffe7d7c9032b606bc226c8260fd97e0976601",
        strip_prefix = "subtle-2.4.1",
        build_file = Label("//cargo/remote:BUILD.subtle-2.4.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__syn__1_0_103",
        url = "https://crates.io/api/v1/crates/syn/1.0.103/download",
        type = "tar.gz",
        sha256 = "a864042229133ada95abf3b54fdc62ef5ccabe9515b64717bcb9a1919e59445d",
        strip_prefix = "syn-1.0.103",
        build_file = Label("//cargo/remote:BUILD.syn-1.0.103.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__synstructure__0_12_6",
        url = "https://crates.io/api/v1/crates/synstructure/0.12.6/download",
        type = "tar.gz",
        sha256 = "f36bdaa60a83aca3921b5259d5400cbf5e90fc51931376a9bd4a0eb79aa7210f",
        strip_prefix = "synstructure-0.12.6",
        build_file = Label("//cargo/remote:BUILD.synstructure-0.12.6.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__tap__1_0_1",
        url = "https://crates.io/api/v1/crates/tap/1.0.1/download",
        type = "tar.gz",
        sha256 = "55937e1799185b12863d447f42597ed69d9928686b8d88a1df17376a097d8369",
        strip_prefix = "tap-1.0.1",
        build_file = Label("//cargo/remote:BUILD.tap-1.0.1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__termcolor__1_1_3",
        url = "https://crates.io/api/v1/crates/termcolor/1.1.3/download",
        type = "tar.gz",
        sha256 = "bab24d30b911b2376f3a13cc2cd443142f0c81dda04c118693e35b3835757755",
        strip_prefix = "termcolor-1.1.3",
        build_file = Label("//cargo/remote:BUILD.termcolor-1.1.3.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__thiserror__1_0_37",
        url = "https://crates.io/api/v1/crates/thiserror/1.0.37/download",
        type = "tar.gz",
        sha256 = "10deb33631e3c9018b9baf9dcbbc4f737320d2b576bac10f6aefa048fa407e3e",
        strip_prefix = "thiserror-1.0.37",
        build_file = Label("//cargo/remote:BUILD.thiserror-1.0.37.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__thiserror_impl__1_0_37",
        url = "https://crates.io/api/v1/crates/thiserror-impl/1.0.37/download",
        type = "tar.gz",
        sha256 = "982d17546b47146b28f7c22e3d08465f6b8903d0ea13c1660d9d84a6e7adcdbb",
        strip_prefix = "thiserror-impl-1.0.37",
        build_file = Label("//cargo/remote:BUILD.thiserror-impl-1.0.37.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__time__0_1_44",
        url = "https://crates.io/api/v1/crates/time/0.1.44/download",
        type = "tar.gz",
        sha256 = "6db9e6914ab8b1ae1c260a4ae7a49b6c5611b40328a735b21862567685e73255",
        strip_prefix = "time-0.1.44",
        build_file = Label("//cargo/remote:BUILD.time-0.1.44.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__tiny_keccak__2_0_2",
        url = "https://crates.io/api/v1/crates/tiny-keccak/2.0.2/download",
        type = "tar.gz",
        sha256 = "2c9d3793400a45f954c52e73d068316d76b6f4e36977e3fcebb13a2721e80237",
        strip_prefix = "tiny-keccak-2.0.2",
        build_file = Label("//cargo/remote:BUILD.tiny-keccak-2.0.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__toml__0_5_9",
        url = "https://crates.io/api/v1/crates/toml/0.5.9/download",
        type = "tar.gz",
        sha256 = "8d82e1a7758622a465f8cee077614c73484dac5b836c02ff6a40d5d1010324d7",
        strip_prefix = "toml-0.5.9",
        build_file = Label("//cargo/remote:BUILD.toml-0.5.9.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__typenum__1_15_0",
        url = "https://crates.io/api/v1/crates/typenum/1.15.0/download",
        type = "tar.gz",
        sha256 = "dcf81ac59edc17cc8697ff311e8f5ef2d99fcbd9817b34cec66f90b6c3dfd987",
        strip_prefix = "typenum-1.15.0",
        build_file = Label("//cargo/remote:BUILD.typenum-1.15.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__uint__0_9_4",
        url = "https://crates.io/api/v1/crates/uint/0.9.4/download",
        type = "tar.gz",
        sha256 = "a45526d29728d135c2900b0d30573fe3ee79fceb12ef534c7bb30e810a91b601",
        strip_prefix = "uint-0.9.4",
        build_file = Label("//cargo/remote:BUILD.uint-0.9.4.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__unicode_ident__1_0_5",
        url = "https://crates.io/api/v1/crates/unicode-ident/1.0.5/download",
        type = "tar.gz",
        sha256 = "6ceab39d59e4c9499d4e5a8ee0e2735b891bb7308ac83dfb4e80cad195c9f6f3",
        strip_prefix = "unicode-ident-1.0.5",
        build_file = Label("//cargo/remote:BUILD.unicode-ident-1.0.5.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__unicode_width__0_1_10",
        url = "https://crates.io/api/v1/crates/unicode-width/0.1.10/download",
        type = "tar.gz",
        sha256 = "c0edd1e5b14653f783770bce4a4dabb4a5108a5370a5f5d8cfe8710c361f6c8b",
        strip_prefix = "unicode-width-0.1.10",
        build_file = Label("//cargo/remote:BUILD.unicode-width-0.1.10.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__unicode_xid__0_2_4",
        url = "https://crates.io/api/v1/crates/unicode-xid/0.2.4/download",
        type = "tar.gz",
        sha256 = "f962df74c8c05a667b5ee8bcf162993134c104e96440b663c8daa176dc772d8c",
        strip_prefix = "unicode-xid-0.2.4",
        build_file = Label("//cargo/remote:BUILD.unicode-xid-0.2.4.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__version_check__0_9_4",
        url = "https://crates.io/api/v1/crates/version_check/0.9.4/download",
        type = "tar.gz",
        sha256 = "49874b5167b65d7193b8aba1567f5c7d93d001cafc34600cee003eda787e483f",
        strip_prefix = "version_check-0.9.4",
        build_file = Label("//cargo/remote:BUILD.version_check-0.9.4.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wasi__0_10_0_wasi_snapshot_preview1",
        url = "https://crates.io/api/v1/crates/wasi/0.10.0+wasi-snapshot-preview1/download",
        type = "tar.gz",
        sha256 = "1a143597ca7c7793eff794def352d41792a93c481eb1042423ff7ff72ba2c31f",
        strip_prefix = "wasi-0.10.0+wasi-snapshot-preview1",
        build_file = Label("//cargo/remote:BUILD.wasi-0.10.0+wasi-snapshot-preview1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wasi__0_11_0_wasi_snapshot_preview1",
        url = "https://crates.io/api/v1/crates/wasi/0.11.0+wasi-snapshot-preview1/download",
        type = "tar.gz",
        sha256 = "9c8d87e72b64a3b4db28d11ce29237c246188f4f51057d65a7eab63b7987e423",
        strip_prefix = "wasi-0.11.0+wasi-snapshot-preview1",
        build_file = Label("//cargo/remote:BUILD.wasi-0.11.0+wasi-snapshot-preview1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wasi__0_9_0_wasi_snapshot_preview1",
        url = "https://crates.io/api/v1/crates/wasi/0.9.0+wasi-snapshot-preview1/download",
        type = "tar.gz",
        sha256 = "cccddf32554fecc6acb585f82a32a72e28b48f8c4c1883ddfeeeaa96f7d8e519",
        strip_prefix = "wasi-0.9.0+wasi-snapshot-preview1",
        build_file = Label("//cargo/remote:BUILD.wasi-0.9.0+wasi-snapshot-preview1.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wasm_bindgen__0_2_83",
        url = "https://crates.io/api/v1/crates/wasm-bindgen/0.2.83/download",
        type = "tar.gz",
        sha256 = "eaf9f5aceeec8be17c128b2e93e031fb8a4d469bb9c4ae2d7dc1888b26887268",
        strip_prefix = "wasm-bindgen-0.2.83",
        build_file = Label("//cargo/remote:BUILD.wasm-bindgen-0.2.83.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wasm_bindgen_backend__0_2_83",
        url = "https://crates.io/api/v1/crates/wasm-bindgen-backend/0.2.83/download",
        type = "tar.gz",
        sha256 = "4c8ffb332579b0557b52d268b91feab8df3615f265d5270fec2a8c95b17c1142",
        strip_prefix = "wasm-bindgen-backend-0.2.83",
        build_file = Label("//cargo/remote:BUILD.wasm-bindgen-backend-0.2.83.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wasm_bindgen_macro__0_2_83",
        url = "https://crates.io/api/v1/crates/wasm-bindgen-macro/0.2.83/download",
        type = "tar.gz",
        sha256 = "052be0f94026e6cbc75cdefc9bae13fd6052cdcaf532fa6c45e7ae33a1e6c810",
        strip_prefix = "wasm-bindgen-macro-0.2.83",
        build_file = Label("//cargo/remote:BUILD.wasm-bindgen-macro-0.2.83.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wasm_bindgen_macro_support__0_2_83",
        url = "https://crates.io/api/v1/crates/wasm-bindgen-macro-support/0.2.83/download",
        type = "tar.gz",
        sha256 = "07bc0c051dc5f23e307b13285f9d75df86bfdf816c5721e573dec1f9b8aa193c",
        strip_prefix = "wasm-bindgen-macro-support-0.2.83",
        build_file = Label("//cargo/remote:BUILD.wasm-bindgen-macro-support-0.2.83.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wasm_bindgen_shared__0_2_83",
        url = "https://crates.io/api/v1/crates/wasm-bindgen-shared/0.2.83/download",
        type = "tar.gz",
        sha256 = "1c38c045535d93ec4f0b4defec448e4291638ee608530863b1e2ba115d4fff7f",
        strip_prefix = "wasm-bindgen-shared-0.2.83",
        build_file = Label("//cargo/remote:BUILD.wasm-bindgen-shared-0.2.83.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wee_alloc__0_4_5",
        url = "https://crates.io/api/v1/crates/wee_alloc/0.4.5/download",
        type = "tar.gz",
        sha256 = "dbb3b5a6b2bb17cb6ad44a2e68a43e8d2722c997da10e928665c72ec6c0a0b8e",
        strip_prefix = "wee_alloc-0.4.5",
        build_file = Label("//cargo/remote:BUILD.wee_alloc-0.4.5.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__winapi__0_3_9",
        url = "https://crates.io/api/v1/crates/winapi/0.3.9/download",
        type = "tar.gz",
        sha256 = "5c839a674fcd7a98952e593242ea400abe93992746761e38641405d28b00f419",
        strip_prefix = "winapi-0.3.9",
        build_file = Label("//cargo/remote:BUILD.winapi-0.3.9.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__winapi_i686_pc_windows_gnu__0_4_0",
        url = "https://crates.io/api/v1/crates/winapi-i686-pc-windows-gnu/0.4.0/download",
        type = "tar.gz",
        sha256 = "ac3b87c63620426dd9b991e5ce0329eff545bccbbb34f3be09ff6fb6ab51b7b6",
        strip_prefix = "winapi-i686-pc-windows-gnu-0.4.0",
        build_file = Label("//cargo/remote:BUILD.winapi-i686-pc-windows-gnu-0.4.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__winapi_util__0_1_5",
        url = "https://crates.io/api/v1/crates/winapi-util/0.1.5/download",
        type = "tar.gz",
        sha256 = "70ec6ce85bb158151cae5e5c87f95a8e97d2c0c4b001223f33a334e3ce5de178",
        strip_prefix = "winapi-util-0.1.5",
        build_file = Label("//cargo/remote:BUILD.winapi-util-0.1.5.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__winapi_x86_64_pc_windows_gnu__0_4_0",
        url = "https://crates.io/api/v1/crates/winapi-x86_64-pc-windows-gnu/0.4.0/download",
        type = "tar.gz",
        sha256 = "712e227841d057c1ee1cd2fb22fa7e5a5461ae8e48fa2ca79ec42cfc1931183f",
        strip_prefix = "winapi-x86_64-pc-windows-gnu-0.4.0",
        build_file = Label("//cargo/remote:BUILD.winapi-x86_64-pc-windows-gnu-0.4.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__wyz__0_2_0",
        url = "https://crates.io/api/v1/crates/wyz/0.2.0/download",
        type = "tar.gz",
        sha256 = "85e60b0d1b5f99db2556934e21937020776a5d31520bf169e851ac44e6420214",
        strip_prefix = "wyz-0.2.0",
        build_file = Label("//cargo/remote:BUILD.wyz-0.2.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__zeroize__1_3_0",
        url = "https://crates.io/api/v1/crates/zeroize/1.3.0/download",
        type = "tar.gz",
        sha256 = "4756f7db3f7b5574938c3eb1c117038b8e07f95ee6718c0efad4ac21508f1efd",
        strip_prefix = "zeroize-1.3.0",
        build_file = Label("//cargo/remote:BUILD.zeroize-1.3.0.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__zeroize_derive__1_3_2",
        url = "https://crates.io/api/v1/crates/zeroize_derive/1.3.2/download",
        type = "tar.gz",
        sha256 = "3f8f187641dad4f680d25c4bfc4225b418165984179f26ca76ec4fb6441d3a17",
        strip_prefix = "zeroize_derive-1.3.2",
        build_file = Label("//cargo/remote:BUILD.zeroize_derive-1.3.2.bazel"),
    )

    maybe(
        http_archive,
        name = "raze__zeropool_bn__0_5_11",
        url = "https://crates.io/api/v1/crates/zeropool-bn/0.5.11/download",
        type = "tar.gz",
        sha256 = "71e61de68ede9ffdd69c01664f65a178c5188b73f78faa21f0936016a888ff7c",
        strip_prefix = "zeropool-bn-0.5.11",
        build_file = Label("//cargo/remote:BUILD.zeropool-bn-0.5.11.bazel"),
    )
