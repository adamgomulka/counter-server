package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "encoding/json"
)

type counterHandler struct {
    counts map[string]int
}

func (c counterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (c *counterHandler) ServeHello(n string, w http.ResponseWriter) {
    if _, p := c.counts[n]; p {
        c.counts[n]++
    } else {
        c.counts[n] = 1
    }
    fmt.Fprintf(w, "Hello, %s!", n)
}

func (c *counterHandler) GetCounter(w http.ResponseWriter) {
    j, err := json.Marshal(c.counts)
    if err == nil {
        w.Write(j)
    }
}

func (c *counterHandler) ResetCounter() {
    c.counts = make(map[string]int)
}

func main() {
    server := http.NewServeMux()
    var counter counterHandler
    server.Handle("/hello/:", counter)
    server.Handle("/counts", counter)
    log.Fatal(http.ListenAndServe(":80", server))
}

