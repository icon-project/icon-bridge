#!/bin/bash
set -e

########################################################

ROOT_DIR=${PWD}
CONTRACTS_DIR="$ROOT_DIR/build/contracts"
JAVASCORE_DIR="$CONTRACTS_DIR/javascore"
SOLIDITY_DIR="$CONTRACTS_DIR/solidity"

mkdir -p "$JAVASCORE_DIR"
mkdir -p "$SOLIDITY_DIR"
mkdir -p "$ROOT_DIR/keys/_ixh/keystore"

mkdir -p bin
cd cmd/iconvalidators
go build .
cd $ROOT_DIR

echo "Build contracts"

cd $ROOT_DIR/javascore
gradle clean
gradle bmc:optimizedJar
gradle bts:optimizedJar
gradle irc2-token:optimizedJar
gradle irc2Tradeable:optimizedJar
cp bmc/build/libs/bmc-optimized.jar $JAVASCORE_DIR/bmc.jar
cp bts/build/libs/bts-optimized.jar $JAVASCORE_DIR/bts.jar
cp irc2Tradeable/build/libs/irc2Tradeable-optimized.jar $JAVASCORE_DIR/irc2Tradeable.jar
cp irc2-token/build/libs/irc2-token-optimized.jar $JAVASCORE_DIR/irc2.jar

echo "Copy solidity"

rm -rf $ROOT_DIR/solidity/bmc/build
rm -rf $ROOT_DIR/solidity/bmc/node_modules
cp -r $ROOT_DIR/solidity/bmc $SOLIDITY_DIR/

rm -rf $ROOT_DIR/solidity/bts/build
rm -rf $ROOT_DIR/solidity/bts/node_modules
cp -r $ROOT_DIR/solidity/bts $SOLIDITY_DIR/

cd $SOLIDITY_DIR
cd bmc && yarn install
cd ..
cd bts && yarn install