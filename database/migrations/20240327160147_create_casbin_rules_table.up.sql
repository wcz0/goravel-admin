CREATE TABLE casbin_rules  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) NULL,
  `v0` varchar(100) NULL,
  `v1` varchar(100) NULL,
  `v2` varchar(100) NULL,
  `v3` varchar(100) NULL,
  `v4` varchar(100) NULL,
  `v5` varchar(100) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX(`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
);
