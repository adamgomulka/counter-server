FROM golang:1.10.4 as build
RUN mkdir /go/src/nyt-server
ADD main.go /go/src/nyt-server
ADD healthcheck.go /go/src/nyt-server
WORKDIR /go/src/nyt-server
RUN go build -o /go/bin/nyt-server .

FROM debian:stretch-slim
COPY --from=build /go/bin/nyt-server /usr/bin/nyt-server
EXPOSE 8080/tcp
CMD ["/usr/bin/nyt-server"]
