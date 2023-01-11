package main

import (
	"fmt"

	"github.com/golang-estonia/structs-to-bytes/stream"
)

func main() {
	person := Person{
		Name: "John",
		Age:  32,
		Path: Points{
			{X: 10, Y: 10},
			{X: 20, Y: 20},
			{X: 30, Y: 30},
		},
	}

	var out stream.Stream
	err := person.Encode(&out)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	data := out.Bytes()
	fmt.Println(data)

	{
		var p Person
		err := p.Decode(stream.FromBytes(data))
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(p)
	}
}
