package main

import (
    "os"
    "net"
)

func main() {
    _, e := net.Dial("tcp", os.Args[1] + ":" + os.Args[2])
    for e != nil {
        _, e = net.Dial("tcp", os.Args[1] + ":" + os.Args[2])
    }

}
