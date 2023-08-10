package parser

import (
	"fmt"
	"go/ast"
	"reflect"
)

type Type interface {
	typeMarker()
}

type IdentType struct {
	Name string
}

type PointerType struct {
	Inner Type
}

type ArrayType struct {
	Inner Type
}

type StructType struct {
	Fields []*Slot
}

type FuncType struct {
	Params  []*Slot
	Results []*Slot
}

type Slot struct {
	Name string
	Type Type
}

func (t *IdentType) typeMarker()   {}
func (t *PointerType) typeMarker() {}
func (t *ArrayType) typeMarker()   {}
func (t *StructType) typeMarker()  {}
func (t *FuncType) typeMarker()    {}

func parseType(expr ast.Expr) (Type, error) {
	switch e := expr.(type) {
	case *ast.Ident:
		return &IdentType{
			Name: e.Name,
		}, nil

	case *ast.StarExpr:
		inner, err := parseType(e.X)
		if err != nil {
			return nil, err
		}

		return &PointerType{
			Inner: inner,
		}, nil

	case *ast.ArrayType:
		if e.Len != nil {
			return nil, fmt.Errorf("%w: array type len not implemented", ErrAstUnsupported)
		}

		inner, err := parseType(e.Elt)
		if err != nil {
			return nil, err
		}

		return &ArrayType{
			Inner: inner,
		}, nil

	case *ast.StructType:
		fields, err := parseFieldList(e.Fields)
		if err != nil {
			return nil, err
		}

		return &StructType{
			Fields: fields,
		}, nil

	case *ast.FuncType:
		if e.TypeParams != nil {
			return nil, fmt.Errorf("%w: type params not implemented", ErrAstUnsupported)
		}

		params, err := parseFieldList(e.Params)
		if err != nil {
			return nil, err
		}

		results, err := parseFieldList(e.Results)
		if err != nil {
			return nil, err
		}

		return &FuncType{
			Params:  params,
			Results: results,
		}, nil

	default:
		return nil, fmt.Errorf("%w: unexpected type %v", ErrAstUnexpected, reflect.TypeOf(expr))
	}
}

func parseFieldList(l *ast.FieldList) ([]*Slot, error) {
	if l == nil {
		return nil, nil
	}

	var fields []*Slot

	for _, f := range l.List {
		inner, err := parseType(f.Type)
		if err != nil {
			return nil, err
		}

		for _, name := range f.Names {
			fields = append(fields, &Slot{
				Name: name.Name,
				Type: inner,
			})
		}
	}

	return fields, nil
}