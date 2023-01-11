package main

import "github.com/golang-estonia/structs-to-bytes/est"

type Person struct {
	Name string
	Age  uint64

	Path Points
}

type Points []Point

type Point struct {
	X, Y uint64
}

func (p *Person) EncodeEst(stream *est.Stream) (err error) {
	if err = stream.WriteString(p.Name); err != nil {
		return err
	}
	if err = stream.WriteUint64(p.Age); err != nil {
		return err
	}
	if err = stream.WriteMessage(p.Path.EncodeEst); err != nil {
		return err
	}
	return nil
}

func (p *Person) DecodeEst(stream *est.Stream) (err error) {
	if p.Name, err = stream.ReadString(); err != nil {
		return err
	}
	if p.Age, err = stream.ReadUint64(); err != nil {
		return err
	}
	if err = stream.ReadMessage(p.Path.DecodeEst); err != nil {
		return err
	}
	return nil
}

func (points Points) EncodeEst(stream *est.Stream) (err error) {
	if err = stream.WriteUint64(uint64(len(points))); err != nil {
		return err
	}
	for i := range points {
		if err = stream.WriteMessage(points[i].EncodeEst); err != nil {
			return err
		}
	}

	return nil
}

func (points *Points) DecodeEst(stream *est.Stream) (err error) {
	n, err := stream.ReadUint64()
	if err != nil {
		return err
	}
	// TODO: verify that `n` size is reasonable

	*points = make(Points, int(n))
	for i := range *points {
		if err = stream.ReadMessage((*points)[i].DecodeEst); err != nil {
			return err
		}
	}

	return nil
}

func (p *Point) EncodeEst(stream *est.Stream) (err error) {
	if err = stream.WriteUint64(p.X); err != nil {
		return err
	}
	if err = stream.WriteUint64(p.Y); err != nil {
		return err
	}
	return nil
}

func (p *Point) DecodeEst(stream *est.Stream) (err error) {
	if p.X, err = stream.ReadUint64(); err != nil {
		return err
	}
	if p.Y, err = stream.ReadUint64(); err != nil {
		return err
	}
	return nil
}
