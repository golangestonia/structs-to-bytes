package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  uint64
	Path Points
}

type Points []Point

type Point struct {
	X, Y uint64
}

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

	data, err := Encode(&person)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	fmt.Println(data)

	{
		var p Person
		err := Decode(data, &p)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(p)
	}
}
