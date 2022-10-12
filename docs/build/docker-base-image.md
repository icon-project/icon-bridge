# Docker Base Image
## Introduction
For making the CI process faster and to include all dependencies of the mono repo inside single image, we have created a base image [iconbridge/build](https://hub.docker.com/r/iconbridge/build).

## Base Image Composition
The base image is built from Dockerfile on [build/base](https://github.com/icon-project/icon-bridge/tree/main/build/base) folder which comprises following major packages:
1. Java/Gradle
2. Nodejs
2. Go
4. Truffle(nodejs global package)
5. Goloop
6. Go-ethereum

## Image Build Trigger
Github Actions workflow file [build-base-image.yml](https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/build-base-image.yaml) handles the build and push to docker hub.
The workflow is triggered when there is commit to main branch and there is change on [build/base](https://github.com/icon-project/icon-bridge/tree/main/build/base) folder.

## Updating Package
If any package needs to be added/deleted/upgraded, change the Dockerfile and merge to main branch. The change will be taken by next builds
