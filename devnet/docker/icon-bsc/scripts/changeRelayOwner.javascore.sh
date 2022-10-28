#! /bin/bash
set -e 

source token.javascore.sh 


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
