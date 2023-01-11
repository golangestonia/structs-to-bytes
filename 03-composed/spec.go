package main

import (
	"github.com/golang-estonia/structs-to-bytes/stream"
	"github.com/zeebo/errs"
)

type Specer interface {
	Spec() Spec
}

type Spec interface {
	Encode(*stream.Stream) error
	Decode(*stream.Stream) error
}

type Uint64 struct {
	V *uint64
}

func (v Uint64) Encode(s *stream.Stream) (err error) {
	return errs.Wrap(s.WriteUint64(*v.V))
}

func (v Uint64) Decode(s *stream.Stream) (err error) {
	*v.V, err = s.ReadUint64()
	return errs.Wrap(err)
}

type String struct {
	V *string
}

func (v String) Encode(s *stream.Stream) (err error) {
	return errs.Wrap(s.WriteString(*v.V))
}

func (v String) Decode(s *stream.Stream) (err error) {
	*v.V, err = s.ReadString()
	return errs.Wrap(err)
}

type Ordered []Spec

func (v Ordered) Encode(s *stream.Stream) (err error) {
	for _, m := range v {
		err = m.Encode(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func (v Ordered) Decode(s *stream.Stream) (err error) {
	for _, m := range v {
		err = m.Decode(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

type Message struct {
	V Spec
}

func (v Message) Encode(s *stream.Stream) (err error) {
	return errs.Wrap(s.WriteMessage(func(s *stream.Stream) error {
		return errs.Wrap(v.V.Encode(s))
	}))
}

func (v Message) Decode(s *stream.Stream) (err error) {
	return errs.Wrap(s.ReadMessage(func(s *stream.Stream) error {
		return errs.Wrap(v.V.Decode(s))
	}))
}

type Slice struct {
	Count    func() int
	SetCount func(v int)
	Elem     func(i int) Spec
}

func (v Slice) Encode(s *stream.Stream) (err error) {
	n := v.Count()
	err = s.WriteUint64(uint64(n))
	if err != nil {
		return errs.Wrap(err)
	}

	for i := 0; i < n; i++ {
		err = Message{v.Elem(i)}.Encode(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func (v Slice) Decode(s *stream.Stream) (err error) {
	var n uint64
	n, err = s.ReadUint64()
	if err != nil {
		return errs.Wrap(err)
	}
	v.SetCount(int(n))

	for i := 0; i < int(n); i++ {
		err = Message{v.Elem(i)}.Decode(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}
