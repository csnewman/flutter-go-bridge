# flutter-go-bridge

Reuse existing Go libraries in Flutter by calling Go from Dart using auto generated bindings.

### Features

- Windows, Linux and Android support.
    - iOS and macOS untested but likely works.
- Sync and Async calls from Dart.
- Passing primitives (various `int`'s), `string`'s and nested structures.
- Automatic `error` to `Exception` mapping.

### Getting Started

1. Install fgb
   ```
   go get -u github.com/csnewman/flutter-go-bridge
   ```

2. Create a wrapper around your Go code.
    - See [Go wrapper](#go-wrapper) below.
    - See [example](https://github.com/csnewman/flutter-go-bridge/blob/master/example/example.go) for an example.

3. Generate bindings
    - `go generate` in the directory containing your wrapper.

4. Use the generated bindings from your Flutter/Dart application

5. Automate library building by integrating into flutter build.
    - See [Platform building](#platform-building) below.
    - See `exampleapp` for an example.

### Go wrapper

The bridge does not support all Go constructs and instead relies on a subset. This allows for a more ergonomic wrapper.

Therefore, you should create a Go package containing a wrapper around more complex Go libraries. The generator will scan
this package and generate a `bridge` package containing the necessary FFI bridging logic.

See [example](https://github.com/csnewman/flutter-go-bridge/blob/master/example/example.go) for a complete example.

The process is as follows:

1. `mkdir example`
2. Create `example/example.go`, containing the following:
   ```go
    //go:generate go run github.com/csnewman/flutter-go-bridge/cmd/flutter-go-bridge generate --src example.go --go bridge/bridge.gen.go --dart ../exampleapp/lib/bridge.gen.dart
    package example
   ```
   Replace the `--dart` path, with a suitable location inside your flutter applications `lib` directory.
3. Write a simplified wrapper, using the [template mappings](#mappings) below.
4. Invoke `go generate` to automatically regenerate the `bridge.gen.go` and the corresponding `bridge.gen.dart` file.

### Dart setup

From Dart:

1. Import generated bridge
   ```dart
   import 'bridge.gen.dart';
   ```
2. Load library
   ```dart
   var lib = DynamicLibrary.open(Platform.isWindows ? "example.dll" : "libexample.so");
   var bridge = Bridge.open(lib);
   ```
   NOTE: The `open` function is planned to be replaced with an autoconfiguring version, which will automatically select
   the correct file.
3. Invoke functions
   ```dart
   bridge.example();
   await bridge.exampleAsync();
   ```

### Mappings

The following patterns are supported:

#### Function calls

```go
func Simple(a int, b int) int {
    return a + b
}
```

Will produce:

```dart
int simple(int a, int b);

Future<int> simpleAsync(int a, int b);
```

#### Function calls, with errors

```go
func SimpleError(a int, b int) (int, error) {
    return 0, errors.New("an example")
}
```

Will produce:

```dart
int simpleError(int a, int b);

Future<int> simpleErrorAsync(int a, int b);
```

If the Go function returns an error, the `simpleError` function will throw a `BridgeException`.

#### Struct passing

```go
type ExampleStruct struct {
    A int
    B string
}

func StructPassing(val ExampleStruct) {
}
```

Will produce:

```dart
final class ExampleStruct {
    int a;
    String b;

    ExampleStruct(this.a, this.b);
}
```
and
```dart
void structPassing(ExampleStruct val);

Future<void> structPassingAsync(ExampleStruct val);
```

Structs passed in this manner will be passed by `value`, meaning the contents will be copied.

### Platform building

The platforms supported by `flutter` use various build tooling, which complicates integrating Go into the build
pipeline. The approach outlined below attempts to unify the configuration where possible.

A complete example can be seen in `exampleapp`.

#### Shared setup

Inside your flutter project, create a directory called `golib` next to the `pubspec.yaml` file.

Inside `golib`, create a file called `CMakeLists.txt` and copy the contents from
the [exampleapp](https://github.com/csnewman/flutter-go-bridge/blob/master/exampleapp/golib/CMakeLists.txt).

Update the `LIBNAME` variable to match the final name you want the library to have on disk, such as `example`
becomes `libexample.so` and `example.dll`.

`GOSRC` should point to the directory containing your `go.mod` file.

`GOMAIN` should be the relative path to the generated Go bridge from the `GOSRC` directory.

#### Linux

Inside `linux/CMakeLists.txt`, after `include(flutter/generated_plugins.cmake)`, insert:

```
# go-flutter-bridge
add_subdirectory(../golib gobuild)
add_dependencies(${BINARY_NAME} libexample)
list(APPEND PLUGIN_BUNDLED_LIBRARIES ${CMAKE_CURRENT_BINARY_DIR}/gobuild/libexample.so)
```

Replacing `example` with your `LIBNAME`.

See the [exampleapp](https://github.com/csnewman/flutter-go-bridge/blob/master/exampleapp/linux/CMakeLists.txt) for
an example.

#### Windows

Inside `windows/CMakeLists.txt`, after `include(flutter/generated_plugins.cmake)`, insert:

```
# go-flutter-bridge
add_subdirectory(../golib gobuild)
add_dependencies(${BINARY_NAME} libexample)
list(APPEND PLUGIN_BUNDLED_LIBRARIES ${CMAKE_CURRENT_BINARY_DIR}/gobuild/example.dll)
```

Replacing `example` with your `LIBNAME`.

Flutter by default will use MSVCC to compile the desktop application. MSVCC is currently not compatible with CGO.

You will need to install and configure a compatible CGO compiler, such as `mingw`:
1. Download `x86_64-[...]-release-posix-seh-msvcrt-[...].7z` from [GitHub](https://github.com/niXman/mingw-builds-binaries/releases).
2. Extract to `C:\mingw`
3. Add `C:\mingw\bin` to your system path.

Alternatively, `zig c` may be suitable. 

See the [exampleapp](https://github.com/csnewman/flutter-go-bridge/blob/master/exampleapp/windows/CMakeLists.txt) for
an example.

#### Android

Add the following to `android/app/build.gradle`:

Under `defaultConfig`:

```
externalNativeBuild {
    cmake {
        targets "libexample"
    }
}
```

Replacing `example` with your `LIBNAME`.

Under `android`:

```
externalNativeBuild {
    cmake {
        path "../../golib/CMakeLists.txt"
    }
}
```

See the [exampleapp](https://github.com/csnewman/flutter-go-bridge/blob/master/exampleapp/android/app/build.gradle) for
an example.

#### iOS/macOS

While the generated code should be compatible with iOS and macOS, no guidelines are available due to a lack of access to
such hardware. If you have integrated with iOS/macOS, please consider opening a pull request with instructions.

#### Manual building

If you are not using Flutter, or wish to customise the build process, you can manually build the Go library and bundle
with your application as necessary:

```
CGO_ENABLED=1 go build -buildmode=c-shared -o libexample.so example/bridge/bridge.gen.go
```

You can specify `GOOS` and `GOARCH`.
