package core

import (
    "fmt"
    "strings"
    "encoding/json"
    "syscall"
    "unsafe"
)

type RpcRequest struct {
    Name string
    Method string
}

type RpcResponse struct {
    StatusCode int
    Message []byte
}

type CounterHandler struct {
    counts map[string]int
}

type Sysinfo_t struct {
    Uptime    int64
    Loads     [3]uint64
    Totalram  uint64
    Freeram   uint64
    Sharedram uint64
    Bufferram uint64
    Totalswap uint64
    Freeswap  uint64
    Procs     uint16
    Pad       uint16
    Totalhigh uint64
    Freehigh  uint64
    Unit      uint32
    // contains filtered or unexported fields
}

func CreateCounterHandler() (h *CounterHandler) {
    h = &CounterHandler{counts: map[string]int{}}
    return
}

func (c *CounterHandler) Execute(r RpcRequest, w *RpcResponse) (e error) {
    if strings.HasPrefix(r.Name[1:], "hello") {
        if r.Method == "GET" {
            n := r.Name[len("/hello/:"):]
            *w = c.serveHello(n)
        } else {
            *w = RpcResponse{StatusCode: 400}
        }
    } else if strings.HasPrefix(r.Name[1:], "counts") {
        if r.Method == "GET" {
            *w = c.getCounter()
        } else if r.Method == "DELETE" {
            *w = c.resetCounter()
        } else {
            *w = RpcResponse{StatusCode: 400}
        }
    } else if strings.HasPrefix(r.Name[1:], "health") {
        if r.Method == "GET" {
            *w = c.getHealth()
        } else {
            *w = RpcResponse{StatusCode: 400}
        }
    }
    return nil
}

func (c *CounterHandler) resetCounter() RpcResponse {
    c.counts = make(map[string]int)
    return RpcResponse{StatusCode: 200}
}

func (c *CounterHandler) serveHello(n string) RpcResponse {
    if _, p := c.counts[n]; p {
        c.counts[n]++
    } else {
        c.counts[n] = 1
    }
    m := []byte(fmt.Sprintf("Hello, %s!", n))
    return RpcResponse{Message: m, StatusCode: 200}
}

func (c *CounterHandler) getCounter() RpcResponse {
    j, _ := json.Marshal(c.counts)
    return RpcResponse{Message: j, StatusCode: 200}
}

func (c *CounterHandler) getHealth() (r RpcResponse) {
    i := &Sysinfo_t{}
    syscall.Syscall(99, uintptr(unsafe.Pointer(i)), 0, 0)
    j, _ := json.Marshal(*i)
    return RpcResponse{Message: j, StatusCode: 200}
}
