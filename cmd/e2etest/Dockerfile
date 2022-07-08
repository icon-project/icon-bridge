FROM go:1.18

RUN apt update && apt upgrade -y && \
    apt install libgmp-dev libssl-dev curl git \
    psmisc dnsutils jq make gcc g++ bash tig tree sudo vim \
    silversearcher-ag unzip emacs-nox nano bash-completion -y

RUN git clone https://github.com/harmony-one/bls.git
RUN git clone https://github.com/harmony-one/mcl.git

RUN cd bls && make -j8 BLS_SWAP_G=1 && make install && cd ..
RUN cd mcl && make install && cd ..

CMD go run main.go
