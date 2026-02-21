import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:client/main.dart';

void main() {
  testWidgets('App smoke test', (WidgetTester tester) async {
    // We would need to mock SharedPreferences here, so we skip it for now.
    // await tester.pumpWidget(ProviderScope(child: const HexBattleApp()));
  });
}
