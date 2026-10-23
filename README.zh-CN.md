<div align="center">
<img src="https://www.goravel.dev/logo.png" width="300" alt="Goravel Admin">

[![Go](https://img.shields.io/github/go-mod/go-version/wcz0/goravel-admin)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/wcz0/goravel-admin)](https://github.com/wcz0/goravel-admin/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/wcz0/goravel-admin)](https://goreportcard.com/report/github.com/wcz0/goravel-admin)
[![License](https://img.shields.io/github/license/wcz0/goravel-admin)](./LICENSE)

</div>

# Goravel Admin

[English](./README.md) | **简体中文**

基于 [amis](https://aisuda.bce.baidu.com/amis/) 低代码渲染器的 [Goravel](https://github.com/goravel/goravel) 后台扩展包。

后端接口与 [owl-admin](https://github.com/slowlyo/owl-admin) v4.7.5 保持同一套契约，因此可以直接使用 owl-admin 的
前端构建产物——不需要写一行前端代码，也不需要额外的前端构建流程。

- **可安装扩展包**：注册一个 ServiceProvider 即可，迁移在 Boot 阶段自注册，宿主无需额外接线。
- **全 amis schema**：列表、表单、详情页都由服务端用 [gamis](https://github.com/wcz0/gamis) 生成，页面即 Go 代码。
- **内置 RBAC**：用户、角色、权限、菜单，不依赖外部授权包。
- **对齐 owl-admin**：验证码、分块上传、xlsx 导出、自定义页面、动态 API、动态关联与扩展名单均已实现，
  差异见[与 owl-admin 的差异](#与-owl-admin-的差异)。

变更记录见 [CHANGELOG.md](./CHANGELOG.md)（英文），各版本发布说明见 [docs/releases](./docs/releases)。

## 版本对应

`goravel-admin` 以**可安装扩展包**的形式提供，版本号与它适配的 Goravel 框架对齐（与 `goravel/gin`、
`goravel/redis` 的做法一致）。

| goravel-admin | goravel/framework | gamis | amis | owl-admin 前端 |
| --- | --- | --- | --- | --- |
| v1.18.x（分支 `main`） | v1.18.x | v6.13.x | 6.13.0 | v4.7.5 |
| v1.17.x（分支 `goravel-1.17`，**本分支**） | v1.17.x | v6.13.x | 6.13.0 | v4.7.5 |

Goravel 每个 minor 版本都会调整契约接口，因此**不同 goravel 版本由不同分支维护**，而不是用一套代码兼容全部：

- `main`：跟随最新的 `goravel/framework`（当前 v1.18），发布 tag 形如 `v1.18.x`
- `goravel-1.17`：适配 v1.17 契约，发布 `v1.17.x`

> 已实测：v1.18 的代码直接降到 v1.17 无法编译（`http.Middleware` 由接口改回 `func(Context)`，
> `Request().Validate` 与 `Validation().Make` 的入参由 `map[string]any` 变回 `map[string]string`），
> 所以低版本必须走对应分支。

### 为什么没有 goravel-1.16

Goravel 1.16 缺少支撑「扩展包」机制的前置能力，无法用同一套安装方式支持：

- `foundation.Setup()` 不存在（1.17 才引入），宿主无法用 `WithProviders / WithMigrations` 组装应用
- `packages.Setup` 没有 `Paths()`，`modify` 包也没有 `RegisterProvider / RegisterMigration`，只有底层的
  `GoFile + match.Providers`
- 1.16 的宿主在 `config/app.go` 里以 `config.Add("app", {"providers": [...]})` 注册 Provider，安装器要改的
  文件位置也不同

因此 1.16 需要独立的应用骨架与安装器实现，属于另一套工程，当前不支持。

## 环境要求

- Go 1.25 及以上
- `main` 分支对应 `goravel/framework` v1.18.x，`goravel-1.17` 分支对应 v1.17.x
- Goravel 支持的任一数据库（SQLite、MySQL、PostgreSQL、SQL Server）

## 快速开始

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

### 控制台命令

| 命令 | 用途 |
| --- | --- |
| `admin:init` | 创建管理员/操作员角色、首页菜单与默认 `admin` 账号 |
| `admin:create-user` | 创建管理员（`--username`、`--password`、`--name`、`--roles`） |
| `admin:reset-password` | 重置密码（`--username`、`--password`；不传密码时随机生成 12 位） |
| `admin:doctor` | 巡检 `APP_KEY`、数据库连通性、13 张核心表、管理员账号、首页菜单唯一性与菜单父级断链 |

## 配置

默认配置由扩展包提供（见 `config/admin.go`），宿主无需复制。所有键都可以用 `ADMIN_*` 环境变量覆盖，例如
`ADMIN_ROUTE_PREFIX`、`ADMIN_ROUTE_PATH`、`ADMIN_ASSETS_PATH`、`ADMIN_AUTH_GUARD`、`ADMIN_LOGIN_CAPTCHA`、
`ADMIN_SHOW_DEVELOPMENT_TOOLS`。

`admin.auth.guard` 默认是 `admin`，登录签发的 token 与中间件校验都读取这个值，宿主不要只改 `config/auth.go`
而漏改这里。

## 示例项目

`example/` 是一个最小宿主应用，通过 `replace` 指向本地包，用来验证 `setup install` 全流程：

```shell
cd example
go run ../setup install
go run . artisan migrate && go run . artisan admin:init
go run .
```

## 可用功能

前端资源直接来自 owl-admin v4.7.5，因此后端接口与 owl-admin 保持同一套契约：

| 能力 | 接口 | 说明 |
| --- | --- | --- |
| 登录 / 登出 | `POST /login`、`GET /logout` | 验证码开启时强制校验 |
| 验证码 | `GET /captcha` | 返回 data URI 图片 + 一次性会话 id，答案存 Cache |
| 当前用户 | `GET /current-user` | 返回姓名、头像与顶部下拉菜单 schema |
| 个人资料 | `GET/PUT /user_setting` | 改头像、姓名、密码（需旧密码） |
| 菜单 | `GET /menus` | 按权限过滤，支持路由/外链/iframe/自定义页面四种类型 |
| 系统设置 | `GET/POST /_settings` | 主题、语言等，含 `enabled_extensions` 与 `show_development_tools` |
| 自定义页面 | `GET /page_schema?sign=` | 页面结构存 `admin_pages`，菜单 `url_type=4` 时使用 |
| 图标搜索 | `GET /_iconify_search` | 内嵌 22 万条 iconify 图标名 |
| 上传 | `upload_image`、`upload_file`、`upload_rich` | 走 `admin.upload.disk` 与按类型分目录 |
| 分块上传 | `upload_chunk_start`、`upload_chunk`、`upload_chunk_finish` | md5 校验 + 目标目录白名单 + 脚本后缀拦截 |
| 导出 | 列表 `_action=export` → `_download_export` | xlsx 生成 + HMAC 签名 + 有效期，防路径穿越 |
| 后台管理 | `system/admin_users`、`admin_roles`、`admin_permissions`、`admin_menus` | 内置 RBAC |
| 开发者工具 | `dev_tools/pages`、`apis`、`relationships`、`extensions` | 需 `ADMIN_SHOW_DEVELOPMENT_TOOLS=true` |

### 与 owl-admin 的差异

owl-admin 是 PHP 运行时，部分能力在 Go 的静态编译模型下没有等价实现，下面是明确的差异：

| owl-admin 能力 | 本包的处理 |
| --- | --- |
| 扩展动态安装（zip 投放 + 运行时注册 PSR-4） | **不可行**。`dev_tools/extensions` 只维护「扩展名单 + 启停状态」，扩展代码需在宿主编译期注册；`enabled_extensions` 会照常下发给前端 |
| 代码生成器（读表 → 生成 PHP Model/Controller/Service/Migration） | **未实现**。产物是 PHP + Eloquent 专属，Go 侧需重写为 struct/控制器/迁移生成器，属独立工程 |
| 可视化编辑器（amis JSON → `amis()->Xxx()` PHP 链式代码） | **不可行**。Go 没有等价的 amis DSL，只保留了 `/dev_tools/editor` 的前端入口 |
| 动态 API 模板（运行时写 PHP 类并即时生效） | **部分实现**。接口记录（路径/模板/参数）可维护，10 个内置模板只作为元信息，运行时分发尚未接入 |
| 动态关联（Eloquent `resolveRelationUsing` 运行时挂载） | **部分实现**。可登记关联配置并给出 `preview_code`，Goravel 的关联需写进 struct，不会即时生效 |
| 多模块 / 多子后台 | **未实现**。单后台场景建议用租户字段或独立部署，复制整套 7 张表成本过高 |
| `admin:ide-helper`（PHPStorm 元数据） | **不适用**。Go 由 gopls 提供补全，无需生成桩文件 |

## 测试

仓库内包含单元测试与端到端契约冒烟脚本：

```shell
go test ./...
```

冒烟脚本会把 owl-admin 前端依赖的接口逐个打一遍：

```shell
cd example
go run . artisan migrate && go run . artisan admin:init
go run . artisan admin:create-user --username smoke --password 'Smoke@12345'
# 另开一个终端
BASE=http://127.0.0.1:3000 python example/smoke_admin.py
```

## 发布

发布说明用英文撰写，按版本存放在 [`docs/releases`](./docs/releases)：

```shell
git tag -a v1.17.0 -m "v1.17.0"
git push origin v1.17.0
gh release create v1.17.0 --title "v1.17.0" --notes-file docs/releases/v1.17.0.md
```

## 致谢

- [owl-admin](https://github.com/slowlyo/owl-admin) — 原始 PHP 实现与前端构建产物
- [gamis](https://github.com/wcz0/gamis) — Go 版 amis 渲染器
- [amis](https://aisuda.bce.baidu.com/amis/) — 低代码渲染器本体

## 许可证

MIT，详见 [LICENSE](./LICENSE)。
