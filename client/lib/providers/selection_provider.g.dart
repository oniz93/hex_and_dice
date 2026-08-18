// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'selection_provider.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(SelectionStateNotifier)
final selectionStateProvider = SelectionStateNotifierProvider._();

final class SelectionStateNotifierProvider
    extends $NotifierProvider<SelectionStateNotifier, SelectionState> {
  SelectionStateNotifierProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'selectionStateProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$selectionStateNotifierHash();

  @$internal
  @override
  SelectionStateNotifier create() => SelectionStateNotifier();

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(SelectionState value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<SelectionState>(value),
    );
  }
}

String _$selectionStateNotifierHash() =>
    r'a2722d49019220011f12eba4993d15f8cc2a3abb';

abstract class _$SelectionStateNotifier extends $Notifier<SelectionState> {
  SelectionState build();
  @$mustCallSuper
  @override
  WhenComplete runBuild() {
    final ref = this.ref as $Ref<SelectionState, SelectionState>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<SelectionState, SelectionState>,
              SelectionState,
              Object?,
              Object?
            >;
    return element.handleCreate(ref, build);
  }
}
