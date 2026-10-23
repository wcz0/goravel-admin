<div align="center">
<img src="https://www.goravel.dev/logo.png" width="300" alt="Logo">

[![Doc](https://pkg.go.dev/badge/github.com/goravel/framework)](https://pkg.go.dev/github.com/goravel/framework)
[![Go](https://img.shields.io/github/go-mod/go-version/goravel/framework)](https://go.dev/)
[![Release](https://img.shields.io/github/release/goravel/framework.svg)](https://github.com/goravel/framework/releases)
[![Test](https://github.com/goravel/framework/actions/workflows/test.yml/badge.svg)](https://github.com/goravel/framework/actions)
[![Report Card](https://goreportcard.com/badge/github.com/goravel/framework)](https://goreportcard.com/report/github.com/goravel/framework)
[![Codecov](https://codecov.io/gh/goravel/framework/branch/master/graph/badge.svg)](https://codecov.io/gh/goravel/framework)
![License](https://img.shields.io/github/license/goravel/framework)</div>

## About Goravel-Admin

使用百度amis 低代码平台, 来快速创建后台

无须你编写任何前端代码, 就能创建一个好看的后台页面

变更记录见 [CHANGELOG.md](./CHANGELOG.md)。

## 版本对应

goravel-admin 以**可安装扩展包**的形式提供，版本号与它适配的 Goravel 框架对齐（与 `goravel/gin`、`goravel/redis` 的做法一致）。

| goravel-admin | goravel/framework | gamis | amis | owl-admin 前端 |
| --- | --- | --- | --- | --- |
| v1.18.x（main） | v1.18.x | v6.13.x | 6.13.0 | v4.7.5 |
| v1.17.x（branch `goravel-1.17`） | v1.17.x | v6.13.x | 6.13.0 | v4.7.5 |

Goravel 每个 minor 版本都会调整契约接口，因此**不同 goravel 版本由不同分支维护**，而不是用一套代码兼容全部：

- `main`：跟随最新的 goravel/framework（当前 v1.18），发布 tag 形如 `v1.18.x`
- `goravel-1.17`：适配 v1.17 契约，发布 `v1.17.x`（已跑通安装流程，与 goravel-authz 的分支约定一致）

> 已实测：v1.18 的代码直接降到 v1.17 无法编译（`http.Middleware` 由接口改回 `func(Context)`，`Request().Validate` 与 `Validation().Make` 的入参由 `map[string]any` 变回 `map[string]string`），所以低版本必须走对应分支。

### 为什么没有 goravel-1.16

Goravel 1.16 缺少支撑"扩展包"机制的前置能力，无法用同一套安装方式支持：

- `foundation.Setup()` 不存在（1.17 才引入），宿主无法用 `WithProviders / WithMigrations` 组装应用
- `packages.Setup` 没有 `Paths()`，`modify` 包也没有 `RegisterProvider / RegisterMigration`，只有底层的 `GoFile + match.Providers`
- 1.16 的宿主在 `config/app.go` 里以 `config.Add("app", {"providers": [...]})` 注册 Provider，安装器要改的文件位置也不同

因此 1.16 需要独立的应用骨架与安装器实现，属于另一套工程，当前不支持。

## Getting started

### 1. 引入扩展包

```shell
go get github.com/wcz0/goravel-admin@latest
```

### 2. 安装（注册 ServiceProvider + 可选配置骨架）

```shell
go run github.com/wcz0/goravel-admin/setup install
```

安装做了什么：

- 在 `bootstrap/providers.go` 里注册 `&admin.ServiceProvider{}`
- 若 `config/admin.go` 不存在，生成一份带说明的覆盖骨架（已存在则不覆盖）
- **不会改动** `bootstrap/migrations.go`：迁移由 Provider 在 Boot 阶段通过 `facades.Schema().Register` 自动挂载

卸载：

```shell
go run github.com/wcz0/goravel-admin/setup uninstall
```

### 3. 跑迁移并初始化后台数据

```shell
go run . artisan key:generate
go run . artisan migrate
go run . artisan admin:init
go run . 
```

打开 `http://host:port/admin`，登录账号 `admin`，密码 `admin`（`admin:init` 创建，请及时修改）。

### 4. 可选：发布资源以便自定义

```shell
go run . artisan vendor:publish --package=github.com/wcz0/goravel-admin --tag=goravel-admin-lang
go run . artisan vendor:publish --package=github.com/wcz0/goravel-admin --tag=goravel-admin-assets
```

## 配置

默认配置由扩展包提供（见 `config/admin.go`），宿主无需复制。所有键都可以用 `ADMIN_*` 环境变量覆盖，例如 `ADMIN_ROUTE_PREFIX`、`ADMIN_ROUTE_PATH`、`ADMIN_ASSETS_PATH`、`ADMIN_AUTH_GUARD`、`ADMIN_LOGIN_CAPTCHA`。

`admin.auth.guard` 默认是 `admin`，登录签发的 token 与中间件校验都读取这个值，宿主不要只改 `config/auth.go` 而漏改这里。

## 示例项目

`example/` 是一个最小宿主应用，通过 `replace` 指向本地包，用来验证 `package:install` 全流程：

```shell
cd example
go run ../setup install
go run . artisan migrate && go run . artisan admin:init
go run .
```

## Roadmap

[For Detail](https://github.com/goravel/goravel/issues?q=is%3Aissue+is%3Aopen)

权限相关由内置 RBAC（角色/权限/菜单）实现，不依赖外部授权包（已移除 goravel-authz）。
