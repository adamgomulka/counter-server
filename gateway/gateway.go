package main

import (
    "fmt"
    "log"
    "net/http"
    "net/rpc"
    "strconv"
    . "../core"
)

var services = map[string][]string{"counter": []string{"/hello/", "/counter"}, "health": []string{"/health"}}
var rpc_port = 1337

type RpcClient struct {
    ServerName string
    client *rpc.Client
}

type HttpRequestHandler struct {
    rpc_client RpcClient
}

func (c *RpcClient) Init() (e error) {
    addr := c.ServerName + ":" + strconv.Itoa(rpc_port)
    c.client, e = rpc.Dial("tcp", addr)
    return
}

func (c *RpcClient) Close() (e error) {
    if c.client != nil {
        e = c.client.Close()
        return
    }
    return
}

func (h HttpRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    rpc_request := &RpcRequest{Name: r.URL.Path, Method: r.Method}
    rpc_response := RpcResponse{}
    handler_name := h.rpc_client.ServerName + ".Execute"
    h.rpc_client.client.Call(handler_name, rpc_request, rpc_response)
    if rpc_response.StatusCode == 200 {
        fmt.Fprintf(w, rpc_response.Message)
    } else {
        http.Error(w, http.StatusText(rpc_response.StatusCode), rpc_response.StatusCode)
    }
}

func CreateRPCClients(services map[string][]string) (h map[string]HttpRequestHandler, e []error) {
    var err error
    h = make(map[string]HttpRequestHandler)
    for s, _ := range services {
        rpc_client := &RpcClient{ServerName: s}
        err = rpc_client.Init()
        if err != nil {
            e = append(e, err)
        }
        handler := &HttpRequestHandler{rpc_client: *rpc_client}
        h[s] = *handler
    }
    return
}

func DefineRoutes(services map[string][]string, handlers map[string]HttpRequestHandler, h *http.ServeMux) {
    for s, rs := range services {
        for _, r := range rs {
            h.Handle(r, handlers[s])
        }
    }
}

func main() {
    var handlers map[string]HttpRequestHandler
    http_server := http.NewServeMux()
    handlers, e := CreateRPCClients(services)
    if len(e) == 0 {
        DefineRoutes(services, handlers, *http_server)
    }
    log.Fatal(http.ListenAndServe(":8080", http_server))
}
