package main

import (
    . "../core"
    "strings"
    "fmt"
    "encoding/json"
    "net/rpc"
    "net"
    "strconv"
    "log"
)

const service string = "counter"

type CounterHandler struct {
    counts map[string]int
}

type RpcServer struct {
    ServerName string
    listener net.Listener
    server *rpc.Server
}

func (s *RpcServer) Init(c *CounterHandler) {
    var e error
    s.server = &rpc.Server{}
    e = s.server.Register(c)
    if e != nil {
        log.Print(e)
    }
    s.listener, e = net.Listen("tcp", ":" + strconv.Itoa(1377))
    if e != nil {
        log.Print(e)
    }
    s.server.Accept(s.listener)
    fmt.Print("init completed")
}

func (s *RpcServer) Close() {
    s.listener.Close()
}

func (c *CounterHandler) Execute(r RpcRequest, w *RpcResponse) (e error) {
    if strings.HasPrefix(r.Name[1:], "hello") {
        if r.Method == "DELETE" {
            *w = c.resetCounter()
        } else if r.Method == "GET" {
            n := r.Name[len("/hello/:"):]
            *w = c.serveHello(n)
        } else {
            *w = RpcResponse{StatusCode: 400}
        }
    } else if strings.HasPrefix(r.Name[1:], "counts") {
        if r.Method == "GET" {
            *w = c.getCounter()
        }
    }
    return nil
}

func (c *CounterHandler) resetCounter() RpcResponse {
    c.counts = make(map[string]int)
    return RpcResponse{Message: "", StatusCode: 200}
}

func (c *CounterHandler) serveHello(n string) RpcResponse {
    if _, p := c.counts[n]; p {
        c.counts[n]++
    } else {
        c.counts[n] = 1
    }
    return RpcResponse{Message: fmt.Sprintf("Hello, %s!", n), StatusCode: 200}
}

func (c *CounterHandler) getCounter() RpcResponse {
    j, _ := json.Marshal(c.counts)
    return RpcResponse{Message: string(j), StatusCode: 200}
}

func main() {
    counts := &CounterHandler{counts: map[string]int{}}
    server := RpcServer{ServerName: service}
    server.Init(counts)
    defer server.Close()
}
