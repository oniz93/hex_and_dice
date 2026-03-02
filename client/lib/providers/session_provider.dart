import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'core_providers.dart';
import '../services/api_service.dart' as api;

part 'session_provider.g.dart';

@Riverpod(keepAlive: true)
class SessionProvider extends _$SessionProvider {
  @override
  Future<api.Session?> build() async {
    final storage = ref.watch(storageServiceProvider);
    final token = storage.token;
    final playerId = storage.playerId;
    final nickname = storage.nickname;

    if (token != null && playerId != null && nickname != null) {
      ref.read(apiServiceProvider).setToken(token);
      return api.Session(id: playerId, nickname: nickname, token: token);
    }
    return null;
  }

  Future<void> registerGuest(String nickname) async {
    state = const AsyncValue.loading();
    try {
      final session =
          await ref.read(apiServiceProvider).registerGuest(nickname);
      ref.read(apiServiceProvider).setToken(session.token);

      // Persist session to storage
      final storage = ref.read(storageServiceProvider);
      await storage.setToken(session.token);
      await storage.setPlayerId(session.id);
      await storage.setNickname(session.nickname);

      state = AsyncValue.data(session);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }
}
