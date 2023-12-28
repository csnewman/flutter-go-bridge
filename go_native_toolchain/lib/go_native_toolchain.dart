import 'dart:io';
import 'dart:math';

import 'package:logging/logging.dart';
import 'package:native_assets_cli/native_assets_cli.dart';

// Not ideal:
import 'package:native_toolchain_c/src/cbuilder/compiler_resolver.dart';
import 'package:native_toolchain_c/src/cbuilder/run_cbuilder.dart';
import 'package:native_toolchain_c/src/utils/run_process.dart';

/// Golang artifact builder
class GoBuilder {
  /// Name of the library to build.
  ///
  /// File will be placed in [BuildConfig.outDir].
  final String name;

  /// Asset identifier.
  ///
  /// Used to output the [BuildOutput.assets].
  ///
  /// If omitted, no asset will be added to the build output.
  final String? assetId;

  /// Golang bridge to build.
  ///
  /// Resolved against [BuildConfig.packageRoot].
  ///
  /// Used to output the [BuildOutput.dependencies].
  final String bridgePath;

  /// Sources to build the library.
  ///
  /// Resolved against [BuildConfig.packageRoot].
  ///
  /// Used to output the [BuildOutput.dependencies].
  final List<String> sources;

  /// The dart files involved in building this artifact.
  ///
  /// Resolved against [BuildConfig.packageRoot].
  ///
  /// Used to output the [BuildOutput.dependencies].
  final List<String> dartBuildFiles;

  GoBuilder({
    required this.name,
    required this.assetId,
    required this.bridgePath,
    this.sources = const [],
    this.dartBuildFiles = const ['build.dart'],
  });

  /// Runs the Golang Compiler.
  ///
  /// Completes with an error if the build fails.
  Future<void> run({
    required BuildConfig buildConfig,
    required BuildOutput buildOutput,
    required Logger? logger,
  }) async {
    final outDir = buildConfig.outDir;
    final packageRoot = buildConfig.packageRoot;
    await Directory.fromUri(outDir).create(recursive: true);
    var linkMode = buildConfig.linkModePreference.preferredLinkMode;
    final libUri = outDir.resolve(
      buildConfig.targetOs.libraryFileName(name, linkMode),
    );
    final bridgePath = packageRoot.resolveUri(Uri.file(this.bridgePath));
    final resolver = CompilerResolver(buildConfig: buildConfig, logger: logger);

    if (!buildConfig.dryRun) {
      final compiler = await resolver.resolveCompiler();
      final target = buildConfig.target;
      var buildMode = linkMode == LinkMode.static ? "c-archive" : "c-shared";
      var goLib = libUri;

      final env = {
        "CGO_ENABLED": "1",
      };

      String buildArch;
      switch(target.architecture) {
        case Architecture.arm:
          buildArch = "arm";
        case Architecture.arm64:
          buildArch = "arm64";
        case Architecture.ia32:
          buildArch = "386";
        case Architecture.x64:
          buildArch = "amd64";
        default:
          throw Exception("Unknown architecture");
      }

      String buildOs;
      switch(target.os) {
        case OS.windows:
          buildOs = "windows";

        case OS.linux:
          buildOs = "linux";

        case OS.android:
          buildOs = "android";

          // The Android Gradle plugin does not honor API level 19 and 20 when
          // invoking clang. Mimic that behavior here.
          // See https://github.com/dart-lang/native/issues/171.
          final minimumApi = target == Target.androidRiscv64 ? 35 : 21;
          final targetAndroidNdkApi = max(buildConfig.targetAndroidNdkApi!, minimumApi);

          final cc = compiler.uri.resolve('./${RunCBuilder.androidNdkClangTargetFlags[target]!}$targetAndroidNdkApi-clang');
          env["CC"] = cc.toFilePath();

        case OS.iOS:
          buildOs = "ios";
          buildMode = "c-archive";
          goLib = outDir.resolve('out.o');

        case OS.macOS:
          buildOs = "darwin";

        default:
          throw Exception("Unknown os");
      }

      env["GOOS"] = buildOs;
      env["GOARCH"] = buildArch;

      await runProcess(
        executable: Uri.file("go"),
        environment: env,
        arguments: [
          "build",
          "-buildmode=$buildMode",
          if (buildConfig.buildMode != BuildMode.debug) '-ldflags=-s -w',
          '-o',
          goLib.toFilePath(),
          bridgePath.toFilePath(),
        ],
        logger: logger,
        captureOutput: false,
        throwOnUnexpectedExitCode: true,
      );

      if (target.os == OS.iOS) {
        //xcrun -sdk iphoneos clang -arch armv7 -fpic -shared -Wl,-all_load
        // libmystatic.a -framework Corefoundation -o libmydynamic.dylib
        await runProcess(
          executable: compiler.uri,
          arguments: [
            '-fpic',
            '-shared',
            '-Wl,-all_load,-force_load',
            goLib.toFilePath(),
            '-framework',
            'CoreFoundation',
            '-o',
            libUri.toFilePath(),
          ],
          logger: logger,
          captureOutput: false,
          throwOnUnexpectedExitCode: true,
        );
      }
    }

    if (assetId != null) {
      final targets = [
        if (!buildConfig.dryRun)
          buildConfig.target
        else
          for (final target in Target.values)
            if (target.os == buildConfig.targetOs) target
      ];
      for (final target in targets) {
        buildOutput.assets.add(Asset(
          id: assetId!,
          linkMode: linkMode,
          target: target,
          path: AssetAbsolutePath(libUri),
        ));
      }
    }

    if (!buildConfig.dryRun) {
      final sources = [
        for (final source in this.sources)
          packageRoot.resolveUri(Uri.file(source)),
      ];
      final dartBuildFiles = [
        for (final source in this.dartBuildFiles) packageRoot.resolve(source),
      ];

      final sourceFiles = await Stream.fromIterable(sources)
          .asyncExpand(
            (path) => Directory(path.toFilePath())
            .list(recursive: true)
            .where((entry) => entry is File)
            .map((file) => file.uri),
      ).toList();

      buildOutput.dependencies.dependencies.addAll({
        ...sourceFiles,
        bridgePath,
        ...dartBuildFiles,
      });
    }
  }
}
