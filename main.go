package main

import (
    "fmt"
    "os"
    "log"
    "errors"
    "io"
)

func main() {
    r, err := os.Open("messages.txt")
    if err != nil {
        log.Fatal(err)
    }
    b := make([]byte, 8)
    for {
        _, err := r.Read(b)
        if errors.Is(err, io.EOF) {
            log.Fatal(err)
        }

        fmt.Printf("read: %s", string(b))
    }
}
