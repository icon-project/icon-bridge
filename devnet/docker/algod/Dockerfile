FROM ubuntu:20.04

RUN apt-get update && \
    apt-get install -y wget && \
    mkdir /tmp/algorand && \
    wget -O /tmp/algorand/algorand.tar.gz https://github.com/algorand/go-algorand/releases/download/v3.13.3-stable/node_stable_linux-amd64_3.13.3.tar.gz && \
    cd /tmp/algorand && \
    tar xf algorand.tar.gz && \
    cd bin && \
    mv algod goal kmd /usr/local/bin && \
    rm -rf /tmp/algorand*

COPY ./template.json /tmp/template.json

RUN goal network create -r /testnet -t /tmp/template.json

COPY ./config.json /testnet/Node
COPY ./algod.token /testnet/Node
COPY ./kmd_config.json /testnet/Node/kmd-v0.5/kmd_config.json
COPY ./kmd.token /testnet/Node/kmd-v0.5/kmd.token

COPY ./scripts /scripts

EXPOSE 4001
EXPOSE 4002

ENTRYPOINT ["/scripts/run_algorand_testnet.sh"]