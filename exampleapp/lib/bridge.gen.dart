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

typedef _FgbEmptyPoint = _FgbCPoint Function();

typedef _FgbDefInit = ffi.Pointer<ffi.Void> Function(ffi.Pointer<ffi.Void>);
typedef _FgbDefIntCAlloc = ffi.Pointer Function(ffi.IntPtr);
typedef _FgbDefIntDartAlloc = ffi.Pointer Function(int);
typedef _FgbDefIntCFree = ffi.Void Function(ffi.Pointer);
typedef _FgbDefIntDartFree = void Function(ffi.Pointer);
typedef _NativeFinalizerFunctionC = ffi.Void Function(ffi.Pointer<ffi.Void>);
typedef _NativeFinalizerFunctionDart = void Function(ffi.Pointer<ffi.Void>);

class _GoAllocator implements ffi.Allocator {
  final _FgbDefIntDartAlloc _allocPtr;
  final _FgbDefIntDartFree _freePtr;

  const _GoAllocator(this._allocPtr, this._freePtr);

  @override
  ffi.Pointer<T> allocate<T extends ffi.NativeType>(int byteCount, {int? alignment}) {
    ffi.Pointer<T> result = _allocPtr(byteCount).cast();
    if (result.address == 0) {
      throw ArgumentError('Could not allocate $byteCount bytes.');
    }
    return result;
  }

  @override
  void free(ffi.Pointer pointer) {
  	_freePtr(pointer);
  }
}

final class _FgbRetAdd extends ffi.Struct {
  @ffi.Int()
  external int res;
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDartAdd = _FgbRetAdd Function(int, int);
typedef _FgbDefCAdd = _FgbRetAdd Function(ffi.Int, ffi.Int);
typedef _FgbAsyncDefDartAdd = void Function(int, int, int);
typedef _FgbAsyncDefCAdd = ffi.Void Function(ffi.Int, ffi.Int, ffi.Uint64);
typedef _FgbAsyncResDefDartAdd = _FgbRetAdd Function(int);
typedef _FgbAsyncResDefCAdd = _FgbRetAdd Function(ffi.Uint64);

final class _FgbRetAddPoints extends ffi.Struct {
  external _FgbCPoint res;
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDartAddPoints = _FgbRetAddPoints Function(_FgbCPoint, _FgbCPoint);
typedef _FgbDefCAddPoints = _FgbRetAddPoints Function(_FgbCPoint, _FgbCPoint);
typedef _FgbAsyncDefDartAddPoints = void Function(_FgbCPoint, _FgbCPoint, int);
typedef _FgbAsyncDefCAddPoints = ffi.Void Function(_FgbCPoint, _FgbCPoint, ffi.Uint64);
typedef _FgbAsyncResDefDartAddPoints = _FgbRetAddPoints Function(int);
typedef _FgbAsyncResDefCAddPoints = _FgbRetAddPoints Function(ffi.Uint64);

final class _FgbRetAddError extends ffi.Struct {
  @ffi.Int()
  external int res;
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDartAddError = _FgbRetAddError Function(int, int);
typedef _FgbDefCAddError = _FgbRetAddError Function(ffi.Int, ffi.Int);
typedef _FgbAsyncDefDartAddError = void Function(int, int, int);
typedef _FgbAsyncDefCAddError = ffi.Void Function(ffi.Int, ffi.Int, ffi.Uint64);
typedef _FgbAsyncResDefDartAddError = _FgbRetAddError Function(int);
typedef _FgbAsyncResDefCAddError = _FgbRetAddError Function(ffi.Uint64);

final class _FgbRetNewObj extends ffi.Struct {
  external ffi.Pointer<ffi.Void> res;
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDartNewObj = _FgbRetNewObj Function(ffi.Pointer<ffi.Void>, int);
typedef _FgbDefCNewObj = _FgbRetNewObj Function(ffi.Pointer<ffi.Void>, ffi.Int);
typedef _FgbAsyncDefDartNewObj = void Function(ffi.Pointer<ffi.Void>, int, int);
typedef _FgbAsyncDefCNewObj = ffi.Void Function(ffi.Pointer<ffi.Void>, ffi.Int, ffi.Uint64);
typedef _FgbAsyncResDefDartNewObj = _FgbRetNewObj Function(int);
typedef _FgbAsyncResDefCNewObj = _FgbRetNewObj Function(ffi.Uint64);

final class _FgbRetModifyObj extends ffi.Struct {
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDartModifyObj = _FgbRetModifyObj Function(ffi.Pointer<ffi.Void>);
typedef _FgbDefCModifyObj = _FgbRetModifyObj Function(ffi.Pointer<ffi.Void>);
typedef _FgbAsyncDefDartModifyObj = void Function(ffi.Pointer<ffi.Void>, int);
typedef _FgbAsyncDefCModifyObj = ffi.Void Function(ffi.Pointer<ffi.Void>, ffi.Uint64);
typedef _FgbAsyncResDefDartModifyObj = _FgbRetModifyObj Function(int);
typedef _FgbAsyncResDefCModifyObj = _FgbRetModifyObj Function(ffi.Uint64);

final class _FgbRetFormatObj extends ffi.Struct {
  external ffi.Pointer<ffi.Void> res;
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDartFormatObj = _FgbRetFormatObj Function(ffi.Pointer<ffi.Void>);
typedef _FgbDefCFormatObj = _FgbRetFormatObj Function(ffi.Pointer<ffi.Void>);
typedef _FgbAsyncDefDartFormatObj = void Function(ffi.Pointer<ffi.Void>, int);
typedef _FgbAsyncDefCFormatObj = ffi.Void Function(ffi.Pointer<ffi.Void>, ffi.Uint64);
typedef _FgbAsyncResDefDartFormatObj = _FgbRetFormatObj Function(int);
typedef _FgbAsyncResDefCFormatObj = _FgbRetFormatObj Function(ffi.Uint64);

final class _FfiBridge implements Bridge {
  late _GoAllocator _allocator;
  late _FgbDefDartAdd _addPtr;
  late _FgbAsyncDefDartAdd _addPtrAsync;
  late _FgbAsyncResDefDartAdd _addPtrAsyncRes;
  late _FgbDefDartAddPoints _addPointsPtr;
  late _FgbAsyncDefDartAddPoints _addPointsPtrAsync;
  late _FgbAsyncResDefDartAddPoints _addPointsPtrAsyncRes;
  late _FgbDefDartAddError _addErrorPtr;
  late _FgbAsyncDefDartAddError _addErrorPtrAsync;
  late _FgbAsyncResDefDartAddError _addErrorPtrAsyncRes;
  late _FgbDefDartNewObj _newObjPtr;
  late _FgbAsyncDefDartNewObj _newObjPtrAsync;
  late _FgbAsyncResDefDartNewObj _newObjPtrAsyncRes;
  late _FgbDefDartModifyObj _modifyObjPtr;
  late _FgbAsyncDefDartModifyObj _modifyObjPtrAsync;
  late _FgbAsyncResDefDartModifyObj _modifyObjPtrAsyncRes;
  late _FgbDefDartFormatObj _formatObjPtr;
  late _FgbAsyncDefDartFormatObj _formatObjPtrAsync;
  late _FgbAsyncResDefDartFormatObj _formatObjPtrAsyncRes;
  late _FgbEmptyPoint _emptyPointPtr;
  late ffi.Pointer<ffi.NativeFinalizerFunction> _freeObjPtr;
  late ffi.NativeFinalizer _objFinalizer;

