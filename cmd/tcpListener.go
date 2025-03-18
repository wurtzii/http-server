package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
    "net"
)

func main() {
    listener, err := net.Listen("tcp", ":42069")
    if err != nil {
        log.Fatal(err)
    }

    defer  listener.Close()
    for {
        conn, err := listener.Accept()
        fmt.Println("connection accepted")
        if err != nil {
            log.Fatal(err)
        }
        ch := getLinesChannel(conn)
        
        for s := range ch {
            fmt.Printf("%s\n", s)
        }
        fmt.Println("connection closed")
    }
}

func getLinesChannel(f io.ReadCloser) <-chan string {
    ch := make(chan string)
    go func() {
        defer f.Close()
        defer close(ch)
        var curr_line string
        for {
            b := make([]byte, 8)
            n, err := f.Read(b)
            if n == 0 {
                if curr_line != "" {
                    ch <- curr_line
                }
                break
            }

            if err != nil {
                if curr_line != "" {
                    ch <- curr_line
                }
                
                if errors.Is(err, io.EOF) {
                    break
                }

                break
            }

            parts := strings.Split(string(b), "\n")
            if len(parts) > 1 {
                for _, s := range parts[:len(parts) -1] {
                    ch <- fmt.Sprintf("%s%s", curr_line, s)
                    curr_line = ""
                }

                curr_line += parts[len(parts) -1]
            } else {
                curr_line += parts[0]
            }
        }
    }()

    return ch
}
