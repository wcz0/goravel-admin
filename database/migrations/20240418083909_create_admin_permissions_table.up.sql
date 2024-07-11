CREATE TABLE admin_permissions (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级权限ID',
  `name` varchar(255) NOT NULL COMMENT '权限名称',
  `value` varchar(255) NOT NULL COMMENT '权限值',
  `method` varchar(255) NOT NULL COMMENT '请求方法',
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
