package core

import (
    "fmt"
    "strings"
    "encoding/json"
    "syscall"
    "unsafe"
    "time"

    "github.com/go-redis/redis"
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
    redis *redis.Client
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
    h = &CounterHandler{}
    h.redis = redis.NewClient(&redis.Options{Addr: "redis:6379", DB: 0})
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
    names := c.redis.Keys("*").Val()
    for _, n := range names {
        c.redis.Del(n)
    }
    return RpcResponse{StatusCode: 200}
}

func (c *CounterHandler) serveHello(n string) RpcResponse {
    if c.redis.Get(n).Err() != redis.Nil {
        c.redis.Incr(n)
    } else {
        c.redis.Set(n, 1)
    }
    m := []byte(fmt.Sprintf("Hello, %s!", n))
    return RpcResponse{Message: m, StatusCode: 200}
}

func (c *CounterHandler) getCounter() RpcResponse {
    names := c.redis.Keys("*").Val()
    counts := map[string]int{}
    for _, n := range names {
        counts[n] , _ = c.redis.Get(n).Int()
    }
    j, _ := json.Marshal(counts)
    return RpcResponse{Message: j, StatusCode: 200}
}

func (c *CounterHandler) getHealth() (r RpcResponse) {
    i := &Sysinfo_t{}
    syscall.Syscall(99, uintptr(unsafe.Pointer(i)), 0, 0)
    j, _ := json.Marshal(*i)
    return RpcResponse{Message: j, StatusCode: 200}
}
