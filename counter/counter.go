package counter

import (
    . "../core"
    "strings"
    "fmt"
    "encoding/json"
)

const service string = "counter"

type CounterHandler struct {
    RpcRequestHandler
    counts map[string]int
}

func (c *CounterHandler) Execute(r RpcRequest, w *RpcResponse) (e error) {
    if strings.HasPrefix(r.Name[1:], "hello") {
        if r.Method == "DELETE" {
            *w = c.ResetCounter()
        } else if r.Method == "GET" {
            n := r.Name[len("/hello/:"):]
            *w = c.ServeHello(n)
        } else {
            *w = RpcResponse{StatusCode: 400}
        }
    } else if strings.HasPrefix(r.Name[1:], "counts") {
        if r.Method == "GET" {
            *w = c.GetCounter()
        }
    }
    return nil
}

func (c *CounterHandler) ResetCounter() RpcResponse {
    c.counts = make(map[string]int)
    return RpcResponse{Message: "", StatusCode: 200}
}

func (c *CounterHandler) ServeHello(n string) RpcResponse {
    if _, p := c.counts[n]; p {
        c.counts[n]++
    } else {
        c.counts[n] = 1
    }
    return RpcResponse{Message: fmt.Sprintf("Hello, %s!", n), StatusCode: 200}
}

func (c *CounterHandler) GetCounter() RpcResponse {
    j, _ := json.Marshal(c.counts)
    return RpcResponse{Message: string(j), StatusCode: 200}
}

func main() {
    counts := CounterHandler{counts: map[string]int{}}
    server := RpcServer{ServerName: service}
    server.Init(counts)
    defer server.Close()
}
