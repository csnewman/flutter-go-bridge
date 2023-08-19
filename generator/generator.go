package generator

import (
	"errors"
	"fmt"
	"log"
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
	TgtPkg       string
	ValueStructs []*ValueStruct
	Functions    []*Func
}

type ValueStruct struct {
	OrigName   string
	SnakeName  string
	CamelName  string
	PascalName string
	Fields     []*Field
}

type Type struct {
	CType     string
	GoType    string
	GoMode    string
	MapName   string
	DartCType string
	DartType  string
	DartMode  string
}

type Field struct {
	Type
	OrigName   string
	SnakeName  string
	PascalName string
	CamelName  string
}

type Func struct {
	SnakeName  string
	CamelName  string
	PascalName string
	TgtName    string
	Params     []*Param
	HasRes     bool
	Res        Type
	HasErr     bool
}

type Param struct {
	Type
	Name string
}

func Generate(goDest string, dartDest string, in *parser.Package) error {
	g := &generator{
		pkg: in,
		unit: &Unit{
			TgtPkg: in.PkgPath,
		},
	}

	log.Println("Generating")
	if err := g.process(in); err != nil {
		return err
	}

	// Generate go
	log.Println("Generating Go")
	f, err := os.Create(goDest)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := GetGoBridgeTemplate().Execute(f, g.unit); err != nil {
		return err
	}

	// Generate dart
	log.Println("Generating Dart")
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
	pkg  *parser.Package
	unit *Unit
}

func (g *generator) process(p *parser.Package) error {
	for _, name := range p.TypeOrder {
		log.Println(" - Type", name)
		if err := g.processTypeDef(p.Types[name]); err != nil {
			return err
		}
	}

	for _, name := range p.FuncOrder {
		log.Println(" - Func", name)
		if err := g.processFunc(p.Funcs[name]); err != nil {
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

			mt, err := g.mapType(res.Type)
			if err != nil {
				return err
			}

			fu.Res = mt
		}

	}

	for _, param := range f.Sig.Params {
		mt, err := g.mapType(param.Type)
		if err != nil {
			return err
		}

		fu.Params = append(fu.Params, &Param{
			Type: mt,
			Name: param.Name,
		})
	}

	g.unit.Functions = append(g.unit.Functions, fu)

	return nil
}

var inbuiltTypes = map[string]Type{
	"error":  {CType: "void*", GoType: "error", GoMode: "todo", DartCType: "todo", DartType: "todo", DartMode: "todo"},
	"string": {CType: "void*", GoType: "string", GoMode: "todo", DartCType: "todo", DartType: "todo", DartMode: "todo"},
	"int8":   {CType: "int8", GoType: "int8", GoMode: "cast", DartCType: "ffi.Int8", DartType: "int", DartMode: "direct"},
	"uint8":  {CType: "uint8", GoType: "uint8", GoMode: "cast", DartCType: "ffi.Uint8", DartType: "int", DartMode: "direct"},
	"byte":   {CType: "byte", GoType: "byte", GoMode: "cast", DartCType: "ffi.Byte", DartType: "int", DartMode: "direct"},
	"int16":  {CType: "int16", GoType: "int16", GoMode: "cast", DartCType: "ffi.Int16", DartType: "int", DartMode: "direct"},
	"uint16": {CType: "uint16", GoType: "uint16", GoMode: "cast", DartCType: "ffi.Uint16", DartType: "int", DartMode: "direct"},
	"int32":  {CType: "int32", GoType: "int32", GoMode: "cast", DartCType: "ffi.Int32", DartType: "int", DartMode: "direct"},
	"uint32": {CType: "uint32", GoType: "uint32", GoMode: "cast", DartCType: "ffi.Uint32", DartType: "int", DartMode: "direct"},
	"int64":  {CType: "int64", GoType: "int64", GoMode: "cast", DartCType: "ffi.Int64", DartType: "int", DartMode: "direct"},
	"uint64": {CType: "uint64", GoType: "uint64", GoMode: "cast", DartCType: "ffi.Uint64", DartType: "int", DartMode: "direct"},
	"int":    {CType: "int", GoType: "int", GoMode: "cast", DartCType: "ffi.Int", DartType: "int", DartMode: "direct"},
	"uint":   {CType: "uint", GoType: "uint", GoMode: "cast", DartCType: "ffi.Uint", DartType: "int", DartMode: "direct"},
}

func (g *generator) processTypeDef(def *parser.TypeDef) error {
	switch t := def.Type.(type) {
	case *parser.IdentType:
		return fmt.Errorf("%w: aliases not supported", ErrUnsupported)
	case *parser.PointerType:
		return fmt.Errorf("%w: aliases not supported", ErrUnsupported)
	case *parser.StructType:
		if def.Usage == parser.UsageModeNone {
			log.Println("   - Type not used")

			return nil
		} else if def.Usage == parser.UsageModeRef {
			log.Println("   - Type used as ref")

			return nil
		}

		log.Println("   - Type used as value")

		vt := &ValueStruct{
			OrigName:   def.Name,
			SnakeName:  strcase.ToSnake(def.Name),
			PascalName: strcase.ToCamel(def.Name),
			CamelName:  strcase.ToLowerCamel(def.Name),
		}

		for _, f := range t.Fields {
			mt, err := g.mapType(f.Type)
			if err != nil {
				return err
			}

			vt.Fields = append(vt.Fields, &Field{
				Type:       mt,
				OrigName:   f.Name,
				SnakeName:  strcase.ToSnake(f.Name),
				PascalName: strcase.ToCamel(f.Name),
				CamelName:  strcase.ToLowerCamel(f.Name),
			})
		}

		g.unit.ValueStructs = append(g.unit.ValueStructs, vt)

		return nil
	default:
		return fmt.Errorf("%w: unexpected type %v", ErrUnexpected, reflect.TypeOf(def.Type))
	}
}

func (g *generator) mapType(pt parser.Type) (Type, error) {
	switch t := pt.(type) {
	case *parser.IdentType:
		ibt, ok := inbuiltTypes[t.Name]
		if ok {
			return ibt, nil
		}

		def, ok := g.pkg.Types[t.Name]
		if !ok {
			return Type{}, fmt.Errorf("%w: type is not defined: %v", ErrUnexpected, t.Name)
		}

		if def.Usage != parser.UsageModeValue {
			return Type{}, fmt.Errorf("%w: field has unexpected usage %v", ErrUnexpected, def.Usage)
		}

		return Type{
			CType:     "fgb_vt_" + strcase.ToSnake(t.Name),
			GoType:    "orig." + t.Name,
			GoMode:    "map",
			MapName:   strcase.ToCamel(t.Name),
			DartCType: "_FgbC" + strcase.ToCamel(t.Name),
			DartType:  strcase.ToCamel(t.Name),
			DartMode:  "map",
		}, nil
	default:
		return Type{}, fmt.Errorf("%w: unexpected type %v", ErrUnexpected, reflect.TypeOf(pt))
	}
}
