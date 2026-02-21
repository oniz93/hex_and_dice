import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'core_providers.dart';
import '../services/api_service.dart' as api;

part 'session_provider.g.dart';

@Riverpod(keepAlive: true)
class SessionProvider extends _$SessionProvider {
  @override
  Future<api.Session?> build() async {
    return null;
  }

  Future<void> registerGuest(String nickname) async {
    state = const AsyncValue.loading();
    try {
      final session = await ref
          .read(apiServiceProvider)
          .registerGuest(nickname);
      ref.read(apiServiceProvider).setToken(session.token);
      state = AsyncValue.data(session);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }
}
