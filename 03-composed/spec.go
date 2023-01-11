package main

import (
	"github.com/golang-estonia/structs-to-bytes/est"
	"github.com/zeebo/errs"
)

type Spec interface {
	EncodeEst(*est.Stream) error
	DecodeEst(*est.Stream) error
}

type Uint64 struct {
	V *uint64
}

func (v Uint64) EncodeEst(s *est.Stream) (err error) {
	return errs.Wrap(s.WriteUint64(*v.V))
}

func (v Uint64) DecodeEst(s *est.Stream) (err error) {
	*v.V, err = s.ReadUint64()
	return errs.Wrap(err)
}

type String struct {
	V *string
}

func (v String) EncodeEst(s *est.Stream) (err error) {
	return errs.Wrap(s.WriteString(*v.V))
}

func (v String) DecodeEst(s *est.Stream) (err error) {
	*v.V, err = s.ReadString()
	return errs.Wrap(err)
}

type Ordered []Spec

func (v Ordered) EncodeEst(s *est.Stream) (err error) {
	for _, m := range v {
		err = m.EncodeEst(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func (v Ordered) DecodeEst(s *est.Stream) (err error) {
	for _, m := range v {
		err = m.DecodeEst(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

type Message struct {
	V Spec
}

func (v Message) EncodeEst(s *est.Stream) (err error) {
	return errs.Wrap(s.WriteMessage(func(s *est.Stream) error {
		return errs.Wrap(v.V.EncodeEst(s))
	}))
}

func (v Message) DecodeEst(s *est.Stream) (err error) {
	return errs.Wrap(s.ReadMessage(func(s *est.Stream) error {
		return errs.Wrap(v.V.DecodeEst(s))
	}))
}

type Slice struct {
	Count    func() int
	SetCount func(v int)
	Elem     func(i int) Spec
}

func (v Slice) EncodeEst(s *est.Stream) (err error) {
	n := v.Count()
	err = s.WriteUint64(uint64(n))
	if err != nil {
		return errs.Wrap(err)
	}

	for i := 0; i < n; i++ {
		err = Message{v.Elem(i)}.EncodeEst(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func (v Slice) DecodeEst(s *est.Stream) (err error) {
	var n uint64
	n, err = s.ReadUint64()
	if err != nil {
		return errs.Wrap(err)
	}
	v.SetCount(int(n))

	for i := 0; i < int(n); i++ {
		err = Message{v.Elem(i)}.DecodeEst(s)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}
