# Release Assets
## Introduction
The collection of assets to deploy to multiple environments.

## Assets Creation
The assets are auto created after the release creation is triggered. The assets are sent to Github Assets section of the specific release by (release.yml workflow file)[https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/release.yml]. It also sends (BMR docker image)[https://hub.docker.com/repository/docker/iconbridge/bmr] to the docker registry.
 
## Assets Location
The assets are list down under "Assets" section of the specific release. Following are the assets created:
1. bmc-javascore-optimized.zip
2. bmc-solidity-contracts.zip
3. bmr-go-iconbridge.zip
4. bts-javascore-optimized.zip
5. bts-solidity-contracts.zip
6. BMR docker image

