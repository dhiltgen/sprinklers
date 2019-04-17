#FROM dhiltgen/golang-armv7l:1.9.2-alpine3.6 as builder
FROM golang:1.12-alpine as builder

#ENV GOPATH /root/go

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
    

# Protobufs...
RUN git clone https://github.com/google/protobuf /tmp/protobuf && \
    cd /tmp/protobuf && git checkout v3.7.1 && \
    ./autogen.sh && \
    ./configure && make install

RUN go get -u github.com/golang/protobuf/protoc-gen-go

COPY . /go/src/github.com/dhiltgen/sprinklers
WORKDIR /go/src/github.com/dhiltgen/sprinklers

# RUN go test -v github.com/dhiltgen/sprinklers/...

# TODO wire this up better
RUN cd api && protoc ./sprinklers.proto --go_out=plugins=grpc:sprinklers

RUN go build -a -tags "netgo static_build" -installsuffix netgo \
    -ldflags "-w -extldflags '-static'" -o /bin/sprinklerd ./api/main/main.go
RUN go build -a -tags "netgo static_build" -installsuffix netgo \
    -ldflags "-w -extldflags '-static'" -o /bin/sprinklers ./client/main.go


FROM alpine:3.9

COPY --from=builder /bin/sprinklerd /bin
COPY --from=builder /bin/sprinklers /bin
COPY --from=builder /go/src/github.com/dhiltgen/sprinklers/circuits.json .

CMD /bin/sprinklerd
