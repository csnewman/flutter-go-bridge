package parser

import (
	"fmt"
	"go/ast"
	"log"
)

type UsageMode string

const (
	UsageModeNone  UsageMode = ""
	UsageModeValue UsageMode = "value"
	UsageModeRef   UsageMode = "ref"
)

type TypeDef struct {
	Name      string
	Type      Type
	Usage     UsageMode
	FuncOrder []string
	Funcs     map[string]*FuncDef
}

func (p *parser) processTypeSpec(ts *ast.TypeSpec) error {
	name := ts.Name.Name

	log.Println(" - Type", name)

	if ts.TypeParams != nil {
		return fmt.Errorf("%w: type params not implemented", ErrAstUnsupported)
	}

	t, err := parseType(ts.Type)
	if err != nil {
		return err
	}

	p.TypeOrder = append(p.TypeOrder, name)
	p.Types[name] = &TypeDef{
		Name:  name,
		Type:  t,
		Funcs: map[string]*FuncDef{},
	}

	return nil
}
