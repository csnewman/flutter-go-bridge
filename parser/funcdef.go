package parser

import (
	"fmt"
	"go/ast"
	"log"
	"reflect"
)

type FuncDef struct {
	Name string
	Sig  *FuncType
	Recv UsageMode
}

func (p *parser) processFuncDecl(d *ast.FuncDecl) error {
	name := d.Name.Name

	log.Println(" - Func", name)

	recvs, err := parseFieldList(d.Recv)
	if err != nil {
		return err
	}

	if len(recvs) > 1 {
		return fmt.Errorf("%w: more than one reciever not implemented", ErrAstUnsupported)
	}

	rawSig, err := parseType(d.Type)
	if err != nil {
		return err
	}

	sig, ok := rawSig.(*FuncType)
	if !ok {
		return fmt.Errorf("%w: unexpected func sig %v", ErrAstUnexpected, reflect.TypeOf(rawSig))
	}

	if _, err := p.processUsage(sig); err != nil {
		return err
	}

	def := &FuncDef{
		Name: name,
		Sig:  sig,
	}

	if len(recvs) == 0 {
		p.FuncOrder = append(p.FuncOrder, name)
		p.Funcs[name] = def
		return nil
	}

	recv := recvs[0]
	var recvType string

	switch r := recv.Type.(type) {
	case *IdentType:
		def.Recv = UsageModeValue
		recvType = r.Name
	case *PointerType:
		ir, ok := r.Inner.(*IdentType)
		if !ok {
			return fmt.Errorf("%w: unexpected nested pointer %v: %v", ErrAstUnexpected, reflect.TypeOf(r.Inner), r.Inner)
		}

		def.Recv = UsageModeRef
		recvType = ir.Name
	default:
		return fmt.Errorf("%w: unexpected recv type %v", ErrAstUnexpected, reflect.TypeOf(recv))
	}

	recvDef, ok := p.Types[recvType]
	if !ok {
		return fmt.Errorf("%w: unknown recv %v", ErrAstUnexpected, recvType)
	}

	if recvDef.Usage != UsageModeNone && recvDef.Usage != def.Recv {
		return fmt.Errorf("%w: %v has both %v and %v usages", ErrAstUnsupported, recvType, recvDef.Usage, def.Recv)
	}

	recvDef.Usage = def.Recv
	recvDef.FuncOrder = append(p.FuncOrder, name)
	recvDef.Funcs[name] = def

	return nil
}
