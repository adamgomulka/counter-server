package main

import (
    "syscall"
    "net/http"
    "encoding/json"
    "unsafe"
)

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

func HealthCheck(w http.ResponseWriter, r *http.Request) {
    i := &Sysinfo_t{}
    syscall.Syscall(99, uintptr(unsafe.Pointer(i)), 0, 0)
    j, err := json.Marshal(*i)
    if err == nil {
        w.Write(j)
    } else {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
}
