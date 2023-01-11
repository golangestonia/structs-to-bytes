package main

import (
	"github.com/golang-estonia/structs-to-bytes/est"
	"github.com/zeebo/errs"
)


type Person struct {
	Name string
	Age uint64
	Path Points
}

func (m Person) EncodeEst(stream *est.Stream) (err error) {
	if err = stream.WriteString(m.Name); err != nil {
		return errs.Wrap(err)
	}
	
	if err = stream.WriteUint64(m.Age); err != nil {
		return errs.Wrap(err)
	}
	
	if err = stream.WriteMessage(m.Path.EncodeEst); err != nil {
		return errs.Wrap(err)
	}
	
	return nil
}

func (m *Person) DecodeEst(stream *est.Stream) (err error) {
	if m.Name, err = stream.ReadString(); err != nil {
		return errs.Wrap(err)
	}
	
	if m.Age, err = stream.ReadUint64(); err != nil {
		return errs.Wrap(err)
	}
	
	if err = stream.ReadMessage(m.Path.DecodeEst); err != nil {
		return errs.Wrap(err)
	}
	
	return nil
}

type Points []Point

func (m Points) EncodeEst(stream *est.Stream) (err error) {
	if err = stream.WriteUint64(uint64(len(m))); err != nil {
		return errs.Wrap(err)
	}
	for i := range m {
		if err = stream.WriteMessage(m[i].EncodeEst); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func (m *Points) DecodeEst(stream *est.Stream) (err error) {
	n, err := stream.ReadUint64()
	if err != nil {
		return errs.Wrap(err)
	}
	*m = make(Points, int(n))
	for i := range *m {
		if err = stream.ReadMessage((*m)[i].DecodeEst); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

type Point struct {
	X uint64
	Y uint64
}

func (m Point) EncodeEst(stream *est.Stream) (err error) {
	if err = stream.WriteUint64(m.X); err != nil {
		return errs.Wrap(err)
	}
	
	if err = stream.WriteUint64(m.Y); err != nil {
		return errs.Wrap(err)
	}
	
	return nil
}

func (m *Point) DecodeEst(stream *est.Stream) (err error) {
	if m.X, err = stream.ReadUint64(); err != nil {
		return errs.Wrap(err)
	}
	
	if m.Y, err = stream.ReadUint64(); err != nil {
		return errs.Wrap(err)
	}
	
	return nil
}
