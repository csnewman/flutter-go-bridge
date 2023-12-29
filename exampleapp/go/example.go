//go:generate go run github.com/csnewman/flutter-go-bridge/cmd/flutter-go-bridge generate --src example.go --go bridge/bridge.gen.go --dart ../lib/bridge.gen.dart
package example

import (
	"fmt"
)

func Add(a int, b int) int {
	return a + b
}

type Point struct {
	X    int
	Y    int
	Name string
}

func AddPoints(a Point, b Point) Point {
	return Point{
		X:    a.X + b.X,
		Y:    a.Y + b.Y,
		Name: a.Name + "+" + b.Name,
	}
}

func AddError(a int, b int) (int, error) {
	return 0, fmt.Errorf("add res was %v", a+b)
}

type Obj struct {
	Name  string
	other int
}

func NewObj(name string, other int) *Obj {
	return &Obj{
		Name:  name,
		other: other,
	}
}

func ModifyObj(o *Obj) {
	o.other *= 2
}

func FormatObj(o *Obj) string {
	return fmt.Sprintf("Obj: Name=%v Other=%v", o.Name, o.other)
}
