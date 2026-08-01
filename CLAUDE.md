# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

contrib.rocks — a service that renders a repository's GitHub contributors as a single SVG image, served at
`GET /image?repo=<owner>/<repo>`. The image is embedded in READMEs, so it is fetched by GitHub's camo proxy far
more often than by browsers. Caching and cheap responses matter more than request-time freshness.

## Commands

Package manager is **pnpm** (`pnpm@10.8.1`); Node version is pinned in `.node-version`. Nx targets are declared
in each `apps/*/project.json`.

```bash
pnpm install                      # required before any nx command

pnpm nx serve api                 # Go API on :3333 (needs GITHUB_AUTH_TOKEN, see below)
pnpm nx serve webapp              # Angular dev server, proxies /image and /api to :3333
pnpm nx serve worker              # worker via tsx — also defaults to :3333, so set PORT to run it
                                  # alongside the API; its one route needs ADC + BigQuery access

pnpm nx test api                  # go test ./... in apps/api
pnpm nx test golib                # go test ./... from the repo root — covers every Go package
pnpm nx test webapp               # Karma + ChromeHeadless, single run
pnpm test:ci                      # all projects

pnpm lint                         # go vet + eslint across projects
pnpm format                       # gofmt -w + prettier -w
pnpm build:all:production         # despite the name, builds only webapp
```

Running a single Go test: `cd apps/api && go test ./internal/service/image/ -run TestName`.

Regenerating the `cupaloy` snapshots under `go/renderer/.snapshots` and `go/dataurl/.snapshots`:

```bash
UPDATE_SNAPSHOTS=true go test -count=1 ./apps/api/go/...
```

Do **not** use the `test:u` package script. It expands to a bare `nx test` with no project and exits 1; supplying
a project makes it worse, because `nx.json` caches the `test` target and `UPDATE_SNAPSHOTS` is not part of the
hash — Nx replays a cached success without ever starting `go test`, so the snapshots are silently not rewritten.
`-count=1` defeats Go's own test cache for the same reason.

The `worker` project's `test` and `lint` targets are `echo` stubs — it has no tests and is not linted.

## Local setup for the API

`config.Load` returns an error and the server exits if `GITHUB_AUTH_TOKEN` is unset. `main.go` calls
`godotenv.Load()` and the `serve` target runs `go run .` with `cwd` = `apps/api`, so put the token in
`apps/api/.env`. Verify `git check-ignore apps/api/.env` exits 0 before writing a real token into it — this is a
GitHub credential, and a leaked one has to be revoked, not un-pushed.

Without Google credentials **and** `CACHE_STORAGE_BUCKET`, `NewServicePack` falls back to an in-memory cache
instead of GCS. That is the normal local configuration.

## Architecture

Three deployables, one Nx workspace:

| Project  | Path           | Stack                        | Deploy target                        |
| -------- | -------------- | ---------------------------- | ------------------------------------ |
| `api`    | `apps/api`     | Go 1.23, Gin, OpenTelemetry  | Cloud Run (image built by `ko`)      |
| `webapp` | `apps/webapp`  | Angular 21 standalone + Material | Firebase Hosting                 |
| `worker` | `apps/worker`  | Node, Hono, BigQuery→Firestore | Cloud Run (`--source`, buildpacks) |

`@nx/angular` has no imports anywhere in this repo and generates nothing here — but it must stay a devDependency
regardless. `nx migrate` reads a package's `packageJsonUpdates` ladder only from a package that is already in
`package.json`; if `@nx/angular` is missing, none of its Angular-version-lockstep entries run, and `@angular/material`,
`@angular/cdk`, `zone.js`, and `angular-eslint` silently stop tracking `@angular/core`. Re-adding it after the fact
with `nx add` does not backfill the skipped ladder — it must already be declared, at the version the workspace is
actually on, before `nx migrate` runs. The price of keeping it is large and moves with every Nx bump.
`@nx/angular` depends on `@nx/rspack`, which drags in `@rspack/binding-*` for **every** published target rather
than the host's alone: at Nx 23 that is 10 platforms (`darwin-{arm64,x64}`, `linux-{x64,arm64}-{gnu,musl}`,
`win32-{x64,ia32,arm64}-msvc`, `wasm32-wasi`), and two minor lines of rspack 1.x (`1.6.8` and `1.7.12`) coexist,
so it lands as 20 store entries. `du -sch node_modules/.pnpm/@rspack+binding-*` reports ~850 MB, but that is the
*logical* size — pnpm clones from the global store, so it overstates what a removal would actually reclaim. Before
quoting any number in a decision about dropping the package, measure the **difference** between installs with and
without the `package.json` line; the absolute figure above is not that difference.

