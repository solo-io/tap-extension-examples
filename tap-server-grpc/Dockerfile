FROM golang:1.20 as build-env

ADD sample-tap-server-grpc/ /go/src/github.com/solo-io/tap-extension-examples/
WORKDIR /go/src/github.com/solo-io/tap-extension-examples/
RUN GOOS=linux go build -o sample-tap-server /go/src/github.com/solo-io/tap-extension-examples/main.go

FROM alpine:latest
RUN apk add --no-cache libc6-compat
COPY --from=build-env /go/src/github.com/solo-io/tap-extension-examples/sample-tap-server /
CMD /sample-tap-server
