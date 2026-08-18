// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'session_provider.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(SessionProvider)
final sessionProviderProvider = SessionProviderProvider._();

final class SessionProviderProvider
    extends $AsyncNotifierProvider<SessionProvider, api.Session?> {
  SessionProviderProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'sessionProviderProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$sessionProviderHash();

  @$internal
  @override
  SessionProvider create() => SessionProvider();
}

String _$sessionProviderHash() => r'412878792a055be744fedede98d82c390d6b89b1';

abstract class _$SessionProvider extends $AsyncNotifier<api.Session?> {
  FutureOr<api.Session?> build();
  @$mustCallSuper
  @override
  WhenComplete runBuild() {
    final ref = this.ref as $Ref<AsyncValue<api.Session?>, api.Session?>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<AsyncValue<api.Session?>, api.Session?>,
              AsyncValue<api.Session?>,
              Object?,
              Object?
            >;
    return element.handleCreate(ref, build);
  }
}
