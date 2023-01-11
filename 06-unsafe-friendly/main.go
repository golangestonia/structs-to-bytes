package main

import "fmt"

func main() {
	stream := NewStream(20)

	{ // Let's write some data
		it := stream.Iterator()
		*it.MoveTo() = MoveTo{
			X: 10,
			Y: 10,
		}
		*it.LineTo() = LineTo{
			X: 20,
			Y: 20,
		}
		*it.Done() = Done{}

		fmt.Println(stream.Data)
	}

	{ // Let's read the data
		it := stream.Iterator()
		for it.Op() != OpDone {
			switch it.Op() {
			case OpMoveTo:
				op := it.MoveTo()
				fmt.Println("move to", *op)
			case OpLineTo:
				op := it.LineTo()
				fmt.Println("line to", *op)
			default:
				panic("unsupported op")
			}
		}
	}
}