TypeScript's version is a separate trap: it looks unmanaged, but `@nx/workspace`'s own `packageJsonUpdates` owns
it. Read the chain out of `node_modules/@nx/workspace/migrations.json` (`packageJsonUpdates`, package `typescript`)
rather than deriving it from `@angular/compiler-cli`'s peer range; the peer range is only a lower bound the ladder
happens to satisfy, not the ladder's own destination. At Nx 23 the chain is two rungs — `21.2.0` → `~5.8.2` and
`21.5.0` → `~5.9.2` — and two independent filters decide whether either fires.

- **The migration window**, applied first: `filterPackageJsonUpdates` (`migrate.js:313`) drops any rung whose own
  version is below the *installed* version of the migrating package, or above the target. Both rungs sit at 21.x,
  so migrating `@nx/workspace` 22 → 23 never reaches them. A rung is reachable only during the migration that
  crosses it — the same mechanism that requires `@nx/angular` to be declared at the version the workspace is on.
- **The `requires` gate**, checked second: `>=5.7.0 <5.8.0` and `>=5.8.0 <5.9.0` against the workspace's *current*
  TypeScript, which is `5.9.3` — outside both.

**So, after every `nx migrate`: read the `typescript` rungs out of the freshly installed
`node_modules/@nx/workspace/migrations.json` and set `package.json`'s `typescript` by hand to the version the
highest reachable rung names. Nothing sets it for you, and a skipped window is not a reason to assume it is
current.** Do not reach for `--interactive` to collect the rungs: it is force-disabled whenever `isCI()` is true,
hangs waiting on a real TTY otherwise, and forcing input through a pseudo-terminal produced actively wrong results
(stopped mid-chain at an intermediate version, and answered an unrelated Nx Cloud telemetry consent prompt it
happened to also raise). Nx 23 deprecates it outright — removal in v24 — in favour of `--include` (`required` /
`optional` / `all`, also `migrate.include` in `nx.json`), which falls through to `all` without a TTY. The
`x-prompt` these rungs carry is not a filter: `migrate.js:192` reads it only under `--interactive`.

The app is not zoneless — but from Angular 21 on, that is no longer the default. v21 flipped the `ZONELESS_ENABLED`
token's default factory to `true`, so `app.config.ts` has to call `provideZoneChangeDetection()` to stay
zone-based; that call is what supplies `ZONELESS_ENABLED: false`. It is not leftover boilerplate from the
migration — deleting it silently makes the app zoneless.

**No test covers that decision.** The Karma builder synthesizes its own `provideZoneChangeDetection()` into a
generated test module whenever `zone.js` appears in the target's `polyfills`
(`@angular/build/src/builders/karma/application_builder.js`), so a bootstrap that lost the provider still passes
the suite. Only loading the deployed page exercises it.

`apps/api/go/` is a separate Nx project named **`golib`** (tag `lib`), holding `apiclient`, `compress`, `dataurl`,
`env`, `httptrace`, `model`, `renderer`, `util`. The boundary is Gin and app config — **not** cloud SDKs:
`go/apiclient` constructs the BigQuery / Firestore / Logging / Storage clients and `go/httptrace` wraps
OpenTelemetry. A new GCP client belongs in `go/apiclient`, not in `internal/`. `apps/api/internal/` holds what
knows about Gin, config, and service wiring — `internal/service/services.go` is where GCS-vs-in-memory cache is
decided. `api` declares `implicitDependencies: ["golib"]`.

One Go module, `contrib.rocks`, rooted at the repo root — not at `apps/api`. Import paths are
`contrib.rocks/apps/api/...`.

### Request flow for `GET /image`

`internal/api/image` (bind + validate params) → `image.Service.GetImage` (cache lookup; on hit, return) →
`contributors.Service.GetContributors` (its own cache, then GitHub API) → `image.Service.RenderImage`
(fetch every avatar, inline as base64 data URL, render SVG, write to cache) → `usage.Service.CollectUsage`.

