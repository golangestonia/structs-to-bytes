package est

import (
	"encoding/binary"
	"errors"

	"github.com/zeebo/errs"
)

var ErrNotEnoughData = errors.New("not enough data")

func AppendUint64(xs []byte, v uint64) []byte {
	return binary.LittleEndian.AppendUint64(xs, v)
}

func AppendBytes(xs []byte, data []byte) []byte {
	xs = AppendUint64(xs, uint64(len(data)))
	xs = append(xs, data...)
	return xs
}

func AppendString(xs []byte, s string) []byte {
	return AppendBytes(xs, []byte(s))
}

func ReadUint64(xs []byte) (v uint64, rest []byte, err error) {
	if len(xs) < 8 {
		return 0, xs, errs.Wrap(ErrNotEnoughData)
	}
	v = binary.LittleEndian.Uint64(xs)
	return v, xs[8:], nil
}

func ReadBytes(xs []byte) (v []byte, rest []byte, err error) {
	n, rest, err := ReadUint64(xs)
	if err != nil {
		return nil, nil, errs.Wrap(ErrNotEnoughData)
	}
	if uint64(len(rest)) < n {
		return nil, nil, errs.Wrap(ErrNotEnoughData)
	}
	return rest[:n], rest[n:], nil
}

func ReadString(xs []byte) (v string, rest []byte, err error) {
	data, rest, err := ReadBytes(xs)
	return string(data), rest, errs.Wrap(err)
}

func AppendMessage(xs []byte, encode func() ([]byte, error)) ([]byte, error) {
	data, err := encode()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return AppendBytes(xs, data), nil
}

func ReadMessage(xs []byte, decode func([]byte) error) (rest []byte, err error) {
	var sub []byte

	sub, rest, err = ReadBytes(xs)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	err = decode(sub)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return rest, nil
}
