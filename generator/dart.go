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
  {{if $f.HasRes}}{{$f.Res.DartType}}{{else}}void{{end}} {{$f.CamelName}}(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  );

  Future<{{if $f.HasRes}}{{$f.Res.DartType}}{{else}}void{{end}}> {{$f.CamelName}}Async(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  );
{{end -}}
}
{{- range $s := $top.ValueStructs}}

final class {{$s.PascalName}} {
{{- range $f := $s.Fields}}
    {{$f.DartType}} {{$f.CamelName}};
{{- end}}

    {{$s.PascalName}}(
{{- range $i, $f := $s.Fields -}}
{{if gt $i 0}}, {{end}}this.{{$f.CamelName}}
{{- end -}}
);
}

final class _FgbC{{$s.PascalName}} extends ffi.Struct {
{{- range $f := $s.Fields}}
{{- if eq $f.DartMode "direct"}}
  @{{$f.DartCType}}()
  external {{$f.DartType}} {{$f.CamelName}};
{{- else if eq $f.DartMode "map"}}
  external {{$f.DartCType}} {{$f.CamelName}};
{{- else}}
  external unknown {{$f.CamelName}};
{{- end}}
{{- end}}
}

typedef _FgbEmpty{{$s.PascalName}} = _FgbC{{$s.PascalName}} Function();
{{- end}}

typedef _FgbDefInit = ffi.Pointer<ffi.Void> Function(ffi.Pointer<ffi.Void>);

{{- range $f := $top.Functions}}

final class _FgbRet{{$f.PascalName}} extends ffi.Struct {
  {{- if $f.HasRes}}
{{- if eq $f.Res.DartMode "direct"}}
  @{{$f.Res.DartCType}}()
  external {{$f.Res.DartType}} res;
{{- else if eq $f.Res.DartMode "map"}}
  external {{$f.Res.DartCType}} res;
{{- else}}
  external unknown res;
{{- end}}
  {{- end}}
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDart{{$f.PascalName}} = _FgbRet{{$f.PascalName}} Function(
  {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}
{{- if eq $p.DartMode "direct" -}}
  {{$p.DartType}}
{{- else if eq $p.DartMode "map" -}}
  {{$p.DartCType}}
{{- else -}}
  unknown
{{- end}}{{end -}}
);
typedef _FgbDefC{{$f.PascalName}} = _FgbRet{{$f.PascalName}} Function(
  {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}
{{- if eq $p.DartMode "direct" -}}
  {{$p.DartCType}}
{{- else if eq $p.DartMode "map" -}}
  {{$p.DartCType}}
{{- else -}}
  unknown
{{- end}}{{end -}}
);
typedef _FgbAsyncDefDart{{$f.PascalName}} = void Function(
  {{- range $p := $f.Params}}{{- if eq $p.DartMode "direct" -}}
  {{$p.DartType}}
{{- else if eq $p.DartMode "map" -}}
  {{$p.DartCType}}
{{- else -}}
  unknown
{{- end}}, {{end -}}
int);
typedef _FgbAsyncDefC{{$f.PascalName}} = ffi.Void Function(
  {{- range $p := $f.Params}}{{- if eq $p.DartMode "direct" -}}
  {{$p.DartCType}}
{{- else if eq $p.DartMode "map" -}}
  {{$p.DartCType}}
{{- else -}}
  unknown
{{- end}}, {{end -}}
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

{{- range $s := $top.ValueStructs}}
  late _FgbEmpty{{$s.PascalName}} _empty{{$s.PascalName}}Ptr;
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
{{range $s := $top.ValueStructs}}
    _empty{{$s.PascalName}}Ptr = lib.lookupFunction<_FgbEmpty{{$s.PascalName}}, _FgbEmpty{{$s.PascalName}}>("fgbempty_{{$s.SnakeName}}");
{{- end}}
  }
{{range $f := $top.Functions}}
  @override
  {{if $f.HasRes}}{{$f.Res.DartType}}{{else}}void{{end}} {{$f.CamelName}}(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  ) {
    {{- range $p := $f.Params}}
    {{- if eq $p.DartMode "direct"}}
    var {{$p.Name}}Dart = {{$p.Name}};
    {{- else if eq $p.DartMode "map"}}
    var {{$p.Name}}Dart = _mapFrom{{$p.DartType}}({{$p.Name}});
    {{- else}}
    var {{$p.Name}}Dart = unknown;
    {{- end}}
    {{end}}
    {{if $f.HasRes}}return {{end}}_process{{$f.PascalName}}(_{{$f.CamelName}}Ptr(
      {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.Name}}Dart{{end -}}
    ));
  }

  @override
  Future<{{if $f.HasRes}}{{$f.Res.DartType}}{{else}}void{{end}}> {{$f.CamelName}}Async(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  ) async {
    {{- range $p := $f.Params}}
    {{- if eq $p.DartMode "direct"}}
    var {{$p.Name}}Dart = {{$p.Name}};
    {{- else if eq $p.DartMode "map"}}
    var {{$p.Name}}Dart = _mapFrom{{$p.DartType}}({{$p.Name}});
    {{- else}}
    var {{$p.Name}}Dart = unknown;
    {{- end}}
    {{end}}
    var recv = ReceivePort('AsyncRecv({{$f.CamelName}})');
    _{{$f.CamelName}}PtrAsync(
      {{- range $p := $f.Params}}{{$p.Name}}Dart, {{end -}}
    recv.sendPort.nativePort);
    var msg = await recv.first;
    recv.close();
    {{if $f.HasRes}}return {{end}}_process{{$f.PascalName}}(_{{$f.CamelName}}PtrAsyncRes(msg[0]));
  }

  {{if $f.HasRes}}{{$f.Res.DartType}}{{else}}void{{end}} _process{{$f.PascalName}}(_FgbRet{{$f.PascalName}} res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      calloc.free(errPtr);

      throw BridgeException(errMsg);
    }
    {{- if $f.HasRes}}

    {{- if eq $f.Res.DartMode "direct"}}
    return res.res;
    {{- else if eq $f.Res.DartMode "map"}}
    return _mapTo{{$f.Res.DartType}}(res.res);
    {{- else}}
    return unknown;
    {{- end}}
    {{- end}}
  }
{{- end -}}
{{- range $s := $top.ValueStructs}}

  {{$s.PascalName}} _mapTo{{$s.PascalName}}(_FgbC{{$s.PascalName}} from) {
    return {{$s.PascalName}}(
      {{- range $i, $f := $s.Fields -}}
      {{if gt $i 0}}, {{end}}
      {{- if eq $f.DartMode "direct" -}}
      from.{{$f.CamelName}}
      {{- else if eq $f.DartMode "map" -}}
      _mapTo{{$f.DartType}}(from.{{$f.CamelName}})
      {{- else -}}
      unknown
      {{- end}}
      {{- end -}}
    );
  }

  _FgbC{{$s.PascalName}} _mapFrom{{$s.PascalName}}({{$s.PascalName}} from) {
    var res = _empty{{$s.PascalName}}Ptr();
    {{- range $f := $s.Fields}}
    {{- if eq $f.DartMode "direct"}}
    res.{{$f.CamelName}} = from.{{$f.CamelName}};
    {{- else if eq $f.DartMode "map"}}
    res.{{$f.CamelName}} = _mapFrom{{$f.DartType}}(from.{{$f.CamelName}});
    {{- else}}
    unknown
    {{- end}}
    {{- end}}
    return res;
  }
{{- end}}
}
`
