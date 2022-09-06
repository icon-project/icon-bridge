
FROM canhlinh/truffle:5.3.0-alpine as truffle
FROM node:18-alpine as node


# go
COPY --from=golang:1.16-alpine /usr/local/go/ /usr/local/go/ 
ENV PATH="/usr/local/go/bin:${PATH}"
ENV PATH="/root/go/bin:${PATH}"

# install pkgs
COPY --from=truffle /usr/local/node /usr/local/node
RUN apk update && apk upgrade
RUN apk add --update nodejs-current npm
RUN apk add --no-cache jq curl inotify-tools git bash gcc libc-dev openssl gradle make binutils-gold python3 py3-pip g++
ENV PATH $PATH:/usr/local/node/bin

#install goloop, ethkey
RUN go install github.com/icon-project/goloop/cmd/goloop@latest
RUN go install github.com/ethereum/go-ethereum/cmd/ethkey@latest

#env
ENV ICONBRIDGE_BASE_DIR=/iconbridge
ENV ICONBRIDGE_CONFIG_DIR=${ICONBRIDGE_BASE_DIR}/config
ENV ICONBRIDGE_CONTRACTS_DIR=${ICONBRIDGE_BASE_DIR}/contracts
ENV ICONBRIDGE_SCRIPTS_DIR=${ICONBRIDGE_BASE_DIR}/js
ENV ICONBRIDGE_BIN_DIR=/iconbridge/bin
RUN mkdir -p ${ICONBRIDGE_CONFIG_DIR} && mkdir -p ${ICONBRIDGE_CONTRACTS_DIR} && mkdir -p ${ICONBRIDGE_SCRIPTS_DIR} && mkdir -p ${ICONBRIDGE_BIN_DIR}

#iconbridge
#RUN git clone https://github.com/icon-project/icon-bridge.git
RUN mkdir -p icon-bridge
COPY . icon-bridge/
RUN cd /icon-bridge/cmd/iconbridge && go build -o iconbridge && cp iconbridge /usr/local/go/bin/
RUN cd /icon-bridge/cmd/iconvalidators && go build -o iconvalidators && cp iconvalidators /usr/local/go/bin/ && cp iconvalidators /iconbridge/bin/

#contracts
RUN mkdir -p ${ICONBRIDGE_CONTRACTS_DIR}/javascore
RUN mkdir -p ${ICONBRIDGE_CONTRACTS_DIR}/solidity

RUN cd /icon-bridge/javascore && gradle clean 
RUN cd /icon-bridge/javascore && gradle bmc:optimizedJar && cp /icon-bridge/javascore/bmc/build/libs/bmc-optimized.jar ${ICONBRIDGE_CONTRACTS_DIR}/javascore/bmc.jar 
RUN cd /icon-bridge/javascore && gradle bts:optimizedJar && cp /icon-bridge/javascore/bts/build/libs/bts-optimized.jar ${ICONBRIDGE_CONTRACTS_DIR}/javascore/bts.jar
RUN cd /icon-bridge/javascore && cp /icon-bridge/javascore/lib/irc2Tradeable-0.1.0-optimized.jar ${ICONBRIDGE_CONTRACTS_DIR}/javascore/irc2Tradeable.jar
RUN cd /icon-bridge && git clone https://github.com/icon-project/java-score-examples.git
RUN cd /icon-bridge/java-score-examples && gradle irc2-token:clean && gradle irc2-token:optimizedJar && cp /icon-bridge/java-score-examples/irc2-token/build/libs/irc2-token-0.9.1-optimized.jar ${ICONBRIDGE_CONTRACTS_DIR}/javascore/irc2.jar

RUN cp -r /icon-bridge/solidity/bmc ${ICONBRIDGE_CONTRACTS_DIR}/solidity/
RUN cp -r /icon-bridge/solidity/bts ${ICONBRIDGE_CONTRACTS_DIR}/solidity/
RUN yarn --cwd $ICONBRIDGE_CONTRACTS_DIR/solidity/bmc --prod
RUN yarn --cwd $ICONBRIDGE_CONTRACTS_DIR/solidity/bts --prod

#scripts
RUN cd /icon-bridge/devnet/docker/icon-bsc && \
     cp ./*.sh /iconbridge/bin/ && \
     cp ./scripts/*.sh /iconbridge/bin/ && \
     cp ./entrypoint.sh / && \
     cp ./scripts/provision.sh / && \
     cp ./scripts/*.js /iconbridge/js/
#replace config.sh with env.variables.sh, so that upgrade scripts can use env.variables.sh without change
RUN cd /icon-bridge/devnet/docker/icon-bsc/scripts && \
     cp config.sh config.sh_testnet && \
     cp ../env.variables.sh config.sh
RUN chmod +x /provision.sh
RUN chmod +x /entrypoint.sh
RUN chmod +x /iconbridge/bin/wait-for-it.sh

WORKDIR /iconbridge/bin
ENTRYPOINT ["/bin/bash","/iconbridge/bin/entrypoint.sh"]