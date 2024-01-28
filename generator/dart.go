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
  factory Bridge.open() {
    return _FfiBridge();
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
{{- range $s := $top.RefStructs}}

abstract interface class {{$s.PascalName}} {
}

final class _Ffi{{$s.PascalName}} implements {{$s.PascalName}}, ffi.Finalizable {
  final ffi.Pointer<ffi.Void> id;

  _Ffi{{$s.PascalName}}(this.id);

  @override
  String toString() {
    return '{{$s.PascalName}}(${id.address})';
  }
}

{{- end}}
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

  @override
  String toString() {
    return '{{$s.PascalName}}{
{{- range $i, $f := $s.Fields -}}
{{if gt $i 0}}, {{end}}{{$f.CamelName}}: ${{$f.CamelName}}
{{- end -}}
}';
  }
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

@ffi.Native<_FgbC{{$s.PascalName}} Function()>(symbol: "fgbempty_{{$s.SnakeName}}")
external _FgbC{{$s.PascalName}} _fgbEmpty{{$s.PascalName}}();
{{- end}}

@ffi.Native<ffi.Pointer<ffi.Void> Function(ffi.Pointer<ffi.Void>)>(symbol: "fgbinternal_init")
external ffi.Pointer<ffi.Void> _fgbInternalInit(ffi.Pointer<ffi.Void> arg0);

@ffi.Native<ffi.Pointer Function(ffi.IntPtr)>(symbol: "fgbinternal_alloc")
external ffi.Pointer _fgbInternalAlloc(int arg0);

@ffi.Native<ffi.Void Function(ffi.Pointer)>(symbol: "fgbinternal_free")
external void _fgbInternalFree(ffi.Pointer arg0);

@ffi.Native<ffi.Void Function(ffi.Pointer<ffi.Void>)>(symbol: "fgbinternal_freepin")
external void _fgbInternalFreePin(ffi.Pointer<ffi.Void> arg0);

ffi.Pointer<ffi.NativeFinalizerFunction> _fgbInternalFreePinPtr = ffi.Native.addressOf(_fgbInternalFreePin);

class _GoAllocator implements ffi.Allocator {
  const _GoAllocator();

  @override
  ffi.Pointer<T> allocate<T extends ffi.NativeType>(int byteCount, {int? alignment}) {
    ffi.Pointer<T> result = _fgbInternalAlloc(byteCount).cast();
    if (result.address == 0) {
      throw ArgumentError('Could not allocate $byteCount bytes.');
    }
    return result;
  }

  @override
  void free(ffi.Pointer pointer) {
  	_fgbInternalFree(pointer);
  }
}

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

@ffi.Native<_FgbRet{{$f.PascalName}} Function(
  {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}
{{- if eq $p.DartMode "direct" -}}
  {{$p.DartCType}}
{{- else if eq $p.DartMode "map" -}}
  {{$p.DartCType}}
{{- else -}}
  unknown
{{- end}}{{end -}}
)>(symbol: "fgb_{{$f.SnakeName}}")
external _FgbRet{{$f.PascalName}} _fgb{{$f.PascalName}}(
  {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}
{{- if eq $p.DartMode "direct" -}}
  {{$p.DartType}}
{{- else if eq $p.DartMode "map" -}}
  {{$p.DartCType}}
{{- else -}}
  unknown
{{- end}} arg{{$i}}{{end -}}
);

@ffi.Native<ffi.Void Function(
{{- range $p := $f.Params}}{{- if eq $p.DartMode "direct" -}}
  {{$p.DartCType}}
{{- else if eq $p.DartMode "map" -}}
  {{$p.DartCType}}
{{- else -}}
  unknown
{{- end}}, {{end -}}
ffi.Uint64)>(symbol: "fgbasync_{{$f.SnakeName}}")
external void _fgbAsync{{$f.PascalName}}(
{{- range $i, $p := $f.Params}}{{- if eq $p.DartMode "direct" -}}
  {{$p.DartType}}
{{- else if eq $p.DartMode "map" -}}
  {{$p.DartCType}}
{{- else -}}
  unknown
{{- end}} arg{{$i}}, {{end -}}
int argPtr);

@ffi.Native<_FgbRet{{$f.PascalName}} Function(ffi.Uint64)>(symbol: "fgbasyncres_{{$f.SnakeName}}")
external _FgbRet{{$f.PascalName}} _fgbAsyncRes{{$f.PascalName}}(int arg0);
{{- end}}

final class _FfiBridge implements Bridge {
  late _GoAllocator _allocator;
  late ffi.NativeFinalizer _pinFinalizer;

