package main

import "unsafe"

// Stream contains the buffer for writing data.
type Stream struct {
	Data []byte
}

func NewStream(max int) Stream {
	return Stream{
		Data: make([]byte, max, max),
	}
}

func (stream *Stream) Iterator() *Iterator {
	return &Iterator{
		Head: unsafe.Pointer(&stream.Data[0]), // use `unsafe.SliceData`
	}
}

// Iterator that can be used to write or read the data.
type Iterator struct {
	Head unsafe.Pointer
	// We have omitted all the safeguards.
}

// Data Types that we wish to encode.
type MoveTo struct{ X, Y int32 }
type LineTo struct{ X, Y int32 }
type Done struct{}

// Op is the byte used to indicate a particular operator.
type Op byte

const (
	OpDone   = Op(0)
	OpMoveTo = Op(1)
	OpLineTo = Op(2)
)

func (it *Iterator) Op() Op { return *(*Op)(it.Head) }

// We return a pointer of the struct inside the original buffer.
func (it *Iterator) MoveTo() (r *MoveTo) {
	// Currently this mixes writing and reading, so wrong usage will corrupt your data.
	*(*Op)(it.Head) = OpMoveTo
	// Grab the pointer to the struct.
	r = (*MoveTo)(unsafe.Add(it.Head, 1))
	// Advance the iterator.
	it.Head = unsafe.Add(it.Head, 1+unsafe.Sizeof(*r))
	return r
}

func (it *Iterator) LineTo() (r *LineTo) {
	*(*Op)(it.Head) = OpLineTo
	r = (*LineTo)(unsafe.Add(it.Head, 1))
	it.Head = unsafe.Add(it.Head, 1+unsafe.Sizeof(*r))
	return r
}

func (it *Iterator) Done() (r *Done) {
	*(*Op)(it.Head) = OpDone
	r = (*Done)(unsafe.Add(it.Head, 1))
	it.Head = unsafe.Add(it.Head, 1+unsafe.Sizeof(*r))
	return r
}
