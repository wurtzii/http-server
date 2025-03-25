package headers

import(
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {
    // Test: Valid single header
    headers := NewHeaders()
    data := []byte("Host: localhost:42069\r\n\r\n")
    n, done, err := headers.Parse(data)
    require.NoError(t, err)
    require.NotNil(t, headers)
    assert.Equal(t, "localhost:42069", headers["host"])
    assert.Equal(t, 23, n)
    assert.False(t, done)

    // Test: Invalid spacing header
    headers = NewHeaders()
    data = []byte("       Host : localhost:42069       \r\n\r\n")
    n, done, err = headers.Parse(data)
    require.Error(t, err)
    assert.Equal(t, 0, n)
    assert.False(t, done)

    // Test: Valid single header with extra whitespace
    headers = NewHeaders()
    data = []byte("     Host:      localhost:42069\r\n\r\n")
    n, done, err = headers.Parse(data)
    require.NoError(t, err)
    assert.Equal(t, 33, n)
    assert.False(t, done)

    // Test: Valid done
    headers = NewHeaders()
    data = []byte("\r\n")
    n, done, err = headers.Parse(data)
    require.NoError(t, err)
    assert.Equal(t, 2, n)
    assert.True(t, done)

    // Test: Invalid character header
    headers = NewHeaders()
    data = []byte("H©st: localhost:42069\r\n")
    n, done, err = headers.Parse(data)
    require.Error(t, err)
    assert.Equal(t, 0, n)
    assert.False(t, done)
    
    // Test: Valid multiple values to one key header
    headers = NewHeaders()
    data = []byte("Set-Person: lane-loves-go\r\n Set-Person: prime-loves-zig\r\n Set-Person: tj-loves-ocaml\r\n")
    td := 0
    n, done, err = headers.Parse(data)
    td += n
    n, done, err = headers.Parse(data[td + 1:])
    td += n
    n, done, err = headers.Parse(data[td +1:])
    td += n
    require.NoError(t, err)
    assert.Equal(t, 85, td)
    assert.False(t, done)
    assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
}
