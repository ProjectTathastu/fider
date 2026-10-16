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

## Implementation and evidence

- Persistence/API committed first; UI/navigation/catalogs committed second.
  `saved_posts` has a unique tenant/user/post primary key. Writes resolve the post
  through `getPostByNumber`; lists reuse `buildSinglePostQuery`, so deleted and
  inaccessible moderated posts are excluded using the original detail rules.
  User/tenant deletion clears saved rows. No vote/subscription writes were added.
- The rebuilt checkout started at `ccd5428`; the deleted attempt's commit was not
  present. Its reported build failure was not reproduced by the focused baseline
  `go build -o /tmp/fider-baseline .`. No speculative application workaround was
  made for that excerpt or the prior S3 panic.
- Prerequisites: installed the existing godotenv executable in `/usr/local/bin`
  and checked discovery using `PATH=/usr/local/bin:/usr/bin:/bin`. `pg_isready -h
  127.0.0.1 -p 5566` reported accepting connections; SQL authenticated as the test
  role to the test database and found `pg_trgm`. MinIO readiness returned HTTP
  200; the fixture identity successfully authenticated to `test-bucket`. Native
  MinIO credentials and tracked `.test.env` were unchanged. The focused unchanged
  `TestBlobStorage/Test_S3_AllOperations` test passed.
- Focused checks passed:
  - `godotenv -f .test.env go test ./app/services/sqlstore/postgres ./app/handlers/apiv1 -run 'TestSaved' -race -count=1`
  - `godotenv -f .test.env go test ./app/handlers -run 'TestSavedPage' -count=1`
  - `TZ=GMT npx jest public/components/PostDetails/SavePostButton.test.tsx public/pages/Saved/Saved.page.test.tsx --runInBand` (8 tests)
  - ESLint on the new components/tests and account menu; `git diff --check`.
  - Go application build, catalog compilation, SSR bundling, and development
    webpack build to a temporary directory. Webpack reported one CSS ordering
    warning between existing TagsSelect and PostDetails styles.
  - All 21 source catalogs retained their existing keys and translations, adding
    only eight Saved messages each. Non-English messages use existing English
    fallback until translated.
- Plain `npx tsc --noEmit --pretty false` reported installed dependency declaration
  errors (Lingui, Tiptap, Node and Playwright). Application source type checking
  with `--skipLibCheck` passed. Compiler settings and dependency files were not
  changed; the plain TypeScript check is not claimed to pass.
- Chromium preview at 1440x1000 and 390x844 used a separate disposable database,
  never the validation database. Two disposable users exercised post creation,
  save/unsave, reload, sign-out/sign-in persistence, private lists/empty states,
  account navigation, original links, search/modal details and anonymous rejection.
  Voting/following remained set after save/unsave and reload; one user's unsave
  left the other's saved row intact. Screenshots were inspected, phone overflow
  was absent and no browser page errors were observed. Preview files and
  generated bundles remain outside Git. Fixture verification links were used;
  email delivery is unverified.
- Full configured validation (including complete `make test`) remains owned by
  SAHR and was not run in this implementation stage. Publication is deferred to
  SAHR's later fork-only draft PR stage; no push, gh, PR or merge was performed.

## Validation attempt 1 prerequisite repair

- SAHR reported `make: godotenv: No such file or directory` at Makefile:76
  for commit `d72d75f1`. Fresh inspection found no executable on PATH or at
  `/usr/local/bin/godotenv`; the earlier prerequisite evidence above describes
  the previous environment and does not establish setup in this session.
- Makefile:76 loads `.test.env` for migrations and line 77 loads it for server
  tests. `.github/workflows/build.yml` installs the CLI separately with
  `go install github.com/joho/godotenv/cmd/godotenv`; go.mod pins v1.5.1.
  The missing preview log supplies no evidence about runtime service health.
  No database, storage, application or test changes were justified by this error.
- Restored the existing prerequisite outside Git with
  `env GOBIN=/tmp/fider-repair-tools go install github.com/joho/godotenv/cmd/godotenv`
  and `sudo -n install -m 0755 /tmp/fider-repair-tools/godotenv /usr/local/bin/godotenv`.
  Direct installation as the workspace user was denied by directory permissions.
  `go version -m /usr/local/bin/godotenv` confirmed v1.5.1.
- A Python subprocess probe with `PATH=/usr/local/bin:/usr/bin:/bin` resolved
  the CLI, loaded a synthetic temporary env file and checked its value in a child
  shell. A failing child and an absent env file both returned nonzero. The initial
  probe incorrectly expected the child's exact exit code (23); the CLI returns
  1. The corrected probe asserts failure propagation and passed all cases.
  Temporary probe inputs were outside Git and contained no credentials.
- No full validation commands were run. Fresh full-suite results, service-backed
  acceptance, independent review and publication remain unverified in this repair
  session; prior preview evidence was not repeated. Email delivery remains
  unverified. SAHR owns full validation and later publication. This machine-local
  installation must also be present if validation runs in a different environment.

## Validation attempt 2 database prerequisite repair

- The supplied excerpt stops during successful catalog extraction and omits the
  actual error. The event log confirms validation failed, but the complete output
  for command `01a9e441-254d-487c-8c47-7f12ff55b3aa` was not available in the
  local session logs. The exact failing operation is therefore unconfirmed.
  Existing binary, compiled catalog and SSR timestamps fall within the reported
  16:21:02–16:21:12 UTC window; this alone does not prove validation success.
- Bounded read-only inspection found a concrete prerequisite mismatch:
  `pg_isready -h 127.0.0.1 -p 5566` reported no response, and
  `godotenv -f .test.env` followed by a Python-managed `psql ... SELECT 1`
  probe returned exit 2 with `Connection refused`. The effective endpoint was
  localhost:5566. The existing PostgreSQL postmaster PID/port and its startup log
  agreed on port 5432. Its log records listening on 5432 and readiness at
  16:19:32 UTC, before the failed command. MinIO readiness returned HTTP 200;
  storage authentication and operations were not repeated.
- Restored a separate disposable PostgreSQL fixture outside Git at
  `/tmp/fider-test-postgres-9vmyc2uo`, using `initdb` and
  `pg_ctl ... -o '-h 127.0.0.1 -p 5566 -k <fixture-directory>' -w start`.
  Created the role/database from the existing test environment and installed
  `pg_trgm`. The existing 5432 service, tracked environment files, application,
  migration behavior and tests were unchanged. No migration or configured
  validation command was run; SAHR still owns schema migration and full tests.
- Focused subprocess probes passed: the unmodified `.test.env` authenticated and
  found `pg_trgm` (`SELECT 1, EXISTS (...)` returned `1|t`); a random incorrect
  password failed authentication; a connection to an unavailable local port
  failed with `Connection refused`. `pg_isready` now reports accepting
  connections on 5566. These establish fixture readiness and failure reporting,
  not application acceptance or a passing validation run. Fixture credentials
  and logs remain outside Git.
- No application repair was supported by the available evidence. Full tests,
  fresh desktop/phone preview, independent review, email delivery and publication
  remain unverified in this repair session. The database fixture is machine-local
  and must be restored if validation uses another environment. No push, gh, PR
  or merge was performed.
