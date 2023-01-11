package main

import (
	"reflect"
	"unsafe"
)

// Ref represents a relative reference in a buffer.
type Ref int32

func (r *Ref) Lock()  {}
func (r *Ref) Clear() { *r = 0 }

func (r *Ref) Set(v unsafe.Pointer) {
	*r = Ref(uintptr(unsafe.Pointer(v)) - uintptr(unsafe.Pointer(r)))
}
func (r *Ref) Get() unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(r)) + uintptr(*r))
}

type Node struct {
	Value  uint32
	_Left  Ref
	_Right Ref
}

func (n *Node) Left() *Node  { return (*Node)(n._Left.Get()) }
func (n *Node) Right() *Node { return (*Node)(n._Right.Get()) }

type Stream struct {
	Data [1 << 10]byte
	Head int
}

func (s *Stream) Write(v interface{}) unsafe.Pointer {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)
	size := int(rt.Elem().Size())

	base := unsafe.Pointer(rv.Pointer())
	data := (*[1 << 10]byte)(base)[:size]

	p := unsafe.Pointer(&s.Data[s.Head])
	s.Head += copy(s.Data[s.Head:], data)
	return p
}
