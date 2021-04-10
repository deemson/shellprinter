package shellprinter_test

import (
	"github.com/deemson/shellprinter"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type mockWriter struct {
	mock.Mock
}

func (w *mockWriter) Write(data []byte) (int, error) {
	args := w.Called(data)
	return args.Int(0), args.Error(1)
}

func TestWritePrefixErrAfterClose(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m).WithPrefixString("bad")
	m.On("Write", []byte("bad")).Return(0, errors.New("prefix error"))
	bytesWritten, err := p.Write([]byte("hello"))
	require.Nil(t, err)
	require.Equal(t, 5, bytesWritten)
	err = p.Close()
	require.NotNil(t, err)
	require.Equal(t, "failed to write prefix: prefix error", err.Error())
}

func TestWritePrefixErr(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m).WithPrefixString("bad")
	m.On("Write", []byte("bad")).Return(0, errors.New("prefix error"))
	bytesWritten, err := p.Write([]byte("hello\nworld"))
	require.NotNil(t, err)
	require.Equal(t, "failed to write prefix: prefix error", err.Error())
	require.Equal(t, 0, bytesWritten)
}

func TestWritePrefixInconsistencyErr(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m).WithPrefixString("bad")
	m.On("Write", []byte("bad")).Return(1, nil)
	bytesWritten, err := p.Write([]byte("hello\nworld"))
	require.NotNil(t, err)
	require.Equal(t, "inconsistency when writing prefix: prefix len = 3, actually written = 1", err.Error())
	require.Equal(t, 0, bytesWritten)
}

func TestWriteDataErr(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m)
	m.On("Write", []byte("hello")).Return(0, errors.New("data error"))
	bytesWritten, err := p.Write([]byte("hello\nworld"))
	require.NotNil(t, err)
	require.Equal(t, "failed to write data: data error", err.Error())
	require.Equal(t, 0, bytesWritten)
}

func TestWriteDataInconsistencyErr(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m)
	m.On("Write", []byte("hello")).Return(2, nil)
	bytesWritten, err := p.Write([]byte("hello\nworld"))
	require.NotNil(t, err)
	require.Equal(t, "inconsistency when writing data: data len = 5, actually written = 2", err.Error())
	require.Equal(t, 0, bytesWritten)
}

func TestWriteNewLineErr(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m)
	m.On("Write", []byte("hello")).Return(5, nil)
	m.On("Write", []byte("\n")).Return(0, errors.New("new line error"))
	bytesWritten, err := p.Write([]byte("hello\nworld"))
	require.NotNil(t, err)
	require.Equal(t, "failed to write new line: new line error", err.Error())
	require.Equal(t, 5, bytesWritten)
}

func TestWriteNewLineInconsistencyErr(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m)
	m.On("Write", []byte("hello")).Return(5, nil)
	m.On("Write", []byte("\n")).Return(2, nil)
	bytesWritten, err := p.Write([]byte("hello\nworld"))
	require.NotNil(t, err)
	require.Equal(t, "inconsistency when writing new line: new line len = 1, actually written = 2", err.Error())
	require.Equal(t, 5, bytesWritten)
}

func TestWriteSuffixErr(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m).WithSuffixString("bad")
	m.On("Write", []byte("hello")).Return(5, nil)
	m.On("Write", []byte("bad")).Return(0, errors.New("suffix error"))
	bytesWritten, err := p.Write([]byte("hello\nworld"))
	require.NotNil(t, err)
	require.Equal(t, "failed to write suffix: suffix error", err.Error())
	require.Equal(t, 5, bytesWritten)
}

func TestWriteSuffixErrInconsistency(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m).WithSuffixString("bad")
	m.On("Write", []byte("hello")).Return(5, nil)
	m.On("Write", []byte("bad")).Return(2, nil)
	bytesWritten, err := p.Write([]byte("hello\nworld"))
	require.NotNil(t, err)
	require.Equal(t, "inconsistency when writing suffix: suffix len = 3, actually written = 2", err.Error())
	require.Equal(t, 5, bytesWritten)
}

func TestFlushBufferErr(t *testing.T) {
	m := new(mockWriter)
	p := shellprinter.New(m)
	bytesWritten, err := p.Write([]byte("hell"))
	require.Nil(t, err)
	require.Equal(t, 4, bytesWritten)
	m.On("Write", []byte("hell")).Return(4, errors.New("flush buffer error"))
	bytesWritten, err = p.Write([]byte("o\nworld"))
	require.NotNil(t, err)
	require.Equal(t, "failed to flush buffer: flush buffer error", err.Error())
	require.Equal(t, 0, bytesWritten)
}
