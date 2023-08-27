package example

import (
	"log"
)

type Inner struct {
	A int
	D string
}

type SomeVal struct {
	V1 int
	V2 int
	I1 Inner
}

func Example(v SomeVal) error {
	log.Println("Call Example", v)
	log.Println("> ", v.I1.D)

	return nil
}

func Other() SomeVal {
	return SomeVal{}
}

func CallMe() (int, error) {
	log.Println("CallMe")

	// return 123, fmt.Errorf("Hell")
	return 123, nil
}
