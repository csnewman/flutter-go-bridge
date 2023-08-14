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
}

abstract interface class Bridge {
  factory Bridge.open(ffi.DynamicLibrary lib) {
    // TODO: Auto configure
    return _FfiBridge(lib);
  }
{{range $f := $top.Functions}}
  {{if $f.HasRes}}{{$f.ResDartType}} {{else}}void {{end}}{{$f.CamelName}}();
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
{{- end}}

final class _FfiBridge implements Bridge {
{{- range $f := $top.Functions}}
  late _FgbRet{{$f.PascalName}} Function() _{{$f.CamelName}}Ptr;
{{- end}}

  _FfiBridge(ffi.DynamicLibrary lib) {
{{- range $f := $top.Functions}}
    _{{$f.CamelName}}Ptr = lib.lookupFunction<
        _FgbRet{{$f.PascalName}} Function(),
        _FgbRet{{$f.PascalName}} Function()
    >("fgb_{{$f.SnakeName}}");
{{- end}}
  }
{{range $f := $top.Functions}}
  @override
  {{if $f.HasRes}}{{$f.ResDartType}} {{else}}void {{end}}{{$f.CamelName}}() {
    var res = _{{$f.CamelName}}Ptr();

    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      calloc.free(errPtr);

      throw BridgeException(errMsg);
    }

    // TODO
    throw BridgeException("todo process response");
  }
{{end -}}
}
`
