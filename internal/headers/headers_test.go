package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("     Host:      localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 40, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["User-Agent"] = "boots"
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.Equal(t, 2, len(headers))
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid char in header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid header with multiple occurences
	headers = NewHeaders()
	data = []byte("Set-Person: lane loves go\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "lane loves go", headers["set-person"])
	assert.Equal(t, 27, n)
	assert.False(t, done)

	data2 := []byte("Set-Person: prime loves zig\r\n")
	n, done, err = headers.Parse(data2)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "lane loves go, prime loves zig", headers["set-person"])
	assert.Equal(t, 29, n)
	assert.False(t, done)

	data3 := []byte("Set-Person: tj loves ocaml\r\n\r\n")
	n, done, err = headers.Parse(data3)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "lane loves go, prime loves zig, tj loves ocaml", headers["set-person"])
	assert.Equal(t, 28, n)
	assert.False(t, done)
}
