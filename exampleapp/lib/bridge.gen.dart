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

  int add(int a, int b);

  Future<int> addAsync(int a, int b);

  Point addPoints(Point a, Point b);

  Future<Point> addPointsAsync(Point a, Point b);

  int addError(int a, int b);

  Future<int> addErrorAsync(int a, int b);

  Obj newObj(String name, int other);

  Future<Obj> newObjAsync(String name, int other);

  void modifyObj(Obj o);

  Future<void> modifyObjAsync(Obj o);

  String formatObj(Obj o);

  Future<String> formatObjAsync(Obj o);
}

abstract interface class Obj {
}

final class _FfiObj implements Obj, ffi.Finalizable {
  final ffi.Pointer<ffi.Void> id;

  _FfiObj(this.id);

  @override
  String toString() {
    return 'Obj(${id.address})';
  }
}

final class Point {
  int x;
  int y;
  String name;

  Point(this.x, this.y, this.name);

  @override
  String toString() {
    return 'Point{x: $x, y: $y, name: $name}';
  }
}

final class _FgbCPoint extends ffi.Struct {
  @ffi.Int()
  external int x;
  @ffi.Int()
  external int y;
  external ffi.Pointer<ffi.Void> name;
}

@ffi.Native<_FgbCPoint Function()>(symbol: "fgbempty_point")
external _FgbCPoint _fgbEmptyPoint();

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

final class _FgbRetAdd extends ffi.Struct {
  @ffi.Int()
  external int res;
  external ffi.Pointer<ffi.Void> err;
}

@ffi.Native<_FgbRetAdd Function(ffi.Int, ffi.Int)>(symbol: "fgb_add")
external _FgbRetAdd _fgbAdd(int arg0, int arg1);

@ffi.Native<ffi.Void Function(ffi.Int, ffi.Int, ffi.Uint64)>(symbol: "fgbasync_add")
external void _fgbAsyncAdd(int arg0, int arg1, int argPtr);

@ffi.Native<_FgbRetAdd Function(ffi.Uint64)>(symbol: "fgbasyncres_add")
external _FgbRetAdd _fgbAsyncResAdd(int arg0);

final class _FgbRetAddPoints extends ffi.Struct {
  external _FgbCPoint res;
  external ffi.Pointer<ffi.Void> err;
}

@ffi.Native<_FgbRetAddPoints Function(_FgbCPoint, _FgbCPoint)>(symbol: "fgb_add_points")
external _FgbRetAddPoints _fgbAddPoints(_FgbCPoint arg0, _FgbCPoint arg1);

@ffi.Native<ffi.Void Function(_FgbCPoint, _FgbCPoint, ffi.Uint64)>(symbol: "fgbasync_add_points")
external void _fgbAsyncAddPoints(_FgbCPoint arg0, _FgbCPoint arg1, int argPtr);

@ffi.Native<_FgbRetAddPoints Function(ffi.Uint64)>(symbol: "fgbasyncres_add_points")
external _FgbRetAddPoints _fgbAsyncResAddPoints(int arg0);

final class _FgbRetAddError extends ffi.Struct {
  @ffi.Int()
  external int res;
  external ffi.Pointer<ffi.Void> err;
}

@ffi.Native<_FgbRetAddError Function(ffi.Int, ffi.Int)>(symbol: "fgb_add_error")
external _FgbRetAddError _fgbAddError(int arg0, int arg1);

@ffi.Native<ffi.Void Function(ffi.Int, ffi.Int, ffi.Uint64)>(symbol: "fgbasync_add_error")
external void _fgbAsyncAddError(int arg0, int arg1, int argPtr);

@ffi.Native<_FgbRetAddError Function(ffi.Uint64)>(symbol: "fgbasyncres_add_error")
external _FgbRetAddError _fgbAsyncResAddError(int arg0);

final class _FgbRetNewObj extends ffi.Struct {
  external ffi.Pointer<ffi.Void> res;
  external ffi.Pointer<ffi.Void> err;
}

@ffi.Native<_FgbRetNewObj Function(ffi.Pointer<ffi.Void>, ffi.Int)>(symbol: "fgb_new_obj")
external _FgbRetNewObj _fgbNewObj(ffi.Pointer<ffi.Void> arg0, int arg1);

@ffi.Native<ffi.Void Function(ffi.Pointer<ffi.Void>, ffi.Int, ffi.Uint64)>(symbol: "fgbasync_new_obj")
external void _fgbAsyncNewObj(ffi.Pointer<ffi.Void> arg0, int arg1, int argPtr);

@ffi.Native<_FgbRetNewObj Function(ffi.Uint64)>(symbol: "fgbasyncres_new_obj")
external _FgbRetNewObj _fgbAsyncResNewObj(int arg0);

final class _FgbRetModifyObj extends ffi.Struct {
  external ffi.Pointer<ffi.Void> err;
}

@ffi.Native<_FgbRetModifyObj Function(ffi.Pointer<ffi.Void>)>(symbol: "fgb_modify_obj")
external _FgbRetModifyObj _fgbModifyObj(ffi.Pointer<ffi.Void> arg0);

@ffi.Native<ffi.Void Function(ffi.Pointer<ffi.Void>, ffi.Uint64)>(symbol: "fgbasync_modify_obj")
external void _fgbAsyncModifyObj(ffi.Pointer<ffi.Void> arg0, int argPtr);

@ffi.Native<_FgbRetModifyObj Function(ffi.Uint64)>(symbol: "fgbasyncres_modify_obj")
external _FgbRetModifyObj _fgbAsyncResModifyObj(int arg0);

