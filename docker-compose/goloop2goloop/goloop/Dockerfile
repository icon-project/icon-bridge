ARG ICONBRIDGE_VERSION=latest
ARG GOLOOP_IMAGE=goloop:latest

FROM iconbridge:${ICONBRIDGE_VERSION} AS iconbridge
FROM ${GOLOOP_IMAGE}

RUN apk add --no-cache jq curl

ENV GOLOOP_PROVISION=/goloop/provisioning
ENV GOLOOP_PROVISION_CONFIG=${GOLOOP_PROVISION}/config
ENV GOLOOP_PROVISION_DATA=${GOLOOP_PROVISION}/data
# copy files for provisioning
COPY --from=iconbridge /iconbridge/contracts/pyscore ${GOLOOP_PROVISION_CONFIG}/pyscore/
COPY ./*.sh /goloop/bin/
COPY ./entrypoint /

RUN provision.sh
WORKDIR /goloop/config
