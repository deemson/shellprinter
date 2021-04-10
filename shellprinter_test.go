package shellprinter_test

import (
	"bytes"
	"github.com/deemson/shellprinter"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNothing_Trivial(t *testing.T) {
	buf := bytes.NewBufferString("")
	p := shellprinter.New(buf)
	bytesWritten, err := p.Write([]byte("hello"))
	require.Nil(t, err)
	require.Equal(t, 5, bytesWritten)
	err = p.Close()
	require.Nil(t, err)
	require.Equal(t, "hello", buf.String())
}

func TestPrefix_Trivial(t *testing.T) {
	buf := bytes.NewBufferString("")
	p := shellprinter.New(buf).WithPrefixString("=>")
	bytesWritten, err := p.Write([]byte("hello"))
	require.Nil(t, err)
	require.Equal(t, 5, bytesWritten)
	err = p.Close()
	require.Nil(t, err)
	require.Equal(t, "=>hello", buf.String())
}

func TestPrefix_3LinesAtOnce(t *testing.T) {
	buf := bytes.NewBufferString("")
	p := shellprinter.New(buf).WithPrefixString("=>")
	data := []byte("hello\nmagnificent\nworld")
	bytesWritten, err := p.Write(data)
	require.Nil(t, err)
	require.Equal(t, len(data), bytesWritten)
	err = p.Close()
	require.Nil(t, err)
	require.Equal(t, "=>hello\n=>magnificent\n=>world", buf.String())
}

func TestPrefix_3Lines2Writes(t *testing.T) {
	buf := bytes.NewBufferString("")
	p := shellprinter.New(buf).WithPrefixString("=>")
	for _, data := range [][]byte{[]byte("hello\nmagnif"), []byte("icent\nworld")} {
		bytesWritten, err := p.Write(data)
		require.Nil(t, err)
		require.Equal(t, len(data), bytesWritten)
	}
	err := p.Close()
	require.Nil(t, err)
	require.Equal(t, "=>hello\n=>magnificent\n=>world", buf.String())
}

func TestSuffix_Trivial(t *testing.T) {
	buf := bytes.NewBufferString("")
	p := shellprinter.New(buf).WithSuffixString("<=")
	data := []byte("hello\n")
	bytesWritten, err := p.Write(data)
	require.Nil(t, err)
	require.Equal(t, len(data), bytesWritten)
	err = p.Close()
	require.Nil(t, err)
	require.Equal(t, "hello<=\n", buf.String())
}
