FROM ubuntu:18.04

ARG TARGETARCH
ARG GOLANG_VERSION="1.16.3"

SHELL ["/bin/bash", "-c"]

ENV GOPATH=/root/go
ENV GO111MODULE=on

RUN apt update && apt upgrade -y && \
    apt install libgmp-dev libssl-dev curl git \
    psmisc dnsutils jq make gcc g++ bash tig tree sudo vim \
    silversearcher-ag unzip emacs-nox nano bash-completion -y

RUN git clone https://github.com/harmony-one/bls.git
RUN git clone https://github.com/harmony-one/mcl.git

RUN cd bls && make -j8 BLS_SWAP_G=1 && make install
RUN cd mcl && make install
RUN rm -rf bls mcl
