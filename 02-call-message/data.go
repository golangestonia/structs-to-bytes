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
	data, err = est.AppendMessage(data, p.Path.EncodeEst)
	if err != nil {
		return nil, errs.Wrap(err)
	}

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

	data, err = est.ReadMessage(data, p.Path.DecodeEst)
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
		data, err = est.AppendMessage(data, points[i].EncodeEst)
		if err != nil {
			return nil, errs.Wrap(err)
		}
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
		data, err = est.ReadMessage(data, (*points)[i].DecodeEst)
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
