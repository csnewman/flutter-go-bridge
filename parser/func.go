package parser

import (
	"go/ast"
	"log"
)

type Func struct{}

func (p *parser) processFuncDecl(d *ast.FuncDecl) {
	log.Println("Func", d)
}
