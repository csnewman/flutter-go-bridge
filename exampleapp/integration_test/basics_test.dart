import 'package:exampleapp/main.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

const notCalledText = "not called";

List<Map<String, String>> basicTests = [
  {
    "name": "add sync",
    "elem-text": "value-add-sync",
    "elem-button": "button-add-sync",
    "result": "3",
  },
  {
    "name": "add async",
    "elem-text": "value-add-async",
    "elem-button": "button-add-async",
    "result": "244444",
  },
  {
    "name": "add points sync",
    "elem-text": "value-add-points-sync",
    "elem-button": "button-add-points-sync",
    "result": "Point{x: 457, y: 791, name: PointA+PointB}",
  },
  {
    "name": "add points async",
    "elem-text": "value-add-points-async",
    "elem-button": "button-add-points-async",
    "result": "Point{x: 457, y: 791, name: PointA+PointB}",
  },
  {
    "name": "add errors sync",
    "elem-text": "value-add-errors-sync",
    "elem-button": "button-add-errors-sync",
    "result": "error=BridgeException: add res was 2",
  },
  {
    "name": "add errors async",
    "elem-text": "value-add-errors-async",
    "elem-button": "button-add-errors-async",
    "result": "error=BridgeException: add res was 4",
  },
  {
    "name": "obj sync",
    "elem-text": "value-obj-sync",
    "elem-button": "button-obj-sync",
    "result": "Obj: Name=test1 Other=200",
  },
  {
    "name": "obj async",
    "elem-text": "value-obj-async",
    "elem-button": "button-obj-async",
    "result": "Obj: Name=test2 Other=246",
  },
];

void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();
  for (var t in basicTests) {
    print(t);
    testWidgets(t["name"]!, (tester) async {
      await tester.pumpWidget(const MaterialApp(
        home: Scaffold(body: Example1()),
      ));

      await tester.pumpAndSettle();

      expect(getText(t["elem-text"]!), notCalledText);

      await tester.tap(find.byKey(Key(t["elem-button"]!)));
      await tester.pumpAndSettle();

      expect(getText(t["elem-text"]!), t["result"]!);
    });
  }
}

String getText(String name) {
  final elem = find.byKey(Key(name)).evaluate().single.widget as Text;
  return elem.data!;
}
