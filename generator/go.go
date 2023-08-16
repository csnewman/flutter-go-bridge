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
    "fmt"
	"sync"
	"sync/atomic"
    "unsafe"

	orig "{{$top.TgtPkg}}"
	"flutter-go-bridge/runtime"
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

var (
	handles   = sync.Map{}
	handleIdx uint64
)

// Required by cgo
func main() {}

//export fgb_internal_init
func fgb_internal_init(p unsafe.Pointer) unsafe.Pointer {
	err := runtime.InitializeApi(p)

    var cerr unsafe.Pointer
    if err != nil {
        cerr = unsafe.Pointer(C.CString(err.Error()))
    }

	return cerr
}
{{range $f := $top.Functions}}
//export fgb_{{$f.SnakeName}}
func fgb_{{$f.SnakeName}}({{range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.Name}} C.{{$p.CType}}{{end}}) (resw C.fgb_ret_{{$f.SnakeName}}) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_{{$f.SnakeName}} {
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	{{range $i, $p := $f.Params}}
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


//export fgbasync_{{$f.SnakeName}}
func fgbasync_{{$f.SnakeName}}({{range $p := $f.Params}}{{$p.Name}} C.{{$p.CType}}, {{end}}fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_{{$f.SnakeName}}({{range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.Name}}{{end}}))

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
        if !sent {
            handles.LoadAndDelete(h)
        }
	}()
}


//export fgbasyncres_{{$f.SnakeName}}
func fgbasyncres_{{$f.SnakeName}}(h uint64) C.fgb_ret_{{$f.SnakeName}} {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_{{$f.SnakeName}}{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_{{$f.SnakeName}})
}
{{end}}`
