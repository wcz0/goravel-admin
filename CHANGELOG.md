# Changelog

本项目遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/) 和[语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### Added

- 新增登录验证码：`GET /captcha` 返回 `captcha_img`（data URI）与 `sys_captcha`，答案存 Cache 并一次性消费；`admin.auth.login_captcha` 开启时登录接口强制校验，amis 登录页同步渲染验证码组件。
- 新增个人资料保存 `PUT /user_setting`：支持改头像、姓名与密码，改密必须校验旧密码，两次密码需一致。
- 新增 `GET /page_schema?sign=` 与「开发者工具 - 页面管理」：页面结构以 amis JSON 存 `admin_pages`，菜单 `url_type=4` 时由前端按 sign 拉取渲染。
- 新增「开发者工具 - 动态 API / 动态关联 / 扩展」三组接口：提供列表、详情、增删改、模板下拉、关联类型下拉与 `preview_code` 预览。
- 新增 `GET /_iconify_search`：内嵌 iconify 全量图标名（22 万条），首次访问解析后常驻内存，单次最多返回 1000 条。
- 新增分块上传三接口 `upload_chunk_start` / `upload_chunk` / `upload_chunk_finish`：会话 10 分钟过期，分片按 md5 校验，合并目标目录走白名单，脚本类后缀直接拒绝。
- 新增导出能力：列表「导出全部 / 本页 / 选中」按钮 + xlsx 生成 + `_download_export` 按 HMAC 签名与有效期下载，路径做穿越校验；用户列表已接入。
- 新增 `admin:create-user`、`admin:reset-password`、`admin:doctor` 命令。`admin:doctor` 巡检 APP_KEY、数据库连通性、13 张核心表、管理员账号、首页菜单唯一性与菜单父级断链。
- 新增 `core/captcha_test.go`，覆盖验证码「生成 → 正确校验 → 一次性失效」与空值/错误答案拒绝。
- 新增 `example/` 最小宿主应用，用于验证 `package:install` 全流程（注册 Provider、迁移、`admin:init`、登录与列表接口）。
- 新增 `setup install` / `setup uninstall`：安装器只注册 ServiceProvider 并生成可选配置骨架，不再改动宿主 `bootstrap/migrations.go`。
- 新增 `migrations.All()` 与 `tools.AuthGuard()`。

### Fixed

- 修复图片/文件上传的配置键拼接错误：原代码读 `admin.upload.directoryimage`，导致文件落到存储根目录，现按 `admin.upload.directory.<type>` 读取，图片进 `images/`、文件进 `files/`。
- 修复个人资料页保存地址错误：表单 `api` 原指向 `PUT /current-user`（该路由不存在），改为 `PUT /user_setting`。
- 修复 `_settings` 的 `enabled_extensions` 恒为空数组，改为读取 `admin_extensions` 中已启用的扩展名。
- 修复主键为空或记录不存在时仍会写库的问题：Goravel ORM 的 `Find`/`First` 在「id=0」和「查无记录」两种情况下都不报错（`First` 会吞掉 `ErrRecordNotFound`），原代码拿到零值结构体后直接 `Save`，会 INSERT 出一条空记录——实测 `POST /dev_tools/extensions/enable` 传 `id=0` 时写出了一条 `name` 为空的扩展并下发到 `enabled_extensions`。新增 `services.firstByID`（先 `Count` 再取），页面、动态 API、动态关联、扩展、菜单、权限、用户的单条读写统一走它，并新增 `services/preview_test.go` 覆盖。
- 修复 `POST /dev_tools/pages` 在 `schema` 直接以 JSON 对象提交时被判为"页面结构不可为空"：现在 `page.schema` 嵌套对象、`schema` 对象、`schema` JSON 文本三种形态都接受。
- 修复动态关联 `preview_code` 生成的片段无法编译：`model` / `title` 原样拼进标识符（`Smoke Roles` 带空格），且返回了 Goravel 并不存在的 `orm.Relation`。现统一转成合法 Go 标识符，并改为 Goravel ORM 实际支持的显式查询骨架（中间表走 `Join`，一对多走 `Where`），表单默认参数同步为 `related` / `table` / `foreign_key` / `related_key`。

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

- 本分支对应 goravel/framework v1.18.x。
- Goravel 1.17 由 `goravel-1.17` 分支提供（已跑通安装流程），与 main 的差异：中间件改为 `func(ctx http.Context)` 闭包、校验规则字面量改为 `map[string]string`、移除 `anyRules()` 转换层。
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
