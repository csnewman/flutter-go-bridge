package generator

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"flutter-go-bridge/parser"
	"github.com/iancoleman/strcase"
)

var (
	ErrUnexpected  = errors.New("parser contained unexpected data")
	ErrUnsupported = errors.New("unsupported")
)

type Unit struct {
	TgtPkg    string
	Functions []*Func
}

type Func struct {
	SnakeName   string
	CamelName   string
	PascalName  string
	TgtName     string
	Params      []*Param
	HasRes      bool
	ResCType    string
	ResGoType   string
	ResGoMode   string
	ResDartType string
	HasErr      bool
}

type Param struct {
	Name      string
	CType     string
	GoType    string
	GoMode    string
	DartCType string
	DartType  string
}

func Generate(goDest string, dartDest string, in *parser.Package) error {
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

	// Generate dart
	f, err = os.Create(dartDest)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := GetDartBridgeTemplate().Execute(f, g.unit); err != nil {
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
		SnakeName:  strcase.ToSnake(f.Name),
		PascalName: strcase.ToCamel(f.Name),
		CamelName:  strcase.ToLowerCamel(f.Name),
		TgtName:    f.Name,
	}

	if len(f.Sig.Results) > 2 {
		return fmt.Errorf("%w: %v contains more than 2 results", ErrUnsupported, len(f.Sig.Results))
	}

	if len(f.Sig.Results) > 1 {
		res := f.Sig.Results[1]
		fu.HasErr = true

		if !parser.IsErrorType(res.Type) {
			return fmt.Errorf("%w: second result must be error %v", ErrUnsupported, res.Type)
		}
	}

	if len(f.Sig.Results) > 0 {
		res := f.Sig.Results[0]

		if parser.IsErrorType(res.Type) {
			if fu.HasErr {
				return fmt.Errorf("%w: multiple error types", ErrUnsupported)
			}

			fu.HasErr = true
		} else {
			fu.HasRes = true

			switch t := res.Type.(type) {
			case *parser.IdentType:
				fu.ResCType = t.Name
				fu.ResGoType = t.Name
				fu.ResGoMode = "cast"
				fu.ResDartType = t.Name
			default:
				return fmt.Errorf("%w: unexpected type %v", ErrUnexpected, reflect.TypeOf(res.Type))
			}
		}

	}

	for _, param := range f.Sig.Params {
		p := &Param{
			Name: param.Name,
		}

		switch t := param.Type.(type) {
		case *parser.IdentType:
			p.CType = t.Name
			p.GoType = t.Name
			p.GoMode = "cast"
			p.DartCType = "ffi.Int32"
			p.DartType = t.Name
		default:
			return fmt.Errorf("%w: unexpected type %v", ErrUnexpected, reflect.TypeOf(param.Type))
		}

		fu.Params = append(fu.Params, p)
	}

	g.unit.Functions = append(g.unit.Functions, fu)

	return nil
}
