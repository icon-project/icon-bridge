FROM node:15.12.0-alpine AS node

FROM alpine:3.16

# Dependencies version pin
ENV GRADLE_VERSION 6.7.1
ENV GOLANG_VERSION 1.18
ENV GOLOOP_VERSION v1.2.11
ENV ETHKEY_VERSION v1.10.21
ENV TRUFFLE_VERSION 5.5.5

RUN apk add git bash openjdk11 jq aws-cli

RUN apk -U add --no-cache curl; \
    curl https://downloads.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip > gradle.zip; \
    unzip gradle.zip; \
    rm gradle.zip; \
    apk del curl;
ENV PATH "$PATH:/gradle-${GRADLE_VERSION}/bin/"

RUN echo "installing go" && \
        wget -O go.tgz https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
        tar -C /usr/local -xzf go.tgz && rm go.tgz && apk add go
RUN go version
ENV GOPATH="/go"
ENV PATH="${PATH}:/go/bin"
RUN go install github.com/icon-project/goloop/cmd/goloop@${GOLOOP_VERSION} && \
    go install github.com/ethereum/go-ethereum/cmd/ethkey@${ETHKEY_VERSION}
# For Github Actions user permission
RUN chmod -R 777 /go

# Install Nodejs
COPY --from=node /usr/lib /usr/lib
COPY --from=node /usr/local/share /usr/local/share
COPY --from=node /usr/local/lib /usr/local/lib
COPY --from=node /usr/local/include /usr/local/include
COPY --from=node /usr/local/bin /usr/local/bin

RUN npm install -g truffle@${TRUFFLE_VERSION}
RUN rm /usr/local/bin/yarn && rm /usr/local/bin/yarnpkg && npm install -g yarn

ENV LD_LIBRARY_PATH /usr/local/lib
RUN apk add git make g++ gmp-dev libressl-dev && \
 git clone https://github.com/harmony-one/bls.git && git clone https://github.com/harmony-one/mcl.git && \
 cd bls && make -j8 BLS_SWAP_G=1 && make install && make BLS_ETH=1 /usr/local/lib/libbls384_256.so && \
 cd ../mcl && make install && \
 cd .. && rm -rf mcl bls

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
