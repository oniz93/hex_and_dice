import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../services/storage_service.dart';
import '../services/api_service.dart';
import '../services/ws_service.dart';
import '../services/audio_service.dart';

// Derive the server base URL from the browser's current location.
String _httpBaseUrl() {
  final host = Uri.base.host;
  if (host == 'localhost' || host == '127.0.0.1') {
    final base = Uri.base;
    return '${base.scheme}://${base.host}${base.hasPort ? ':${base.port}' : ''}';
  }
  return 'http://api.hexdice.teomiscia.com';
}

String _wsBaseUrl() {
  final host = Uri.base.host;
  if (host == 'localhost' || host == '127.0.0.1') {
    final base = Uri.base;
    final scheme = base.scheme == 'https' ? 'wss' : 'ws';
    return '$scheme://${base.host}${base.hasPort ? ':${base.port}' : ''}';
  }
  return 'ws://api.hexdice.teomiscia.com';
}

final sharedPreferencesProvider = Provider<SharedPreferences>((ref) {
  throw UnimplementedError();
});

final storageServiceProvider = Provider<StorageService>((ref) {
  final prefs = ref.watch(sharedPreferencesProvider);
  return StorageService(prefs);
});

final apiServiceProvider = Provider<ApiService>((ref) {
  return ApiService(baseUrl: _httpBaseUrl());
});

final wsServiceProvider = Provider<WsService>((ref) {
  ref.keepAlive();
  return WsService(baseUrl: _wsBaseUrl());
});

final audioServiceProvider = Provider<AudioService>((ref) {
  return AudioService();
});
