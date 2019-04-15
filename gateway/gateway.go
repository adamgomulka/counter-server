package main

import (
    "fmt"
    "log"
    "net/http"
    "net/rpc"
    "strconv"
    . "../core"
)

var services = map[string][]string{"CounterHandler": []string{"/hello/", "/counter"}, "health": []string{"/health"}}
var rpc_port = 1377

type RequestHandler struct {
    ServerName string
    client *rpc.Client
}

func (h *RequestHandler) init() (e error) {
    //addr := c.ServerName + ":" + strconv.Itoa(rpc_port)
    fmt.Printf("Initializing request handler%s", "\n")
    addr := "127.0.0.1:" + strconv.Itoa(rpc_port)
    h.client, e = rpc.Dial("tcp", addr)
    fmt.Printf("Request handler initialized: %s%s", h.ServerName, "\n")
    if e != nil {
        log.Print(e)
    } else if e == nil {
        fmt.Printf("No errors encountered%s", "\n")
    }
    return
}


/*
func (h *RequestHandler) close() (e error) {
    if c.client != nil {
        e = c.client.Close()
        return
    }
    return
}
*/

func (h RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    rpc_request := &RpcRequest{Name: r.URL.Path, Method: r.Method}
    rpc_response := new(RpcResponse)
    handler_name := h.ServerName + ".Execute"
    h.client.Call(handler_name, rpc_request, rpc_response)
    /* if e != nil {
        fmt.Printf("Error handling request in ServeHTTP->Call%s", "\n")
        log.Print(e)
    }
    if rpc_response.StatusCode == 200 {
        fmt.Fprintf(w, rpc_response.Message)
    } else {
        http.Error(w, http.StatusText(rpc_response.StatusCode), rpc_response.StatusCode)
    }*/
    fmt.Printf("Status Code: %s%s", strconv.Itoa(rpc_response.StatusCode), "\n")
    fmt.Printf("Message: %s%s", rpc_response.Message, "\n")
}

func createRPCClients(services map[string][]string) (h map[string]*RequestHandler) {
    h = make(map[string]*RequestHandler)
    for s, _ := range services {
        handler := &RequestHandler{ServerName: s}
        handler.init()
        h[s] = handler
    }
    return
}

func defineRoutes(services map[string][]string, handlers map[string]*RequestHandler, h *http.ServeMux) {
    for s, rs := range services {
        for _, r := range rs {
            h.Handle(r, handlers[s])
        }
    }
}

func main() {
    http_server := http.NewServeMux()
    handlers := createRPCClients(services)
    defineRoutes(services, handlers, http_server)
    log.Fatal(http.ListenAndServe(":8080", http_server))
}
