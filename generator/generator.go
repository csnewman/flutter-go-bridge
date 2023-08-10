package generator

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"flutter-go-bridge/parser"
	"github.com/iancoleman/strcase"
)

var ErrUnexpected = errors.New("parser contained unexpected data")

type Unit struct {
	TgtPkg    string
	Functions []*Func
}

type Func struct {
	Name    string
	TgtName string
	Params  []*Param
}

type Param struct {
	Name   string
	CType  string
	GoType string
	GoMode string
}

func Generate(goDest string, in *parser.Package) error {
	g := &generator{
		unit: &Unit{
			TgtPkg: in.PkgPath,
		},
	}

	if err := g.process(in); err != nil {
		return err
	}

	// Generate go
	f, err := os.Create(goDest)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := GetGoBridgeTemplate().Execute(f, g.unit); err != nil {
		return err
	}

	return nil
}

type generator struct {
	unit *Unit
}

func (g *generator) process(p *parser.Package) error {
	for _, fName := range p.FuncOrder {
		if err := g.processFunc(p.Funcs[fName]); err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) processFunc(f *parser.FuncDef) error {
	fu := &Func{
		Name:    strcase.ToSnake(f.Name),
		TgtName: f.Name,
	}

	for _, param := range f.Sig.Params {
		p := &Param{
			Name: param.Name,
		}

		switch t := param.Type.(type) {
		case *parser.IdentType:
			p.CType = "C." + t.Name
			p.GoType = t.Name
			p.GoMode = "cast"
		default:
			return fmt.Errorf("%w: unexpected type %v", ErrUnexpected, reflect.TypeOf(param.Type))
		}

		fu.Params = append(fu.Params, p)
	}

	g.unit.Functions = append(g.unit.Functions, fu)

	return nil
}
