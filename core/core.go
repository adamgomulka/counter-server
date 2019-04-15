package core

import (
    "fmt"
    "strings"
    "encoding/json"
)

type RpcRequest struct {
    Name string
    Method string
}

type RpcResponse struct {
    StatusCode int
    Message string
}

type CounterHandler struct {
    counts map[string]int
}


func CreateCounterHandler() (h *CounterHandler) {
    h = &CounterHandler{counts: map[string]int{}}
    return
}

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
