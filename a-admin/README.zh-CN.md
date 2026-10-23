# amis 后台前端实现

[English](./README.md) | **简体中文**

本目录是 [owl-admin](https://github.com/slowlyo/owl-admin) v4.7.5 的前端源码（amis 6.13.0），保留在仓库里是为了让
构建过程可复现。

正常使用 `goravel-admin` **不需要**构建前端：编译产物已提交在 [`../assets/dist`](../assets/dist)，并通过
[`../resources.go`](../resources.go) 的 `//go:embed` 打进 Go 二进制。只有改动前端代码时才需要重新构建：

```shell
# 安装依赖
pnpm install

# 开发模式
pnpm run dev

# 生产构建，并把产物同步到被 embed 的目录
pnpm run build
rm -rf ../assets/dist && cp -r dist ../assets/dist
```

## Owl-admin 后台思路

使用 amis 低代码框架实现后台快速开发：后端只负责输出 JSON schema，前端负责渲染，因此任何语言都能实现同一套后台
UI。

### 后端库使用

- Go：[gamis](https://github.com/wcz0/gamis) —— 本包使用的实现
- PHP：[amis-renderer](https://github.com/slowlyo/amis-renderer)，或完整项目
  [owl-admin](https://github.com/slowlyo/owl-admin)

### 接口文档

- Apifox：<https://owl-admin.apifox.cn>
- 使用文档：<https://doc.owladmin.com/#/examples/custom-login>
- Demo：<http://demo.owladmin.com/admin>
- amis 组件文档：<https://aisuda.bce.baidu.com/amis/zh-CN/components/page>

### 返回格式

所有 JSON 接口统一返回：

```json
{
  "status": 0,
  "msg": "Login Success",
  "doNotDisplayToast": 0,
  "data": null
}
```

`status = 0` 表示成功，其余值按错误处理并把 `msg` 以 toast 形式提示；`doNotDisplayToast = 1` 表示不弹提示，
`page_schema` 这类「查不到也算正常」的接口会用到。
