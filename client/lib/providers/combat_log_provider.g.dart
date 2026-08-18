// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'combat_log_provider.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(CombatLog)
final combatLogProvider = CombatLogProvider._();

final class CombatLogProvider
    extends $NotifierProvider<CombatLog, List<CombatLogEntry>> {
  CombatLogProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'combatLogProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$combatLogHash();

  @$internal
  @override
  CombatLog create() => CombatLog();

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(List<CombatLogEntry> value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<List<CombatLogEntry>>(value),
    );
  }
}

String _$combatLogHash() => r'6c4f59b944a9f4afe5bc524c2d18004e94722b42';

abstract class _$CombatLog extends $Notifier<List<CombatLogEntry>> {
  List<CombatLogEntry> build();
  @$mustCallSuper
  @override
  WhenComplete runBuild() {
    final ref = this.ref as $Ref<List<CombatLogEntry>, List<CombatLogEntry>>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<List<CombatLogEntry>, List<CombatLogEntry>>,
              List<CombatLogEntry>,
              Object?,
              Object?
            >;
    return element.handleCreate(ref, build);
  }
}
