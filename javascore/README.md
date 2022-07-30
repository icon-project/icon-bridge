# ICON BTP BSC

ICON Blockchain Transmission Protocol for Binance Smart Chain (WIP)

## Requirements

- [ICON javaee](https://github.com/icon-project/goloop/tree/master/javaee)
- [ICON sdk](https://github.com/icon-project/goloop/tree/master/sdk/java) 
- [ICON goloop](https://github.com/icon-project/goloop) (for running the integration tests).
#### Obtain the goloop repo
```
$ git clone git@github.com:icon-project/goloop.git
$ GOLOOP_ROOT=/path/to/goloop
```

## Build
##### Dependencies
This project currently depends on building local maven snapshot for the ICON Java SCORE. 
###### Build & publish javaee
1. goto the api folder
    ``` cd ${GOLOOP_ROOT}/javaee/ ```
2. run 
    ``` ./gradlew publishToMavenLocal -x:signMavenJavaPublication ```

###### Build & publish sdk
1. goto the icon-sdk folder
    ``` cd ${GOLOOP_ROOT}/sdk/java ```
2. run 
    ``` ./gradlew publishToMavenLocal -x:signMavenJavaPublication ```

##### Build bmc, bsr and bts
``` ./gradlew bmc:build ```
``` ./gradlew bsr:build ```
``` ./gradlew bts:build ```

## Run Integration Tests with deployment using integration test cases
Follow local gochain setup guide:
[gochain_icon_local_node_guide](https://github.com/icon-project/goloop/blob/master/doc/gochain_icon_local_node_guide.md)

From the integration-tests project, run the following:

``` ./gradlew testJavaScore -DNO_SERVER=true -DCHAIN_ENV=./data/env.properties ```

For a specific test, use --tests <testname>

``` ./gradlew <project_name>:testJavaScore -DNO_SERVER=true --tests MTATest -DCHAIN_ENV=./data/env.properties ```


### Deployment in a local node using scripts & integration test for local node

Run integration tests from local deployment:

steps:

1. clean & create optimizeJar of BMC, BSH, BMV

``` gradle <project_name>:clean ```

``` gradle <project_name>:optimizedJar```

2. run ```deploy-script.sh deployToLocal```
4. substitute the BMC score address, BMV score address & BMV score deploy Txn address from the output of above script at setup() method in BMVLocalTest.java file
5. Also pass appropriate keystore file & password in setup() method in BMVLocalTest.java file
3. run BMV local Test from integration test
   ``` gradle bsh:testJavaScore -DNO_SERVER=true --tests BMVLocalTest -DCHAIN_ENV=./data/env.properties -PkeystoreName=keystore -PkeystorePass=Admin@123```
   
Other commands:


BMC:

``` gradle bmc:deployToLocal -PkeystoreName=../keys/keystore_god.json -PkeystorePass=gochain ```

Deploying BMC.jar (replace with proper parameters):

``` goloop rpc --uri http://btp.net.solidwallet.io/api/v3 sendtx deploy bmc-0.1.0.jar \
    --key_store keystore --key_password Admin@123 \
    --nid 0x42 --step_limit 13610920001 \
    --content_type application/zip \
    --param _net="0x07.icon" ```
BTS:

``` gradle bts:deployToLocal -DBMC_ADDRESS=<BMC_SCORE_ADDRESS> -PkeystoreName=../keys/keystore_god.json -PkeystorePass=gochain ```

``` goloop rpc --uri http://btp.net.solidwallet.io/api/v3 sendtx deploy ./bts/build/libs/bts-0.1.0.jar \
    --key_store ./keys/keystore.json --key_password gochain \
    --nid 0x42 --step_limit 13610920001 \
    --content_type application/zip \
    --param _bmc=cx14579031817b2973f50b78bc1507e9c2d446e0f7 \
    --param _net="0x07.icon" 
    ```



//@Deprecated (this module was moved in to the bts)
BSH:
``` gradle token-bsh:deployToLocal -DBMC_ADDRESS=<BMC_SCORE_ADDRESS> -PkeystoreName=../keys/keystore_god.json -PkeystorePass=gochain ```

