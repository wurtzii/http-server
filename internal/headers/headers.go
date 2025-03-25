package headers

import (
	"bytes"
	"errors"
	"fmt"
    "strings"
)

type Headers map[string]string

const (
    crlf = "\r\n"
)

func NewHeaders() Headers {
    return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
    idx := bytes.Index(data, []byte(crlf))
    if idx == -1 { 
        return 0, false, nil
    }

    if idx == 0 {
        return 2, true, nil // account for consumed crlf
    }

    header := data[:idx]
    idx = bytes.Index(header, []byte(":"))
    if idx == -1 || idx == len(header) - 1 {
        return 0, false, errors.New("invalid header syntax")
    }

    key := string(header[:idx])
    key = strings.TrimLeft(key, " ")
    key = strings.ToLower(key)

    if len(key) < 1 {
        return 0, false, fmt.Errorf("found 0 length key")
    }

    if spcIdx := bytes.Index([]byte(key), []byte(" ")); spcIdx != -1 { // check for remaining spaces
        return 0, false, fmt.Errorf("space separated key. found '%s'", key)
    }

    if !OnlyValidChar([]byte(key)) {
        return 0, false, fmt.Errorf("invalid characters in key. found '%s'", key)
    }

    value := string(header[idx + 1:])
    value = strings.TrimSpace(value)

    if len(value) < 1 {
        return 0, false, fmt.Errorf("0 length value")
    }

    if val, ok := h[key]; ok {
        h[key] = fmt.Sprintf("%s, %s", val, value)
    } else {
        h[key] = value
    }
    return len(header) + 2, false, nil // account for the extra 2 bytes because of crlf
}

func OnlyValidChar(dat []byte) bool {
    validSymbols := ("!#$%&'*+-.^_`|~")
    for _, b := range dat {
        if b >= '0' && b <= '9' {
            continue
        }

        if b >= 'A' && b <= 'z' {
            continue
        }

        if strings.Contains(validSymbols, string(b)) {
            continue
        }

        return false
    }

    return true
}
