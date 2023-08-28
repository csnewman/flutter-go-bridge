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

  void example(SomeVal v);

  Future<void> exampleAsync(SomeVal v);

  SomeVal other();

  Future<SomeVal> otherAsync();

  int callMe();

  Future<int> callMeAsync();
}

final class Inner {
    int a;
    String d;

    Inner(this.a, this.d);
}

final class _FgbCInner extends ffi.Struct {
  @ffi.Int()
  external int a;
  external ffi.Pointer<ffi.Void> d;
}

typedef _FgbEmptyInner = _FgbCInner Function();

final class SomeVal {
    int v1;
    int v2;
    Inner i1;

    SomeVal(this.v1, this.v2, this.i1);
}

final class _FgbCSomeVal extends ffi.Struct {
  @ffi.Int()
  external int v1;
  @ffi.Int()
  external int v2;
  external _FgbCInner i1;
}

typedef _FgbEmptySomeVal = _FgbCSomeVal Function();

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

final class _FgbRetExample extends ffi.Struct {
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDartExample = _FgbRetExample Function(_FgbCSomeVal);
typedef _FgbDefCExample = _FgbRetExample Function(_FgbCSomeVal);
typedef _FgbAsyncDefDartExample = void Function(_FgbCSomeVal, int);
typedef _FgbAsyncDefCExample = ffi.Void Function(_FgbCSomeVal, ffi.Uint64);
typedef _FgbAsyncResDefDartExample = _FgbRetExample Function(int);
typedef _FgbAsyncResDefCExample = _FgbRetExample Function(ffi.Uint64);

final class _FgbRetOther extends ffi.Struct {
  external _FgbCSomeVal res;
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDartOther = _FgbRetOther Function();
typedef _FgbDefCOther = _FgbRetOther Function();
typedef _FgbAsyncDefDartOther = void Function(int);
typedef _FgbAsyncDefCOther = ffi.Void Function(ffi.Uint64);
typedef _FgbAsyncResDefDartOther = _FgbRetOther Function(int);
typedef _FgbAsyncResDefCOther = _FgbRetOther Function(ffi.Uint64);

final class _FgbRetCallMe extends ffi.Struct {
  @ffi.Int()
  external int res;
  external ffi.Pointer<ffi.Void> err;
}

typedef _FgbDefDartCallMe = _FgbRetCallMe Function();
typedef _FgbDefCCallMe = _FgbRetCallMe Function();
typedef _FgbAsyncDefDartCallMe = void Function(int);
typedef _FgbAsyncDefCCallMe = ffi.Void Function(ffi.Uint64);
typedef _FgbAsyncResDefDartCallMe = _FgbRetCallMe Function(int);
typedef _FgbAsyncResDefCCallMe = _FgbRetCallMe Function(ffi.Uint64);

final class _FfiBridge implements Bridge {
  late _GoAllocator _allocator;
  late _FgbDefDartExample _examplePtr;
  late _FgbAsyncDefDartExample _examplePtrAsync;
  late _FgbAsyncResDefDartExample _examplePtrAsyncRes;
  late _FgbDefDartOther _otherPtr;
  late _FgbAsyncDefDartOther _otherPtrAsync;
  late _FgbAsyncResDefDartOther _otherPtrAsyncRes;
  late _FgbDefDartCallMe _callMePtr;
  late _FgbAsyncDefDartCallMe _callMePtrAsync;
  late _FgbAsyncResDefDartCallMe _callMePtrAsyncRes;
  late _FgbEmptyInner _emptyInnerPtr;
  late _FgbEmptySomeVal _emptySomeValPtr;

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

    _examplePtr = lib.lookupFunction<_FgbDefCExample, _FgbDefDartExample>("fgb_example");
    _examplePtrAsync = lib.lookupFunction<_FgbAsyncDefCExample, _FgbAsyncDefDartExample>("fgbasync_example");
    _examplePtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefCExample, _FgbAsyncResDefDartExample>("fgbasyncres_example");
    _otherPtr = lib.lookupFunction<_FgbDefCOther, _FgbDefDartOther>("fgb_other");
    _otherPtrAsync = lib.lookupFunction<_FgbAsyncDefCOther, _FgbAsyncDefDartOther>("fgbasync_other");
    _otherPtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefCOther, _FgbAsyncResDefDartOther>("fgbasyncres_other");
    _callMePtr = lib.lookupFunction<_FgbDefCCallMe, _FgbDefDartCallMe>("fgb_call_me");
    _callMePtrAsync = lib.lookupFunction<_FgbAsyncDefCCallMe, _FgbAsyncDefDartCallMe>("fgbasync_call_me");
    _callMePtrAsyncRes = lib.lookupFunction<_FgbAsyncResDefCCallMe, _FgbAsyncResDefDartCallMe>("fgbasyncres_call_me");

    _emptyInnerPtr = lib.lookupFunction<_FgbEmptyInner, _FgbEmptyInner>("fgbempty_inner");
    _emptySomeValPtr = lib.lookupFunction<_FgbEmptySomeVal, _FgbEmptySomeVal>("fgbempty_some_val");
  }

  @override
  void example(SomeVal v) {
    var vDart = _mapFromSomeVal(v);
    
    _processExample(_examplePtr(vDart));
  }

  @override
  Future<void> exampleAsync(SomeVal v) async {
    var vDart = _mapFromSomeVal(v);
    
    var recv = ReceivePort('AsyncRecv(example)');
    _examplePtrAsync(vDart, recv.sendPort.nativePort);
    var msg = await recv.first;
    recv.close();
    _processExample(_examplePtrAsyncRes(msg[0]));
  }

  void _processExample(_FgbRetExample res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
  }
  @override
  SomeVal other() {
    return _processOther(_otherPtr());
  }

  @override
  Future<SomeVal> otherAsync() async {
    var recv = ReceivePort('AsyncRecv(other)');
    _otherPtrAsync(recv.sendPort.nativePort);
    var msg = await recv.first;
    recv.close();
    return _processOther(_otherPtrAsyncRes(msg[0]));
  }

  SomeVal _processOther(_FgbRetOther res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
    return _mapToSomeVal(res.res);
  }
  @override
  int callMe() {
    return _processCallMe(_callMePtr());
  }

  @override
  Future<int> callMeAsync() async {
    var recv = ReceivePort('AsyncRecv(callMe)');
    _callMePtrAsync(recv.sendPort.nativePort);
    var msg = await recv.first;
    recv.close();
    return _processCallMe(_callMePtrAsyncRes(msg[0]));
  }

  int _processCallMe(_FgbRetCallMe res) {
    if (res.err != ffi.nullptr) {
      var errPtr = ffi.Pointer<Utf8>.fromAddress(res.err.address);
      var errMsg = errPtr.toDartString(); 
      _allocator.free(errPtr);

      throw BridgeException(errMsg);
    }
    return res.res;
  }

  Inner _mapToInner(_FgbCInner from) {
    return Inner(from.a, _mapToString(from.d));
  }

  _FgbCInner _mapFromInner(Inner from) {
    var res = _emptyInnerPtr();
    res.a = from.a;
    res.d = _mapFromString(from.d);
    return res;
  }

  SomeVal _mapToSomeVal(_FgbCSomeVal from) {
    return SomeVal(from.v1, from.v2, _mapToInner(from.i1));
  }

  _FgbCSomeVal _mapFromSomeVal(SomeVal from) {
    var res = _emptySomeValPtr();
    res.v1 = from.v1;
    res.v2 = from.v2;
    res.i1 = _mapFromInner(from.i1);
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
