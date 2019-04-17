FROM dhiltgen/golang-armv7l:1.9.2-alpine3.6 as builder

ENV GOPATH /root/go

RUN go get \
    github.com/stianeikeland/go-rpio \
    github.com/stretchr/testify/require \
    github.com/emicklei/go-restful \
    github.com/prometheus/client_golang/prometheus


COPY . /go/src/github.com/dhiltgen/sprinklers

# RUN go test -v github.com/dhiltgen/sprinklers/...

RUN go build -a -tags "netgo static_build" -installsuffix netgo \
    -ldflags "-w -extldflags '-static'" -o /bin/sprinklers github.com/dhiltgen/sprinklers/cmd


FROM alpine:3.6

COPY --from=builder /bin/sprinklers /bin
COPY --from=builder /go/src/github.com/dhiltgen/sprinklers/circuits.json .

CMD /bin/sprinklers
