package est

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/zeebo/errs"
)

const maxBytesLength = 1024

// Stream implements writing and reading values for some basic types.
type Stream struct {
	data bytes.Buffer
}

func StreamFromBytes(data []byte) *Stream {
	return &Stream{data: *bytes.NewBuffer(data)}
}
func (stream *Stream) Bytes() []byte { return stream.data.Bytes() }

func (stream *Stream) WriteUint64(v uint64) error {
	var data [8]byte
	binary.LittleEndian.PutUint64(data[:], v)
	_, err := stream.data.Write(data[:])
	if err != nil {
		return errs.New("failed to write uint64: %w", err)
	}
	return nil
}

func (stream *Stream) WriteBytes(data []byte) error {
	if len(data) > maxBytesLength {
		return errs.New("too many bytes (%v)", len(data))
	}
	err := stream.WriteUint64(uint64(len(data)))
	if err != nil {
		return errs.New("failed to bytes length: %w", err)
	}
	_, err = stream.data.Write(data)
	if err != nil {
		return errs.New("failed to bytes data: %w", err)
	}
	return nil
}

func (stream *Stream) WriteString(data string) error {
	err := stream.WriteBytes([]byte(data))
	if err != nil {
		return errs.New("failed to write string: %w", err)
	}
	return nil
}

func (stream *Stream) ReadUint64() (uint64, error) {
	var data [8]byte
	_, err := io.ReadFull(&stream.data, data[:])
	if err != nil {
		return 0, errs.New("failed to read uint64: %w", err)
	}
	v := binary.LittleEndian.Uint64(data[:])
	return v, nil
}

func (stream *Stream) ReadBytes() ([]byte, error) {
	size, err := stream.ReadUint64()
	if err != nil {
		return nil, errs.New("failed to read bytes length: %w", err)
	}
	// it's useful to check for reasonable values to avoid potential OOM issues.
	if size > maxBytesLength {
		return nil, errs.New("too many bytes %v", size)
	}

	data := make([]byte, size)
	_, err = io.ReadFull(&stream.data, data[:])
	if err != nil {
		return nil, errs.New("failed to read bytes data: %w", err)
	}

	return data, nil
}

func (stream *Stream) ReadString() (string, error) {
	data, err := stream.ReadBytes()
	if err != nil {
		return "", errs.New("failed to read string: %w", err)
	}
	return string(data), nil
}

func (stream *Stream) WriteMessage(fn func(*Stream) error) error {
	var sub Stream
	err := fn(&sub)
	if err != nil {
		return errs.Wrap(err)
	}

	err = stream.WriteBytes(sub.Bytes())
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func (stream *Stream) ReadMessage(fn func(*Stream) error) error {
	data, err := stream.ReadBytes()
	if err != nil {
		return errs.Wrap(err)
	}

	err = fn(StreamFromBytes(data))
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}
