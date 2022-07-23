#!/bin/sh
set -e

source provision.sh

if [ "$ICONBRIDGE_OFFSET" != "" ] && [ -f "$ICONBRIDGE_OFFSET" ]; then
    export ICONBRIDGE_OFFSET=$(cat ${ICONBRIDGE_OFFSET})
fi

if [ "$ICONBRIDGE_CONFIG" != "" ] && [ ! -f "$ICONBRIDGE_CONFIG" ]; then
    echo "Setup configuration $ICONBRIDGE_CONFIG"
    UNSET="ICONBRIDGE_CONFIG"
    CMD="iconbridge save $ICONBRIDGE_CONFIG"
    if [ "$ICONBRIDGE_KEY_SECRET" != "" ] && [ ! -f "$ICONBRIDGE_KEY_SECRET" ]; then
        mkdir -p $(dirname $ICONBRIDGE_KEY_SECRET)
        printf $(date | md5sum | head -c16) >$ICONBRIDGE_KEY_SECRET
    fi
    if [ "$ICONBRIDGE_KEY_STORE" != "" ] && [ ! -f "$ICONBRIDGE_KEY_STORE" ]; then
        echo "Save keystore $ICONBRIDGE_CONFIG"
        UNSET="$UNSET ICONBRIDGE_KEY_STORE"
        CMD="$CMD --save_key_store=$ICONBRIDGE_KEY_STORE"
    fi
    if [ "$ICONBRIDGE_OFFSET" != "" ] && [ -f "$ICONBRIDGE_OFFSET" ]; then
        export ICONBRIDGE_OFFSET=$(cat ${ICONBRIDGE_OFFSET})
    fi

    if [ "$ICONBRIDGE_SRC_ADDRESS" != "" ] && [ -f "$ICONBRIDGE_SRC_ADDRESS" ]; then
        export ICONBRIDGE_SRC_ADDRESS=$(cat ${ICONBRIDGE_SRC_ADDRESS})
    fi
    if [ "$ICONBRIDGE_SRC_ENDPOINT" != "" ] && [ -f "$ICONBRIDGE_SRC_ENDPOINT" ]; then
        export ICONBRIDGE_SRC_ENDPOINT=$(cat ${ICONBRIDGE_SRC_ENDPOINT})
    fi
    if [ "$ICONBRIDGE_DST_ADDRESS" != "" ] && [ -f "$ICONBRIDGE_DST_ADDRESS" ]; then
        export ICONBRIDGE_DST_ADDRESS=$(cat ${ICONBRIDGE_DST_ADDRESS})
    fi
    if [ "$ICONBRIDGE_DST_ENDPOINT" != "" ] && [ -f "$ICONBRIDGE_DST_ENDPOINT" ]; then
        export ICONBRIDGE_DST_ENDPOINT=$(cat ${ICONBRIDGE_DST_ENDPOINT})
    fi
    sh -c "unset $UNSET ; $CMD"
fi

exec "$@"
