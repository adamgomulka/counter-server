FROM golang:1.10.4 as build
ADD core/core.go /go/src/counter-server/core/core.go
ADD rpcserv/rpcserv.go /go/src/counter-server/rpcserv/rpcserv.go
ADD gateway/gateway.go /go/src/counter-server/gateway/gateway.go
ADD tcp_wait/tcp_wait.go /go/src/counter-server/tcp_wait/tcp_wait.go
WORKDIR /go/src/counter-server
RUN go get github.com/go-redis/redis
RUN go build -o /go/bin/gateway gateway/gateway.go
RUN go build -o /go/bin/rpcserv rpcserv/rpcserv.go
RUN go build -o /go/bin/tcp_wait tcp_wait/tcp_wait.go

FROM debian:stretch-slim
COPY --from=build /go/bin/gateway /usr/bin/gateway
COPY --from=build /go/bin/rpcserv /usr/bin/rpcserv
COPY --from=build /go/bin/tcp_wait /usr/bin/tcp_wait
ADD start-server.sh /usr/local/bin/start-server.sh
CMD ["/usr/local/bin/start-server.sh"]
