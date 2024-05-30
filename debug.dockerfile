# This version should match that in .nvmrc
FROM node:22.2.0 AS nodebuilder

WORKDIR /go/src/github.com/elielamora/kratos-selfservice-ui-go

ADD . .

RUN make clean build-css

ADD . .

# This version should match the version in go.mod
FROM golang:1.22 AS gobuilder

WORKDIR /go/src/github.com/elielamora/kratos-selfservice-ui-go

ADD go.mod go.mod
ADD go.sum go.sum

ENV GO111MODULE on
ENV CGO_ENABLED 0

RUN go mod download

ADD . .

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN go build -ldflags="-extldflags=-static" -gcflags="all=-N -l" -o /usr/bin/kratos-selfservice-ui-go

FROM scratch
COPY --from=gobuilder /usr/bin/kratos-selfservice-ui-go /
COPY --from=gobuilder /go/bin/dlv /

# Expose the default port that we will be listening to
EXPOSE 4455 40000

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/kratos-selfservice-ui-go"]