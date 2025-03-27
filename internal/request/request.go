package request

import (
    "io"
    "strings"
    "errors"
    "fmt"
    "bytes"
    "httpfromtcp/internal/headers"
)

const (
    crlf = "\r\n"
    requestStateInitialized = iota
    requestStateParsingHeaders
    requestStateDone
)

var (
    ErrInvalidRequest = errors.New("invalid request")
)

type Request struct {
    RequestLine RequestLine
    Headers     headers.Headers
    State       int 
}

type RequestLine struct {
    HttpVersion    string
    RequestTarget   string
    Method          string
}

func RequestFromReader(r io.Reader) (*Request, error) {
    req := &Request{
        State:  requestStateInitialized,
        Headers: headers.NewHeaders(),
    }

    bytesRead := 0
    bytesParsed := 0
    var buf []byte
    for req.State !=  requestStateDone {
        // make sure we can always add at least 8 bytes of data to the buffer
        dat := make([]byte, 8) 
        n, rerr := r.Read(dat)

        buf = append(buf, dat[:n]...)
        bytesRead += n

        n, err := req.Parse(buf)
        bytesParsed += n
        buf = buf[n:]
        if err != nil {
            return nil, err
        }

        if rerr != nil {
            if errors.Is(rerr, io.EOF) {
                if req.State != requestStateDone {
                    return nil, fmt.Errorf("eof but parsing not complete")
                }
            } else {
                return nil, err
            }
        }
    }

    return req, nil
}

func (r *Request) Parse(dat []byte) (int, error) {
    // we know were not done request line until found CRLF
    totalBytesParsed := 0
    switch r.State {
    case requestStateInitialized:
        n, reqline, err := ParseRequestLine(dat)
        if err != nil {
            return 0, err
        }

        if n > 0 {
            r.State = requestStateParsingHeaders 
            r.RequestLine = *reqline
        }
        
        totalBytesParsed += n

    case requestStateParsingHeaders:
        for {
            n, done, err := r.Headers.Parse(dat[totalBytesParsed:])
            totalBytesParsed += n
            if err != nil {
                return n, err
            }
            
            if n == 0 { // not enough data available to parse
                return totalBytesParsed, nil
            }

            if done { // headers are parsed
                r.State = requestStateDone
                break
            }
        }
    }

    return totalBytesParsed, nil 
}

func (r *Request) ParseSingleHeader(data []byte) (n int, err error) {
    n, done, err := r.Headers.Parse(data)

    if err != nil {
        return 0, err
    }

    if done {
        return n, nil
    }

    return n, nil;
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
    return idx + 2, reqline, nil
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
