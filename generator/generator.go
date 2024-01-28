package generator

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/csnewman/flutter-go-bridge/parser"
	"github.com/iancoleman/strcase"
)

var (
	ErrUnexpected  = errors.New("parser contained unexpected data")
	ErrUnsupported = errors.New("unsupported")
)

type Unit struct {
	TgtPkg       string
	ValueStructs []*ValueStruct
	RefStructs   []*RefStruct
	Functions    []*Func
}

type ValueStruct struct {
	OrigName   string
	SnakeName  string
	CamelName  string
	PascalName string
	Fields     []*Field
}

type RefStruct struct {
	OrigName   string
	SnakeName  string
	CamelName  string
	PascalName string
}

type Type struct {
	CType     string
	GoType    string
	GoCType   string
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

	log.Println("Preparing")
	if err := g.process(in); err != nil {
		return err
	}

	log.Println("Generating")

	// Generate go
	if goDest != "" {
		log.Println(" - Go", goDest)
		f, err := os.Create(goDest)
		if err != nil {
			return err
		}

		defer f.Close()

		if err := GetGoBridgeTemplate().Execute(f, g.unit); err != nil {
			return err
		}
	}

	// Generate dart
	if dartDest != "" {
		log.Println(" - Dart", dartDest)
		f, err := os.Create(dartDest)
		if err != nil {
			return err
		}

		defer f.Close()

		if err := GetDartBridgeTemplate().Execute(f, g.unit); err != nil {
			return err
		}
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
	"error":  {CType: "void*", GoType: "error", GoCType: "unsafe.Pointer", GoMode: "map", MapName: "Error", DartCType: "ffi.Pointer<ffi.Void>", DartType: "String", DartMode: "map"},
	"string": {CType: "void*", GoType: "string", GoCType: "unsafe.Pointer", GoMode: "map", MapName: "String", DartCType: "ffi.Pointer<ffi.Void>", DartType: "String", DartMode: "map"},
	"int8":   {CType: "int8_t", GoType: "int8", GoCType: "C.int8_t", GoMode: "cast", DartCType: "ffi.Int8", DartType: "int", DartMode: "direct"},
	"uint8":  {CType: "uint8_t", GoType: "uint8", GoCType: "C.uint8_t", GoMode: "cast", DartCType: "ffi.Uint8", DartType: "int", DartMode: "direct"},
	"byte":   {CType: "byte", GoType: "byte", GoCType: "C.byte", GoMode: "cast", DartCType: "ffi.Byte", DartType: "int", DartMode: "direct"},
	"int16":  {CType: "int16_t", GoType: "int16", GoCType: "C.int16_t", GoMode: "cast", DartCType: "ffi.Int16", DartType: "int", DartMode: "direct"},
	"uint16": {CType: "uint16_t", GoType: "uint16", GoCType: "C.uint16_t", GoMode: "cast", DartCType: "ffi.Uint16", DartType: "int", DartMode: "direct"},
	"int32":  {CType: "int32_t", GoType: "int32", GoCType: "C.int32_t", GoMode: "cast", DartCType: "ffi.Int32", DartType: "int", DartMode: "direct"},
	"uint32": {CType: "uint32_t", GoType: "uint32", GoCType: "C.uint32_t", GoMode: "cast", DartCType: "ffi.Uint32", DartType: "int", DartMode: "direct"},
	"int64":  {CType: "int64_t", GoType: "int64", GoCType: "C.int64_t", GoMode: "cast", DartCType: "ffi.Int64", DartType: "int", DartMode: "direct"},
	"uint64": {CType: "uint64_t", GoType: "uint64", GoCType: "C.uint64_t", GoMode: "cast", DartCType: "ffi.Uint64", DartType: "int", DartMode: "direct"},
	"int":    {CType: "int", GoType: "int", GoCType: "C.int", GoMode: "cast", DartCType: "ffi.Int", DartType: "int", DartMode: "direct"},
	"uint":   {CType: "uint", GoType: "uint", GoCType: "C.uint", GoMode: "cast", DartCType: "ffi.Uint", DartType: "int", DartMode: "direct"},
}

func (g *generator) processTypeDef(def *parser.TypeDef) error {
	switch t := def.Type.(type) {
	case *parser.IdentType:
		return fmt.Errorf("%w: aliases not supported", ErrUnsupported)
	case *parser.PointerType:
		return fmt.Errorf("%w: aliases not supported", ErrUnsupported)
	case *parser.StructType:
		switch def.Usage {
		case parser.UsageModeNone:
			log.Println("   - Type not used")
			return nil
		case parser.UsageModeValue:
			log.Println("   - Type used as value")
			return g.processValueStruct(def, t)
		case parser.UsageModeRef:
			log.Println("   - Type used as ref")
			return g.processRefStruct(def, t)
		default:
			return fmt.Errorf("%w: unexpected struct usage %v", ErrUnexpected, def.Usage)
		}
	default:
		return fmt.Errorf("%w: unexpected type def %v", ErrUnexpected, reflect.TypeOf(def.Type))
	}
}

func (g *generator) processValueStruct(def *parser.TypeDef, t *parser.StructType) error {
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
}

func (g *generator) processRefStruct(def *parser.TypeDef, t *parser.StructType) error {
	rt := &RefStruct{
		OrigName:   def.Name,
		SnakeName:  strcase.ToSnake(def.Name),
		PascalName: strcase.ToCamel(def.Name),
		CamelName:  strcase.ToLowerCamel(def.Name),
	}

	g.unit.RefStructs = append(g.unit.RefStructs, rt)

	return nil
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
			GoCType:   "C.fgb_vt_" + strcase.ToSnake(t.Name),
			GoMode:    "map",
			MapName:   strcase.ToCamel(t.Name),
			DartCType: "_FgbC" + strcase.ToCamel(t.Name),
			DartType:  strcase.ToCamel(t.Name),
			DartMode:  "map",
		}, nil
	case *parser.PointerType:
		ident, ok := t.Inner.(*parser.IdentType)
		if !ok {
			return Type{}, fmt.Errorf("%w: pointer to %v (%v) is not supported", ErrUnsupported, t.Inner, reflect.TypeOf(t.Inner))
		}

		if _, ok := inbuiltTypes[ident.Name]; ok {
			return Type{}, fmt.Errorf("%w: pointer to primitive %v is not supported", ErrUnsupported, ident.Name)
		}

		def, ok := g.pkg.Types[ident.Name]
		if !ok {
			return Type{}, fmt.Errorf("%w: type is not defined: %v", ErrUnexpected, ident.Name)
		}

		if def.Usage != parser.UsageModeRef {
			return Type{}, fmt.Errorf("%w: pointer type has unexpected usage %v", ErrUnexpected, def.Usage)
		}

		return Type{
			CType:     "uintptr_t",
			GoType:    "*orig." + ident.Name,
			GoCType:   "C.uintptr_t",
			GoMode:    "map",
			MapName:   strcase.ToCamel(ident.Name),
			DartCType: "ffi.Pointer<ffi.Void>",
			DartType:  strcase.ToCamel(ident.Name),
			DartMode:  "map",
		}, nil
	default:
		return Type{}, fmt.Errorf("%w: mapping unexpected type %v: %v", ErrUnexpected, reflect.TypeOf(pt), pt)
	}
}
