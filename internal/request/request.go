package request

import (
    "io"
    "strings"
    "errors"
    "fmt"
    "bytes"
)

const (
    crlf = "\r\n"
    initialized = 0
    done = 1
)

var (
    ErrInvalidRequest = errors.New("invalid request")
)

type Request struct {
    RequestLine RequestLine
    State       int 
}

type RequestLine struct {
    HttpVersion    string
    RequestTarget   string
    Method          string
}

func RequestFromReader(r io.Reader) (*Request, error) {
    req := &Request{
        State:  initialized,
    }

    bytesRead := 0
    bytesParsed := 0
    var buf []byte
    for req.State != done {
        // make sure we can always add at least 8 bytes of data to the buffer
        dat := make([]byte, 8) 
        n, err := r.Read(dat)
        buf = append(buf, dat[:n]...)
        bytesRead += n

        n, err = req.Parse(buf)
        if err != nil {
            return nil, err
        }

        bytesParsed += n
    }
    return req, nil
}

func (r *Request) Parse(dat []byte) (int, error) {
    // we know were not done request line until found CRLF
    n, reqline, err := ParseRequestLine(dat)
    if err != nil {
        return 0, err
    }

    if n > 0 {
        r.State = done
        r.RequestLine = *reqline
    }


    return n, nil 
}


func ParseRequestLine(dat []byte) (int, *RequestLine, error) {
    idx := bytes.Index(dat, []byte(crlf))
    if idx == -1 {
        return 0, nil, nil
    }

    reqlineText := string(dat[:idx])
    reqline, err := RequestLineFromString(reqlineText)
    if err != nil {
        return idx, nil, err
    }
    return idx, reqline, nil
}

func RequestLineFromString(str string) (*RequestLine, error) {
    parts := strings.Split(str, " ")   
    if len(parts) != 3 {
        return nil, fmt.Errorf("request line requires 3 parts found %d", len(parts))
    }

    method := parts[0]
    if method != strings.ToUpper(method) {
        return nil, fmt.Errorf("method declaration has lowercase characters")
    }

    target := parts[1]
    versionParts := strings.Split(parts[2], "/")
    if len(versionParts) != 2 {
        return nil, fmt.Errorf("malformed start-line %s", str)
    }

    httpPart := versionParts[0]
    if httpPart != "HTTP" {
        return nil, fmt.Errorf("invalid http version %s", httpPart)
    }

    version := versionParts[1]
    if version != "1.1" {
        return nil, fmt.Errorf("invalid http version %s", version)
    }

    reqline := &RequestLine{
        HttpVersion:    version,
        RequestTarget:  target,
        Method:         method,
    }

    return reqline, nil
}
