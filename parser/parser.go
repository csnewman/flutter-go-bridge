package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"reflect"
	"unicode"

	"golang.org/x/tools/go/packages"
)

var (
	ErrAstUnexpected  = errors.New("ast contained unexpected data")
	ErrAstUnsupported = errors.New("ast contained unsupported data")
	ErrInternal       = errors.New("internal error")
)

type Package struct {
	PkgPath   string
	TypeOrder []string
	Types     map[string]*TypeDef
	FuncOrder []string
	Funcs     map[string]*FuncDef
}

func Parse(path string, printAST bool) (*Package, error) {
	c := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports |
			packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypesSizes,
	}

	pkgs, err := packages.Load(c, "file="+path)
	if err != nil {
		return nil, err
	}

	p := &parser{
		Types: map[string]*TypeDef{},
		Funcs: map[string]*FuncDef{},
	}

	var pkgPath string

	log.Println("Parsing")
	for _, pkg := range pkgs {
		log.Println(" - Package", pkg)
		pkgPath = pkg.PkgPath

		for i, file := range pkg.CompiledGoFiles {
			log.Println("   - File", file)

			if printAST {
				ast.Print(pkg.Fset, pkg.Syntax[i])
			}

			if err := p.parse(pkg.Syntax[i]); err != nil {
				return nil, err
			}
		}
	}

	if err := p.process(); err != nil {
		return nil, err
	}

	if err := p.validate(); err != nil {
		return nil, err
	}

	return &Package{
		PkgPath:   pkgPath,
		TypeOrder: p.TypeOrder,
		Types:     p.Types,
		FuncOrder: p.FuncOrder,
		Funcs:     p.Funcs,
	}, nil
}

type parser struct {
	typeSpecs []*ast.TypeSpec
	funcDecls []*ast.FuncDecl
	TypeOrder []string
	Types     map[string]*TypeDef
	FuncOrder []string
	Funcs     map[string]*FuncDef
}

func (p *parser) parse(file *ast.File) error {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			// IMPORT, CONST, TYPE, or VAR
			switch d.Tok {
			case token.IMPORT, token.CONST, token.VAR:
				// ignore
				continue
			case token.TYPE:
				if len(d.Specs) != 1 {
					return fmt.Errorf("%w: type specs len != 1", ErrAstUnexpected)
				}

				ts, ok := (d.Specs[0]).(*ast.TypeSpec)
				if !ok {
					return fmt.Errorf("%w: type spec not *ast.TypeSpec", ErrAstUnexpected)
				}

				p.typeSpecs = append(p.typeSpecs, ts)
			default:
				return fmt.Errorf("%w: unexpected decl token %v", ErrAstUnexpected, d.Tok)
			}

		case *ast.FuncDecl:
			p.funcDecls = append(p.funcDecls, d)

		default:
			return fmt.Errorf("%w: unexpected decl %v", ErrAstUnexpected, reflect.TypeOf(decl))
		}
	}

	return nil
}

func (p *parser) process() error {
	log.Println("Processing")

	for _, spec := range p.typeSpecs {
		if err := p.processTypeSpec(spec); err != nil {
			return err
		}
	}

	for _, decl := range p.funcDecls {
		if err := p.processFuncDecl(decl); err != nil {
			return err
		}
	}

	changed := true

	for changed {
		changed = false

		for _, def := range p.Types {
			if def.Usage != UsageModeValue {
				continue
			}

			switch t := def.Type.(type) {
			case *StructType:
				for _, f := range t.Fields {
					c, err := p.processUsage(f.Type)
					if err != nil {
						return err
					}

					changed = changed || c
				}
			}
		}
	}

	return nil
}

func (p *parser) validate() error {
	log.Println("Validating")

	for _, def := range p.Types {
		log.Println(" - Type", def.Name)

		switch t := def.Type.(type) {
		case *StructType:
			switch def.Usage {
			case UsageModeNone:
				log.Printf("  - WARNING: %v is unused", def.Name)
			case UsageModeValue:
				// Ensure no private fields exist, as these can't be exposed
				for _, f := range t.Fields {
					for _, c := range f.Name {
						if !unicode.IsUpper(c) {
							return fmt.Errorf("%w: private field %v in value type %v", ErrAstUnsupported, f.Name, def.Name)
						}

						break
					}
				}
			case UsageModeRef:
				// No validation
			default:
				return fmt.Errorf("%w: unexpected struct usage when validating %v", ErrInternal, def.Usage)
			}
		default:
			return fmt.Errorf("%w: Unexpected type when validating %v", ErrInternal, reflect.TypeOf(def.Type))
		}
	}

	return nil
}
