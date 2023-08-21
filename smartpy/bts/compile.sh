
#!/usr/bin/env bash

set -e -o pipefail

echo "----------------------------------------"
echo "Compiling contracts ... "
echo "----------------------------------------"

# Expected location of SmartPy CLI.
SMART_PY_CLI=~/smartpy-cli/SmartPy.sh

# Build artifact directory.
OUT_DIR=./contracts/build/.contract_build

# Array of SmartPy files to compile.
# CONTRACTS_ARRAY=(counter)

# Exit if SmartPy is not installed. 
if [ ! -f "$SMART_PY_CLI" ]; then
    echo "Fatal: Please install SmartPy CLI at $SMART_PY_CLI" && exit
fi

function processContract {
    CONTRACT_NAME=$1
    OUT_DIR=$2
    CONTRACT_IN="./contracts/src/${CONTRACT_NAME}.py"
    CONTRACT_OUT="${CONTRACT_NAME}.json"
    STORAGE_OUT="${CONTRACT_NAME}_storage.json"
    CONTRACT_COMPILED="${CONTRACT_NAME}/step_000_cont_0_contract.json"
    STORAGE_COMPILED="${CONTRACT_NAME}/step_000_cont_0_storage.json"

    echo ">> Processing ${CONTRACT_NAME}"

    # Ensure file exists.
    if [ ! -f "$CONTRACT_IN" ]; then
        echo "Fatal: $CONTRACT_IN not found. Running from wrong dir?" && exit
    fi

#    echo ">>> [1 / 3] Testing ${CONTRACT_NAME} ... "
#    $SMART_PY_CLI test $CONTRACT_IN $OUT_DIR --html

    echo ">>> [1 / 2] Compiling ${CONTRACT_NAME} ..."
    $SMART_PY_CLI compile $CONTRACT_IN $OUT_DIR --html

    echo ">>> [2 / 2] Extracting Michelson contract ... "
    cp $OUT_DIR/$CONTRACT_COMPILED ./contracts/build/$CONTRACT_OUT
    cp $OUT_DIR/$STORAGE_COMPILED ./contracts/build/$STORAGE_OUT

    echo ">>> Michelson contract written to ${CONTRACT_OUT}"
}

export PYTHONPATH=$PWD


echo "> [1 / 2] Compiling Contracts."
# Use if you want to pass a contract or more as arguments. 
for n in $(seq 1 $#); do
  processContract $1 $OUT_DIR
  shift
done

# Use if you want to compile all contracts in CONTRACTS_ARRAY. No arguments needed. 
# for i in ${!CONTRACTS_ARRAY[@]}; do
#     processContract ${CONTRACTS_ARRAY[$i]} $OUT_DIR
#  done

# Remove build artifacts.
echo "> [2 / 2] Cleaning up ..."
rm -rf $OUT_DIR
rm -rf ./contracts/__pycache__
rm -rf ./__pycache__


echo "> Removed artifacts."

echo "> Compilation successful."