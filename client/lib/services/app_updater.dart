import 'dart:js_interop';

import 'package:flutter/foundation.dart';
import 'package:web/web.dart' as web;

/// Unregisters any active service workers and reloads the page.
///
/// This forces browsers to fetch the newest Flutter web build instead of
/// serving a stale service-worker cache.
Future<void> forceUpdateApp() async {
  if (!kIsWeb) return;

  try {
    final registrations =
        (await web.window.navigator.serviceWorker.getRegistrations().toDart)
            .toDart;
    for (final registration in registrations) {
      await registration.unregister().toDart;
    }
  } catch (_) {
    // Fall through: reload even if service worker cleanup fails.
  }

  web.window.location.reload();
}
