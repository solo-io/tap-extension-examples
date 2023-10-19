FROM golang:1.20 as build-env

ARG GO_MAIN_FILE

# ADD ./ /go/src/github.com/solo-io/tap-extension-examples/
RUN mkdir -p /go/src/github.com/solo-io/tap-extension-examples/
WORKDIR /go/src/github.com/solo-io/tap-extension-examples/

ADD pkg/ ./pkg
ADD tap-server-grpc ./tap-server-grpc
ADD tap-server-http ./tap-server-http
ADD go.mod go.sum ./

RUN echo $GO_MAIN_FILE
RUN GOOS=linux go build -o sample-tap-server $GO_MAIN_FILE

FROM alpine:latest
RUN apk add --no-cache libc6-compat
COPY --from=build-env /go/src/github.com/solo-io/tap-extension-examples/sample-tap-server /
CMD /sample-tap-server
