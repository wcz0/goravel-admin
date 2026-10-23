<div align="center">
<img src="https://www.goravel.dev/logo.png" width="300" alt="Goravel Admin">

[![Go](https://img.shields.io/github/go-mod/go-version/wcz0/goravel-admin)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/wcz0/goravel-admin)](https://github.com/wcz0/goravel-admin/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/wcz0/goravel-admin)](https://goreportcard.com/report/github.com/wcz0/goravel-admin)
[![License](https://img.shields.io/github/license/wcz0/goravel-admin)](./LICENSE)

</div>

# Goravel Admin

**English** | [简体中文](./README.zh-CN.md)

An installable admin panel for the [Goravel](https://github.com/goravel/goravel) framework, built on the
[amis](https://aisuda.bce.baidu.com/amis/) low-code renderer.

The backend speaks the exact same HTTP contract as [owl-admin](https://github.com/slowlyo/owl-admin) v4.7.5, so the
shipped owl-admin frontend build runs unchanged — no frontend code to write, no separate build pipeline.

- **Installable package** — register one service provider, migrations self-register, nothing else to wire up.
- **amis schema everywhere** — every list, form and detail page is generated on the server with
  [gamis](https://github.com/wcz0/gamis), so screens are pure Go.
- **Built-in RBAC** — users, roles, permissions and menus, no external authorization package.
- **Owl-admin parity** — captcha, chunked upload, xlsx export, custom pages, dynamic APIs, relationships and an
  extension registry are all implemented (see [Feature parity](#feature-parity)).

See [CHANGELOG.md](./CHANGELOG.md) for changes, and [docs/releases](./docs/releases) for per-version release notes.

## Version matrix

`goravel-admin` follows the Goravel framework version it targets (the same convention as `goravel/gin` and
`goravel/redis`).

| goravel-admin | goravel/framework | gamis | amis | owl-admin frontend |
| --- | --- | --- | --- | --- |
| v1.18.x (branch `main`) | v1.18.x | v6.13.x | 6.13.0 | v4.7.5 |
| v1.17.x (branch `goravel-1.17`, **this branch**) | v1.17.x | v6.13.x | 6.13.0 | v4.7.5 |

Each Goravel minor release changes its contract interfaces, so **every Goravel minor gets its own branch** instead of
one codebase trying to satisfy all of them:

- `main` tracks the latest `goravel/framework` (currently v1.18) and releases tags like `v1.18.x`.
- `goravel-1.17` targets the v1.17 contracts and releases `v1.17.x`.

> Verified: the v1.18 code base does not compile against v1.17 (`http.Middleware` flips back from an interface to
> `func(Context)`, and `Request().Validate` / `Validation().Make` take `map[string]string` again instead of
> `map[string]any`). Older minors must use their own branch.

### Why there is no `goravel-1.16` branch

Goravel 1.16 lacks the primitives the package-install mechanism depends on:

- `foundation.Setup()` does not exist (introduced in 1.17), so a host app cannot be assembled with
  `WithProviders` / `WithMigrations`.
- `packages.Setup` has no `Paths()`, and the `modify` package has no `RegisterProvider` / `RegisterMigration` — only
  the low-level `GoFile + match.Providers`.
- A 1.16 host registers providers in `config/app.go` via `config.Add("app", {"providers": [...]})`, so the installer
  would have to patch a different file.

Supporting 1.16 would therefore mean a separate application skeleton and a separate installer — a different project.
It is not supported.

## Requirements

- Go 1.25 or newer
- `goravel/framework` v1.18.x on `main`, v1.17.x on `goravel-1.17`
- Any database supported by Goravel (SQLite, MySQL, PostgreSQL, SQL Server)

## Getting started

### 1. Require the package

```shell
go get github.com/wcz0/goravel-admin@latest
```

### 2. Install (register the service provider, optionally publish config)

```shell
go run github.com/wcz0/goravel-admin/setup install
```

The installer:

- registers `&admin.ServiceProvider{}` in `bootstrap/providers.go`
- writes an annotated `config/admin.go` override skeleton **if** the file does not exist yet (never overwrites)
- **does not touch** `bootstrap/migrations.go` — migrations are mounted by the provider during `Boot` through
  `facades.Schema().Register`

To uninstall:

```shell
go run github.com/wcz0/goravel-admin/setup uninstall
```

### 3. Migrate and seed

```shell
go run . artisan key:generate
go run . artisan migrate
go run . artisan admin:init
go run .
```

Open `http://host:port/admin` and sign in with `admin` / `admin` (created by `admin:init` — change it immediately).

### 4. Optional: publish assets for customization

```shell
go run . artisan vendor:publish --package=github.com/wcz0/goravel-admin --tag=goravel-admin-lang
go run . artisan vendor:publish --package=github.com/wcz0/goravel-admin --tag=goravel-admin-assets
```

### Console commands

| Command | Purpose |
| --- | --- |
| `admin:init` | Create the administrator/operator roles, home menu and the default `admin` account |
| `admin:create-user` | Create an administrator (`--username`, `--password`, `--name`, `--roles`) |
| `admin:reset-password` | Reset a password (`--username`, `--password`; a random 12-char password is generated when omitted) |
| `admin:doctor` | Sanity-check `APP_KEY`, database connectivity, the 13 core tables, admin accounts, home menu uniqueness and broken menu parents |

## Configuration

Defaults ship with the package (`config/admin.go`), so hosts do not need to copy anything. Every key can be overridden
with an `ADMIN_*` environment variable, for example `ADMIN_ROUTE_PREFIX`, `ADMIN_ROUTE_PATH`, `ADMIN_ASSETS_PATH`,
`ADMIN_AUTH_GUARD`, `ADMIN_LOGIN_CAPTCHA`, `ADMIN_SHOW_DEVELOPMENT_TOOLS`.

`admin.auth.guard` defaults to `admin`. Login token issuing and middleware verification both read it, so do not change
`config/auth.go` alone and leave this behind.

## Example application

`example/` is a minimal host app that points at the local package through `replace`; it is how the whole
`setup install` flow is verified:

```shell
cd example
go run ../setup install
go run . artisan migrate && go run . artisan admin:init
go run .
```

## Feature parity

The frontend is the owl-admin v4.7.5 build, so the backend implements the same contract:

| Capability | Endpoint | Notes |
| --- | --- | --- |
| Login / logout | `POST /login`, `GET /logout` | Captcha enforced when enabled |
| Captcha | `GET /captcha` | Data-URI image plus a single-use session id; the answer lives in the cache |
| Current user | `GET /current-user` | Name, avatar and the top-bar dropdown schema |
| Profile | `GET/PUT /user_setting` | Avatar, name, password (old password required) |
| Menus | `GET /menus` | Permission filtered; route, external link, iframe and custom page types |
| Settings | `GET/POST /_settings` | Theme, locale, `enabled_extensions`, `show_development_tools` |
| Custom pages | `GET /page_schema?sign=` | Schema stored in `admin_pages`, used by menus with `url_type=4` |
| Icon search | `GET /_iconify_search` | 220k iconify names embedded in the binary |
| Upload | `upload_image`, `upload_file`, `upload_rich` | Uses `admin.upload.disk` with per-type directories |
| Chunked upload | `upload_chunk_start`, `upload_chunk`, `upload_chunk_finish` | md5 per part, target directory allow-list, script extensions rejected |
| Export | list `_action=export` → `_download_export` | xlsx generation, HMAC-signed expiring URL, traversal-safe |
| Admin management | `system/admin_users`, `admin_roles`, `admin_permissions`, `admin_menus` | RBAC built in |
| Dev tools | `dev_tools/pages`, `apis`, `relationships`, `extensions` | Requires `ADMIN_SHOW_DEVELOPMENT_TOOLS=true` |

### Differences from owl-admin

owl-admin runs on PHP; some of its features have no equivalent under Go's static compilation model. The honest list:

| owl-admin capability | This package |
| --- | --- |
| Runtime extension install (drop a zip, register PSR-4 on the fly) | **Not possible.** `dev_tools/extensions` only keeps a registry with enable/disable state; extension code must be registered at host compile time. `enabled_extensions` is still reported to the frontend. |
| Code generator (read tables → generate PHP model/controller/service/migration) | **Not implemented.** The output is PHP + Eloquent specific; a Go version would need struct/controller/migration generators — a project of its own. |
| Visual editor (amis JSON → `amis()->Xxx()` PHP chains) | **Not possible.** Go has no equivalent amis DSL; only the `/dev_tools/editor` frontend entry point is kept. |
| Dynamic API templates (write PHP classes at runtime and take effect instantly) | **Partial.** API records (path/template/args) are manageable, and the 10 built-in templates are metadata only; runtime dispatch is not wired yet. |
| Dynamic relationships (`resolveRelationUsing` at runtime) | **Partial.** Relationship configs are stored and `preview_code` is generated, but Goravel relationships live on the struct and do not take effect at runtime. |
| Multi-module / multiple sub-admins | **Not implemented.** For multi-tenant cases prefer a tenant column or separate deployments; duplicating all 7 tables per module is too costly. |
| `admin:ide-helper` (PHPStorm metadata) | **Not applicable.** gopls provides completion in Go; no stub files needed. |

## Testing

Unit tests plus an end-to-end contract smoke test live in the repo:

```shell
go test ./...
```

The smoke test drives every endpoint the owl-admin frontend relies on against a running instance:

```shell
cd example
go run . artisan migrate && go run . artisan admin:init
go run . artisan admin:create-user --username smoke --password 'Smoke@12345'
# in another shell
BASE=http://127.0.0.1:3000 python example/smoke_admin.py
```

## Releasing

Release notes are written in English and stored per version under [`docs/releases`](./docs/releases):

```shell
git tag -a v1.17.0 -m "v1.17.0"
git push origin v1.17.0
gh release create v1.17.0 --title "v1.17.0" --notes-file docs/releases/v1.17.0.md
```

## Credits

- [owl-admin](https://github.com/slowlyo/owl-admin) — the original PHP implementation and the frontend build
- [gamis](https://github.com/wcz0/gamis) — amis renderers for Go
- [amis](https://aisuda.bce.baidu.com/amis/) — the low-code renderer itself

## License

MIT. See [LICENSE](./LICENSE).
