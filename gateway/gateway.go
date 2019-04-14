package gateway

import (
    "fmt"
    "log"
    "net/http"
    "net/rpc"
    "strings"
    "encoding/json"
    "../counter"
    . "../core"
)

const (
    services map[string][]string{"counter": ["/hello/", "/counter"], "health": ["/health"]}
    RpcPort = 1337
)

func (c *RpcClient) Init() (e error) {
    addr := c.ServerName + ":" + strconv.Itoa(RpcPort)
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

func (s *RpcServer) Init(handler) (e error) {
    s.server = rpc.NewServer()
    s.server.Register(handler)
    s.listener, e = net.Listen("tcp", ":" + strconv.Itoa(RpcPort))
    if e != nil {
        return
    }
    s.server.Accept(s.listener)
    return
}

func (s *RpcServer) Close() (e error) {
    if s.listener != nil {
        e = s.listener.Close()
    }
    return
}

func (h *HttpRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (e error) {
    rpc_request := &RpcRequest{Name: r.URL.Path, Method: r.Method}
    rpc_response := RpcResponse{}
    handler_name = h.RpcClient.ServerName + ".Execute"
    e = h.rpc_client.client.Call(handler_name, rpc_request, rpc_response)
    if rpc_response.StatusCode == 200 {
        fmt.Fprintf(w, rpc_response.Message)
    } else {
        http.Error(w, http.StatusText(rpc_response.StatusCode), rpc_response.StatusCode)
    }
    return
}

func CreateRPCClients(services map[string][]string) (h map[string]RequestHandler, e []error) {
    var err error
    for s, _ := range services {
        rpc_client := RpcClient{ServerName: s, Port: rpc_port}
        err = rpc_client.Init()
        if err != nil {
            e = append(e, err)
        handler := RequestHandler{rpc_client: rpc_client}
        h[s] := handler
    }
    return
}

func (h *ServeMux) DefineRoutes(services map[string][]string, handlers map[string]RequestHandler) {
    for s, rs := range services {
        for _, r := range rs {
            h.Handle(r, handlers[s])
        }
    }
}

func main() {
    var handlers map[string]RequestHandler
    http_server := http.NewServeMux()
    handlers, e := CreateRPCClients()
    if len(e) == 0 {
        http_server.DefineRoutes(handlers)
    }
    log.Fatal(http.ListenAndServe(":80", http_server))
}

