import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../services/ws_service.dart';
import 'core_providers.dart';
import 'dart:async';

part 'connection_provider.g.dart';

@Riverpod(keepAlive: true)
class ConnectionStateNotifier extends _$ConnectionStateNotifier {
  StreamSubscription? _sub;

  @override
  WsConnectionState build() {
    final wsService = ref.watch(wsServiceProvider);

    _sub?.cancel();
    _sub = wsService.connectionState.listen((state) {
      this.state = state;
    });

    ref.onDispose(() {
      _sub?.cancel();
    });

    return WsConnectionState.disconnected;
  }

  Future<void> connect(String token) async {
    await ref.read(wsServiceProvider).connect(token);
  }

  void disconnect() {
    ref.read(wsServiceProvider).disconnect();
  }
}
