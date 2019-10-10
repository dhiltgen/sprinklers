FROM golang:1.12-alpine as builder

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

COPY . /go/src/github.com/dhiltgen/sprinklers
WORKDIR /go/src/github.com/dhiltgen/sprinklers


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

# Target for easily extracing the coverage log
FROM scratch as coverage
COPY --from=test /cover.* /

# Minimal image for the daemon
FROM alpine:3.9 as daemon
COPY --from=builder /bin/sprinklerd /bin
COPY --from=builder /go/src/github.com/dhiltgen/sprinklers/circuits.json .
ENTRYPOINT ["/bin/sprinklerd"]

# Minimal image for the client
FROM alpine:3.9 as client
COPY --from=builder /bin/sprinklers /bin
ENTRYPOINT ["/bin/sprinklers"]
