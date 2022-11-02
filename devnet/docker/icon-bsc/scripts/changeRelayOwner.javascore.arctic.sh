#! /bin/bash
set -e 

source config/icon-ice-arctic-config.sh
source token.javascore.sh 

if [ "$(basename $CONFIG_DIR)" != "_ixh_arctic" ]; then 
    echo "The scipt is currently meant for Arctic chain only. Got config dir ${CONFIG_DIR}"
    exit 0
fi

if [ $# -eq 0 ]; then
    echo "No arguments supplied: Pass --remove addr or --add"
elif [ $1 == "--remove" ]; then
    echo "Removing relay " $2
    configure_bmc_javascore_removeRelay $2
elif [ $1 == "--add" ]; then
    echo "Adding relay"
    configure_bmc_javascore_addRelay
else
    echo "Invalid argument: Pass --remove addr or --add"
fi
