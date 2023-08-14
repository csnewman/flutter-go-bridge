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
    "unsafe"

	orig "{{$top.TgtPkg}}"
)

/*
#include <stdint.h>
{{- range $f := $top.Functions}}

typedef struct {
    {{- if $f.HasRes}}
    {{$f.ResCType}} res;
    {{- end}}
    void* err;
} fgb_ret_{{$f.SnakeName}};
{{- end}}
*/
import "C"

// Required by cgo
func main() {}
{{range $f := $top.Functions}}
//export fgb_{{$f.SnakeName}}
func fgb_{{$f.SnakeName}}({{range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.Name}} C.{{$p.CType}}{{end}}) C.fgb_ret_{{$f.SnakeName}} {
	{{- range $i, $p := $f.Params}}
	{{- if eq $p.GoMode "cast"}}
	{{$p.Name}}Go := ({{$p.GoType}})({{$p.Name}})
	{{- end}}
	{{- end}}
	{{if $f.HasRes}}gres{{if $f.HasErr}}, {{end}}{{end}}{{if $f.HasErr}}gerr{{end -}}
	{{if or $f.HasRes $f.HasErr}} := {{end -}}
	orig.{{$f.TgtName}}({{range $i, $p := $f.Params}}
		{{- if gt $i 0}}, {{end}}{{$p.Name}}Go
	{{- end}})
    {{- if $f.HasRes}}
    {{- if eq $f.ResGoMode "cast"}}
	cres := (C.{{$f.ResCType}})(gres)
	{{- end}}
    {{- end}}

    var cerr unsafe.Pointer
    {{- if $f.HasErr}}
    if gerr != nil {
        cerr = unsafe.Pointer(C.CString(gerr.Error()))
    }
    {{- end}}

    return C.fgb_ret_{{$f.SnakeName}} {
        {{- if $f.HasRes}}
        res: cres,
        {{- end}}
        err: cerr,
    }
}
{{end}}`
