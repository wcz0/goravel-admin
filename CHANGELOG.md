# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Release notes per version live in [`docs/releases`](./docs/releases).

## [Unreleased]

### Added

- **Added the `goravel-1.17` support line (this branch)**, targeting `goravel/framework` v1.17.3; `example/go.mod`
  pins `gin v1.17.1`, `redis v1.17.1` and `sqlite v1.17.0` accordingly.
- Login captcha: `GET /captcha` returns `captcha_img` (a data URI) and `sys_captcha`; the answer is stored in the cache
  and consumed once. When `admin.auth.login_captcha` is enabled the login endpoint enforces it, and the amis login page
  renders the captcha input group.
- Profile update `PUT /user_setting`: change avatar, name and password. Changing the password requires the old one and
  both new entries must match.
- `GET /page_schema?sign=` plus "Dev Tools → Pages": page schemas are stored as amis JSON in `admin_pages` and fetched
  by sign when a menu uses `url_type=4`.
- "Dev Tools → Dynamic APIs / Relationships / Extensions": list, show, create, update, delete, template dropdown,
  relationship type dropdown and `preview_code`.
- `GET /_iconify_search`: all iconify icon names (220k) embedded in the binary, parsed once on first use and kept in
  memory; returns at most 1000 matches per call.
- Chunked upload endpoints `upload_chunk_start` / `upload_chunk` / `upload_chunk_finish`: 10-minute sessions, md5
  verification per part, target directory allow-list, script extensions rejected.
- Export: "export all / this page / selected" toolbar actions, xlsx generation, and an HMAC-signed, expiring
  `_download_export` URL with path traversal protection. Wired up for the user list.
- `admin:create-user`, `admin:reset-password` and `admin:doctor` commands. `admin:doctor` checks `APP_KEY`, database
  connectivity, the 12 core tables, admin accounts, home menu uniqueness and broken menu parents.
- `core/captcha_test.go` covering "generate → verify → single-use expiry" plus rejection of empty and wrong answers.
- `example/`: a minimal host application that verifies the full `setup install` flow (provider registration,
  migrations, `admin:init`, login and list endpoints).
- `setup install` / `setup uninstall`: the installer only registers the service provider and writes an optional config
  skeleton; it no longer modifies the host `bootstrap/migrations.go`.
- `migrations.All()` and `tools.AuthGuard()`.
- `services/preview_test.go` covering identifier sanitization, relationship `preview_code`, page schema extraction and
  primary-key guards.
- `example/smoke_admin.py`: an end-to-end contract smoke test covering every endpoint the owl-admin frontend calls.
- Layered test suite alongside the smoke test: package-level unit tests that need neither container nor database
  (response envelope, export signing and traversal guards, permission tree ordering, identifier conversion, upload
  validation, icon search, legacy `StringSlice` decoding, migration inventory), plus container-backed integration
  tests in `example/tests/feature` that render a real xlsx, write it to the local disk and verify the signature round
  trip.
- `.github/workflows/tests.yml`: build, `go vet` and unit tests with coverage on every push and pull request, plus a
  second job that boots the example application and runs the integration tests.

### Fixed

- `enums.CodeEnum.Message()` returned `Unknown error` for `401 Unauthorized` and `403 PermissionDenied` — both were
  missing from the switch.
- Permission tree children came back in random order: sibling nodes were collected in a `map`, and Go randomises map
  iteration, so "Roles" could appear above "Users" on one refresh and below it on the next. Siblings are now sorted by
  `custom_order` and then by `id`.
- Export worksheet names were truncated by **bytes**, so a Chinese sheet name longer than 32 bytes was cut mid-rune and
  produced invalid UTF-8 that Excel refused to open. Truncation is now rune-based and respects Excel's 31-character
  limit.
- `response.Unauthorized` / `TokenExpired` / `FormError` / `LimitError` / `PermissionError` had no `status` field. The
  frontend decides whether a body is an AMis-style envelope with `typeof backend.status === 'number'`, so missing
  `status` made these error responses fall through as raw axios responses.
