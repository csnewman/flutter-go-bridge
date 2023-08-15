package generator

import "text/template"

func GetDartBridgeTemplate() *template.Template {
	t, err := template.New("dart-bridge").Parse(dartBridgeTemplateSrc)
	if err != nil {
		panic(err)
	}

	return t
}

var dartBridgeTemplateSrc = `{{$top := . -}}
import 'dart:ffi' as ffi;
import 'package:ffi/ffi.dart';

final class BridgeException implements Exception {
  String cause;

  BridgeException(this.cause);

  @override
  String toString() {
    return 'BridgeException: $cause';
  }
}

abstract interface class Bridge {
  factory Bridge.open(ffi.DynamicLibrary lib) {
    // TODO: Auto configure
    return _FfiBridge(lib);
  }
{{range $f := $top.Functions}}
  {{if $f.HasRes}}{{$f.ResDartType}} {{else}}void {{end}}{{$f.CamelName}}(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  );
{{end -}}
}

{{- range $f := $top.Functions}}

final class _FgbRet{{$f.PascalName}} extends ffi.Struct {
  {{- if $f.HasRes}}
  @ffi.Int32()
  external {{$f.ResCType}} res;
  {{- end}}
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDart{{$f.PascalName}} = _FgbRet{{$f.PascalName}} Function(
  {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}}{{end -}}
);
typedef _FgbDefC{{$f.PascalName}} = _FgbRet{{$f.PascalName}} Function(
  {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartCType}}{{end -}}
);
{{- end}}

final class _FfiBridge implements Bridge {
{{- range $f := $top.Functions}}
  late _FgbDefDart{{$f.PascalName}} _{{$f.CamelName}}Ptr;
{{- end}}

  _FfiBridge(ffi.DynamicLibrary lib) {
{{- range $f := $top.Functions}}
    _{{$f.CamelName}}Ptr = lib.lookupFunction<_FgbDefC{{$f.PascalName}}, _FgbDefDart{{$f.PascalName}}>("fgb_{{$f.SnakeName}}");
{{- end}}
  }
{{range $f := $top.Functions}}
  @override
  {{if $f.HasRes}}{{$f.ResDartType}} {{else}}void {{end}}{{$f.CamelName}}(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  ) {
    var res = _{{$f.CamelName}}Ptr(
      {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.Name}}{{end -}}
    );

    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      calloc.free(errPtr);

      throw BridgeException(errMsg);
    }
    {{- if $f.HasRes}}

    return res.res;
    {{- end}}
  }
{{end -}}
}
`
