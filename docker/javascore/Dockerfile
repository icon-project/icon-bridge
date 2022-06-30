ARG GOLOOP_IMAGE=iconloop/goloop-icon:v0.9.7
FROM ${GOLOOP_IMAGE} as build-stage

RUN apk add git
RUN apk add gradle

ENV ROOT=/goloop
ENV JAVASCORE_SRC=/goloop/btp/javascore
ENV DIST_DIR=/dist
WORKDIR ${ROOT}

#copy javascore code & solidity code base into the container for future build & provision
ADD javascore ${JAVASCORE_SRC}

#build javascore jars
WORKDIR ${JAVASCORE_SRC}
RUN gradle bsh:clean
RUN gradle bsh:optimizedJar
RUN gradle nativecoin:clean
RUN gradle nativecoin:optimizedJar
RUN gradle bsr:clean
RUN gradle bsr:optimizedJar

#copy jars to dist
WORKDIR ${ROOT}
RUN mkdir /dist

RUN cp ${JAVASCORE_SRC}/bsh/build/libs/bsh-optimized.jar ${DIST_DIR}
RUN cp ${JAVASCORE_SRC}/nativecoin/build/libs/nativecoin-0.1.0-optimized.jar ${DIST_DIR}/nativecoin-optimized.jar
RUN cp ${JAVASCORE_SRC}/bsr/build/libs/restrictions-0.1.0-optimized.jar ${DIST_DIR}/restrictions-optimized.jar

RUN cp ${JAVASCORE_SRC}/lib/BMC-0.1.0-optimized.jar ${DIST_DIR}/bmc-optimized.jar
RUN cp ${JAVASCORE_SRC}/lib/irc2Tradeable-0.1.0-optimized.jar ${DIST_DIR}/irc2Tradeable-optimized.jar

RUN git clone https://github.com/icon-project/java-score-examples.git
WORKDIR /goloop/java-score-examples
RUN gradle irc2-token:clean
RUN gradle irc2-token:optimizedJar
WORKDIR ${ROOT}
RUN cp java-score-examples/irc2-token/build/libs/irc2-token-0.9.1-optimized.jar ${DIST_DIR}/irc2-token-optimized.jar

