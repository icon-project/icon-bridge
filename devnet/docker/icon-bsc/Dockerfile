ARG ICONBRIDGE_VERSION=latest
ARG GOLOOP_IMAGE=iconloop/goloop-icon:v0.9.7

FROM ${GOLOOP_IMAGE} as goloop
FROM canhlinh/truffle:5.3.0-alpine as truffle
FROM iconbridge:${ICONBRIDGE_VERSION} AS iconbridge
FROM node:16.1-alpine as node

# install truffle cli
COPY --from=truffle /usr/local/node /usr/local/node
RUN apk add --update nodejs-current npm
RUN apk add --no-cache jq curl inotify-tools git
ENV PATH $PATH:/usr/local/node/bin
RUN truffle version
#install yarn
RUN yarn --version

#install eth-cli 
RUN yarn global add eth-cli

# install iconbridge cli
COPY --from=iconbridge . .
ENV PATH $PATH:/iconbridge/bin

# install goloop cli
COPY --from=goloop /goloop/bin /goloop/bin
ENV PATH $PATH:/goloop/bin

ENV ICONBRIDGE_BASE_DIR=/iconbridge
ENV ICONBRIDGE_CONFIG_DIR=${ICONBRIDGE_BASE_DIR}/config
ENV ICONBRIDGE_CONTRACTS_DIR=${ICONBRIDGE_BASE_DIR}/contracts
ENV ICONBRIDGE_SCRIPTS_DIR=${ICONBRIDGE_BASE_DIR}/js
ENV ICONBRIDGE_BIN_DIR=/iconbridge/bin

COPY ./*.sh /iconbridge/bin/
COPY ./scripts/*.sh /iconbridge/bin/
COPY ./entrypoint.sh /
COPY ./scripts/provision.sh /
COPY ./iconvalidators /iconbridge/bin/


#RUN mkdir {ICONBRIDGE_SCRIPTS_DIR}
#COPY ./*.js /iconbridge/js/
COPY ./scripts/*.js /iconbridge/js/

RUN chmod +x /provision.sh
RUN chmod +x /entrypoint.sh
RUN chmod +x /iconbridge/bin/wait-for-it.sh
RUN chmod +x /iconbridge/bin/iconvalidators

RUN yarn --cwd $ICONBRIDGE_CONTRACTS_DIR/solidity/bmc --prod
RUN yarn --cwd $ICONBRIDGE_CONTRACTS_DIR/solidity/bts --prod

WORKDIR /iconbridge/config
ENTRYPOINT ["entrypoint.sh"]