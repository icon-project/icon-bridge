cp ./res/nep141_service.wasm ${PWD}/bsh/bts/res/NEP141_CONTRACT.wasm
MANIFEST_PATH=${PWD}/bsh/bts/Cargo.toml
RUSTFLAGS='-C link-arg=-s' cargo build --manifest-path ${MANIFEST_PATH} --target wasm32-unknown-unknown --release --out-dir ./res -Z unstable-options