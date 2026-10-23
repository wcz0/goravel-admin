# amis admin frontend

**English** | [简体中文](./README.zh-CN.md)

This directory holds the [owl-admin](https://github.com/slowlyo/owl-admin) v4.7.5 frontend sources (amis 6.13.0), kept
in-repo so the build can be reproduced.

`goravel-admin` does **not** require a frontend build for normal use: the compiled output is committed under
[`../assets/dist`](../assets/dist) and embedded into the Go binary through `//go:embed` in [`../resources.go`](../resources.go).
Rebuild only when you change frontend code:

```shell
# install dependencies
pnpm install

# dev server
pnpm run dev

# production build, then copy the result into the embedded directory
pnpm run build
rm -rf ../assets/dist && cp -r dist ../assets/dist
```

## Owl-admin in a nutshell

A backend built with amis is a set of JSON schemas served by an API; the frontend renders them, so the same admin UI can
be produced by any backend language.

### Backend libraries

- Go: [gamis](https://github.com/wcz0/gamis) — used by this package
- PHP: [amis-renderer](https://github.com/slowlyo/amis-renderer), or the full project
  [owl-admin](https://github.com/slowlyo/owl-admin)

### API reference

- Apifox: <https://owl-admin.apifox.cn>
- Usage docs: <https://doc.owladmin.com/#/examples/custom-login>
- Demo: <http://demo.owladmin.com/admin>
- amis component docs: <https://aisuda.bce.baidu.com/amis/zh-CN/components/page>

### Response envelope

Every JSON endpoint answers with this shape:

```json
{
  "status": 0,
  "msg": "Login Success",
  "doNotDisplayToast": 0,
  "data": null
}
```

`status = 0` means success; anything else is treated as an error and `msg` is shown as a toast. `doNotDisplayToast = 1`
suppresses that toast — used by endpoints such as `page_schema`, where a missing page is expected.