  _FfiBridge(ffi.DynamicLibrary lib) {
    var allocPtr = lib.lookupFunction<_FgbDefIntCAlloc, _FgbDefIntDartAlloc>("fgbinternal_alloc");
    var freePtr = lib.lookupFunction<_FgbDefIntCFree, _FgbDefIntDartFree>("fgbinternal_free");
    _allocator = _GoAllocator(allocPtr, freePtr);

    var initPtr = lib.lookupFunction<_FgbDefInit, _FgbDefInit>("fgbinternal_init");
    var initRes = initPtr(ffi.NativeApi.initializeApiDLData);
    if (initRes != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(initRes.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }

    _addPtr = lib.lookupFunction<_FgbDefCAdd, _FgbDefDartAdd>("fgb_add");
    _addPtrAsync = lib.lookupFunction<_FgbAsyncDefCAdd, _FgbAsyncDefDartAdd>("fgbasync_add");
    _addPtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefCAdd, _FgbAsyncResDefDartAdd>("fgbasyncres_add");
    _addPointsPtr = lib.lookupFunction<_FgbDefCAddPoints, _FgbDefDartAddPoints>("fgb_add_points");
    _addPointsPtrAsync = lib.lookupFunction<_FgbAsyncDefCAddPoints, _FgbAsyncDefDartAddPoints>("fgbasync_add_points");
    _addPointsPtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefCAddPoints, _FgbAsyncResDefDartAddPoints>("fgbasyncres_add_points");
    _addErrorPtr = lib.lookupFunction<_FgbDefCAddError, _FgbDefDartAddError>("fgb_add_error");
    _addErrorPtrAsync = lib.lookupFunction<_FgbAsyncDefCAddError, _FgbAsyncDefDartAddError>("fgbasync_add_error");
    _addErrorPtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefCAddError, _FgbAsyncResDefDartAddError>("fgbasyncres_add_error");
    _newObjPtr = lib.lookupFunction<_FgbDefCNewObj, _FgbDefDartNewObj>("fgb_new_obj");
    _newObjPtrAsync = lib.lookupFunction<_FgbAsyncDefCNewObj, _FgbAsyncDefDartNewObj>("fgbasync_new_obj");
    _newObjPtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefCNewObj, _FgbAsyncResDefDartNewObj>("fgbasyncres_new_obj");
    _modifyObjPtr = lib.lookupFunction<_FgbDefCModifyObj, _FgbDefDartModifyObj>("fgb_modify_obj");
    _modifyObjPtrAsync = lib.lookupFunction<_FgbAsyncDefCModifyObj, _FgbAsyncDefDartModifyObj>("fgbasync_modify_obj");
    _modifyObjPtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefCModifyObj, _FgbAsyncResDefDartModifyObj>("fgbasyncres_modify_obj");
    _formatObjPtr = lib.lookupFunction<_FgbDefCFormatObj, _FgbDefDartFormatObj>("fgb_format_obj");
    _formatObjPtrAsync = lib.lookupFunction<_FgbAsyncDefCFormatObj, _FgbAsyncDefDartFormatObj>("fgbasync_format_obj");
    _formatObjPtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefCFormatObj, _FgbAsyncResDefDartFormatObj>("fgbasyncres_format_obj");

    _emptyPointPtr = lib.lookupFunction<_FgbEmptyPoint, _FgbEmptyPoint>("fgbempty_point");

    _freeObjPtr = lib.lookup<ffi.NativeFinalizerFunction>("fgbfree_obj");
    _objFinalizer = ffi.NativeFinalizer(_freeObjPtr);
  }

  @override
  int add(int a, int b) {
    var __Dart__a = a;
    
    var __Dart__b = b;
    
    return _processAdd(_addPtr(__Dart__a, __Dart__b));
  }

  @override
  Future<int> addAsync(int a, int b) async {
    var __Dart__a = a;
    
    var __Dart__b = b;
    
    var __DartRecv__ = ReceivePort('AsyncRecv(add)');
    _addPtrAsync(__Dart__a, __Dart__b, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processAdd(_addPtrAsyncRes(__DartMsg__[0]));
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
    
    return _processAddPoints(_addPointsPtr(__Dart__a, __Dart__b));
  }

  @override
  Future<Point> addPointsAsync(Point a, Point b) async {
    var __Dart__a = _mapFromPoint(a);
    
    var __Dart__b = _mapFromPoint(b);
    
    var __DartRecv__ = ReceivePort('AsyncRecv(addPoints)');
    _addPointsPtrAsync(__Dart__a, __Dart__b, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processAddPoints(_addPointsPtrAsyncRes(__DartMsg__[0]));
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
    
    return _processAddError(_addErrorPtr(__Dart__a, __Dart__b));
  }

  @override
  Future<int> addErrorAsync(int a, int b) async {
    var __Dart__a = a;
    
    var __Dart__b = b;
    
    var __DartRecv__ = ReceivePort('AsyncRecv(addError)');
    _addErrorPtrAsync(__Dart__a, __Dart__b, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processAddError(_addErrorPtrAsyncRes(__DartMsg__[0]));
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
    
    return _processNewObj(_newObjPtr(__Dart__name, __Dart__other));
  }

  @override
  Future<Obj> newObjAsync(String name, int other) async {
    var __Dart__name = _mapFromString(name);
    
    var __Dart__other = other;
    
    var __DartRecv__ = ReceivePort('AsyncRecv(newObj)');
    _newObjPtrAsync(__Dart__name, __Dart__other, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processNewObj(_newObjPtrAsyncRes(__DartMsg__[0]));
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
    
    _processModifyObj(_modifyObjPtr(__Dart__o));
  }

  @override
  Future<void> modifyObjAsync(Obj o) async {
    var __Dart__o = _mapFromObj(o);
    
    var __DartRecv__ = ReceivePort('AsyncRecv(modifyObj)');
    _modifyObjPtrAsync(__Dart__o, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    _processModifyObj(_modifyObjPtrAsyncRes(__DartMsg__[0]));
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
    
    return _processFormatObj(_formatObjPtr(__Dart__o));
  }

  @override
  Future<String> formatObjAsync(Obj o) async {
    var __Dart__o = _mapFromObj(o);
    
    var __DartRecv__ = ReceivePort('AsyncRecv(formatObj)');
    _formatObjPtrAsync(__Dart__o, __DartRecv__.sendPort.nativePort);
    var __DartMsg__ = await __DartRecv__.first;
    __DartRecv__.close();
    return _processFormatObj(_formatObjPtrAsyncRes(__DartMsg__[0]));
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
    var res = _emptyPointPtr();
    res.x = from.x;
    res.y = from.y;
    res.name = _mapFromString(from.name);
    return res;
  }

  Obj _mapToObj(ffi.Pointer<ffi.Void> from) {
    var res = _FfiObj(from);
    _objFinalizer.attach(res, from);
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
