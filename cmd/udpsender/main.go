package main

import (
    "net"
    "log"
    "bufio"
    "os"
    "fmt"
)

const (
    port = ":42069"
)

func main() {
    udpaddr, err := net.ResolveUDPAddr("udp", port)
    if err != nil {
        log.Fatal(err)
    }

    conn, err := net.DialUDP("udp", nil, udpaddr)
    if err != nil {
        log.Fatal(err)
    }

    defer conn.Close()
    r := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("> ")
        str, err := r.ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }

        _, err = conn.Write([]byte(str))
        if err != nil {
            log.Fatalf("could not write to %s err %s", conn.RemoteAddr(), err.Error())
        }
    }
}
