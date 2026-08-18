// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'connection_provider.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(ConnectionStateNotifier)
final connectionStateProvider = ConnectionStateNotifierProvider._();

final class ConnectionStateNotifierProvider
    extends $NotifierProvider<ConnectionStateNotifier, WsConnectionState> {
  ConnectionStateNotifierProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'connectionStateProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$connectionStateNotifierHash();

  @$internal
  @override
  ConnectionStateNotifier create() => ConnectionStateNotifier();

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(WsConnectionState value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<WsConnectionState>(value),
    );
  }
}

String _$connectionStateNotifierHash() =>
    r'a96c3ad8fd77ec0960aa177c949f281833cfd26b';

abstract class _$ConnectionStateNotifier extends $Notifier<WsConnectionState> {
  WsConnectionState build();
  @$mustCallSuper
  @override
  WhenComplete runBuild() {
    final ref = this.ref as $Ref<WsConnectionState, WsConnectionState>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<WsConnectionState, WsConnectionState>,
              WsConnectionState,
              Object?,
              Object?
            >;
    return element.handleCreate(ref, build);
  }
}
