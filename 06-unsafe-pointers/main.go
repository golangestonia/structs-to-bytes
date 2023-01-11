package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

// The implementation is based on https://www.youtube.com/watch?v=qZM2B4D7hZs
//
// This merely a rough demonstration how to handle pointers together with
// unsafe casting. You probably would write a nicer wrapper around all of this.

func serialize() []byte {
	var stream Stream
	var emptyRef Ref
	index := stream.Write(&emptyRef)
	left := stream.Write(&Node{4129422, 0, 0})
	right := stream.Write(&Node{1238471, 0, 0})
	root := stream.Write(&Node{329501, 0, 0})

	(*Node)(root)._Left.Set(left)
	(*Node)(root)._Right.Set(right)
	(*Ref)(index).Set(root)

	return stream.Data[:stream.Head]
}

func main() {
	data := serialize()
	defer runtime.KeepAlive(data)

	index := (*Ref)(unsafe.Pointer(&data[0]))
	root := (*Node)(index.Get())
	left := root.Left()
	right := root.Right()

	fmt.Println(root)
	fmt.Println(left)
	fmt.Println(right)
}
