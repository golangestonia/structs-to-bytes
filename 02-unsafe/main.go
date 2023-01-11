package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/exp/slices"
)

type Point struct {
	X, Y int32
}

func main() {
	points := make([]Point, 8)
	for i := range points {
		points[i].X = int32(i)
		points[i].Y = int32(i)
	}

	slicedata := unsafe.Pointer(&points[0])
	// With Go 1.20 use:
	//   slicedata := unsafe.SliceData(points)
	data := unsafe.Slice((*byte)(slicedata), uintptr(len(points))*unsafe.Sizeof(Point{}))
	fmt.Println(data)

	{ // Converting back
		input := slices.Clone(data)

		slicedata := unsafe.Pointer(&input[0])
		// With Go 1.20 use:
		//   slicedata := unsafe.SliceData(input)
		data := unsafe.Slice((*Point)(slicedata), uintptr(len(input))/unsafe.Sizeof(Point{}))
		fmt.Println(data)
	}
}
