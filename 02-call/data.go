package main

import (
	"errors"

	"github.com/golang-estonia/structs-to-bytes/stream"
	"github.com/zeebo/errs"
)

var ErrTooMuchData = errors.New("too much data")

type Person struct {
	Name string
	Age  uint64
	Path Points
}

type Points []Point

type Point struct {
	X, Y uint64
}

func (p *Person) Encode() (data []byte, err error) {
	data = stream.AppendString(data, p.Name)
	data = stream.AppendUint64(data, p.Age)

	sub, err := p.Path.Encode()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	data = stream.AppendBytes(data, sub)

	return data, nil
}

func (p *Person) Decode(data []byte) (err error) {
	p.Name, data, err = stream.ReadString(data)
	if err != nil {
		return errs.Wrap(err)
	}

	p.Age, data, err = stream.ReadUint64(data)
	if err != nil {
		return errs.Wrap(err)
	}

	var sub []byte
	sub, data, err = stream.ReadBytes(data)
	if err != nil {
		return errs.Wrap(err)
	}
	err = p.Path.Decode(sub)
	if err != nil {
		return errs.Wrap(err)
	}

	if len(data) != 0 {
		return errs.Wrap(ErrTooMuchData)
	}

	return nil
}

func (points Points) Encode() (data []byte, err error) {
	data = stream.AppendUint64(data, uint64(len(points)))
	for i := range points {
		sub, err := points[i].Encode()
		if err != nil {
			return nil, errs.Wrap(err)
		}
		data = stream.AppendBytes(data, sub)
	}
	return data, nil
}

func (points *Points) Decode(data []byte) (err error) {
	var n uint64
	n, data, err = stream.ReadUint64(data)
	if err != nil {
		return errs.Wrap(err)
	}

	*points = make(Points, int(n))
	for i := range *points {
		var sub []byte
		sub, data, err = stream.ReadBytes(data)
		if err != nil {
			return errs.Wrap(err)
		}
		err = (*points)[i].Decode(sub)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	if len(data) != 0 {
		return errs.Wrap(ErrTooMuchData)
	}

	return nil
}

func (p *Point) Encode() (data []byte, err error) {
	data = stream.AppendUint64(data, p.X)
	data = stream.AppendUint64(data, p.Y)
	return data, nil
}

func (p *Point) Decode(data []byte) (err error) {
	p.X, data, err = stream.ReadUint64(data)
	if err != nil {
		return errs.Wrap(err)
	}
	p.Y, data, err = stream.ReadUint64(data)
	if err != nil {
		return errs.Wrap(err)
	}
	if len(data) != 0 {
		return errs.Wrap(ErrTooMuchData)
	}
	return nil
}
