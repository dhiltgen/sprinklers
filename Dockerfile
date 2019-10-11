ARG GO_VERSION=1.13
ARG ALPINE_VERSION=3.10
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as builder

# Get the baseline dependencies installed
RUN apk --update add \
    git \
    autoconf \
    automake \
    bison \
    flex \
    libtool \
    libc-dev \
    gcc \
    g++ \
    libgcc \
    zlib-dev \
    file \
    make \
    tar

# Build Protobufs...
RUN git clone https://github.com/google/protobuf /tmp/protobuf && \
    cd /tmp/protobuf && git checkout v3.7.1 && \
    ./autogen.sh && \
    ./configure && \
    make install

RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u golang.org/x/lint/golint

COPY . /go/src/github.com/dhiltgen/sprinklers
WORKDIR /go/src/github.com/dhiltgen/sprinklers

FROM builder as build

# TODO this needs some cleanup...
RUN cd api && protoc ./sprinklers.proto --go_out=plugins=grpc:sprinklers

# Build the daemon and client
RUN go build -a -tags "netgo static_build" -installsuffix netgo \
    -ldflags "-w -extldflags '-static'" -o /bin/sprinklerd ./api/main/main.go
RUN go build -a -tags "netgo static_build" -installsuffix netgo \
    -ldflags "-w -extldflags '-static'" -o /bin/sprinklers ./client/main.go

# Unit test setup
FROM builder as test
RUN go test -coverprofile=/cover.out -v github.com/dhiltgen/sprinklers/...
RUN go tool cover -html=./cover.out -o /cover.html
RUN go list github.com/dhiltgen/sprinklers/... | grep -v /vendor/ | xargs golint -set_exit_status

# Target for easily extracing the coverage log
FROM scratch as coverage
COPY --from=test /cover.* /

# Minimal image for the daemon
ARG ALPINE_VERSION=3.10
FROM alpine:${ALPINE_VERSION} as daemon
COPY --from=build /bin/sprinklerd /bin
COPY --from=build /go/src/github.com/dhiltgen/sprinklers/circuits.json .
ENTRYPOINT ["/bin/sprinklerd"]

# Minimal image for the client
ARG ALPINE_VERSION=3.10
FROM alpine:${ALPINE_VERSION} as client
COPY --from=build /bin/sprinklers /bin
ENTRYPOINT ["/bin/sprinklers"]
