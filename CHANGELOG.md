# Changelog

本项目遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/) 和[语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### Changed

- 升级 `github.com/wcz0/gamis/v6` 到 v6.13.0，后台页面切换到 amis 6.13.0 映射，并同步控制器 import 路径。
- 清理数据库敏感文件：密码哈希、备份脚本和过时的权限系统文档。
- 权限表 `value` 字段重命名为 `slug`，同步更新数据库结构和测试数据，并补充服务协议相关表。

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
