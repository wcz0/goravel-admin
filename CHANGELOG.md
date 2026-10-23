# Changelog

本项目遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/) 和[语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### Added

- **新增 `goravel-1.17` 支持线（本分支）**，适配 goravel/framework v1.17.3；`example/go.mod` 同步 `gin v1.17.1`、`redis v1.17.1`、`sqlite v1.17.0`。
- 新增 `example/` 最小宿主应用，用于验证 `package:install` 全流程（注册 Provider、迁移、`admin:init`、登录与列表接口）。
- 新增 `setup install` / `setup uninstall`：安装器只注册 ServiceProvider 并生成可选配置骨架，不再改动宿主 `bootstrap/migrations.go`。
- 新增 `migrations.All()` 与 `tools.AuthGuard()`。

### Changed

- **破坏性：仓库形态由完整 skeleton app 改为可安装扩展包。** `app/**` 下的控制器、中间件、模型、迁移、命令、路由整体上移到包根目录，import 路径统一为 `github.com/wcz0/goravel-admin/...`；宿主的 `main.go`、`bootstrap/`、`config/` 骨架移入 `example/`。
- 默认配置改由包内 `config.New()` 提供，并在 `ServiceProvider.Register` 阶段与宿主配置深度合并（宿主优先），宿主不再需要拷贝配置。
- 后台路由、静态资源、内嵌入口页在 `ServiceProvider.Boot` 注册，`/admin` 与 `/admin-assets`、`/admin-api` 前缀与原 Gin 路由保持一致。
- 迁移改为 `ServiceProvider.Boot` 阶段通过 `facades.Schema().Register()` 自注册，修复安装器为宿主生成的 `bootstrap/migrations.go` 中包内 `migrations` 与宿主 `migrations` 同名导致无法编译的问题。
- 登录、注销、token 解析统一读取 `admin.auth.guard`（默认 `admin`），修复登录走默认 guard 导致签发 token 的 `sub` 为 `user`、带 token 请求被判未认证的问题。
- 升级 `github.com/wcz0/gamis/v6` 到 v6.13.0，后台页面切换到 amis 6.13.0 映射，并同步控制器 import 路径。
- 同步 owl-admin v4.7.5 前端源码与构建产物（amis 6.13.0），更新后台静态资源。
- 登录页提交按钮与用户下拉菜单按钮从 `VanillaAction`（输出 `type: button`）改为 `Button`/`Action`（输出 `type: action`），与 amis 6.13 的按钮类型对齐。
- 清理数据库敏感文件：密码哈希、备份脚本和过时的权限系统文档。
- 权限表 `value` 字段重命名为 `slug`，同步更新数据库结构和测试数据，并补充服务协议相关表。
- README 重写为扩展包安装方式，并补充 goravel 版本对应矩阵与分支策略。

### Removed

- 移除 `goravel-authz` 依赖，以及遗留的 Casbin 配置、`admin_rule` 模型、`casbin_rules` 相关迁移和数据库导出文件；权限完全由内置 RBAC 实现。

### Compatibility

- 本分支对应 goravel/framework v1.17.x。与 `main`（v1.18）的差异仅在框架契约层：
  - `middleware/*.go` 改为返回 `func(ctx http.Context)`（1.17 的 `http.Middleware` 是函数类型）
  - 校验规则字面量改为 `map[string]string`，并移除 `anyRules()` 转换层
- 不提供 v1.16 支持线：Goravel 1.16 没有 `foundation.Setup()`，也没有 `packages.Setup.Paths()` 与 `modify.RegisterProvider`，扩展包安装机制在 1.16 上不存在，需要独立的应用骨架与安装器实现。

## [1.0.0] - 2025-11-27

### Added

- 完整的 RBAC 权限管理：用户、角色、权限、菜单。
- 权限自动生成与批量删除。
- 管理员设置表与菜单初始化配置。
- 文件上传与前端静态资源。
- 基于 goravel-authz 的权限控制；登录使用 Goravel 自带能力。

### Changed

- 数据库迁移由 SQL 重构为 Go 迁移文件。
- 更新 README 启动指南。
- 优化菜单验证规则和删除按钮逻辑。

## [0.0.2] - 2024-12-03

### Added

- 用户列表。

### Fixed

- 修复若干问题。

## [0.0.1] - 2024-10-02

### Added

- 初始版本：Goravel 后台基础骨架。
