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

RUN gradle clean 
RUN gradle bmc:optimizedJar
#RUN gradle bsr:optimizedJar
RUN gradle bts:optimizedJar
#copy jars to dist
WORKDIR ${ROOT}
RUN mkdir /dist

RUN cp ${JAVASCORE_SRC}/bmc/build/libs/bmc-0.1.0-optimized.jar ${DIST_DIR}/bmc.jar
#RUN cp ${JAVASCORE_SRC}/bsr/build/libs/restrictions-0.1.0-optimized.jar ${DIST_DIR}/bsr.jar
RUN cp ${JAVASCORE_SRC}/bts/build/libs/bts-0.1.0-optimized.jar ${DIST_DIR}/bts.jar
RUN cp ${JAVASCORE_SRC}/lib/irc2Tradeable-0.1.0-optimized.jar ${DIST_DIR}/irc2Tradeable.jar

RUN git clone https://github.com/icon-project/java-score-examples.git
WORKDIR /goloop/java-score-examples
RUN gradle irc2-token:clean
RUN gradle irc2-token:optimizedJar
WORKDIR ${ROOT}
RUN cp java-score-examples/irc2-token/build/libs/irc2-token-0.9.1-optimized.jar ${DIST_DIR}/irc2.jar