  _FfiBridge() {
    _allocator = const _GoAllocator();
    _pinFinalizer = ffi.NativeFinalizer(_fgbInternalFreePinPtr);

    var initRes = _fgbInternalInit(ffi.NativeApi.initializeApiDLData);
    if (initRes != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(initRes.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
  }
{{range $f := $top.Functions}}
  @override
  {{if $f.HasRes}}{{$f.Res.DartType}}{{else}}void{{end}} {{$f.CamelName}}(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  ) {
    {{- range $p := $f.Params -}}
    {{- if eq $p.DartMode "direct"}}
    var __Dart__{{$p.Name}} = {{$p.Name}};
    {{- else if eq $p.DartMode "map"}}
    var __Dart__{{$p.Name}} = _mapFrom{{$p.MapName}}({{$p.Name}});
    {{- else}}
    var __Dart__{{$p.Name}} = unknown;
    {{- end -}}
    {{end}}
    {{if $f.HasRes}}return {{end}}_process{{$f.PascalName}}(_fgb{{$f.PascalName}}(
      {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}__Dart__{{$p.Name}}{{end -}}
    ));
  }

  @override
  Future<{{if $f.HasRes}}{{$f.Res.DartType}}{{else}}void{{end}}> {{$f.CamelName}}Async(
    {{- range $i, $p := $f.Params}}{{if gt $i 0}}, {{end}}{{$p.DartType}} {{$p.Name}}{{end -}}
  ) async {
    {{- range $p := $f.Params -}}
    {{- if eq $p.DartMode "direct"}}
    var __Dart__{{$p.Name}} = {{$p.Name}};
    {{- else if eq $p.DartMode "map"}}
    var __Dart__{{$p.Name}} = _mapFrom{{$p.DartType}}({{$p.Name}});
    {{- else}}
    var __Dart__{{$p.Name}} = unknown;
    {{- end -}}
    {{end}}
    var __DartRecv__ = ReceivePort('AsyncRecv({{$f.CamelName}})');
    _fgbAsync{{$f.PascalName}}(
      {{- range $p := $f.Params}}__Dart__{{$p.Name}}, {{end -}}
    __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    {{if $f.HasRes}}return {{end}}_process{{$f.PascalName}}(_fgbAsyncRes{{$f.PascalName}}(__DartMsg__[0]));
  }

  {{if $f.HasRes}}{{$f.Res.DartType}}{{else}}void{{end}} _process{{$f.PascalName}}(_FgbRet{{$f.PascalName}} res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
    {{- if $f.HasRes}}

    {{- if eq $f.Res.DartMode "direct"}}
    return res.res;
    {{- else if eq $f.Res.DartMode "map"}}
    return _mapTo{{$f.Res.MapName}}(res.res);
    {{- else}}
    return unknown;
    {{- end}}
    {{- end}}
  }
{{end -}}
{{- range $s := $top.ValueStructs}}

  {{$s.PascalName}} _mapTo{{$s.PascalName}}(_FgbC{{$s.PascalName}} from) {
    return {{$s.PascalName}}(
      {{- range $i, $f := $s.Fields -}}
      {{if gt $i 0}}, {{end}}
      {{- if eq $f.DartMode "direct" -}}
      from.{{$f.CamelName}}
      {{- else if eq $f.DartMode "map" -}}
      _mapTo{{$f.MapName}}(from.{{$f.CamelName}})
      {{- else -}}
      unknown
      {{- end}}
      {{- end -}}
    );
  }

  _FgbC{{$s.PascalName}} _mapFrom{{$s.PascalName}}({{$s.PascalName}} from) {
    var res = _fgbEmpty{{$s.PascalName}}();
    {{- range $f := $s.Fields}}
    {{- if eq $f.DartMode "direct"}}
    res.{{$f.CamelName}} = from.{{$f.CamelName}};
    {{- else if eq $f.DartMode "map"}}
    res.{{$f.CamelName}} = _mapFrom{{$f.MapName}}(from.{{$f.CamelName}});
    {{- else}}
    unknown
    {{- end}}
    {{- end}}
    return res;
  }
{{- end}}
{{- range $s := $top.RefStructs}}

  {{$s.PascalName}} _mapTo{{$s.PascalName}}(ffi.Pointer<ffi.Void> from) {
    var res = _Ffi{{$s.PascalName}}(from);
    _pinFinalizer.attach(res, from);
    return res;
  }

  ffi.Pointer<ffi.Void> _mapFrom{{$s.PascalName}}({{$s.PascalName}} from) {
    if (from is! _Ffi{{$s.PascalName}}) {
      throw 'Mismatched reference struct instance type';
    }

    return from.id;
  }
{{- end}}

  String _mapToString(ffi.Pointer<ffi.Void> from) {
    var res = ffi.Pointer<Utf8>.fromAddress(from.address).toDartString();
    _allocator.free(from);
    return res;
  }

  ffi.Pointer<ffi.Void> _mapFromString(String from) {
    var res = from.toNativeUtf8(allocator: _allocator);
    return ffi.Pointer<ffi.Void>.fromAddress(res.address);
  }
}
`
