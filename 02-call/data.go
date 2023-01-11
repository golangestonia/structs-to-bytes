package main

import (
	"errors"

	"github.com/golang-estonia/structs-to-bytes/est"
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

func (p *Person) EncodeEst() (data []byte, err error) {
	data = est.AppendString(data, p.Name)
	data = est.AppendUint64(data, p.Age)

	sub, err := p.Path.EncodeEst()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	data = est.AppendBytes(data, sub)

	return data, nil
}

func (p *Person) DecodeEst(data []byte) (err error) {
	p.Name, data, err = est.ReadString(data)
	if err != nil {
		return errs.Wrap(err)
	}

	p.Age, data, err = est.ReadUint64(data)
	if err != nil {
		return errs.Wrap(err)
	}

	var sub []byte
	sub, data, err = est.ReadBytes(data)
	if err != nil {
		return errs.Wrap(err)
	}
	err = p.Path.DecodeEst(sub)
	if err != nil {
		return errs.Wrap(err)
	}

	if len(data) != 0 {
		return errs.Wrap(ErrTooMuchData)
	}

	return nil
}

func (points Points) EncodeEst() (data []byte, err error) {
	data = est.AppendUint64(data, uint64(len(points)))
	for i := range points {
		sub, err := points[i].EncodeEst()
		if err != nil {
			return nil, errs.Wrap(err)
		}
		data = est.AppendBytes(data, sub)
	}
	return data, nil
}

func (points *Points) DecodeEst(data []byte) (err error) {
	var n uint64
	n, data, err = est.ReadUint64(data)
	if err != nil {
		return errs.Wrap(err)
	}

	*points = make(Points, int(n))
	for i := range *points {
		var sub []byte
		sub, data, err = est.ReadBytes(data)
		if err != nil {
			return errs.Wrap(err)
		}
		err = (*points)[i].DecodeEst(sub)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	if len(data) != 0 {
		return errs.Wrap(ErrTooMuchData)
	}

	return nil
}

func (p *Point) EncodeEst() (data []byte, err error) {
	data = est.AppendUint64(data, p.X)
	data = est.AppendUint64(data, p.Y)
	return data, nil
}

func (p *Point) DecodeEst(data []byte) (err error) {
	p.X, data, err = est.ReadUint64(data)
	if err != nil {
		return errs.Wrap(err)
	}
	p.Y, data, err = est.ReadUint64(data)
	if err != nil {
		return errs.Wrap(err)
	}
	if len(data) != 0 {
		return errs.Wrap(ErrTooMuchData)
	}
	return nil
}
