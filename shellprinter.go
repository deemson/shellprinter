package shellprinter

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"io"
)

type ShellPrinter struct {
	writer io.Writer
	// The buffer holds data for the incomplete line from previous writes
	buffer *bytes.Buffer
	prefix []byte
	suffix []byte
}

func New(writer io.Writer) *ShellPrinter {
	return &ShellPrinter{
		writer: writer,
		buffer: new(bytes.Buffer),
		prefix: nil,
		suffix: nil,
	}
}

func (p *ShellPrinter) WithPrefix(prefix []byte) *ShellPrinter {
	p.prefix = prefix
	return p
}

func (p *ShellPrinter) WithPrefixString(prefix string) *ShellPrinter {
	return p.WithPrefix([]byte(prefix))
}

func (p *ShellPrinter) WithSuffix(suffix []byte) *ShellPrinter {
	p.suffix = suffix
	return p
}

func (p *ShellPrinter) WithSuffixString(suffix string) *ShellPrinter {
	return p.WithSuffix([]byte(suffix))
}

func (p *ShellPrinter) Write(data []byte) (int, error) {
	bytesWrittenSoFar := 0
	dataRemainder := data
	for {
		// Ignoring the error here as it's always nil judging by the implementation
		afterNewLineIndex, dataBeforeNewLine, err := bufio.ScanLines(dataRemainder, false)
		if err != nil {
			return bytesWrittenSoFar, errors.Wrap(err, "failed to scan lines")
		}
		if afterNewLineIndex == 0 {
			// New line was not found
			break
		}
		err = p.writePrefix()
		if err != nil {
			return bytesWrittenSoFar, err
		}
		err = p.flushBuffer()
		if err != nil {
			return bytesWrittenSoFar, err
		}
		dataBytesWritten, err := p.writer.Write(dataBeforeNewLine)
		if err != nil {
			return bytesWrittenSoFar, errors.Wrap(err, "failed to write data")
		}
		if dataBytesWritten != len(dataBeforeNewLine) {
			return bytesWrittenSoFar, errors.Errorf(
				`inconsistency when writing data: data len = %d, actually written = %d`,
				len(dataBeforeNewLine),
				dataBytesWritten,
			)
		}
		bytesWrittenSoFar += dataBytesWritten
		err = p.writeSuffix()
		if err != nil {
			return bytesWrittenSoFar, err
		}
		newLineData := dataRemainder[len(dataBeforeNewLine):afterNewLineIndex]
		newLineBytesWritten, err := p.writer.Write(newLineData)
		if err != nil {
			return bytesWrittenSoFar, errors.Wrap(err, "failed to write new line")
		}
		if newLineBytesWritten != len(newLineData) {
			return bytesWrittenSoFar, errors.Errorf(
				`inconsistency when writing new line: new line len = %d, actually written = %d`,
				len(newLineData),
				newLineBytesWritten,
			)
		}
		bytesWrittenSoFar += newLineBytesWritten
		// Advance data by cutting everything before the new line
		dataRemainder = dataRemainder[afterNewLineIndex:]
	}
	if len(dataRemainder) > 0 {
		// Being here means there was no more new lines detected and there's data left
		// Put remaining data into buffer and wait the new calls to Write to supply new lines
		bytesWrittenToBuffer, err := p.buffer.Write(dataRemainder)
		if err != nil {
			return bytesWrittenSoFar, errors.Wrap(err, "failed to write to buffer")
		}
		if len(dataRemainder) != bytesWrittenToBuffer {
			return bytesWrittenSoFar, errors.Errorf(
				`inconsistency when writing to buffer: data len = %d, actually written = %d`,
				len(dataRemainder),
				bytesWrittenToBuffer,
			)
		}
		bytesWrittenSoFar += len(dataRemainder)
	}
	return bytesWrittenSoFar, nil
}

func (p *ShellPrinter) flushBuffer() error {
	if p.buffer.Len() > 0 {
		expectedBytesWritten := int64(p.buffer.Len())
		actualBytesWritten, err := p.buffer.WriteTo(p.writer)
		if err != nil {
			return errors.Wrap(err, "failed to flush buffer")
		}
		if expectedBytesWritten != actualBytesWritten {
			return errors.Errorf(
				`inconsistency when flushing buffer: buffer len = %d, actually written = %d`,
				expectedBytesWritten,
				actualBytesWritten,
			)
		}
		p.buffer.Reset()
	}
	return nil
}

func (p *ShellPrinter) writePrefix() error {
	if p.prefix != nil && len(p.prefix) > 0 {
		expectedBytesWritten := len(p.prefix)
		actualBytesWritten, err := p.writer.Write(p.prefix)
		if err != nil {
			return errors.Wrap(err, "failed to write prefix")
		}
		if expectedBytesWritten != actualBytesWritten {
			return errors.Errorf(
				`inconsistency when writing prefix: prefix len = %d, actually written = %d`,
				expectedBytesWritten,
				actualBytesWritten,
			)
		}
	}
	return nil
}

func (p *ShellPrinter) writeSuffix() error {
	if p.suffix != nil && len(p.suffix) > 0 {
		expectedBytesWritten := len(p.suffix)
		actualBytesWritten, err := p.writer.Write(p.suffix)
		if err != nil {
			return errors.Wrap(err, "failed to write suffix")
		}
		if expectedBytesWritten != actualBytesWritten {
			return errors.Errorf(
				`inconsistency when writing suffix: suffix len = %d, actually written = %d`,
				expectedBytesWritten,
				actualBytesWritten,
			)
		}
	}
	return nil
}

func (p *ShellPrinter) Close() error {
	if p.buffer.Len() > 0 {
		err := p.writePrefix()
		if err != nil {
			return err
		}
		return p.flushBuffer()
	}
	return nil
}
