package main

import (
    . "../core"
    "fmt"
    "net/rpc"
    "net"
    "strconv"
    "log"
)

const service string = "counter"
type RpcServer struct {
    ServerName string
    listener net.Listener
    server *rpc.Server
}

func (s *RpcServer) init(c *CounterHandler) {
    var e error
    fmt.Printf("Intializing RPC Server%s", "\n")
    s.server = rpc.NewServer()
    fmt.Printf("Registering Counter Handler on RPC Server%s", "\n")
    s.server.Register(c)
    fmt.Printf("Opening TCP listener on port 1377%s", "\n")
    s.listener, e = net.Listen("tcp", ":" + strconv.Itoa(1377))
    if e != nil {
        log.Print(e)
    } else if e == nil {
        fmt.Printf("TCP listener opened successfully%s", "\n")
    }
    fmt.Printf("Accepting RPC connections on TCP listnener%s", "\n")
    s.server.Accept(s.listener)
}

func (s *RpcServer) Close() {
    fmt.Printf("Closing listener")
    s.listener.Close()
}

func main() {
    counts := CreateCounterHandler()
    server := RpcServer{ServerName: service}
    server.init(counts)
    defer server.Close()
}
