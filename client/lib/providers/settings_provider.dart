import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'core_providers.dart';

part 'settings_provider.g.dart';

class SettingsState {
  final String nickname;
  final double musicVol;
  final double sfxVol;
  final bool muted;
  final String? gameId;
  final String? token;

  const SettingsState({
    required this.nickname,
    required this.musicVol,
    required this.sfxVol,
    required this.muted,
    this.gameId,
    this.token,
  });

  SettingsState copyWith({
    String? nickname,
    double? musicVol,
    double? sfxVol,
    bool? muted,
    String? gameId,
    String? token,
  }) {
    return SettingsState(
      nickname: nickname ?? this.nickname,
      musicVol: musicVol ?? this.musicVol,
      sfxVol: sfxVol ?? this.sfxVol,
      muted: muted ?? this.muted,
      gameId: gameId ?? this.gameId,
      token: token ?? this.token,
    );
  }
}

@Riverpod(keepAlive: true)
class Settings extends _$Settings {
  @override
  SettingsState build() {
    final storage = ref.watch(storageServiceProvider);
    return SettingsState(
      nickname: storage.nickname ?? 'Guest',
      musicVol: storage.musicVol,
      sfxVol: storage.sfxVol,
      muted: storage.muted,
      gameId: storage.gameId,
      token: storage.token,
    );
  }

  Future<void> setNickname(String nickname) async {
    await ref.read(storageServiceProvider).setNickname(nickname);
    state = state.copyWith(nickname: nickname);
  }

  Future<void> setMusicVol(double vol) async {
    await ref.read(storageServiceProvider).setMusicVol(vol);
    state = state.copyWith(musicVol: vol);
  }

  Future<void> setSfxVol(double vol) async {
    await ref.read(storageServiceProvider).setSfxVol(vol);
    state = state.copyWith(sfxVol: vol);
  }

  Future<void> setMuted(bool muted) async {
    await ref.read(storageServiceProvider).setMuted(muted);
    state = state.copyWith(muted: muted);
  }

  Future<void> setReconnectData(String gameId, String token) async {
    await ref.read(storageServiceProvider).setGameId(gameId);
    await ref.read(storageServiceProvider).setToken(token);
    state = state.copyWith(gameId: gameId, token: token);
  }

  Future<void> clearReconnectData() async {
    await ref.read(storageServiceProvider).clearReconnectData();
    state = state.copyWith(gameId: null, token: null);
  }
}