final class _FgbRetFormatObj extends ffi.Struct {
  external ffi.Pointer<ffi.Void> res;
  external ffi.Pointer<ffi.Void> err;
}

@ffi.Native<_FgbRetFormatObj Function(ffi.Pointer<ffi.Void>)>(symbol: "fgb_format_obj")
external _FgbRetFormatObj _fgbFormatObj(ffi.Pointer<ffi.Void> arg0);

@ffi.Native<ffi.Void Function(ffi.Pointer<ffi.Void>, ffi.Uint64)>(symbol: "fgbasync_format_obj")
external void _fgbAsyncFormatObj(ffi.Pointer<ffi.Void> arg0, int argPtr);

@ffi.Native<_FgbRetFormatObj Function(ffi.Uint64)>(symbol: "fgbasyncres_format_obj")
external _FgbRetFormatObj _fgbAsyncResFormatObj(int arg0);

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

  @override
  int add(int a, int b) {
    var __Dart__a = a;
    var __Dart__b = b;
    return _processAdd(_fgbAdd(__Dart__a, __Dart__b));
  }

  @override
  Future<int> addAsync(int a, int b) async {
    var __Dart__a = a;
    var __Dart__b = b;
    var __DartRecv__ = ReceivePort('AsyncRecv(add)');
    _fgbAsyncAdd(__Dart__a, __Dart__b, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processAdd(_fgbAsyncResAdd(__DartMsg__[0]));
  }

  int _processAdd(_FgbRetAdd res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
    return res.res;
  }

  @override
  Point addPoints(Point a, Point b) {
    var __Dart__a = _mapFromPoint(a);
    var __Dart__b = _mapFromPoint(b);
    return _processAddPoints(_fgbAddPoints(__Dart__a, __Dart__b));
  }

  @override
  Future<Point> addPointsAsync(Point a, Point b) async {
    var __Dart__a = _mapFromPoint(a);
    var __Dart__b = _mapFromPoint(b);
    var __DartRecv__ = ReceivePort('AsyncRecv(addPoints)');
    _fgbAsyncAddPoints(__Dart__a, __Dart__b, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processAddPoints(_fgbAsyncResAddPoints(__DartMsg__[0]));
  }

  Point _processAddPoints(_FgbRetAddPoints res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
    return _mapToPoint(res.res);
  }

  @override
  int addError(int a, int b) {
    var __Dart__a = a;
    var __Dart__b = b;
    return _processAddError(_fgbAddError(__Dart__a, __Dart__b));
  }

  @override
  Future<int> addErrorAsync(int a, int b) async {
    var __Dart__a = a;
    var __Dart__b = b;
    var __DartRecv__ = ReceivePort('AsyncRecv(addError)');
    _fgbAsyncAddError(__Dart__a, __Dart__b, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processAddError(_fgbAsyncResAddError(__DartMsg__[0]));
  }

  int _processAddError(_FgbRetAddError res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
    return res.res;
  }

  @override
  Obj newObj(String name, int other) {
    var __Dart__name = _mapFromString(name);
    var __Dart__other = other;
    return _processNewObj(_fgbNewObj(__Dart__name, __Dart__other));
  }

  @override
  Future<Obj> newObjAsync(String name, int other) async {
    var __Dart__name = _mapFromString(name);
    var __Dart__other = other;
    var __DartRecv__ = ReceivePort('AsyncRecv(newObj)');
    _fgbAsyncNewObj(__Dart__name, __Dart__other, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processNewObj(_fgbAsyncResNewObj(__DartMsg__[0]));
  }

  Obj _processNewObj(_FgbRetNewObj res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
    return _mapToObj(res.res);
  }

  @override
  void modifyObj(Obj o) {
    var __Dart__o = _mapFromObj(o);
    _processModifyObj(_fgbModifyObj(__Dart__o));
  }

  @override
  Future<void> modifyObjAsync(Obj o) async {
    var __Dart__o = _mapFromObj(o);
    var __DartRecv__ = ReceivePort('AsyncRecv(modifyObj)');
    _fgbAsyncModifyObj(__Dart__o, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    _processModifyObj(_fgbAsyncResModifyObj(__DartMsg__[0]));
  }

  void _processModifyObj(_FgbRetModifyObj res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
  }

  @override
  String formatObj(Obj o) {
    var __Dart__o = _mapFromObj(o);
    return _processFormatObj(_fgbFormatObj(__Dart__o));
  }

  @override
  Future<String> formatObjAsync(Obj o) async {
    var __Dart__o = _mapFromObj(o);
    var __DartRecv__ = ReceivePort('AsyncRecv(formatObj)');
    _fgbAsyncFormatObj(__Dart__o, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processFormatObj(_fgbAsyncResFormatObj(__DartMsg__[0]));
  }

  String _processFormatObj(_FgbRetFormatObj res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
    return _mapToString(res.res);
  }


  Point _mapToPoint(_FgbCPoint from) {
    return Point(from.x, from.y, _mapToString(from.name));
  }

  _FgbCPoint _mapFromPoint(Point from) {
    var res = _fgbEmptyPoint();
    res.x = from.x;
    res.y = from.y;
    res.name = _mapFromString(from.name);
    return res;
  }

  Obj _mapToObj(ffi.Pointer<ffi.Void> from) {
    var res = _FfiObj(from);
    _pinFinalizer.attach(res, from);
    return res;
  }

  ffi.Pointer<ffi.Void> _mapFromObj(Obj from) {
    if (from is! _FfiObj) {
      throw 'Mismatched reference struct instance type';
    }

    return from.id;
  }

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
