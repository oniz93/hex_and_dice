// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'game_state_provider.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(GameStateNotifier)
final gameStateProvider = GameStateNotifierProvider._();

final class GameStateNotifierProvider
    extends $NotifierProvider<GameStateNotifier, GameState?> {
  GameStateNotifierProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'gameStateProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$gameStateNotifierHash();

  @$internal
  @override
  GameStateNotifier create() => GameStateNotifier();

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(GameState? value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<GameState?>(value),
    );
  }
}

String _$gameStateNotifierHash() => r'e830085df3c9999484aa45524628f323616851ef';

abstract class _$GameStateNotifier extends $Notifier<GameState?> {
  GameState? build();
  @$mustCallSuper
  @override
  WhenComplete runBuild() {
    final ref = this.ref as $Ref<GameState?, GameState?>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<GameState?, GameState?>,
              GameState?,
              Object?,
              Object?
            >;
    return element.handleCreate(ref, build);
  }
}
