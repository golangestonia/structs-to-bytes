package main

import (
	"encoding/binary"
	"fmt"
)

func Example() {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, 123456789)
	fmt.Println(data)
	value := binary.LittleEndian.Uint64(data)
	fmt.Println(value)

	// Output:
	// [21 205 91 7 0 0 0 0]
	// 123456789
}
