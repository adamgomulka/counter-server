FROM golang:1.10.4 as build
RUN go get github.com/adamgomulka/counter-server/core
ADD gateway.go /go/src/counter-server/gateway/gateway.go
WORKDIR /go/src/counter-server/gateway
RUN go build -o /go/bin/gateway .

FROM debian:stretch-slim
COPY --from=build /go/bin/gateway /usr/bin/gateway
EXPOSE 8080/tcp
CMD ["/usr/bin/gateway"]