Every avatar is inlined as a base64 data URL rather than referenced by URL, so the SVG is self-contained. That
makes rendering expensive — one HTTP request per contributor, parallelised with `errgroup` — which is why the
cache is load-bearing rather than an optimisation. Keep this in mind before adding work to the render path.

Handlers depend on service *interfaces* declared in the consuming package (`api/image/api.go`), not on the
concrete service structs. Wiring happens in `internal/service/services.go`.

### Caching

Both layers share one `AppCache` (GCS bucket in deployed environments, in-memory locally). Keys are built only in
`internal/service/internal/cachekey`:

- contributors JSON — `contributors/v1.2/{owner}--{repo}.json`
- rendered image — `image/{owner}--{repo}--{anon|noanon}_{max}_{columns}.svg`

**The image key carries no renderer version.** Changing SVG output does not invalidate cached images; the
contributors key has an explicit `v1.2` for exactly that reason. If you change `go/renderer` output, plan for
stale entries — responses also carry `cache-control: max-age=259200` (3 days) plus an ETag.

Renderer defaults live in `internal/service/image/options.go` (max 100, 12 columns, 64px items) and are applied by
`normalizeRendererOptions` before the cache key is computed, so `?max=0` and no `max` hit the same entry.

### webapp ↔ worker coupling

The worker has a single route, `/update-featured-repositories`. It queries the BigQuery table
`contributors-img.repository_usage.weekly_repository_usage` and writes `featured_repositories` and `usage_stats`
documents into the Firestore collection named after `APP_ENV`. That table is not defined in this repo; the API
emits the matching per-request records as structured logs under the `repository-usage` log group
(`internal/service/usage`), and the pipeline between the two lives in GCP config, not here.

The webapp reads only one of those two documents: `{collection}/featured_repositories`, from the single Firestore
call site in `app/shared/featured-repository/firestore.ts`. **`usage_stats` has no reader in this repo** — the
worker writes it and nothing here consumes it, so changing its shape does not break the webapp.

The collection comes from `environment.firestoreRootCollectionName`. All three `src/environments/environment*.ts`
import the same `firebaseConfig.prod` (`firebase-config.ts` defines no other key), so every environment points at
the one `contributors-img` Firebase project; they differ in `firestoreRootCollectionName`
(`development` / `staging` / `production`) and in `production`, which is `false` only in `environment.ts`.
Environment separation here is by collection, not by project. `firebase/firestore.rules` makes all three
collections world-readable and client-write-denied.

## Deployment

- Push to `main` → `deploy-staging.yml` deploys all three to staging.
- `release-please.yml` maintains a release PR; merging it triggers `deploy-production.yml` at the release SHA.
- `deploy-production.yml` also accepts `workflow_dispatch` with an arbitrary `ref`.
- Firebase Hosting rewrites `/image` and `/api/**` to the Cloud Run API service, so the webapp and API share an
  origin in production. Locally, `proxy.conf.json` reproduces this.

Conventional commits are required — release-please derives versions and CHANGELOG from them.

## Conventions

`.github/copilot-instructions.md` also exists, but do not adopt it wholesale. Its coding rules hold (English
comments, "why" over "what", conventional commits, strict TS). Three things in it do not:

- its Commands section says `npx nx …` — this repo is pnpm, use `pnpm nx`;
- it offers `nx format:write`, which runs Prettier only and skips the `gofmt -w` in the `api` and `golib` `format`
  targets. Use `pnpm format`. No CI job checks Go formatting, so unformatted Go can reach `main` uncaught;
- "write unit tests for all new functionality" does not apply to `worker`, which has no test runner.

Ignore its opening instruction to edit itself when it disagrees with the implementation. Commands are canonical
here, in this file.

Additionally:

- Comments and identifiers in English. Some existing Go comments are in Japanese; do not add more.
- Angular: standalone components, `OnPush` (Nx generator default), SCSS. There is no `NgModule` and no
  `angular.json` — project config lives in `apps/webapp/project.json`. `workspace.json` is **not** the project
  list: it names only `api`, `webapp`, and `worker`, while `golib` is discovered by Nx from
  `apps/api/go/project.json`.
- Tests: Go uses table tests plus `cupaloy` snapshots under `.snapshots/`; webapp uses Karma with
  `@testing-library/angular`.
- `.npmrc` sets `shared-workspace-lockfile=false`, so `apps/worker/pnpm-lock.yaml` exists separately. It must stay
  self-contained — production deploys `--source ./apps/worker` and resolves dependencies there.
