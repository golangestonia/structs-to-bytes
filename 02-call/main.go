package main

import "fmt"

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

	data, err := person.Encode()
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	fmt.Println(data)

	{
		var p Person
		err := p.Decode(data)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(p)
	}
}
