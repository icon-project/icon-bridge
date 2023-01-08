#-------------------------------------------------------------------------------
#
# 	Makefile for building target binaries.
#


# Configuration
ROOT_DIR=${PWD}
CONTRACTS_DIR=$(ROOT_DIR)/build/contracts
JAVASCORE_DIR=$(CONTRACTS_DIR)/javascore
SOLIDITY_DIR=$(CONTRACTS_DIR)/solidity

icon-clean:
	rm -rf $(ROOT_DIR)/build ; \
	cd ./javascore ; \
	gradle clean ; \



# TODO split ./build.sh into smaller processes and move to the make file
build-contracts:
	cd ./javascore ; \
    gradle clean ; \
    gradle bmc:optimizedJar ; \
    gradle bts:optimizedJar ; \
    gradle irc2-token:optimizedJar ; \
    gradle irc2Tradeable:optimizedJar ; \
    cp bmc/build/libs/bmc-optimized.jar $(JAVASCORE_DIR)/bmc.jar ; \
    cp bts/build/libs/bts-optimized.jar $(JAVASCORE_DIR)/bts.jar ; \
    cp irc2Tradeable/build/libs/irc2Tradeable-optimized.jar $(JAVASCORE_DIR)/irc2Tradeable.jar ; \
    cp irc2-token/build/libs/irc2-token-optimized.jar $(JAVASCORE_DIR)/irc2.jar ; \
