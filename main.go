package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
    f, err := os.Open("messages.txt")
    ch := getLinesChannel(f)
    if err != nil {
        log.Fatal(err)
    }

    for s := range ch {
        fmt.Printf("read: %s\n", s)
    }
}

func getLinesChannel(f io.ReadCloser) <-chan string {
    ch := make(chan string)
    go func() {
        defer f.Close()
        defer close(ch)
        var curr_line string
        for {
            b := make([]byte, 8, 8)
            _, err := f.Read(b)
            if err != nil {
                if errors.Is(err, io.EOF) {
                    break
                }

                fmt.Println(err)
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
