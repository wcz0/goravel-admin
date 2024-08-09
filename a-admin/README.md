# amis 后台前端实现

**Owl-admin 后台思路**

使用amis低代码框架 实现 后台快速开发

可以使用各种语言实现后台api

### 后端库使用
Go amis实现 [https://github.com/wcz0/gamis](https://github.com/wcz0/gamis)
PHP amis 实现 [https://github.com/slowlyo/amis-renderer](https://github.com/slowlyo/amis-renderer)
    php 推荐使用完整项目 [https://github.com/slowlyo/owl-admin](https://github.com/slowlyo/owl-admin)

### 后端需要实现的api接口文档
Apifox地址 [https://owl-admin.apifox.cn](https://owl-admin.apifox.cn)

具体使用如何使用, 文档地址 [https://doc.owladmin.com/#/examples/custom-login](https://doc.owladmin.com/#/examples/custom-login)

demo地址 [http://demo.owladmin.com/admin](http://demo.owladmin.com/admin)

amis 组件文档地址 [https://aisuda.bce.baidu.com/amis/zh-CN/components/page](https://aisuda.bce.baidu.com/amis/zh-CN/components/page)


### 其他定义
可以修改前端代码


返回格式

```json
{
  "status": 0,
  "msg": "Login Success",
  "doNotDisplayToast": 0,
  "data": null
}
```

```


// 初始化项目
pnpm install

// 开发模式
pnpm run dev

// 构建
pnpm run build
```
