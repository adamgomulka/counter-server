package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "encoding/json"
)

type CounterHandler struct {
    counts map[string]int
}

func (c *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if strings.HasPrefix(r.URL.Path[1:], "hello") {
        c.ServeHello(r.URL.Path[len("/hello/:"):], w)
    } else if strings.HasPrefix(r.URL.Path[1:], "counts") {
        if r.Method == "GET" {
            c.GetCounter(w)
        } else if r.Method == "DELETE" {
            c.ResetCounter()
        } else {
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        }
    } else {
         http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    }
}

func (c *CounterHandler) ServeHello(n string, w http.ResponseWriter) {
    if _, p := c.counts[n]; p {
        c.counts[n]++
    } else {
        c.counts[n] = 1
    }
    fmt.Fprintf(w, "Hello, %s!", n)
}

func (c *CounterHandler) GetCounter(w http.ResponseWriter) {
    j, err := json.Marshal(c.counts)
    if err == nil {
        w.Write(j)
    }
}

func (c *CounterHandler) ResetCounter() {
    c.counts = make(map[string]int)
}

func main() {
    server := http.NewServeMux()
    counter := &CounterHandler{counts: map[string]int{}}
    server.Handle("/hello/", counter)
    server.Handle("/counts", counter)
    server.HandleFunc("/health", HealthCheck)
    log.Fatal(http.ListenAndServe(":8080", server))
}

