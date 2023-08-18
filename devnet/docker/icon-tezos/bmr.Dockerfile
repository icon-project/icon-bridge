# dev build
FROM ubuntu:18.04

ARG TARGETARCH
ARG GOLANG_VERSION="1.20"

SHELL ["/bin/bash", "-c"]

ENV GOPATH=/root/go
ENV GO111MODULE=on
ENV GIMME_GO_VERSION=${GOLANG_VERSION}
ENV PATH="/root/bin:${PATH}"

RUN apt update && apt upgrade -y && \
    apt install libgmp-dev libssl-dev curl git \
    psmisc dnsutils jq make gcc g++ bash tig tree sudo vim \
    silversearcher-ag unzip emacs-nox nano bash-completion -y

RUN mkdir ~/bin && \
    curl -sL -o ~/bin/gimme \
    https://raw.githubusercontent.com/travis-ci/gimme/master/gimme && \
    chmod +x ~/bin/gimme

RUN eval "$(~/bin/gimme ${GIMME_GO_VERSION})"
RUN touch /root/.bash_profile && \
    gimme ${GIMME_GO_VERSION} >> /root/.bash_profile && \
    echo "GIMME_GO_VERSION='${GIMME_GO_VERSION}'" >> /root/.bash_profile && \
    echo "GO111MODULE='on'" >> /root/.bash_profile && \
    echo ". ~/.bash_profile" >> /root/.profile && \
    echo ". ~/.bash_profile" >> /root/.bashrc
ENV PATH="/root/.gimme/versions/go${GIMME_GO_VERSION}.linux.${TARGETARCH:-amd64}/bin:${GOPATH}/bin:${PATH}"
RUN . ~/.bash_profile

COPY . bmr
RUN cd bmr/cmd/iconbridge 

# prod build
FROM ubuntu:18.04
SHELL ["/bin/bash", "-c"]
RUN apt update -y && apt install -y make ca-certificates libssl-dev
COPY --from=0 /bmr/cmd/iconbridge/iconbridge /bin/iconbridge
RUN rm -rf /var/lib/apt/lists/*