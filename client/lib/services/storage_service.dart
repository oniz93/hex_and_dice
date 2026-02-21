import 'package:shared_preferences/shared_preferences.dart';

class StorageService {
  final SharedPreferences _prefs;

  StorageService(this._prefs);

  static Future<StorageService> init() async {
    final prefs = await SharedPreferences.getInstance();
    return StorageService(prefs);
  }

  String? get nickname => _prefs.getString('nickname');
  Future<void> setNickname(String value) => _prefs.setString('nickname', value);

  double get musicVol => _prefs.getDouble('musicVol') ?? 1.0;
  Future<void> setMusicVol(double value) => _prefs.setDouble('musicVol', value);

  double get sfxVol => _prefs.getDouble('sfxVol') ?? 1.0;
  Future<void> setSfxVol(double value) => _prefs.setDouble('sfxVol', value);

  bool get muted => _prefs.getBool('muted') ?? false;
  Future<void> setMuted(bool value) => _prefs.setBool('muted', value);

  String? get gameId => _prefs.getString('gameId');
  Future<void> setGameId(String value) => _prefs.setString('gameId', value);

  String? get token => _prefs.getString('token');
  Future<void> setToken(String value) => _prefs.setString('token', value);

  Future<void> clearReconnectData() async {
    await _prefs.remove('gameId');
    await _prefs.remove('token');
  }
}
