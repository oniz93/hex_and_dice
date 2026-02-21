import 'package:flame_audio/flame_audio.dart';

class AudioService {
  double _musicVol = 1.0;
  double _sfxVol = 1.0;
  bool _muted = false;

  AudioService();

  Future<void> preload() async {
    // try to preload but don't crash if assets aren't there yet
    try {
      await FlameAudio.audioCache.loadAll([
        'music/menu_theme.ogg',
        'music/battle_theme.ogg',
        'sfx/attack_hit.wav',
        'sfx/attack_miss.wav',
        'sfx/troop_move.wav',
        'sfx/troop_death.wav',
        'sfx/dice_roll.wav',
        'sfx/turn_start.wav',
        'sfx/structure_capture.wav',
        'sfx/coin_gain.wav',
        'sfx/purchase.wav',
        'sfx/emote_pop.wav',
      ]);
    } catch (e) {
      print('Audio preload failed: $e');
    }
  }

  void setMusicVolume(double vol) {
    _musicVol = vol;
    // update current BGM vol
  }

  void setSfxVolume(double vol) {
    _sfxVol = vol;
  }

  void setMuted(bool muted) {
    _muted = muted;
    if (muted) {
      FlameAudio.bgm.stop();
    }
  }

  Future<void> playMenuMusic() async {
    if (_muted) return;
    try {
      await FlameAudio.bgm.play('music/menu_theme.ogg', volume: _musicVol);
    } catch (_) {}
  }

  Future<void> playBattleMusic() async {
    if (_muted) return;
    try {
      await FlameAudio.bgm.play('music/battle_theme.ogg', volume: _musicVol);
    } catch (_) {}
  }

  Future<void> stopMusic() async {
    FlameAudio.bgm.stop();
  }

  void _playSfx(String file) {
    if (_muted) return;
    try {
      FlameAudio.play(file, volume: _sfxVol);
    } catch (_) {}
  }

  void playAttackHit() => _playSfx('sfx/attack_hit.wav');
  void playAttackMiss() => _playSfx('sfx/attack_miss.wav');
  void playTroopMove() => _playSfx('sfx/troop_move.wav');
  void playTroopDeath() => _playSfx('sfx/troop_death.wav');
  void playDiceRoll() => _playSfx('sfx/dice_roll.wav');
  void playTurnStart() => _playSfx('sfx/turn_start.wav');
  void playStructureCapture() => _playSfx('sfx/structure_capture.wav');
  void playCoinGain() => _playSfx('sfx/coin_gain.wav');
  void playPurchase() => _playSfx('sfx/purchase.wav');
  void playEmotePop() => _playSfx('sfx/emote_pop.wav');
}
