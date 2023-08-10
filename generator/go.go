package generator

import "text/template"

func GetGoBridgeTemplate() *template.Template {
	t, err := template.New("go-bridge").Parse(goBridgeTemplateSrc)
	if err != nil {
		panic(err)
	}

	return t
}

var goBridgeTemplateSrc = `{{$top := . -}}
package main

import (
	orig "{{$top.TgtPkg}}"
)

/*
#include <stdint.h>
*/
import "C"

// Required by cgo
func main() {}
{{range $f := $top.Functions}}
//export fgb_{{$f.Name}}
func fgb_{{$f.Name}}({{range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.Name}} {{$p.CType}}{{end}}) {
	{{- range $i, $p := $f.Params}}
	{{- if eq $p.GoMode "cast"}}
	{{$p.Name}}Go := {{$p.GoType}}({{$p.Name}})
	{{- end}}
	{{- end}}
	orig.{{$f.TgtName}}({{range $i, $p := $f.Params}}
		{{- if gt $i 0}}, {{end}}{{$p.Name}}Go
	{{- end}})
}
{{end}}`
