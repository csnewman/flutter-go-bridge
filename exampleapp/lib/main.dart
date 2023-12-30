import 'package:flutter/material.dart';
import 'bridge.gen.dart';

Bridge? _bridge;

Bridge getBridge() {
  if (_bridge != null) {
    return _bridge!;
  }

  _bridge = Bridge.open();
  return _bridge!;
}

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'flutter-go-bridge Demo',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: const MyHomePage(title: 'flutter-go-bridge demo'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 1,
      child: Scaffold(
        appBar: AppBar(
          backgroundColor: Theme.of(context).colorScheme.inversePrimary,
          title: Text(widget.title),
          bottom: const TabBar(
            isScrollable: true,
            tabs: [
              Tab(text: "1: basics"),
            ],
          ),
        ),
        body: const TabBarView(
          children: [
            Example1(),
          ],
        ),
      ),
    );
  }
}

class Example1 extends StatefulWidget {
  const Example1({super.key});

  @override
  State<Example1> createState() => _Example1State();
}

class _Example1State extends State<Example1> {
  String _addSyncState = 'not called';
  String _addAsyncState = 'not called';
  String _addPointsSyncState = 'not called';
  String _addPointsAsyncState = 'not called';
  String _addErrorSyncState = 'not called';
  String _addErrorAsyncState = 'not called';
  String _objSyncState = 'not called';
  String _objAsyncState = 'not called';

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.start,
      children: <Widget>[
        Row(
          children: [
            OutlinedButton(
              onPressed: () {
                setState(() {
                  _addSyncState = getBridge().add(1, 2).toString();
                });
              },
              key: const Key('button-add-sync'),
              child: const Text("Add (sync)"),
            ),
            Expanded(
              child: Text(
                _addSyncState,
                key: const Key('value-add-sync'),
              ),
            ),
          ],
        ),
        Row(
          children: [
            OutlinedButton(
              onPressed: () {
                getBridge().addAsync(121321, 123123).then((value) {
                  setState(() {
                    _addAsyncState = value.toString();
                  });
                });
              },
              key: const Key('button-add-async'),
              child: const Text("Add (async)"),
            ),
            Expanded(
              child: Text(
                _addAsyncState,
                key: const Key('value-add-async'),
              ),
            ),
          ],
        ),
        Row(
          children: [
            OutlinedButton(
              onPressed: () {
                setState(() {
                  var a = Point(1, 2, "PointA");
                  var b = Point(456, 789, "PointB");
                  _addPointsSyncState = getBridge().addPoints(a, b).toString();
                });
              },
              key: const Key('button-add-points-sync'),
              child: const Text("AddPoints (sync)"),
            ),
            Expanded(
              child: Text(
                _addPointsSyncState,
                key: const Key('value-add-points-sync'),
              ),
            ),
          ],
        ),
        Row(
          children: [
            OutlinedButton(
              onPressed: () {
                var a = Point(1, 2, "PointA");
                var b = Point(456, 789, "PointB");
                getBridge().addPointsAsync(a, b).then((value) {
                  setState(() {
                    _addPointsAsyncState = value.toString();
                  });
                });
              },
              key: const Key('button-add-points-async'),
              child: const Text("AddPoints (async)"),
            ),
            Expanded(
              child: Text(
                _addPointsAsyncState,
                key: const Key('value-add-points-async'),
              ),
            ),
          ],
        ),
        Row(
          children: [
            OutlinedButton(
              onPressed: () {
                setState(() {
                  try {
                    var res = getBridge().addError(1, 1);
                    _addErrorSyncState = 'success=$res';
                  } catch (e) {
                    _addErrorSyncState = 'error=$e';
                  }
                });
              },
              key: const Key('button-add-errors-sync'),
              child: const Text("AddErrors (sync)"),
            ),
            Expanded(
              child: Text(
                _addErrorSyncState,
                key: const Key('value-add-errors-sync'),
              ),
            ),
          ],
        ),
        Row(
          children: [
            OutlinedButton(
              onPressed: () {
                getBridge().addErrorAsync(2, 2).then((value) {
                  setState(() {
                    _addErrorAsyncState = 'success=$value';
                  });
                }, onError: (e) {
                  setState(() {
                    _addErrorAsyncState = 'error=$e';
                  });
                });
              },
              key: const Key('button-add-errors-async'),
              child: const Text("AddError (async)"),
            ),
            Expanded(
              child: Text(
                _addErrorAsyncState,
                key: const Key('value-add-errors-async'),
              ),
            ),
          ],
        ),
        Row(
          children: [
            OutlinedButton(
              onPressed: () {
                setState(() {
                  var b = getBridge();
                  var o = b.newObj("test1", 100);
                  b.modifyObj(o);
                  _objSyncState = b.formatObj(o);
                });
              },
              key: const Key('button-obj-sync'),
              child: const Text("Obj format (sync)"),
            ),
            Expanded(
              child: Text(
                _objSyncState,
                key: const Key('value-obj-sync'),
              ),
            ),
          ],
        ),
        Row(
          children: [
            OutlinedButton(
              onPressed: () {
                getBridge().newObjAsync("test2", 123).then(
                      (o) => {
                        getBridge().modifyObjAsync(o).then((_) => {
                              getBridge().formatObjAsync(o).then((value) => {
                                    setState(() {
                                      _objAsyncState = value;
                                    })
                                  }),
                            }),
                      },
                    );
              },
              key: const Key('button-obj-async'),
              child: const Text("Obj format (async)"),
            ),
            Expanded(
              child: Text(
                _objAsyncState,
                key: const Key('value-obj-async'),
              ),
            ),
          ],
        ),
      ],
    );
  }
}
