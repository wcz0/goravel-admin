<div align="center">

# 项目还处于wip状态. 

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

## Getting started

运行项目, 打开 localhost:8080/admin

```
// Generate APP_KEY
go run . artisan key:generate
// 运行项目
go run .
```



执行命令, 更新前端资源 [skip]

```shell
// make admin frontend static file
cd a-admin
// install front library
pnpm i
// generate static file
pnpm build
// copy file to Public folder
cp -r dist/* ../public/admin/
```

## Roadmap

[For Detail](https://github.com/goravel/goravel/issues?q=is%3Aissue+is%3Aopen)
使用 goravel-authz 来编写权限相关
