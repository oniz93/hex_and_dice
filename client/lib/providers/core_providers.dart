import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../services/storage_service.dart';
import '../services/api_service.dart';
import '../services/ws_service.dart';
import '../services/audio_service.dart';

// These providers will be overridden in main() once SharedPreferences is initialized
final sharedPreferencesProvider = Provider<SharedPreferences>((ref) {
  throw UnimplementedError();
});

final storageServiceProvider = Provider<StorageService>((ref) {
  final prefs = ref.watch(sharedPreferencesProvider);
  return StorageService(prefs);
});

final apiServiceProvider = Provider<ApiService>((ref) {
  return ApiService(baseUrl: 'http://localhost:8080');
});

final wsServiceProvider = Provider<WsService>((ref) {
  return WsService(baseUrl: 'ws://localhost:8080');
});

final audioServiceProvider = Provider<AudioService>((ref) {
  return AudioService();
});