- `FormError` / `LimitError` / `UnauthorizedError` / `PermissionError` discarded the caller's message (`if msg != ""
  { msg = "fixed text" }`, with the variable never used). They now return the template default when no message is
  given and the custom message otherwise, and never mutate the shared package-level template.
- `FailMsg` / `FailData` / `FailMsgData` omitted the `code` field that every other failure response carries; they now
  include `code: 400`.
- `safeExtension` compared extensions case-sensitively, so `PHP` / `.EXE` slipped through if a future upload entry
  point forgot to lowercase. The check now normalizes case and a leading dot itself.
- `VerifyExport` accepted any signature when `APP_KEY` was empty (empty HMAC key = anyone can forge a link). It now
  fails closed with an explicit error.

- Image and file upload configuration keys were concatenated incorrectly (`admin.upload.directoryimage`), so uploads
  landed in the storage root. Files now resolve `admin.upload.directory.<type>`: images go to `images/`, files to
  `files/`.
- The profile form posted to `PUT /current-user`, which does not exist; it now targets `PUT /user_setting`.
- `enabled_extensions` in `_settings` was always an empty array; it now reports the names enabled in `admin_extensions`.
- Rows could be written when the primary key was empty or the record did not exist. Goravel's `Find` and `First` both
  stay silent on `id=0` and on missing rows (`First` swallows `ErrRecordNotFound`), so the old code received a
  zero-value struct and `Save` inserted a blank row — `POST /dev_tools/extensions/enable` with `id=0` created an
  extension with an empty name and pushed it into `enabled_extensions`. Added `services.firstByID` (a `Count` before
  `First`) and routed single-record reads of pages, dynamic APIs, relationships, extensions, menus, permissions and
  users through it.
- `POST /dev_tools/pages` rejected a `schema` submitted as a JSON object with "page schema cannot be empty". The nested
  `page.schema` object, a top-level `schema` object and a `schema` JSON string are all accepted now.
- `preview_code` for relationships produced code that would not compile: `model` / `title` were pasted into identifiers
  verbatim (a title such as `Smoke Roles` broke it) and the snippet returned `orm.Relation`, which Goravel does not
  have. Identifiers are sanitized to valid Go names, and the snippet is now an explicit-query skeleton the ORM really
  supports (pivot tables via `Join`, one-to-many via `Where`). Form defaults are `related` / `table` / `foreign_key` /
  `related_key`.

- `query.Where(...)` returns a **new** query object in Goravel (the conditions are copied into a fresh instance), so
  every call that ignored the return value silently dropped its filter. Thirteen sites were affected across pages,
  dynamic APIs, relationships, menus, permissions, roles and users — the search boxes in those lists returned the whole
  table, and the `id != ?` guard in the uniqueness checks was never applied (editing a page without touching its `sign`
  reported a duplicate). All of them now assign the result, and OR pairs go through `WhereAny` so they stay grouped.
- `admin:init` inserted the admin↔role pivot twice per run: a raw `INSERT` at step 4 and a `Create` at step 10. A single
  run therefore left two identical rows, and the user list rendered duplicate role tags. The command is now idempotent
  end to end (`FirstOrCreate` for roles, the user, permissions, menus, settings and every pivot), so running it again
  changes nothing.
- Seeded permissions carried paths that do not exist in this package (`/auth/roles`, `/auth/menu`, …); no request could
  ever match them, and `admin_permission_menu` stayed empty, so any non-superadmin user logged in to an empty sidebar.
  Permissions now describe the real routes, they are mapped to menus, and the Operator role receives a working subset.
- `AdminUser.AllPermissions()` asked the `AdminRoles` association to fill a `[]AdminPermission` slice. The types never
  matched, the call always errored and returned `nil`, so **every** non-superadmin request was answered with 403.
  It now aggregates permissions across the user's roles with de-duplication.
- Permission path matching stripped the leading slash from the pattern before comparing it with
  `ctx.Request().Path()` (`/admin-api/...`), so wildcards never matched anything. `tools.MatchPath` now normalises both
  sides consistently and has unit tests.
- `admin.auth.except` and `admin.permission.except` whitelisted `system/admin_users`, which left the user list readable
  **without logging in**. Both lists are empty now; only the login-related endpoints are exempt.
- `AdminUserService.List` returned every user with no pagination metadata (`total` / `page` / `perPage`), so amis could
  not render the pager, and it issued one extra query per row to load roles. Pagination is implemented, `perPage` is
  clamped to 500, and roles are preloaded.
- `admin:init` seeded the child menus with a hard-coded `parent_id = 2`; children are now attached to the `/system`
  menu by its actual ID.
- `admin:init` used `FirstOrCreate` for the built-in permissions, so an existing database kept the old `http_path`
  values (`/auth/roles`, …) after an upgrade and non-superadmins stayed blocked. It now compares and updates `name` and
  `http_path` for the eight built-in slugs when they differ.

### Changed

- **Breaking: the repository changed from a full skeleton app into an installable package.** Controllers, middleware,
  models, migrations, commands and routes moved from `app/**` to the package root, imports became
  `github.com/wcz0/goravel-admin/...`, and the host `main.go`, `bootstrap/` and `config/` skeletons moved into
  `example/`.
- Default configuration now ships with the package via `config.New()` and is deep-merged with host configuration during
  `ServiceProvider.Register` (host wins), so hosts no longer need to copy it.
- Admin routes, static assets and the embedded entry page are registered in `ServiceProvider.Boot`; the `/admin`,
  `/admin-assets` and `/admin-api` prefixes match the previous Gin routes.
- Migrations self-register in `ServiceProvider.Boot` via `facades.Schema().Register()`, fixing the compile failure the
  installer caused when it generated a host `bootstrap/migrations.go` where the package `migrations` clashed with the
  host `migrations`.
- Login, logout and token parsing all read `admin.auth.guard` (default `admin`), fixing tokens issued under the default
  guard (`sub` was `user`) being treated as unauthenticated.
- Upgraded `github.com/wcz0/gamis/v6` to v6.13.0, moved admin pages to the amis 6.13.0 mapping and updated controller
  imports.
- Synced the owl-admin v4.7.5 frontend sources and build output (amis 6.13.0) and refreshed the bundled assets.
- The login submit button and the user dropdown actions changed from `VanillaAction` (`type: button`) to
  `Button`/`Action` (`type: action`) to match amis 6.13.
- Removed sensitive database files: password hashes, backup scripts and outdated permission-system documents.
- Renamed the permission table column `value` to `slug`, updated schema and seed data accordingly, and added the
  service-agreement tables.
- Rewrote the README for the package installation flow and documented the Goravel version matrix and branch policy.
  The README is now available in English and Chinese.

### Removed

- Removed the `goravel-authz` dependency along with the leftover Casbin configuration, the `admin_rule` model, the
  `casbin_rules` migrations and the database dump. Permissions are fully handled by the built-in RBAC.
- Removed the "Get PHP Code" button and the `/dev_tools/editor` page. Both called `/dev_tools/editor_parse`, which had
  no backend implementation: Go has no equivalent of owl-admin's `amis()->Xxx()` PHP codegen, so the request returned
  an empty `200`. The amis schema editor itself (`custom-amis-editor`) is unchanged and still edits `admin_pages`
  schemas.
- Removed the code generator leftovers: the `admin_code_generators` table had no model, controller or route behind it,
  and came with a dead `/dev_tools/code_generator` menu entry and ~100 lines of unused copy in both language packs. The
  migration, the `admin:doctor` check and the menu entry are gone, so the migration count drops from 13 to 12.
- Removed `core.Menu.Extra()` and `core.Menu.devToolMenus()`. Nothing ever called `Extra()`, and the routes it built
  pointed at frontend pages that do not exist (`user_setting`, `/dev_tools/editor`, `/dev_tools/code_generator`).
- README no longer carries the "Differences from owl-admin" table. Runtime extension install, the PHP code generator
  and multi-module sub-admins are out of scope rather than gaps in the package, so the scope is stated in one
  paragraph instead of an inventory of missing features.

### Compatibility

- This branch targets `goravel/framework` v1.17.x. Differences from `main` (v1.18) are confined to the framework
  contract:
  - `middleware/*.go` return `func(ctx http.Context)` because `http.Middleware` is a function type in 1.17
  - validation rule literals use `map[string]string`, and the `anyRules()` conversion layer is gone
- No v1.16 line: Goravel 1.16 has no `foundation.Setup()`, no `packages.Setup.Paths()` and no
  `modify.RegisterProvider`, so the package-install mechanism does not exist there. Supporting it would require a
  separate app skeleton and installer.

## [1.0.0] - 2025-11-27

### Added

- Complete RBAC: users, roles, permissions and menus.
- Automatic permission generation and bulk deletion.
- Admin settings table and seeded menu configuration.
- File uploads and bundled frontend assets.
- Authorization through `goravel-authz`; authentication uses the built-in Goravel capabilities.

### Changed

- Database migrations were rewritten from SQL into Go migration files.
- Updated the README startup guide.
- Refined menu validation rules and delete button behaviour.

## [0.0.2] - 2024-12-03

### Added

- User list.

### Fixed

- Various fixes.

## [0.0.1] - 2024-10-02

### Added

- Initial release: base Goravel admin skeleton.
