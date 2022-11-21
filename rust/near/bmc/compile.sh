MANIFEST_PATH=${PWD}/bmc/Cargo.toml
RUSTFLAGS='-C link-arg=-s' cargo build --manifest-path ${MANIFEST_PATH} --target wasm32-unknown-unknown --release --out-dir ./res -Z unstable-options