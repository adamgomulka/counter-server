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
/*
type CounterHandler struct {
    counts map[string]int
}
*/
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
    s.listener.Close()
}
/*
func (c *CounterHandler) Execute(r RpcRequest, w *RpcResponse) (e error) {
    if strings.HasPrefix(r.Name[1:], "hello") {
        if r.Method == "DELETE" {
            fmt.Printf("Method DELETE%sRunning resetCounter()%s", "\n", "\n")
            *w = c.resetCounter()
            fmt.Printf("Finished running resetCounter()%s", "\n")
        } else if r.Method == "GET" {
            fmt.Printf("Method GET%s", "\n")
            n := r.Name[len("/hello/:"):]
            fmt.Printf("Running serveHello(%s)%s", n, "\n")
            *w = c.serveHello(n)
            fmt.Printf("Finished running serveHello()%s", "\n")
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
*/
func main() {
    //counts := &CounterHandler{counts: map[string]int{}}
    counts := CreateCounterHandler()
    server := RpcServer{ServerName: service}
    server.init(counts)
    defer server.Close()
}
