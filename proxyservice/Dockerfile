FROM golang:1.9.2-alpine3.6 as go_builder

RUN apk update
RUN apk add git swig make bash

RUN go get -u github.com/go-swagger/go-swagger/cmd/swagger
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/kardianos/govendor

WORKDIR /go/src/ProxyService/proxyservice

COPY . .

ARG BUILD

RUN make clean && make swagger
RUN make prepare
RUN make build

FROM alpine:3.6
RUN apk update && apk add ca-certificates
RUN mkdir ProxyService

COPY --from=go_builder /go/src/ProxyService/proxyservice/main ProxyService/
WORKDIR /ProxyService

CMD ["./main"]
EXPOSE 3000 9090
