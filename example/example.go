//go:generate go run github.com/csnewman/flutter-go-bridge/cmd/flutter-go-bridge generate --src example.go --go bridge/bridge.gen.go --dart ../exampleapp/lib/bridge.gen.dart
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

	// return 123, fmt.Errorf("Hello")
	return 123, nil
}
