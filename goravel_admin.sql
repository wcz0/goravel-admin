/*
 Navicat Premium Dump SQL

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80040 (8.0.40)
 Source Host           : localhost:3306
 Source Schema         : goravel_admin

 Target Server Type    : MySQL
 Target Server Version : 80040 (8.0.40)
 File Encoding         : 65001

 Date: 03/12/2024 00:11:46
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
  `http_method_bak` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `http_path_bak` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `http_method` json DEFAULT NULL,
  `http_path` json DEFAULT NULL,
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
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method_bak`, `http_path_bak`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (1, '首页', 'home', NULL, '[\'/home*\']', NULL, '[\"/home*\"]', 0, 0, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method_bak`, `http_path_bak`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (2, '系统', 'system', NULL, '', NULL, NULL, 0, 0, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method_bak`, `http_path_bak`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (3, '管理员', 'admin_users', NULL, '[\'/admin_users*\']', NULL, '[\"/admin_users*\"]', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method_bak`, `http_path_bak`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (4, '角色', 'roles', NULL, '[\'/roles*\']', NULL, '[\"/roles*\"]', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method_bak`, `http_path_bak`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (5, '权限', 'permissions', NULL, '[\'/permissions*\']', NULL, '[\"/permissions*\"]', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method_bak`, `http_path_bak`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (6, '菜单', 'menus', NULL, '[\'/menus*\']', NULL, '[\"/menus*\"]', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
INSERT INTO `admin_permissions` (`id`, `name`, `slug`, `http_method_bak`, `http_path_bak`, `http_method`, `http_path`, `custom_order`, `parent_id`, `created_at`, `updated_at`) VALUES (7, '设置', 'settings', NULL, '[\'/settings*\']', NULL, '[\"/settings*\"]', 0, 2, '2024-08-10 07:04:53', '2024-08-10 07:04:53');
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
INSERT INTO `admin_settings` (`key`, `values`, `created_at`, `updated_at`) VALUES ('admin_locale', '\"zh_CN\"', NULL, '2024-12-03 00:00:22');
INSERT INTO `admin_settings` (`key`, `values`, `created_at`, `updated_at`) VALUES ('system_theme_setting', '{\"accordionMenu\":false,\"animateInDuration\":600,\"animateInType\":\"alpha\",\"animateOutDuration\":600,\"animateOutType\":\"alpha\",\"breadcrumb\":true,\"darkTheme\":false,\"enableTab\":true,\"footer\":false,\"keepAlive\":false,\"layoutMode\":\"default\",\"loginTemplate\":\"default\",\"siderTheme\":\"light\",\"tabIcon\":true,\"themeColor\":\"#1677ff\",\"topTheme\":\"light\"}', '2024-12-02 23:54:15', '2024-12-02 23:58:50');
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
INSERT INTO `admin_users` (`id`, `username`, `password`, `enabled`, `name`, `avatar`, `remember_token`, `created_at`, `updated_at`) VALUES (1, 'admin', '$2a$12$qbXT0QDJh5PYU1WrpeT3.ufPNkSc.YaoIAuTrSPUee/mvSvgquIZm', 1, 'Administrator', NULL, NULL, '2024-07-24 16:51:57', '2024-07-24 16:51:57');
COMMIT;

-- ----------------------------
-- Table structure for casbin_rules
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rules`;
CREATE TABLE `casbin_rules` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `ptype` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

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
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v0` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v1` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v2` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v3` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v4` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `v5` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `ptype` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of casbin_rules_second
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for migrations
-- ----------------------------
DROP TABLE IF EXISTS `migrations`;
CREATE TABLE `migrations` (
  `version` bigint NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of migrations
-- ----------------------------
BEGIN;
INSERT INTO `migrations` (`version`, `dirty`) VALUES (20240423045835, 0);
COMMIT;

-- ----------------------------
-- Table structure for password_reset_tokens
-- ----------------------------
DROP TABLE IF EXISTS `password_reset_tokens`;
CREATE TABLE `password_reset_tokens` (
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`email`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of password_reset_tokens
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for personal_access_tokens
-- ----------------------------
DROP TABLE IF EXISTS `personal_access_tokens`;
CREATE TABLE `personal_access_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `tokenable_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `tokenable_id` bigint unsigned NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `abilities` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `last_used_at` timestamp NULL DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `personal_access_tokens_token_unique` (`token`) USING BTREE,
  KEY `personal_access_tokens_tokenable_type_tokenable_id_index` (`tokenable_type`,`tokenable_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of personal_access_tokens
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email_verified_at` timestamp NULL DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `remember_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `users_email_unique` (`email`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
