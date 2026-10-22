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

## Getting started

复制.env.example 为 .env
修改 .env 中的配置
创建数据库, 然后执行迁移和初始化数据
```shell
go run . artisan migrate
go run . artisan db:seed
```


```
// Generate APP_KEY
go run . artisan key:generate
// 运行项目
go run .
```

运行项目, 打开 host:port/admin

登录
账号:admin
密码:admin


## Roadmap

[For Detail](https://github.com/goravel/goravel/issues?q=is%3Aissue+is%3Aopen)
权限相关由内置 RBAC（角色/权限/菜单）实现，不依赖外部授权包
