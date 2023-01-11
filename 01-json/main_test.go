package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type Person struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address,omitempty"`
	Password string `json:"-"`
}

func ExampleEncode() {
	data, err := json.Marshal(Person{
		Name:     "John",
		Email:    "john@email.test",
		Password: "hunter2",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// Output: {"name":"John","email":"john@email.test"}
}

func ExampleDecode() {
	var person Person
	data := []byte(`{"name":"John","email":"john@email.test"}`)
	err := json.Unmarshal(data, &person)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", person)

	// Output: main.Person{Name:"John", Email:"john@email.test", Address:"", Password:""}
}

func ExampleInlineDecode() {
	var person struct {
		Name  string
		Email string
	}
	data := []byte(`{"name":"John","email":"john@email.test"}`)
	err := json.Unmarshal(data, &person)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", person)

	// Output: {Name:John Email:john@email.test}
}

func ExampleEncoder() {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	_ = enc.Encode(Person{Name: "Alice"})
	_ = enc.Encode(Person{Name: "Bob"})
	_ = enc.Encode(Person{Name: "Charlie"})
	fmt.Println(buf.String())

	// Output:
	// {"name":"Alice","email":""}
	// {"name":"Bob","email":""}
	// {"name":"Charlie","email":""}
}

// START ENCODINGDEFINITIONS OMIT
// Default text encoding for different encoding packages.
type TextMarshaler interface {
	MarshalText() (text []byte, err error)
}
type TextUnmarshaler interface {
	UnmarshalText(text []byte) error
}

// Default binary encoding for different encoding packages.
type BinaryMarshaler interface {
	MarshalBinary() (data []byte, err error)
}
type BinaryUnmarshaler interface {
	UnmarshalBinary(data []byte) error
}

// END ENCODINGDEFINITIONS OMIT

// START POINT OMIT
type Point struct {
	X, Y int32
}

// Custom encoding for `encoding/json`, `encoding/xml` and `encoding/gob`.
func (p Point) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("x%d y%d", p.X, p.Y)), nil
}

// Custom encoding for `encoding/gob`.
func (p Point) MarshalBinary() ([]byte, error) {
	var data [8]byte
	binary.BigEndian.PutUint32(data[0:4], uint32(p.X))
	binary.BigEndian.PutUint32(data[4:8], uint32(p.Y))
	return data[:], nil
}

// END POINT OMIT

func (p *Point) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("expected length 8, but got %d", len(data))
	}
	data = data[0:8]
	p.X = int32(binary.BigEndian.Uint32(data[0:4]))
	p.Y = int32(binary.BigEndian.Uint32(data[4:8]))
	return nil
}

func ExamplePoint() {
	type Drone struct {
		Name     string
		Location Point
	}

	drone := Drone{
		Name:     "Johnny",
		Location: Point{X: 100, Y: 100},
	}

	data, err := json.Marshal(drone)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	// Output: {"Name":"Johnny","Location":"x100 y100"}
}

// START JSONDEFINITIONS OMIT
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

// END JSONDEFINITIONS OMIT

// START QUOTEDPERSON OMIT
type QuotedPerson struct {
	Name string
}

func (p QuotedPerson) MarshalJSON() ([]byte, error) {
	qname, err := json.Marshal("!" + p.Name + "!")
	if err != nil {
		return nil, err
	}

	return []byte(`{"quoted": ` + string(qname) + `}`), nil
}

func ExampleQuotedPerson() {
	data, err := json.Marshal(QuotedPerson{Name: "Hello"})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"quoted":"!Hello!"}
}

// END QUOTEDPERSON OMIT

// START TEMPPERSON OMIT
type TempPerson struct {
	Name string
}

func (p TempPerson) MarshalJSON() ([]byte, error) {
	var temp struct {
		Quoted string `json:"quoted"`
	}
	temp.Quoted = "!" + p.Name + "!"
	return json.Marshal(temp)
}

func ExampleTempPerson() {
	data, err := json.Marshal(TempPerson{Name: "Hello"})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"quoted":"!Hello!"}
}

// END TEMPPERSON OMIT
