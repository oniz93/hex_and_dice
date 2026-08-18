// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'turn_timer_provider.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(TurnTimer)
final turnTimerProvider = TurnTimerProvider._();

final class TurnTimerProvider extends $NotifierProvider<TurnTimer, int> {
  TurnTimerProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'turnTimerProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$turnTimerHash();

  @$internal
  @override
  TurnTimer create() => TurnTimer();

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(int value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<int>(value),
    );
  }
}

String _$turnTimerHash() => r'b52aa314c144d2b360f6d0c8da371cbe8c7267ee';

abstract class _$TurnTimer extends $Notifier<int> {
  int build();
  @$mustCallSuper
  @override
  WhenComplete runBuild() {
    final ref = this.ref as $Ref<int, int>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<int, int>,
              int,
              Object?,
              Object?
            >;
    return element.handleCreate(ref, build);
  }
}
