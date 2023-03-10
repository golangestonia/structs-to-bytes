package main

// See More Details in https://egonelbre.com/composed-serialization/

import (
	"fmt"

	"github.com/golang-estonia/structs-to-bytes/est"
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

	var s est.Stream
	err := person.Est().EncodeEst(&s)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	data := s.Bytes()
	fmt.Println(data)

	{
		var p Person
		err := p.Est().DecodeEst(est.StreamFromBytes(data))
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(p)
	}
}
