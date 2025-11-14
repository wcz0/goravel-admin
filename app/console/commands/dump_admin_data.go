package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"goravel/app/models/admin"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

type DumpAdminData struct{}

// Signature 命令的名称和签名
func (receiver *DumpAdminData) Signature() string {
	return "admin:dump-data"
}

// Description 命令的描述
func (receiver *DumpAdminData) Description() string {
	return "导出数据库中的admin数据为初始化文件"
}

func (receiver *DumpAdminData) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
	}
}

// Handle 执行命令的逻辑
func (receiver *DumpAdminData) Handle(ctx console.Context) error {
	// 生成初始化数据文件
	if err := receiver.generateSeedDataFile(); err != nil {
		return err
	}

	fmt.Println("admin数据导出成功")
	return nil
}

// generateSeedDataFile 生成初始化数据文件
func (receiver *DumpAdminData) generateSeedDataFile() error {
	// 创建目录
	seedDir := "database/seeders"
	if err := os.MkdirAll(seedDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 查询所有数据
	users, err := receiver.getAllUsers()
	if err != nil {
		return fmt.Errorf("查询用户数据失败: %w", err)
	}

	roles, err := receiver.getAllRoles()
	if err != nil {
		return fmt.Errorf("查询角色数据失败: %w", err)
	}

	permissions, err := receiver.getAllPermissions()
	if err != nil {
		return fmt.Errorf("查询权限数据失败: %w", err)
	}

	menus, err := receiver.getAllMenus()
	if err != nil {
		return fmt.Errorf("查询菜单数据失败: %w", err)
	}

	settings, err := receiver.getAllSettings()
	if err != nil {
		return fmt.Errorf("查询设置数据失败: %w", err)
	}

	roleUsers, err := receiver.getAllRoleUsers()
	if err != nil {
		return fmt.Errorf("查询角色用户关联数据失败: %w", err)
	}

	rolePermissions, err := receiver.getAllRolePermissions()
	if err != nil {
		return fmt.Errorf("查询角色权限关联数据失败: %w", err)
	}

	permissionMenus, err := receiver.getAllPermissionMenus()
	if err != nil {
		return fmt.Errorf("查询权限菜单关联数据失败: %w", err)
	}

	// 生成代码文件
	return receiver.generateSeedFile(users, roles, permissions, menus, settings, roleUsers, rolePermissions, permissionMenus)
}

// getAllUsers 获取所有用户
func (receiver *DumpAdminData) getAllUsers() ([]*admin.AdminUser, error) {
	var users []*admin.AdminUser
	if err := facades.Orm().Query().Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// getAllRoles 获取所有角色
func (receiver *DumpAdminData) getAllRoles() ([]*admin.AdminRole, error) {
	var roles []*admin.AdminRole
	if err := facades.Orm().Query().Find(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

// getAllPermissions 获取所有权限
func (receiver *DumpAdminData) getAllPermissions() ([]*admin.AdminPermission, error) {
	// 使用原生SQL查询，手动处理Python格式的JSON
	query := `
		SELECT 
			id, name, slug, 
			http_method, http_path, custom_order, parent_id, 
			created_at, updated_at 
		FROM admin_permissions`
	
	var results []map[string]interface{}
	if err := facades.Orm().Query().Raw(query).Scan(&results); err != nil {
		return nil, err
	}
	
	var permissions []*admin.AdminPermission
	for _, result := range results {
		permission := &admin.AdminPermission{
			Name:        result["name"].(string),
			Slug:        result["slug"].(string),
		}
		
		// 安全地转换和设置字段
		if customOrder, ok := result["custom_order"]; ok && customOrder != nil {
			permission.CustomOrder = getIntFromInterface(customOrder)
		}
		if parentId, ok := result["parent_id"]; ok && parentId != nil {
			permission.ParentId = getUintFromInterface(parentId)
		}
		
		// 设置ID
		if id, ok := result["id"]; ok && id != nil {
			// 通过反射设置ID字段
			setField(permission, "ID", id)
		}
		
		// 设置时间字段
		if createdAt, ok := result["created_at"]; ok && createdAt != nil {
			setField(permission, "CreatedAt", createdAt)
		}
		if updatedAt, ok := result["updated_at"]; ok && updatedAt != nil {
			setField(permission, "UpdatedAt", updatedAt)
		}
		
		// 转换Python格式的JSON为Go格式
		if httpMethod, ok := result["http_method"]; ok && httpMethod != nil {
			if str, ok := httpMethod.(string); ok {
				permission.HttpMethod = receiver.convertPythonListToGoSlice(str)
			}
		}
		if httpPath, ok := result["http_path"]; ok && httpPath != nil {
			if str, ok := httpPath.(string); ok {
				permission.HttpPath = receiver.convertPythonListToGoSlice(str)
			}
		}
		
		permissions = append(permissions, permission)
	}
	
	return permissions, nil
}

// getAllMenus 获取所有菜单
func (receiver *DumpAdminData) getAllMenus() ([]*admin.AdminMenu, error) {
	var menus []*admin.AdminMenu
	if err := facades.Orm().Query().Find(&menus); err != nil {
		return nil, err
	}
	return menus, nil
}

// convertPythonListToGoSlice 将Python格式的列表转换为Go slice
func (receiver *DumpAdminData) convertPythonListToGoSlice(pythonList string) []string {
	if pythonList == "" {
		return []string{}
	}
	
	// 使用正则表达式提取Python列表中的字符串
	re := regexp.MustCompile(`'([^']*)'`)
	matches := re.FindAllStringSubmatch(pythonList, -1)
	
	var result []string
	for _, match := range matches {
		if len(match) > 1 {
			result = append(result, match[1])
		}
	}
	
	return result
}

// setField 使用反射设置结构体私有字段，并处理类型转换
func setField(obj interface{}, fieldName string, value interface{}) {
	v := reflect.ValueOf(obj).Elem()
	field := v.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		// 跳过时间类型字段，因为它们需要特殊处理
		if fieldName == "CreatedAt" || fieldName == "UpdatedAt" {
			return
		}
		// 处理类型转换
		convertedValue := convertValueType(value, field.Type())
		field.Set(convertedValue)
	}
}

// convertValueType 将值转换为目标类型
func convertValueType(value interface{}, targetType reflect.Type) reflect.Value {
	// 如果已经是目标类型，直接返回
	if reflect.TypeOf(value) == targetType {
		return reflect.ValueOf(value)
	}
	
	// 获取目标类型的种类
	targetKind := targetType.Kind()
	
	// 处理数字类型之间的转换
	switch targetKind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var convertedVal uint64
		switch v := value.(type) {
		case int:
			convertedVal = uint64(v)
		case int32:
			convertedVal = uint64(v)
		case int64:
			convertedVal = uint64(v)
		case uint:
			convertedVal = uint64(v)
		case uint8:
			convertedVal = uint64(v)
		case uint16:
			convertedVal = uint64(v)
		case uint32:
			convertedVal = uint64(v)
		case uint64:
			convertedVal = v
		case float32:
			convertedVal = uint64(v)
		case float64:
			convertedVal = uint64(v)
		default:
			return reflect.ValueOf(value)
		}
		
		// 根据目标类型进行正确的转换
		switch targetKind {
		case reflect.Uint:
			return reflect.ValueOf(uint(convertedVal))
		case reflect.Uint8:
			return reflect.ValueOf(uint8(convertedVal))
		case reflect.Uint16:
			return reflect.ValueOf(uint16(convertedVal))
		case reflect.Uint32:
			return reflect.ValueOf(uint32(convertedVal))
		case reflect.Uint64:
			return reflect.ValueOf(convertedVal)
		}
	}
	
	// 默认返回原值
	return reflect.ValueOf(value)
}

// getAllSettings 获取所有设置
func (receiver *DumpAdminData) getAllSettings() ([]*admin.AdminSetting, error) {
	var settings []*admin.AdminSetting
	if err := facades.Orm().Query().Find(&settings); err != nil {
		return nil, err
	}
	return settings, nil
}

// getIntFromInterface 安全地从interface{}转换为int
func getIntFromInterface(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

// getUintFromInterface 安全地从interface{}转换为uint
func getUintFromInterface(value interface{}) uint {
	switch v := value.(type) {
	case uint:
		return v
	case uint8:
		return uint(v)
	case uint16:
		return uint(v)
	case uint32:
		return uint(v)
	case uint64:
		return uint(v)
	case int:
		return uint(v)
	case int32:
		return uint(v)
	case int64:
		return uint(v)
	case float32:
		return uint(v)
	case float64:
		return uint(v)
	default:
		return 0
	}
}

// getAllRoleUsers 获取所有角色用户关联
func (receiver *DumpAdminData) getAllRoleUsers() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := "SELECT role_id, user_id, created_at, updated_at FROM admin_role_users"
	if err := facades.Orm().Query().Raw(query).Scan(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// getAllRolePermissions 获取所有角色权限关联
func (receiver *DumpAdminData) getAllRolePermissions() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := "SELECT role_id, permission_id, created_at, updated_at FROM admin_role_permissions"
	if err := facades.Orm().Query().Raw(query).Scan(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// getAllPermissionMenus 获取所有权限菜单关联
func (receiver *DumpAdminData) getAllPermissionMenus() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := "SELECT permission_id, menu_id, created_at, updated_at FROM admin_permission_menu"
	if err := facades.Orm().Query().Raw(query).Scan(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// generateSeedFile 生成种子文件
func (receiver *DumpAdminData) generateSeedFile(users []*admin.AdminUser, roles []*admin.AdminRole, permissions []*admin.AdminPermission, menus []*admin.AdminMenu, settings []*admin.AdminSetting, roleUsers []map[string]interface{}, rolePermissions []map[string]interface{}, permissionMenus []map[string]interface{}) error {
	// 生成文件内容
	code := fmt.Sprintf(`package seeders

import (
	"encoding/json"
	"fmt"
	"time"

	"goravel/app/models/admin"

	"github.com/goravel/framework/facades"
)

const %s = "%s"

// AdminSeeder 管理员数据种子
type AdminSeeder struct{}

// Run 执行种子
func (seeder *AdminSeeder) Run() error {
	return seeder.seedAdminData()
}

// seedAdminData 填充管理员数据
func (seeder *AdminSeeder) seedAdminData() error {
	// 1. 清理现有数据（如果需要）
	// 注意：这里不清理数据，只是插入新数据

	// 2. 插入用户数据
	users := %s
	for _, user := range users {
		if err := facades.Orm().Query().Create(user); err != nil {
			return fmt.Errorf("创建用户失败: %%w", err)
		}
	}

	// 3. 插入角色数据
	roles := %s
	for _, role := range roles {
		if err := facades.Orm().Query().Create(role); err != nil {
			return fmt.Errorf("创建角色失败: %%w", err)
		}
	}

	// 4. 插入权限数据
	permissions := %s
	for _, permission := range permissions {
		if err := facades.Orm().Query().Create(permission); err != nil {
			return fmt.Errorf("创建权限失败: %%w", err)
		}
	}

	// 5. 插入菜单数据
	menus := %s
	for _, menu := range menus {
		if err := facades.Orm().Query().Create(menu); err != nil {
			return fmt.Errorf("创建菜单失败: %%w", err)
		}
	}

	// 6. 插入设置数据
	settings := %s
	for _, setting := range settings {
		if err := facades.Orm().Query().Create(setting); err != nil {
			return fmt.Errorf("创建设置失败: %%w", err)
		}
	}

	// 7. 关联角色和用户
	roleUsers := %s
	for _, ru := range roleUsers {
		// 通过原生SQL执行关联插入
		query := "INSERT INTO admin_role_users (role_id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)"
		if err := facades.Orm().Query().Exec(query, ru["role_id"], ru["user_id"], ru["created_at"], ru["updated_at"]); err != nil {
			return fmt.Errorf("关联角色用户失败: %%w", err)
		}
	}

	// 8. 关联角色和权限
	rolePermissions := %s
	for _, rp := range rolePermissions {
		query := "INSERT INTO admin_role_permissions (role_id, permission_id, created_at, updated_at) VALUES (?, ?, ?, ?)"
		if err := facades.Orm().Query().Exec(query, rp["role_id"], rp["permission_id"], rp["created_at"], rp["updated_at"]); err != nil {
			return fmt.Errorf("关联角色权限失败: %%w", err)
		}
	}

	// 9. 关联权限和菜单
	permissionMenus := %s
	for _, pm := range permissionMenus {
		query := "INSERT INTO admin_permission_menu (permission_id, menu_id, created_at, updated_at) VALUES (?, ?, ?, ?)"
		if err := facades.Orm().Query().Exec(query, pm["permission_id"], pm["menu_id"], pm["created_at"], pm["updated_at"]); err != nil {
			return fmt.Errorf("关联权限菜单失败: %%w", err)
		}
	}

	return nil
}
`, 
		"SEED_COMMENT", 
		time.Now().Format("2006-01-02 15:04:05"),
		receiver.generateUsersCode(users),
		receiver.generateRolesCode(roles),
		receiver.generatePermissionsCode(permissions),
		receiver.generateMenusCode(menus),
		receiver.generateSettingsCode(settings),
		receiver.generateRoleUsersCode(roleUsers),
		receiver.generateRolePermissionsCode(rolePermissions),
		receiver.generatePermissionMenusCode(permissionMenus),
	)

	// 写入文件
	filePath := filepath.Join("database/seeders", "admin_seeder.go")
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// generateUsersCode 生成用户代码
func (receiver *DumpAdminData) generateUsersCode(users []*admin.AdminUser) string {
	if len(users) == 0 {
		return "[]*admin.AdminUser{}"
	}

	var code strings.Builder
	code.WriteString("[]*admin.AdminUser{\n")
	for _, user := range users {
		code.WriteString(fmt.Sprintf(`{
ID: %d,
Username: %s,
Password: %s,
Enabled: %d,
Name: %s,
Avatar: %s,
RememberToken: %s,
CreatedAt: %s,
UpdatedAt: %s,
},`, 
			user.ID,
			receiver.formatStringValue(user.Username),
			receiver.formatStringValue(user.Password),
			user.Enabled,
			receiver.formatStringValue(user.Name),
			receiver.formatStringValue(user.Avatar),
			receiver.formatStringValue(user.RememberToken),
			receiver.formatTime(user.CreatedAt),
			receiver.formatTime(user.UpdatedAt),
		))
		code.WriteString("\n")
	}
	code.WriteString("}")
	return code.String()
}

// generateRolesCode 生成角色代码
func (receiver *DumpAdminData) generateRolesCode(roles []*admin.AdminRole) string {
	if len(roles) == 0 {
		return "[]*admin.AdminRole{}"
	}

	var code strings.Builder
	code.WriteString("[]*admin.AdminRole{\n")
	for _, role := range roles {
		code.WriteString(fmt.Sprintf(`{
ID: %d,
Name: %s,
Slug: %s,
CreatedAt: %s,
UpdatedAt: %s,
},`, 
			role.ID,
			receiver.formatStringValue(role.Name),
			receiver.formatStringValue(role.Slug),
			receiver.formatTime(role.CreatedAt),
			receiver.formatTime(role.UpdatedAt),
		))
		code.WriteString("\n")
	}
	code.WriteString("}")
	return code.String()
}

// generatePermissionsCode 生成权限代码
func (receiver *DumpAdminData) generatePermissionsCode(permissions []*admin.AdminPermission) string {
	if len(permissions) == 0 {
		return "[]*admin.AdminPermission{}"
	}

	var code strings.Builder
	code.WriteString("[]*admin.AdminPermission{\n")
	for _, permission := range permissions {
		httpMethodJSON, _ := json.Marshal(permission.HttpMethod)
		httpPathJSON, _ := json.Marshal(permission.HttpPath)
		
		code.WriteString(fmt.Sprintf(`{
ID: %d,
Name: %s,
Slug: %s,
HttpMethod: %s,
HttpPath: %s,
CustomOrder: %d,
CreatedAt: %s,
UpdatedAt: %s,
},`, 
			permission.ID,
			receiver.formatStringValue(permission.Name),
			receiver.formatStringValue(permission.Slug),
			receiver.formatStringValue(string(httpMethodJSON)),
			receiver.formatStringValue(string(httpPathJSON)),
			permission.CustomOrder,
			receiver.formatTime(permission.CreatedAt),
			receiver.formatTime(permission.UpdatedAt),
		))
		code.WriteString("\n")
	}
	code.WriteString("}")
	return code.String()
}

// generateMenusCode 生成菜单代码
func (receiver *DumpAdminData) generateMenusCode(menus []*admin.AdminMenu) string {
	if len(menus) == 0 {
		return "[]*admin.AdminMenu{}"
	}

	var code strings.Builder
	code.WriteString("[]*admin.AdminMenu{\n")
	for _, menu := range menus {
		code.WriteString(fmt.Sprintf(`{
ID: %d,
ParentId: %d,
CustomOrder: %d,
Title: %s,
Icon: %s,
Url: %s,
UrlType: %d,
Visible: %d,
IsHome: %d,
Component: %s,
IsFull: %d,
Extension: %s,
CreatedAt: %s,
UpdatedAt: %s,
},`, 
			menu.ID,
			menu.ParentId,
			menu.CustomOrder,
			receiver.formatStringValue(menu.Title),
			receiver.formatStringValue(menu.Icon),
			receiver.formatStringValue(menu.Url),
			menu.UrlType,
			menu.Visible,
			menu.IsHome,
			receiver.formatStringValue(menu.Component),
			menu.IsFull,
			receiver.formatStringValue(menu.Extension),
			receiver.formatTime(menu.CreatedAt),
			receiver.formatTime(menu.UpdatedAt),
		))
		code.WriteString("\n")
	}
	code.WriteString("}")
	return code.String()
}

// generateSettingsCode 生成设置代码
func (receiver *DumpAdminData) generateSettingsCode(settings []*admin.AdminSetting) string {
	if len(settings) == 0 {
		return "[]*admin.AdminSetting{}"
	}

	var code strings.Builder
	code.WriteString("[]*admin.AdminSetting{\n")
	for _, setting := range settings {
		code.WriteString(fmt.Sprintf(`{
Key: %s,
Values: %s,
CreatedAt: %s,
UpdatedAt: %s,
},`, 
			receiver.formatStringValue(setting.Key),
			receiver.formatStringValue(setting.Values),
			receiver.formatTime(setting.CreatedAt),
			receiver.formatTime(setting.UpdatedAt),
		))
		code.WriteString("\n")
	}
	code.WriteString("}")
	return code.String()
}

// generateRoleUsersCode 生成角色用户关联代码
func (receiver *DumpAdminData) generateRoleUsersCode(roleUsers []map[string]interface{}) string {
	if len(roleUsers) == 0 {
		return "[]map[string]interface{}{}"
	}

	var code strings.Builder
	code.WriteString("[]map[string]interface{}{\n")
	for _, ru := range roleUsers {
		code.WriteString(fmt.Sprintf(`{
"role_id": %v,
"user_id": %v,
"created_at": %s,
"updated_at": %s,
},`, 
			ru["role_id"],
			ru["user_id"],
			receiver.formatTime(ru["created_at"]),
			receiver.formatTime(ru["updated_at"]),
		))
		code.WriteString("\n")
	}
	code.WriteString("}")
	return code.String()
}

// generateRolePermissionsCode 生成角色权限关联代码
func (receiver *DumpAdminData) generateRolePermissionsCode(rolePermissions []map[string]interface{}) string {
	if len(rolePermissions) == 0 {
		return "[]map[string]interface{}{}"
	}

	var code strings.Builder
	code.WriteString("[]map[string]interface{}{\n")
	for _, rp := range rolePermissions {
		code.WriteString(fmt.Sprintf(`{
"role_id": %v,
"permission_id": %v,
"created_at": %s,
"updated_at": %s,
},`, 
			rp["role_id"],
			rp["permission_id"],
			receiver.formatTime(rp["created_at"]),
			receiver.formatTime(rp["updated_at"]),
		))
		code.WriteString("\n")
	}
	code.WriteString("}")
	return code.String()
}

// generatePermissionMenusCode 生成权限菜单关联代码
func (receiver *DumpAdminData) generatePermissionMenusCode(permissionMenus []map[string]interface{}) string {
	if len(permissionMenus) == 0 {
		return "[]map[string]interface{}{}"
	}

	var code strings.Builder
	code.WriteString("[]map[string]interface{}{\n")
	for _, pm := range permissionMenus {
		code.WriteString(fmt.Sprintf(`{
"permission_id": %v,
"menu_id": %v,
"created_at": %s,
"updated_at": %s,
},`, 
			pm["permission_id"],
			pm["menu_id"],
			receiver.formatTime(pm["created_at"]),
			receiver.formatTime(pm["updated_at"]),
		))
		code.WriteString("\n")
	}
	code.WriteString("}")
	return code.String()
}

// formatString 格式化字符串指针
func (receiver *DumpAdminData) formatString(s *string) string {
	if s == nil {
		return "nil"
	}
	// 转义双引号
	escaped := strings.ReplaceAll(*s, `"`, `\"`)
	return fmt.Sprintf("\"%s\"", escaped)
}

// formatStringValue 格式化字符串值
func (receiver *DumpAdminData) formatStringValue(s string) string {
	// 转义双引号
	escaped := strings.ReplaceAll(s, `"`, `\"`)
	return fmt.Sprintf("\"%s\"", escaped)
}

// formatTime 格式化时间
func (receiver *DumpAdminData) formatTime(t interface{}) string {
	if t == nil {
		return "nil"
	}
	
	switch v := t.(type) {
	case time.Time:
		if v.IsZero() {
			return "nil"
		}
		return fmt.Sprintf("time.Date(%d, time.%s, %d, %d, %d, %d, 0, time.UTC)", 
			v.Year(), v.Month().String(), v.Day(), v.Hour(), v.Minute(), v.Second())
	case *time.Time:
		if v == nil || v.IsZero() {
			return "nil"
		}
		return fmt.Sprintf("time.Date(%d, time.%s, %d, %d, %d, %d, 0, time.UTC)", 
			v.Year(), v.Month().String(), v.Day(), v.Hour(), v.Minute(), v.Second())
	case string:
		if v == "" {
			return "nil"
		}
		return fmt.Sprintf("\"%s\"", v)
	default:
		return "nil"
	}
}