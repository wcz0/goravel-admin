CREATE TABLE admin_permission_menu (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `permission` int NOT NULL,
  `menu` int NOT NULL,
  `created_at` datetime NULL DEFAULT NULL,
  `updated_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
