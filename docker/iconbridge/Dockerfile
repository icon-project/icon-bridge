ARG BASE_IMAGE
FROM ${BASE_IMAGE}
LABEL MAINTAINER="t_arch@iconloop.com"

# install dependencies
#RUN apk add --no-cache jq

ARG VERSION
LABEL ICONBRIDGE_VERSION="$VERSION"

# install
COPY dist/bin/iconbridge /iconbridge/bin/iconbridge
ENV PATH $PATH:/iconbridge/bin

# copy extras
#COPY dist/contracts/pyscore/*.zip /iconbridge/contracts/pyscore/
COPY dist/contracts/solidity /iconbridge/contracts/solidity
COPY dist/contracts/javascore /iconbridge/contracts/javascore
WORKDIR /iconbridge

# container configuration
VOLUME ["/iconbridge/data"]

# goloop entrypoint
ENV ICONBRIDGE_BASE_DIR=/iconbridge/data
ENV ICONBRIDGE_CONFIG=/iconbridge/config/config.json
ENV ICONBRIDGE_KEY_STORE=/iconbridge/config/keystore.json
ENV ICONBRIDGE_KEY_SECRET=/iconbridge/config/keysecret
ENV ICONBRIDGE_LOG_WRITER_FILENAME=/iconbridge/data/iconbridge.log

COPY ./entrypoint /iconbridge

RUN chmod +x entrypoint
ENTRYPOINT ["/iconbridge/entrypoint"]
