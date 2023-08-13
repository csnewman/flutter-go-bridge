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
{{- if or $f.HasRes $f.HasErr}}

typedef struct {
    {{- if $f.HasRes}}
    {{$f.ResCType}} res;
    {{- end}}
    {{- if $f.HasErr}}
    void* err;
    {{- end}}
} fgb_ret_{{$f.Name}};
{{- end}}
{{- end}}
*/
import "C"

// Required by cgo
func main() {}
{{range $f := $top.Functions}}
//export fgb_{{$f.Name}}
func fgb_{{$f.Name}}({{range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.Name}} C.{{$p.CType}}{{end}}) 
{{- if or $f.HasRes $f.HasErr}} C.fgb_ret_{{$f.Name}} {{- end}} {
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
    {{- if or $f.HasRes $f.HasErr}}
    {{- if $f.HasRes}}
    {{- if eq $f.ResGoMode "cast"}}
	cres := (C.{{$f.ResCType}})(gres)
	{{- end}}
    {{- end}}

    {{- if $f.HasErr}}

    var cerr unsafe.Pointer
    if gerr != nil {
        cerr = unsafe.Pointer(C.CString(gerr.Error()))
    }
    {{- end}}

    return C.fgb_ret_{{$f.Name}} {
        {{- if $f.HasRes}}
        res: cres,
        {{- end}}
        {{- if $f.HasErr}}
        err: cerr,
        {{- end}}
    }
    {{- end}} 
}
{{end}}`
