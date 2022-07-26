#!/bin/bash
bmr_dir=$PWD
bmr_src_dir="$bmr_dir"
root_dir="$bmr_dir/../../../.."
env
DOCKER_REPO=${DOCKER_REPO:-localhost:5000}
DOCKER_IMG=${DOCKER_IMG:-bmr-hmny}
DOCKER_TAG=${DOCKER_TAG:-latest}
DOCKER_IMGTAG=${DOCKER_IMGTAG:-"$DOCKER_REPO/$DOCKER_IMG:$DOCKER_TAG"}

cd $root_dir
docker build -f $bmr_src_dir/Dockerfile -t $DOCKER_IMGTAG .
