#!/bin/bash
set -e 
# docker-compose rm
# docker image rm icon-bsc_goloop
# rm -rf work/*
# #docker-compose -f docker-compose.yml -f docker-compose.provision.yml up -d   --force-recreate && docker-compose stop
# #docker inspect iconbridge_src -f '{{ json .State.Health.Log }}' | jq .
# docker-compose build
export DOCKER_DEFAULT_PLATFORM=linux/amd64
echo "Build BMR"
cd ../../../../icon-bridge/
docker build -f ./devnet/docker/icon-bsc/Dockerfile -t iconbridge_bsc:latest .
echo "Build BSC"
docker build --tag bsc-node ./devnet/docker/bsc-node --build-arg KEYSTORE_PASS="Perlia0"
echo "Build ICON"
cd ./devnet/docker/goloop
docker build --tag icon-bsc_goloop:latest .