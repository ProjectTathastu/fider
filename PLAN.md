# Private Saved feedback

Rebuild on clean main after the previous sandbox was deleted. Keep voting,
subscriptions, search, tracked translation catalogs and complete tests intact.

1. Add tenant/user scoped saved-post persistence with unique idempotent writes,
   authenticated handlers, visibility filtering and deletion cleanup. Cover
   persistence, isolation, authorization and failure propagation.
2. Add Save/Unsave to post details and private Saved navigation/list/empty state,
   with original-post links, accessible controls and frontend regression tests.
3. Reconcile source catalogs, run focused checks, inspect clean Git status and
   commit. SAHR owns full configured validation and later fork-only draft PR
   publication; do not push, invoke gh, merge or run the full validation suite.

Prerequisite probes are separate from validation: verify godotenv in a minimal
validation PATH, fixture PostgreSQL authentication on 5566 with pg_trgm, and
MinIO fixture identity/bucket on 9000 without changing native root credentials.
Do not commit runtime configuration, credentials, generated bundles or preview
helpers. Desktop/phone preview with two disposable users remains an acceptance
criterion; record evidence and explicitly name anything not verified.
