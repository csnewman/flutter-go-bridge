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
import 'dart:isolate';
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
  {{if $f.HasRes}}{{$f.ResDartType}}{{else}}void{{end}} {{$f.CamelName}}(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  );

  Future<{{if $f.HasRes}}{{$f.ResDartType}}{{else}}void{{end}}> {{$f.CamelName}}Async(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  );
{{end -}}
}

typedef _FgbDefInit = ffi.Pointer<ffi.Void> Function(ffi.Pointer<ffi.Void>);

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
typedef _FgbAsyncDefDart{{$f.PascalName}} = void Function(
  {{- range $p := $f.Params}}{{$p.DartType}}, {{end -}}
int);
typedef _FgbAsyncDefC{{$f.PascalName}} = ffi.Void Function(
  {{- range $p := $f.Params}}{{$p.DartCType}}, {{end -}}
ffi.Uint64);
typedef _FgbAsyncResDefDart{{$f.PascalName}} = _FgbRet{{$f.PascalName}} Function(int);
typedef _FgbAsyncResDefC{{$f.PascalName}} = _FgbRet{{$f.PascalName}} Function(ffi.Uint64);
{{- end}}

final class _FfiBridge implements Bridge {
{{- range $f := $top.Functions}}
  late _FgbDefDart{{$f.PascalName}} _{{$f.CamelName}}Ptr;
  late _FgbAsyncDefDart{{$f.PascalName}} _{{$f.CamelName}}PtrAsync;
  late _FgbAsyncResDefDart{{$f.PascalName}} _{{$f.CamelName}}PtrAsyncRes;
{{- end}}

  _FfiBridge(ffi.DynamicLibrary lib) {
    var initPtr = lib.lookupFunction<_FgbDefInit, _FgbDefInit>("fgb_internal_init");
    var initRes = initPtr(ffi.NativeApi.initializeApiDLData);
    if (initRes != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(initRes.address);
      var errMsg = errPtr.toDartString(); 
      calloc.free(errPtr);

      throw BridgeException(errMsg);
    }
{{range $f := $top.Functions}}
    _{{$f.CamelName}}Ptr = lib.lookupFunction<_FgbDefC{{$f.PascalName}}, _FgbDefDart{{$f.PascalName}}>("fgb_{{$f.SnakeName}}");
    _{{$f.CamelName}}PtrAsync = lib.lookupFunction<_FgbAsyncDefC{{$f.PascalName}}, _FgbAsyncDefDart{{$f.PascalName}}>("fgbasync_{{$f.SnakeName}}");
    _{{$f.CamelName}}PtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefC{{$f.PascalName}}, _FgbAsyncResDefDart{{$f.PascalName}}>("fgbasyncres_{{$f.SnakeName}}");
{{- end}}
  }
{{range $f := $top.Functions}}
  @override
  {{if $f.HasRes}}{{$f.ResDartType}}{{else}}void{{end}} {{$f.CamelName}}(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  ) {
    {{if $f.HasRes}}return {{end}}_process{{$f.PascalName}}(_{{$f.CamelName}}Ptr(
      {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.Name}}{{end -}}
    ));
  }

  @override
  Future<{{if $f.HasRes}}{{$f.ResDartType}}{{else}}void{{end}}> {{$f.CamelName}}Async(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  ) async {
    var recv = ReceivePort('AsyncRecv({{$f.CamelName}})');
    _{{$f.CamelName}}PtrAsync(
      {{- range $p := $f.Params}}{{$p.Name}}, {{end -}}
    recv.sendPort.nativePort);
    var msg = await recv.first;
    recv.close();
    {{if $f.HasRes}}return {{end}}_process{{$f.PascalName}}(_{{$f.CamelName}}PtrAsyncRes(msg[0]));
  }

  {{if $f.HasRes}}{{$f.ResDartType}}{{else}}void{{end}} _process{{$f.PascalName}}(_FgbRet{{$f.PascalName}} res) {
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
