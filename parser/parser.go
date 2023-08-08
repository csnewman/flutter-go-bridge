package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/packages"
)

var ErrAstUnexpected = errors.New("ast contained unexpected data")

func Parse(path string) error {
	c := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports |
			packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypesSizes,
	}

	pkgs, err := packages.Load(c, "file="+path)
	if err != nil {
		return err
	}

	p := &parser{}

	log.Println("Parsing")
	for _, pkg := range pkgs {
		log.Println(" - Package", pkg)

		for i, file := range pkg.CompiledGoFiles {
			log.Println("   - File", file)

			if err := p.parse(pkg.Syntax[i]); err != nil {
				return err
			}
		}
	}

	if err := p.process(); err != nil {
		return err
	}

	return nil
}

type parser struct {
	typeSpecs []*ast.TypeSpec
	funcDecls []*ast.FuncDecl
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
			return fmt.Errorf("%w: unexpected decl %v", ErrAstUnexpected, decl)
		}
	}

	return nil
}

func (p *parser) process() error {
	for _, spec := range p.typeSpecs {
		p.processTypeSpec(spec)
	}

	for _, decl := range p.funcDecls {
		p.processFuncDecl(decl)
	}

	return nil
}

func (p *parser) processTypeSpec(ts *ast.TypeSpec) {
	log.Println("TypeSpec", ts)
}
