FROM iconloop/goloop-icon:latest

ARG CONFIG_JSON=""
#RUN [[ ! -z "${CONFIG_JSON}" ]] && echo "${CONFIG_JSON}" > /goloop/config.json

ENV GOCHAIN_DB_TYPE=rocksdb
ENV GOCHAIN_CLEAN_DATA=true
ENV JAVAEE_BIN=/goloop/execman/bin/execman
ENV PYEE_VERIFY_PACKAGE=true
