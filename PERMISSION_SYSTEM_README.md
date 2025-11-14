# Goravel Admin 权限系统

这是一个基于 Goravel 框架构建的完整权限管理系统，支持基于角色的访问控制（RBAC）。

## 系统架构

### 数据库设计

#### 表结构

1. **admin_permissions** - 权限表
   - `id` - 主键
   - `parent_id` - 父权限ID（用于构建权限树）
   - `name` - 权限名称
   - `value` - 权限标识符
   - `http_method` - HTTP方法（JSON格式）
   - `http_path` - HTTP路径（JSON格式）
   - `custom_order` - 自定义排序
   - `created_at`, `updated_at` - 时间戳

2. **admin_roles** - 角色表
   - `id` - 主键
   - `name` - 角色名称
   - `slug` - 角色标识符
   - `created_at`, `updated_at` - 时间戳

3. **admin_role_permissions** - 角色权限关联表
   - `role_id` - 角色ID
   - `permission_id` - 权限ID
   - 唯一约束：(role_id, permission_id)

4. **admin_user_roles** - 用户角色关联表
   - `user_id` - 用户ID
   - `role_id` - 角色ID
   - 唯一约束：(user_id, role_id)

### 模型层

#### AdminPermission（权限模型）
```go
// 核心方法
func (p *AdminPermission) GetHttpMethods() []string
func (p *AdminPermission) GetHttpPaths() []string
```

#### AdminRole（角色模型）
```go
// 核心方法
func (r *AdminRole) HasUser(user *AdminUser) bool
func (r *AdminRole) PrimaryKey() interface{}
```

#### AdminUser（用户模型）
```go
// 核心方法
func (u *AdminUser) IsAdministrator() bool
func (u *AdminUser) AllPermissions() []AdminPermission
func (u *AdminUser) HasRole(role *AdminRole) bool
func (u *AdminUser) HasPermission(permission *AdminPermission) bool
```

### 服务层

#### AdminPermissionService（权限服务）
- `Store(ctx)` - 创建权限
- `Update(ctx)` - 更新权限
- `Destroy(ctx)` - 删除权限
- `List(ctx)` - 获取权限列表（树形结构）
- `Show(ctx)` - 获取权限详情

#### AdminRoleService（角色服务）
- `Store(ctx)` - 创建角色
- `Update(ctx)` - 更新角色
- `Destroy(ctx)` - 删除角色
- `List(ctx)` - 获取角色列表
- `Show(ctx)` - 获取角色详情

#### AdminUserService（用户服务）
- `Login(ctx)` - 用户登录
- `List(ctx)` - 获取用户列表
- `RoleOptions(ctx)` - 获取角色选项（用于下拉框）
- `Store(ctx)` - 创建用户

### 控制器层

#### PermissionController（权限控制器）
- 完整的CRUD操作
- HTTP方法和路径配置
- 权限树形展示
- 角色关联查询

#### RoleController（角色控制器）
- 完整的CRUD操作
- 权限分配和管理
- 用户角色关联

#### UserController（用户控制器）
- 完整的CRUD操作
- 角色选择和管理
- 用户设置

### 中间件

#### AdminPermissionMiddleware（权限中间件）
- 自动根据路由生成权限标识
- 支持手动权限检查
- 角色继承权限检查
- 超级管理员权限跳过

## 使用方法

### 1. 数据库迁移

```bash
# 运行迁移文件
go run artisan migrate
```

### 2. 权限验证

#### 在控制器中使用中间件
```go
func (c *Controller) GetUsers(ctx http.Context) error {
    // 中间件会自动验证权限
    // 权限标识格式：模块.操作（如：user.index, user.store）
    return ctx.Response().Json(200, map[string]interface{}{
        "users": []interface{}{},
    })
}
```

#### 手动权限检查
```go
// 在控制器或服务中
func (c *Controller) CheckPermission(user *AdminUser, permissionValue string) bool {
    return middleware.HasPermission(user, permissionValue)
}

// 在视图中
func (u *AdminUser) HasPermission(permission *AdminPermission) bool {
    return middleware.HasPermission(u, permission.Value)
}
```

### 3. 权限分配流程

#### 创建权限
```go
// 通过API创建权限
POST /api/admin/permissions
{
    "parent_id": 0,
    "name": "用户管理",
    "value": "user",
    "http_method": ["GET"],
    "http_path": ["/admin/users"],
    "custom_order": 1
}
```

#### 创建角色并分配权限
```go
// 创建角色
POST /api/admin/roles
{
    "name": "编辑员",
    "slug": "editor"
}

// 分配权限给角色
PUT /api/admin/roles/{id}/permissions
{
    "permission_ids": [1, 2, 3]
}
```

#### 分配角色给用户
```go
// 分配角色给用户
POST /api/admin/users
{
    "username": "editor1",
    "name": "编辑员用户",
    "password": "password",
    "admin_role_ids": [2] // 角色ID列表
}
```

### 4. 权限标识规范

#### 常见权限格式
```
# 模块操作权限
user.index - 查看用户列表
user.show - 查看用户详情
user.store - 创建用户
user.update - 更新用户
user.destroy - 删除用户

# 权限管理
permission.index - 查看权限列表
permission.show - 查看权限详情
permission.store - 创建权限
permission.update - 更新权限
permission.destroy - 删除权限

# 角色管理
role.index - 查看角色列表
role.show - 查看角色详情
role.store - 创建角色
role.update - 更新角色
role.destroy - 删除角色
```

### 5. 前端页面功能

#### 权限管理页面
- 权限树形结构展示
- 拖拽排序功能
- HTTP方法和路径配置
- 批量操作支持

#### 角色管理页面
- 角色信息维护
- 权限分配界面
- 用户角色关联

#### 用户管理页面
- 用户信息维护
- 角色选择（支持搜索和多选）
- 权限预览

## 特性说明

### 1. 权限继承
- 权限支持父子关系，形成树形结构
- 角色可以继承多个权限
- 用户通过角色获得权限

### 2. 超级管理员
- ID为1的用户默认为超级管理员
- 超级管理员拥有所有权限，跳过权限验证

### 3. 唯一性约束
- 权限名称和标识符必须唯一
- 角色名称和标识符必须唯一
- 用户名必须唯一

### 4. 数据验证
- HTTP方法和路径存储为JSON格式
- 支持数组格式的输入验证
- 完整的错误处理和提示

### 5. 灵活的配置
- 支持自定义排序
- 支持权限树形展示
- 支持多选权限分配

## 测试

运行测试：
```bash
# 运行权限系统测试
go test -v ./tests/feature/admin_permission_test.go
```

测试覆盖：
- 权限API功能
- 角色权限分配
- 用户角色分配
- 权限验证逻辑
- 中间件功能

## 注意事项

1. **数据库备份**：在进行权限结构调整前，请先备份数据库
2. **权限测试**：在生产环境中部署前，请充分测试权限功能
3. **性能考虑**：大量权限数据时，考虑使用缓存策略
4. **安全加固**：定期审查权限配置，确保安全

## 扩展建议

1. **缓存机制**：可以添加权限缓存，提高性能
2. **审计日志**：记录权限变更操作
3. **动态权限**：支持运行时动态加载权限
4. **权限组**：支持权限组管理，提高可维护性

---

*本权限系统基于 Goravel 框架构建，提供完整的RBAC权限管理解决方案。*