package main

type Person struct {
	Name string
	Age  uint64
	Path Points
}

type Points []Point

type Point struct {
	X, Y uint64
}

func (p *Person) Est() Spec {
	return Ordered{
		String{&p.Name},
		Uint64{&p.Age},
		Message{p.Path.Est()},
	}
}

func (p *Person) EstName() Spec {
	return Ordered{
		String{&p.Name},
	}
}

func (p *Points) Est() Spec {
	return Slice{
		Count:    func() int { return len(*p) },
		SetCount: func(n int) { *p = make(Points, n) },
		Elem: func(i int) Spec {
			return (*p)[i].Est()
		},
	}
}

func (p *Point) Est() Spec {
	return Ordered{
		Uint64{&p.X},
		Uint64{&p.Y},
	}
}
