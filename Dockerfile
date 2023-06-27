FROM golang:1.17-alpine

ENV GOPATH /go

ENV env=docker

RUN mkdir -p "$GOPATH/src/github.com/alvinatthariq/walletsvc" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

ADD . ${GOPATH}/src/github.com/alvinatthariq/walletsvc/

WORKDIR ${GOPATH}/src/github.com/alvinatthariq/walletsvc

COPY go.mod go.sum ./

RUN go get ./...

COPY *.go *.json ./

RUN apk update && apk add --no-cache git

RUN CGO_ENABLED=0 GOOS=linux go build -tags dynamic -o walletsvc

EXPOSE 8080

ENTRYPOINT ["/go/src/github.com/alvinatthariq/walletsvc/walletsvc"]