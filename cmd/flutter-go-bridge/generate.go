package main

import (
	"fmt"
	"log"

	"github.com/csnewman/flutter-go-bridge/generator"
	"github.com/csnewman/flutter-go-bridge/parser"
	"github.com/spf13/cobra"
)

var cmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "Generate bindings",
	RunE:  executeGen,
}

var (
	genPrintAST bool
	genSrc      string
	genGoDst    string
	genDartDst  string
)

func init() {
	cmdGenerate.Flags().BoolVar(&genPrintAST, "print-ast", false, "Print the parsed AST (for debug use)")
	cmdGenerate.Flags().StringVar(&genSrc, "src", "", "The Go source to generate bindings for")
	cmdGenerate.Flags().StringVar(&genGoDst, "go", "", "The destination folder to store the Go bindings")
	cmdGenerate.Flags().StringVar(&genDartDst, "dart", "", "The destination folder to store the Dart bindings")

	_ = cmdGenerate.MarkFlagRequired("src")
}

func executeGen(cmd *cobra.Command, args []string) error {
	if genGoDst == "" && genDartDst == "" {
		return fmt.Errorf(`atleast one of the "go" or "dart" flags must be provided`)
	}

	log.Println("flutter-go-bridge generator")

	p, err := parser.Parse(genSrc, genPrintAST)
	if err != nil {
		log.Fatalln(err)
	}

	err = generator.Generate(genGoDst, genDartDst, p)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}
