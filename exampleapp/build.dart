import 'package:go_native_toolchain/go_native_toolchain.dart';
import 'package:logging/logging.dart';
import 'package:native_assets_cli/native_assets_cli.dart';

const packageName = 'exampleapp';

void main(List<String> args) async {
  final buildConfig = await BuildConfig.fromArgs(args);
  final buildOutput = BuildOutput();

  final gobuulder = GoBuilder(
    name: packageName,
    assetId: 'package:$packageName/bridge.gen.dart',
    bridgePath: 'go/bridge',
    sources: ['go/']
  );

  await gobuulder.run(
    buildConfig: buildConfig,
    buildOutput: buildOutput,
    logger: Logger('')..onRecord.listen((record) => print(record.message)),
  );
  await buildOutput.writeToFile(outDir: buildConfig.outDir);
}
