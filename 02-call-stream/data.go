package main

import "github.com/golang-estonia/structs-to-bytes/stream"

type Person struct {
	Name string
	Age  uint64

	Path Points
}

type Points []Point

type Point struct {
	X, Y uint64
}

func (p *Person) Encode(stream *stream.Stream) (err error) {
	if err = stream.WriteString(p.Name); err != nil {
		return err
	}
	if err = stream.WriteUint64(p.Age); err != nil {
		return err
	}
	if err = stream.WriteMessage(p.Path.Encode); err != nil {
		return err
	}
	return nil
}

func (p *Person) Decode(stream *stream.Stream) (err error) {
	if p.Name, err = stream.ReadString(); err != nil {
		return err
	}
	if p.Age, err = stream.ReadUint64(); err != nil {
		return err
	}
	if err = stream.ReadMessage(p.Path.Decode); err != nil {
		return err
	}
	return nil
}

func (points Points) Encode(stream *stream.Stream) (err error) {
	if err = stream.WriteUint64(uint64(len(points))); err != nil {
		return err
	}
	for i := range points {
		if err = stream.WriteMessage(points[i].Encode); err != nil {
			return err
		}
	}

	return nil
}

func (points *Points) Decode(stream *stream.Stream) (err error) {
	n, err := stream.ReadUint64()
	if err != nil {
		return err
	}
	// TODO: verify that `n` size is reasonable

	*points = make(Points, int(n))
	for i := range *points {
		if err = stream.ReadMessage((*points)[i].Decode); err != nil {
			return err
		}
	}

	return nil
}

func (p *Point) Encode(stream *stream.Stream) (err error) {
	if err = stream.WriteUint64(p.X); err != nil {
		return err
	}
	if err = stream.WriteUint64(p.Y); err != nil {
		return err
	}
	return nil
}

func (p *Point) Decode(stream *stream.Stream) (err error) {
	if p.X, err = stream.ReadUint64(); err != nil {
		return err
	}
	if p.Y, err = stream.ReadUint64(); err != nil {
		return err
	}
	return nil
}
