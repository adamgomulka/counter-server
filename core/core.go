package core

import (
    "net/rpc"
    "net"
    "strconv"
)

const rpc_port int = 1337

type HttpRequestHandler struct {
    rpc_client RpcClient
}

type RpcRequestHandler struct{}

type RpcRequest struct {
    Name string
    Method string
}

type RpcResponse struct {
    StatusCode int
    Message string
}

type RpcClient struct {
    ServerName string
    client *rpc.Client
}

type RpcServer struct {
    ServerName string
    listener net.Listener
    server *rpc.Server
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

func (s *RpcServer) Init(h *RpcRequestHandler) (e error) {
    s.server = rpc.NewServer()
    s.server.Register(h)
    s.listener, e = net.Listen("tcp", ":" + strconv.Itoa(rpc_port))
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

