/*
 Navicat Premium Dump SQL

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80043 (8.0.43)
 Source Host           : localhost:3306
 Source Schema         : goravel_admin

 Target Server Type    : MySQL
 Target Server Version : 80043 (8.0.43)
 File Encoding         : 65001

 Date: 12/12/2025 14:25:43
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_apis
-- ----------------------------
DROP TABLE IF EXISTS `admin_apis`;
CREATE TABLE `admin_apis` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '接口名称',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '接口路径',
  `template` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '接口模板',
  `enabled` tinyint NOT NULL DEFAULT '1' COMMENT '是否启用',
  `args` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '接口参数',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_apis
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_code_generators
-- ----------------------------
DROP TABLE IF EXISTS `admin_code_generators`;
CREATE TABLE `admin_code_generators` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '名称',
  `table_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '表名',
  `primary_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'id' COMMENT '主键名',
  `model_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '模型名',
  `controller_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '控制器名',
  `service_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '服务名',
  `columns` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字段信息',
  `need_timestamps` tinyint NOT NULL DEFAULT '0' COMMENT '是否需要时间戳',
  `soft_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否需要软删除',
  `needs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '需要生成的代码',
  `menu_info` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '菜单信息',
  `page_info` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '页面信息',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_code_generators
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_extensions
-- ----------------------------
DROP TABLE IF EXISTS `admin_extensions`;
CREATE TABLE `admin_extensions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_enabled` tinyint NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `admin_extensions_name_unique` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_extensions
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_menus
-- ----------------------------
DROP TABLE IF EXISTS `admin_menus`;
CREATE TABLE `admin_menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` int NOT NULL DEFAULT '0',
  `custom_order` int NOT NULL DEFAULT '0',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单名称',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单图标',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单路由',
  `url_type` tinyint NOT NULL DEFAULT '1' COMMENT '路由类型(1:路由,2:外链,3:iframe)',
  `visible` tinyint NOT NULL DEFAULT '1' COMMENT '是否可见',
  `is_home` tinyint NOT NULL DEFAULT '0' COMMENT '是否为首页',
  `keep_alive` tinyint DEFAULT NULL COMMENT '页面缓存',
  `iframe_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'iframe_url',
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单组件',
  `is_full` tinyint NOT NULL DEFAULT '0' COMMENT '是否是完整页面',
  `extension` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '扩展',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_menus
-- ----------------------------
BEGIN;
INSERT INTO `admin_menus` (`id`, `parent_id`, `custom_order`, `title`, `icon`, `url`, `url_type`, `visible`, `is_home`, `keep_alive`, `iframe_url`, `component`, `is_full`, `extension`, `created_at`, `updated_at`) VALUES (1, 0, 0, 'dashboard', 'mdi:chart-line', '/dashboard', 1, 1, 1, NULL, NULL, NULL, 0, NULL, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_menus` (`id`, `parent_id`, `custom_order`, `title`, `icon`, `url`, `url_type`, `visible`, `is_home`, `keep_alive`, `iframe_url`, `component`, `is_full`, `extension`, `created_at`, `updated_at`) VALUES (2, 0, 0, 'admin_system', 'material-symbols:settings-outline', '/system', 1, 1, 0, NULL, NULL, NULL, 0, NULL, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_menus` (`id`, `parent_id`, `custom_order`, `title`, `icon`, `url`, `url_type`, `visible`, `is_home`, `keep_alive`, `iframe_url`, `component`, `is_full`, `extension`, `created_at`, `updated_at`) VALUES (3, 2, 0, 'admin_users', 'ph:user-gear', '/system/admin_users', 1, 1, 0, NULL, NULL, NULL, 0, NULL, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_menus` (`id`, `parent_id`, `custom_order`, `title`, `icon`, `url`, `url_type`, `visible`, `is_home`, `keep_alive`, `iframe_url`, `component`, `is_full`, `extension`, `created_at`, `updated_at`) VALUES (4, 2, 0, 'admin_roles', 'carbon:user-role', '/system/admin_roles', 1, 1, 0, NULL, NULL, NULL, 0, NULL, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_menus` (`id`, `parent_id`, `custom_order`, `title`, `icon`, `url`, `url_type`, `visible`, `is_home`, `keep_alive`, `iframe_url`, `component`, `is_full`, `extension`, `created_at`, `updated_at`) VALUES (5, 2, 0, 'admin_permission', 'fluent-mdl2:permissions', '/system/admin_permissions', 1, 1, 0, NULL, NULL, NULL, 0, NULL, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_menus` (`id`, `parent_id`, `custom_order`, `title`, `icon`, `url`, `url_type`, `visible`, `is_home`, `keep_alive`, `iframe_url`, `component`, `is_full`, `extension`, `created_at`, `updated_at`) VALUES (6, 2, 0, 'admin_menu', 'ant-design:menu-unfold-outlined', '/system/admin_menus', 1, 1, 0, NULL, NULL, NULL, 0, NULL, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_menus` (`id`, `parent_id`, `custom_order`, `title`, `icon`, `url`, `url_type`, `visible`, `is_home`, `keep_alive`, `iframe_url`, `component`, `is_full`, `extension`, `created_at`, `updated_at`) VALUES (7, 2, 0, 'admin_setting', 'akar-icons:settings-horizontal', '/system/settings', 1, 1, 0, NULL, NULL, NULL, 0, NULL, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
COMMIT;

-- ----------------------------
-- Table structure for admin_pages
-- ----------------------------
DROP TABLE IF EXISTS `admin_pages`;
CREATE TABLE `admin_pages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '页面名称',
  `sign` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '页面标识',
  `schema` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '页面结构',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_pages
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_permission_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_permission_menu`;
CREATE TABLE `admin_permission_menu` (
  `permission_id` int NOT NULL,
  `menu_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  KEY `admin_permission_menu_permission_id_menu_id_index` (`permission_id`,`menu_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_permission_menu
-- ----------------------------
BEGIN;
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (2, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (3, 3, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (2, 3, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (4, 4, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (2, 4, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (5, 5, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (2, 5, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (6, 6, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (2, 6, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (7, 7, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`) VALUES (2, 7, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
COMMIT;

-- ----------------------------
-- Table structure for admin_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_permissions`;
CREATE TABLE `admin_permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `http_method` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `http_path` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `custom_order` int NOT NULL DEFAULT '0',
  `parent_id` int NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `admin_permissions_name_unique` (`name`) USING BTREE,
  UNIQUE KEY `admin_permissions_slug_unique` (`slug`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_permissions
-- ----------------------------
BEGIN;
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (1, '首页', 'home', NULL, '[\'/home*\']', 0, 0, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (2, '系统', 'system', NULL, '', 0, 0, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (3, '管理员', 'admin_users', NULL, '[\'/admin_users*\']', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (4, '角色', 'roles', NULL, '[\'/roles*\']', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (5, '权限', 'permissions', NULL, '[\'/permissions*\']', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (6, '菜单', 'menus', NULL, '[\'/menus*\']', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (7, '设置', 'settings', NULL, '[\'/settings*\']', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
COMMIT;

-- ----------------------------
-- Table structure for admin_relationships
-- ----------------------------
DROP TABLE IF EXISTS `admin_relationships`;
CREATE TABLE `admin_relationships` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `model` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模型',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联名称',
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联名称',
  `args` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '关联参数',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '额外参数',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_relationships
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_permissions`;
CREATE TABLE `admin_role_permissions` (
  `role_id` int NOT NULL,
  `permission_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  KEY `admin_role_permissions_role_id_permission_id_index` (`role_id`,`permission_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_role_permissions
-- ----------------------------
BEGIN;
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`, `created_at`, `updated_at`) VALUES (1, 1, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`, `created_at`, `updated_at`) VALUES (1, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`, `created_at`, `updated_at`) VALUES (1, 3, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`, `created_at`, `updated_at`) VALUES (1, 4, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`, `created_at`, `updated_at`) VALUES (1, 5, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`, `created_at`, `updated_at`) VALUES (1, 6, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`, `created_at`, `updated_at`) VALUES (1, 7, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
COMMIT;

-- ----------------------------
-- Table structure for admin_role_users
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_users`;
CREATE TABLE `admin_role_users` (
  `role_id` int NOT NULL,
  `user_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  KEY `admin_role_users_role_id_user_id_index` (`role_id`,`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_role_users
-- ----------------------------
BEGIN;
INSERT INTO `admin_role_users` (`role_id`, `user_id`, `created_at`, `updated_at`) VALUES (1, 1, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
COMMIT;

-- ----------------------------
-- Table structure for admin_roles
-- ----------------------------
DROP TABLE IF EXISTS `admin_roles`;
CREATE TABLE `admin_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `admin_roles_name_unique` (`name`) USING BTREE,
  UNIQUE KEY `admin_roles_slug_unique` (`slug`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_roles
-- ----------------------------
BEGIN;
INSERT INTO `admin_roles` (`id`, `name`, `slug`, `created_at`, `updated_at`) VALUES (1, 'Administrator', 'administrator', '2024-08-10 07:04:53', '2024-08-10 07:04:53');
COMMIT;

-- ----------------------------
-- Table structure for admin_settings
-- ----------------------------
DROP TABLE IF EXISTS `admin_settings`;
CREATE TABLE `admin_settings` (
  `key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `values` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_settings
-- ----------------------------
BEGIN;
INSERT INTO `admin_settings` (`key`, `values`, `created_at`, `updated_at`) VALUES ('admin_locale', '\"zh_CN\"', '2024-08-10 07:04:53', '2025-11-12 00:59:19');
INSERT INTO `admin_settings` (`key`, `values`, `created_at`, `updated_at`) VALUES ('system_theme_setting', '{\"darkTheme\":false,\"footer\":false,\"breadcrumb\":true,\"themeColor\":\"#1677ff\",\"layoutMode\":\"default\",\"siderTheme\":\"light\",\"topTheme\":\"light\",\"animateInType\":\"alpha\",\"animateInDuration\":600,\"animateOutType\":\"alpha\",\"animateOutDuration\":600,\"loginTemplate\":\"default\",\"keepAlive\":false,\"enableTab\":false,\"tabIcon\":true,\"accordionMenu\":false}', '2024-09-27 10:07:04', '2024-09-27 10:07:04');
COMMIT;

-- ----------------------------
-- Table structure for admin_users
-- ----------------------------
DROP TABLE IF EXISTS `admin_users`;
CREATE TABLE `admin_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `enabled` tinyint NOT NULL DEFAULT '1',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `remember_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `admin_users_username_unique` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_users
-- ----------------------------
BEGIN;
INSERT INTO `admin_users` (`id`, `username`, `password`, `enabled`, `name`, `avatar`, `remember_token`, `created_at`, `updated_at`) VALUES (1, 'admin', '$2a$10$xHRH5sa52t1MC79h8Ml4p.jyCObM2SNOfXgrHtd1YbGxk1e05HIk.', 1, 'Administrator', NULL, NULL, '2024-08-10 07:04:53', '2025-11-13 00:45:39');
COMMIT;

-- ----------------------------
-- Table structure for casbin_rules
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rules`;
CREATE TABLE `casbin_rules` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v0` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v1` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v2` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v3` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v4` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v5` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `casbin_rules_ptype_v0_v1_v2_v3_v4_v5_unique` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of casbin_rules
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for casbin_rules_second
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rules_second`;
CREATE TABLE `casbin_rules_second` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v0` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v1` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v2` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v3` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v4` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `v5` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `casbin_rules_second_ptype_v0_v1_v2_v3_v4_v5_unique` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of casbin_rules_second
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for failed_jobs
-- ----------------------------
DROP TABLE IF EXISTS `failed_jobs`;
CREATE TABLE `failed_jobs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `connection` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `queue` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `payload` longtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `exception` longtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `failed_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `failed_jobs_uuid_unique` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of failed_jobs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for jobs
-- ----------------------------
DROP TABLE IF EXISTS `jobs`;
CREATE TABLE `jobs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `queue` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `payload` longtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `attempts` tinyint unsigned NOT NULL DEFAULT '0',
  `reserved_at` datetime DEFAULT NULL,
  `available_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `jobs_queue_index` (`queue`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of jobs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for migrations
-- ----------------------------
DROP TABLE IF EXISTS `migrations`;
CREATE TABLE `migrations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `migration` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `batch` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of migrations
-- ----------------------------
BEGIN;
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (1, '20210101000002_create_casbin_rules_table', 1);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (2, '20210101000002_create_jobs_table', 1);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (3, '20240417000001_create_admin_users_table', 1);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (4, '20240417000002_create_users_table', 1);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (5, '20240418000001_create_admin_settings_table', 1);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (6, '20240418000002_create_admin_menus_table', 1);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (7, '20240418000003_create_admin_permission_menu_table', 2);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (8, '20240418000004_create_admin_permissions_table', 3);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (9, '20240418000005_create_admin_extensions_table', 3);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (10, '20240419000001_create_admin_roles_table', 3);
INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES (11, '20240423000001_create_casbin_rules_second_table', 4);
COMMIT;

-- ----------------------------
-- Table structure for service_agreements
-- ----------------------------
DROP TABLE IF EXISTS `service_agreements`;
CREATE TABLE `service_agreements` (
  `id` int NOT NULL AUTO_INCREMENT,
  `version` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `locale` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT 'zh-CN',
  `type` tinyint NOT NULL COMMENT '协议类型：1=用户协议, 2=隐私政策, 3=未成年人保护及家长监护指南, 4=广告与营销内容披露政策, 5=直播与音频内容管理规则, 6=创作者合作协议',
  `is_published` tinyint(1) DEFAULT '0' COMMENT '是否已发布',
  `published_at` timestamp NULL DEFAULT NULL COMMENT '发布时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_version_type` (`version`,`type`),
  KEY `idx_type` (`type`),
  KEY `idx_published` (`is_published`),
  KEY `idx_version_type` (`version`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of service_agreements
-- ----------------------------
BEGIN;
INSERT INTO `service_agreements` (`id`, `version`, `title`, `content`, `locale`, `type`, `is_published`, `published_at`, `created_at`, `updated_at`) VALUES (1, '1.0.0', '用户协议', '<h1>用户协议</h1><p>欢迎使用Enchat平台服务。本协议是您与Enchat平台之间关于使用平台服务的法律协议。</p><h2>一、服务内容</h2><p>Enchat平台为用户提供即时通讯、社交互动等功能。</p><h2>二、用户义务</h2><p>用户应遵守相关法律法规，不得发布违法违规内容。</p><h2>三、免责声明</h2><p>平台不对用户发布的内容承担责任。</p>', 'zh-CN', 1, 1, '2025-11-12 00:57:40', '2025-11-12 00:57:40', '2025-11-12 00:57:40');
INSERT INTO `service_agreements` (`id`, `version`, `title`, `content`, `locale`, `type`, `is_published`, `published_at`, `created_at`, `updated_at`) VALUES (2, '1.0.0', '隐私政策', '<h1>隐私政策</h1><p>本隐私政策说明了我们如何收集、使用、保护您的个人信息。</p><h2>一、信息收集</h2><p>我们会收集您提供的个人信息，包括但不限于：昵称、头像、联系方式等。</p><h2>二、信息使用</h2><p>您的个人信息将用于提供和改进我们的服务。</p><h2>三、信息保护</h2><p>我们采取合理的安全措施保护您的个人信息。</p>', 'zh-CN', 2, 1, '2025-11-12 00:57:40', '2025-11-12 00:57:40', '2025-11-12 00:57:40');
INSERT INTO `service_agreements` (`id`, `version`, `title`, `content`, `locale`, `type`, `is_published`, `published_at`, `created_at`, `updated_at`) VALUES (3, '1.0.0', '未成年人保护及家长监护指南', '<h1>未成年人保护及家长监护指南</h1><p>为了保护未成年人使用平台的安全和身心健康，我们制定了以下保护措施。</p><h2>一、年龄限制</h2><p>未满14周岁的用户需要在监护人同意下使用本服务。</p><h2>二、监护人责任</h2><p>监护人应当教育和引导未成年人正确使用网络服务。</p><h2>三、内容过滤</h2><p>平台将对未成年人可能接触的内容进行适当过滤。</p>', 'zh-CN', 3, 1, '2025-11-12 00:57:40', '2025-11-12 00:57:40', '2025-11-12 00:57:40');
INSERT INTO `service_agreements` (`id`, `version`, `title`, `content`, `locale`, `type`, `is_published`, `published_at`, `created_at`, `updated_at`) VALUES (4, '1.0.0', '广告与营销内容披露政策', '<h1>广告与营销内容披露政策</h1><p>为了保护用户权益，规范平台广告和营销内容的管理。</p><h2>一、广告标识</h2><p>所有广告内容必须明确标识为广告。</p><h2>二、内容标准</h2><p>广告内容应当真实、合法，不得误导用户。</p><h2>三、投诉处理</h2><p>用户可对不当广告内容进行投诉举报。</p>', 'zh-CN', 4, 1, '2025-11-12 00:57:40', '2025-11-12 00:57:40', '2025-11-12 00:57:40');
INSERT INTO `service_agreements` (`id`, `version`, `title`, `content`, `locale`, `type`, `is_published`, `published_at`, `created_at`, `updated_at`) VALUES (5, '1.0.0', '直播与音频内容管理规则', '<h1>直播与音频内容管理规则</h1><p>为了营造健康向上的网络环境，规范直播和音频内容的管理。</p><h2>一、内容规范</h2><p>直播和音频内容应积极健康，符合社会道德规范。</p><h2>二、禁止内容</h2><p>严禁传播暴力、色情、欺诈等违法违规内容。</p><h2>三、审核机制</h2><p>平台对直播和音频内容进行实时监控和事后审核。</p>', 'zh-CN', 5, 1, '2025-11-12 00:57:40', '2025-11-12 00:57:40', '2025-11-12 00:57:40');
INSERT INTO `service_agreements` (`id`, `version`, `title`, `content`, `locale`, `type`, `is_published`, `published_at`, `created_at`, `updated_at`) VALUES (6, '1.0.0', '创作者合作协议', '<h1>创作者合作协议</h1><p>本协议规范了内容创作者与平台之间的权利义务关系。</p><h2>一、合作内容</h2><p>创作者在平台发布原创内容，平台提供展示和推广服务。</p><h2>二、权利归属</h2><p>创作者保留内容著作权，平台获得内容使用许可。</p><h2>三、收益分配</h2><p>平台将按照约定比例向创作者分配收益。</p>', 'zh-CN', 6, 1, '2025-11-12 00:57:40', '2025-11-12 00:57:40', '2025-11-12 00:57:40');
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `password` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `real_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `card_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `mark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `parther_id` int unsigned DEFAULT '0',
  `group_id` int unsigned DEFAULT '0',
  `nickname` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `phone` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `add_ip` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `last_time` datetime DEFAULT NULL,
  `last_ip` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `money` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '0',
  `brokerage_price` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '0',
  `integral` int unsigned DEFAULT '0',
  `exp` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '0',
  `sign_num` int unsigned DEFAULT '0',
  `state` tinyint unsigned DEFAULT '1',
  `level` tinyint unsigned DEFAULT '0',
  `agent_level` int unsigned DEFAULT '0',
  `is_spread` tinyint unsigned DEFAULT '0',
  `spread_uid` bigint unsigned DEFAULT NULL,
  `spread_time` datetime DEFAULT NULL,
  `user_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `is_promoter` tinyint unsigned DEFAULT '0',
  `pay_count` int unsigned DEFAULT '0',
  `spread_count` int unsigned DEFAULT '0',
  `clean_time` datetime DEFAULT NULL,
  `address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `admin_id` int unsigned DEFAULT '0',
  `login_type` tinyint unsigned DEFAULT '0',
  `record_phone` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `member_level` tinyint unsigned DEFAULT '0',
  `member_ever` tinyint unsigned DEFAULT '0',
  `overdue_time` datetime DEFAULT NULL,
  `division_type` tinyint unsigned DEFAULT '0',
  `division_status` tinyint unsigned DEFAULT '0',
  `is_division` tinyint unsigned DEFAULT '0',
  `is_agent` tinyint unsigned DEFAULT '0',
  `is_staff` tinyint unsigned DEFAULT '0',
  `division_id` int unsigned DEFAULT '0',
  `agent_id` int DEFAULT '0',
  `staff_id` int DEFAULT '0',
  `division_percent` tinyint unsigned DEFAULT '0',
  `division_end_time` datetime DEFAULT NULL,
  `division_last_time` datetime DEFAULT NULL,
  `division_invite` int unsigned DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_username_unique` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
