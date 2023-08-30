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
  late _FgbEmptyPoint _emptyPointPtr;

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

    _emptyPointPtr = lib.lookupFunction<_FgbEmptyPoint, _FgbEmptyPoint>("fgbempty_point");
  }

  @override
  int add(int a, int b) {
    var aDart = a;
    
    var bDart = b;
    
    return _processAdd(_addPtr(aDart, bDart));
  }

  @override
  Future<int> addAsync(int a, int b) async {
    var aDart = a;
    
    var bDart = b;
    
    var recv = ReceivePort('AsyncRecv(add)');
    _addPtrAsync(aDart, bDart, recv.sendPort.nativePort);
    var msg = await recv.first;
    recv.close();
    return _processAdd(_addPtrAsyncRes(msg[0]));
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
    var aDart = _mapFromPoint(a);
    
    var bDart = _mapFromPoint(b);
    
    return _processAddPoints(_addPointsPtr(aDart, bDart));
  }

  @override
  Future<Point> addPointsAsync(Point a, Point b) async {
    var aDart = _mapFromPoint(a);
    
    var bDart = _mapFromPoint(b);
    
    var recv = ReceivePort('AsyncRecv(addPoints)');
    _addPointsPtrAsync(aDart, bDart, recv.sendPort.nativePort);
    var msg = await recv.first;
    recv.close();
    return _processAddPoints(_addPointsPtrAsyncRes(msg[0]));
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
    var aDart = a;
    
    var bDart = b;
    
    return _processAddError(_addErrorPtr(aDart, bDart));
  }

  @override
  Future<int> addErrorAsync(int a, int b) async {
    var aDart = a;
    
    var bDart = b;
    
    var recv = ReceivePort('AsyncRecv(addError)');
    _addErrorPtrAsync(aDart, bDart, recv.sendPort.nativePort);
    var msg = await recv.first;
    recv.close();
    return _processAddError(_addErrorPtrAsyncRes(msg[0]));
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
