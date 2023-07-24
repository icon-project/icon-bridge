FROM golang:1.20-alpine
WORKDIR /app

COPY . .
#RUN go mod download

#RUN go get -u blockwatch.cc/tzgo

#RUN go get blockwatch.cc/tzgo/contract

#RUN go get blockwatch.cc/tzgo/micheline

#RUN go get blockwatch.cc/tzgo/rpc

#RUN go get blockwatch.cc/tzgo/signer

#RUN go get blockwatch.cc/tzgo/tezos

#RUN  go get github.com/echa/log

#RUN go get github.com/joho/godotenv

RUN go mod tidy

COPY *.go ./

RUN go build -o tezos-address

EXPOSE 8080

CMD [ "./tezos-address" ]
