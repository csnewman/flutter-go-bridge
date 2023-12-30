# flutter-go-bridge

Reuse existing Go libraries in Flutter by calling Go from Dart using auto generated bindings.

### Features

- MacOS, iOS, Windows, Linux and Android support.
- Sync and Async calls from Dart.
- Passing primitives (various `int`'s), `string`'s and nested structures.
- Automatic `error` to `Exception` mapping.
- Basic object passing.

### Getting Started

1. Install fgb
   ```
   go get -u github.com/csnewman/flutter-go-bridge
   ```

2. Create a wrapper around your Go code.
    - See [Go wrapper](#go-wrapper) below.
    - See [example](https://github.com/csnewman/flutter-go-bridge/blob/master/exampleapp/go/example.go) for an example.

3. Generate bindings
    - Run `go generate` in the directory containing your wrapper.

4. Use the generated bindings from your Flutter/Dart application

5. Automate library building by integrating into flutter build.
    - See [Platform building](#platform-building) below.
    - See the `exampleapp` folder for a full example.

6. When modifying the Go code, you may need to call `flutter clean` to trigger a rebuild, dependent upon your Go source 
   location and configured source directories.

### Go wrapper

The bridge does not support all Go constructs and instead relies on a subset. This allows for a more ergonomic wrapper.

Therefore, you should create a Go package containing a wrapper around more complex Go libraries. The generator will scan
this package and generate a `bridge` package containing the necessary FFI bridging logic.

See [example](https://github.com/csnewman/flutter-go-bridge/blob/master/exampleapp/go/example.go) for a complete example.

The process is as follows:

1. `mkdir go` (inside the Flutter project)
2. Create `go/example.go`, containing the following:
   ```go
    //go:generate go run github.com/csnewman/flutter-go-bridge/cmd/flutter-go-bridge generate --src example.go --go bridge/bridge.gen.go --dart ../lib/bridge.gen.dart
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
   var bridge = Bridge.open();
   ```
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

Value structs can not contain private fields.

#### Objects

TODO: Document

### Platform building

The platforms supported by `flutter` use various build tooling, which complicates integrating Go into the build
pipeline. Originally this project had hooks into the build systems for Windows, Linux and Android, however this had
high maintenance and was not trivial to integrate into the Mac ecosystem.

Flutter (& Dart) currently have an experimental feature called
[Native Assets](https://github.com/flutter/flutter/issues/129757) which greatly simplifies the setup. This does however
mean that for now, this project requires using the flutter `master` channel.

#### Native Assets approach

A complete example can be seen in the `exampleapp` folder.

1. Switch to the `master` flutter channel
   ```bash
   flutter channel master
   ```
2. Enable the [Native Assets](https://github.com/flutter/flutter/issues/129757) experiment
   ```bash
   flutter config --enable-native-asset
   ```
3. Add the required dependencies to `pubspec.yaml`
   ```yaml
   cli_config: ^0.1.2
   logging: ^1.2.0
   native_assets_cli: ^0.3.2
   go_native_toolchain: ^0.0.1
   ffi: ^2.1.0
   ```
4. Fetch dependencies
   ```bash
   flutter pub get
   ```
5. Create a `build.dart` file
   ```dart
   import 'package:go_native_toolchain/go_native_toolchain.dart';
   import 'package:logging/logging.dart';
   import 'package:native_assets_cli/native_assets_cli.dart';
   
   const packageName = 'exampleapp';
   
   void main(List<String> args) async {
     final buildConfig = await BuildConfig.fromArgs(args);
     final buildOutput = BuildOutput();
     
     final gobuilder = GoBuilder(
       name: packageName,
       assetId: 'package:$packageName/bridge.gen.dart',
       bridgePath: 'go/bridge'
     );
   
     await gobuilder.run(
       buildConfig: buildConfig,
       buildOutput: buildOutput,
       logger: Logger('')..onRecord.listen((record) => print(record.message)),
     );
     await buildOutput.writeToFile(outDir: buildConfig.outDir);
   }
   ```
   The `assetId` path needs to match the location of the autogenerated `bridge.gen.dart` file, as flutter uses this
   internally to automate library resolution. You may need to specify a list of source directories to the `GoBuilder`
   to allow automatic rebuilding as necessary.

You should now be able to use your IDE and other tooling as usual.

#### Manual building

If you do not want to use the `master` channel or wish to customise the build process, you can manually build the Go
library and bundle with your application as necessary:

```bash
CGO_ENABLED=1 go build -buildmode=c-shared -o libexample.so example/bridge/bridge.gen.go
```

You can specify `GOOS` and `GOARCH` as necessary.
