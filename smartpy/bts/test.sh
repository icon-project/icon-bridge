#!/usr/bin/env bash

set -e -o pipefail

echo "----------------------------------------"
echo "Compiling contracts ... "
echo "----------------------------------------"

# Expected location of SmartPy CLI.
SMART_PY_CLI=~/smartpy-cli/SmartPy.sh

# Build artifact directory.
TEST_OUT_DIR=./contracts/build/.contract_build/test

# Array of SmartPy files to compile.
# CONTRACTS_ARRAY=(counter)

# Exit if SmartPy is not installed.
if [ ! -f "$SMART_PY_CLI" ]; then
    echo "Fatal: Please install SmartPy CLI at $SMART_PY_CLI" && exit
fi

function processContract {
    CONTRACT_NAME=$1
    TEST_OUT_DIR=$2
    CONTRACT_IN_TEST="./contracts/tests/${CONTRACT_NAME}_test.py"
    CONTRACT_OUT="${CONTRACT_NAME}.json"
    STORAGE_OUT="${CONTRACT_NAME}_storage.json"
    CONTRACT_COMPILED="${CONTRACT_NAME}/step_000_cont_0_contract.json"
    STORAGE_COMPILED="${CONTRACT_NAME}/step_000_cont_0_storage.json"

    echo ">> Processing ${CONTRACT_NAME}"

    # Ensure file exists.
    if [ ! -f "$CONTRACT_IN_TEST" ]; then
        echo "Fatal: $CONTRACT_IN_TEST not found. Running from wrong dir?" && exit
    fi

    echo ">>> Commencing the tests ${CONTRACT_NAME} ... "
    $SMART_PY_CLI test $CONTRACT_IN_TEST $TEST_OUT_DIR --html
}

export PYTHONPATH=$PWD


# Use if you want to pass a contract or more as arguments.
for n in $(seq 1 $#); do
  processContract $1 $TEST_OUT_DIR
  shift
done


echo "> Test Complete."