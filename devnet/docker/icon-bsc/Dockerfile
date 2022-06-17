ARG BTPSIMPLE_VERSION=latest
ARG GOLOOP_IMAGE=iconloop/goloop-icon:v0.9.7

FROM ${GOLOOP_IMAGE}
FROM canhlinh/truffle:5.3.0-alpine as truffle
FROM btpsimple:${BTPSIMPLE_VERSION} AS btpsimple
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

# install btpsimple cli
COPY --from=btpsimple . .
ENV PATH $PATH:/btpsimple/bin

# install goloop cli
COPY --from=goloop /goloop/bin /goloop/bin
ENV PATH $PATH:/goloop/bin

ENV BTPSIMPLE_BASE_DIR=/btpsimple
ENV BTPSIMPLE_CONFIG_DIR=${BTPSIMPLE_BASE_DIR}/config
ENV BTPSIMPLE_CONTRACTS_DIR=${BTPSIMPLE_BASE_DIR}/contracts
ENV BTPSIMPLE_SCRIPTS_DIR=${BTPSIMPLE_BASE_DIR}/js
ENV BTPSIMPLE_BIN_DIR=/btpsimple/bin

COPY ./*.sh /btpsimple/bin/
COPY ./scripts/*.sh /btpsimple/bin/
COPY ./entrypoint.sh /
COPY ./scripts/provision.sh /
COPY ./iconvalidators /btpsimple/bin/


#RUN mkdir {BTPSIMPLE_SCRIPTS_DIR}
#COPY ./*.js /btpsimple/js/
COPY ./scripts/*.js /btpsimple/js/

RUN chmod +x /provision.sh
RUN chmod +x /entrypoint.sh
RUN chmod +x /btpsimple/bin/wait-for-it.sh

RUN yarn --cwd $BTPSIMPLE_CONTRACTS_DIR/solidity/bmc --prod
RUN yarn --cwd $BTPSIMPLE_CONTRACTS_DIR/solidity/bsh --prod
RUN yarn --cwd $BTPSIMPLE_CONTRACTS_DIR/solidity/TokenBSH --prod

WORKDIR /btpsimple/config
ENTRYPOINT ["entrypoint.sh"]