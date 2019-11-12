SET NAMES utf8;
-- SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;


-- vip 等级
DROP TABLE IF EXISTS `ph_level`;
CREATE TABLE `ph_level` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `create_date` datetime NOT NULL COMMENT '创建日期',
  `modify_date` datetime NOT NULL COMMENT '修改日期',
  `creator` bigint(20) NOT NULL COMMENT '创建人id',
  `modifior` bigint(20) NOT NULL COMMENT '修改人id',
  `version` bigint(20) NOT NULL COMMENT '版本号',
  `vip_level` bigint(20)  COMMENT 'vip等级',
  `total_bet` bigint(20)   COMMENT '累计投注',
  `level_gift`bigint(20)   COMMENT '晋级礼金',
  `month_bet` bigint(20 )  COMMENT '每月打码量',
  `month_gift` bigint(20 ) COMMENT '每月好运金',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='vip等级';

-- 期数
DROP TABLE IF EXISTS `ph_period`;
CREATE TABLE `ph_period` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `create_date` datetime NOT NULL COMMENT '创建日期',
  `modify_date` datetime NOT NULL COMMENT '修改日期',
  `creator` bigint(20) NOT NULL COMMENT '创建人id',
  `modifior` bigint(20) NOT NULL COMMENT '修改人id',
  `version` bigint(20) NOT NULL COMMENT '版本号',
  `period_name` VARCHAR(255) COMMENT '期数名称',
  `rank`   bigint(20)   COMMENT '排序字段',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='期数';

-- 会员单期投注
DROP TABLE IF EXISTS `ph_member_single`;
CREATE TABLE `ph_member_single` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `create_date` datetime NOT NULL COMMENT '创建日期',
  `modify_date` datetime NOT NULL COMMENT '修改日期',
  `creator` bigint(20) NOT NULL COMMENT '创建人id',
  `modifior` bigint(20) NOT NULL COMMENT '修改人id',
  `account` VARCHAR(255) BINARY  COMMENT '会员账号',
  `bet` bigint(20) COMMENT '投注金额',
  `level_gift`   bigint(20)   COMMENT '晋级彩金',
  `lucky_gift`   bigint(20)   COMMENT '当天好运金',
  `period`  VARCHAR(255)   COMMENT '期数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='会员单期投注';

-- 会员统计
DROP TABLE IF EXISTS `ph_member_total`;
CREATE TABLE `ph_member_total` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `create_date` datetime NOT NULL COMMENT '创建日期',
  `modify_date` datetime NOT NULL COMMENT '修改日期',
  `creator` bigint(20) NOT NULL COMMENT '创建人id',
  `modifior` bigint(20) NOT NULL COMMENT '修改人id',
  `account` VARCHAR(255) BINARY  COMMENT '会员账号',
  `level` bigint(20)  COMMENT 'vip等级',
  `bet`   bigint(20)   COMMENT '投注额',
  `total_level_gift`   bigint(20)   COMMENT '晋级金总额',
  `total_lucky_gift`   bigint(20)   COMMENT '好运金总额',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='会员统计';


DROP TABLE IF EXISTS `ph_quick_nav`;
CREATE TABLE `ph_quick_nav` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `create_date` datetime NOT NULL COMMENT '创建日期',
  `modify_date` datetime NOT NULL COMMENT '修改日期',
  `creator` bigint(20) NOT NULL COMMENT '创建人id',
  `modifior` bigint(20) NOT NULL COMMENT '修改人id',
  `version` bigint(20) NOT NULL COMMENT '版本号',
  `name` VARCHAR(255) NOT NULL COMMENT '网址名称',
  `web_site` VARCHAR(255) NOT NULL COMMENT '网址',
  `icon` VARCHAR(255) DEFAULT NULL COMMENT '图标',
  `seq` INT DEFAULT 0 COMMENT '排序(升序显示)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='快捷导航';

DROP TABLE IF EXISTS `ph_site_config`;
CREATE TABLE `ph_site_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `create_date` datetime NOT NULL COMMENT '创建日期',
  `modify_date` datetime NOT NULL COMMENT '修改日期',
  `creator` bigint(20) NOT NULL COMMENT '创建人id',
  `modifior` bigint(20) NOT NULL COMMENT '修改人id',
  `version` bigint(20) NOT NULL COMMENT '版本号',
  `code` VARCHAR(20) NOT NULL COMMENT '代码',
  `value` VARCHAR(255) NOT NULL COMMENT '值',
  `is_system` tinyint(1) NOT NULL COMMENT '是否内置(内置不可删除)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;

BEGIN;
INSERT INTO `ph_site_config` VALUES ('1', '2017-10-15 14:22:56', '2017-10-16 09:39:39', 1, 1, '1', 'NAME', '平台名称', 1);
COMMIT;

DROP TABLE IF EXISTS `ph_admin`;
CREATE TABLE `ph_admin` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `create_date` datetime NOT NULL COMMENT '创建日期',
  `modify_date` datetime NOT NULL COMMENT '修改日期',
  `creator` bigint(20) NOT NULL COMMENT '创建人id',
  `modifior` bigint(20) NOT NULL COMMENT '修改人id',
  `version` bigint(20) NOT NULL COMMENT '版本号',
  `ref_id` bigint(20) NOT NULL COMMENT '关联ID(公司id等)',
  `enabled` tinyint(4) NOT NULL COMMENT '是否启用',
  `locked` tinyint(4) NOT NULL COMMENT '是否锁定',
  `is_system` tinyint(4) NOT NULL COMMENT '是否内置(内置不显示)',
  `locked_date` datetime DEFAULT NULL COMMENT '锁定日期',
  `login_date` datetime DEFAULT NULL COMMENT '最后登录时间',
  `login_failure_count` int(11) NOT NULL COMMENT '连续登录失败次数',
  `login_ip` varchar(255) DEFAULT NULL COMMENT '最后登录IP',
  `salt` varchar(255) DEFAULT NULL COMMENT '加密盐',
  `name` varchar(255) DEFAULT NULL COMMENT '姓名',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `email` varchar(255) DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(255) DEFAULT '' COMMENT '手机',
  `login_verify` tinyint(4) DEFAULT 0 COMMENT '登录验证',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_admin_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Records of `ph_admin`
-- ----------------------------
BEGIN;
INSERT INTO `ph_admin` VALUES (1,'2017-10-01 20:59:59','2017-10-01 22:08:51',1,1,1,0,1,0,1,NULL,'2017-10-26 14:24:03',0,'127.0.0.1','17b007bdb8e7af362a1167bcce7277c9','超级管理员','9b16a6a8b524be91d0f440f61ed76fab','superadmin','','',0);
INSERT INTO `ph_admin` VALUES (2,'2017-10-13 22:24:06','2017-10-13 22:24:31',1,1,0,0,1,0,0,NULL,'2017-10-13 22:54:08',0,'127.0.0.1','253da8a9583fccd5645690aa25a71d20','管理员','ec7d8fc2e0093ffec5f39fede8e0bdd6','admin','','',0);
COMMIT;


-- ----------------------------
--  Table structure for ``
-- ----------------------------
DROP TABLE IF EXISTS `ph_role`;
CREATE TABLE `ph_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `create_date` datetime NOT NULL COMMENT '创建日期',
  `modify_date` datetime NOT NULL COMMENT '修改日期',
  `creator` bigint(20) NOT NULL COMMENT '创建人id',
  `modifior` bigint(20) NOT NULL COMMENT '修改人id',
  `version` bigint(20) NOT NULL COMMENT '版本号',
  `enabled` tinyint(4) NOT NULL COMMENT '是否启用',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `is_system` tinyint(4) NOT NULL COMMENT '是否内置(内置不可选择)',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
--  Records of `role`
-- ----------------------------
BEGIN;
INSERT INTO `ph_role` VALUES ('1', '2017-10-01 20:59:59', '2017-10-01 20:59:59', 1, 1, '1', 1, '拥有最高后台管理权限', 1, '超级管理员');
INSERT INTO `ph_role` VALUES ('2', '2017-10-01 20:59:59', '2017-10-01 20:59:59', 1, 1, '1', 1, '拥有后台管理权限', 0, '管理员');

COMMIT;

-- ----------------------------
--  Table structure for `ph_admin_role`
-- ----------------------------
DROP TABLE IF EXISTS `ph_admin_role`;
CREATE TABLE `ph_admin_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `admin_id` bigint(20) NOT NULL,
  `role_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_admin_role` (`admin_id`,`role_id`),
  KEY `ind_admin_role_role_id` (`role_id`),
  CONSTRAINT `fk_admin_role_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `ph_admin` (`id`),
  CONSTRAINT `fk_admin_role_role_id` FOREIGN KEY (`role_id`) REFERENCES `ph_role` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Records of `ph_admin_role`
-- ----------------------------
BEGIN;
INSERT INTO `ph_admin_role` VALUES ('1', '1', '1');
INSERT INTO `ph_admin_role` VALUES ('2', '2', '2');
COMMIT;

-- ----------------------------
--  Table structure for `ph_permission`
-- ----------------------------
SET FOREIGN_KEY_CHECKS = 0;


DROP TABLE IF EXISTS `ph_permission`;
CREATE TABLE `ph_permission` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `create_date` datetime NOT NULL COMMENT '创建日期',
  `modify_date` datetime NOT NULL COMMENT '修改日期',
  `creator` bigint(20) NOT NULL COMMENT '创建人id',
  `modifior` bigint(20) NOT NULL COMMENT '修改人id',
  `version` bigint(20) NOT NULL COMMENT '版本号',
  `pid` bigint(20) NOT NULL COMMENT '父节点id',
  `enabled` tinyint(4) NOT NULL COMMENT '是否启用',
  `display` tinyint(4) NOT NULL COMMENT '是否显示',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '链接地址',
  `name` varchar(255) NOT NULL COMMENT '名称',
  `icon` varchar(16) DEFAULT NULL COMMENT '图标',
  `sort` INT DEFAULT 100 COMMENT '排序',
  PRIMARY KEY (`id`),
  UNIQUE KEY `value` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='权限菜单,最多4级';

-- ----------------------------
--  Records of `permission`
-- ----------------------------
BEGIN;
INSERT INTO `ph_permission` VALUES ('1', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '0', '1', '0', '系统框架', 'BaseController.Index', '系统框架', '', '1');
INSERT INTO `ph_permission` VALUES ('2', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '0', '1', '0', '修改密码', 'ChangePwdController.Get', '修改密码', '', '2');
INSERT INTO `ph_permission` VALUES ('3', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '0', '1', '0', '系统信息', 'SysIndexController.Get', '系统信息', '', '3');
INSERT INTO `ph_permission` VALUES ('10', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '0', '1', '0', '系统通用-文件上传', 'SyscommonController.Upload', '系统通用-文件上传', '', '10');
INSERT INTO `ph_permission` VALUES ('20', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '0', '1', '1', '系统设置', '', '系统设置', '#xe614;', '100');
INSERT INTO `ph_permission` VALUES ('21', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '20', '1', '1', '管理员', 'AdminIndexController.Get', '管理员', '', '100');
INSERT INTO `ph_permission` VALUES ('22', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '21', '1', '0', '添加管理员', 'AdminAddController.Get', '添加管理员', '', '100');
INSERT INTO `ph_permission` VALUES ('23', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '21', '1', '0', '编辑管理员', 'AdminEditController.Get', '编辑管理员', '', '100');
INSERT INTO `ph_permission` VALUES ('24', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '21', '1', '0', '删除管理员', 'AdminIndexController.Delone', '删除管理员', '', '100');
INSERT INTO `ph_permission` VALUES ('25', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '21', '1', '0', '锁定解锁管理员', 'AdminIndexController.Locked', '锁定解锁管理员', '', '100');
INSERT INTO `ph_permission` VALUES ('26', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '21', '1', '0', '管理员登录验证', 'AdminIndexController.LoginVerify', '管理员登录验证', '', '100');
INSERT INTO `ph_permission` VALUES ('30', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '20', '1', '1', '角色管理', 'RoleIndexController.Get', '角色管理', '', '100');
INSERT INTO `ph_permission` VALUES ('31', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '30', '1', '0', '添加角色', 'RoleAddController.Get', '添加角色', '', '100');
INSERT INTO `ph_permission` VALUES ('32', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '30', '1', '0', '编辑角色', 'RoleEditController.Get', '编辑角色', '', '100');
INSERT INTO `ph_permission` VALUES ('33', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '30', '1', '0', '删除角色', 'RoleIndexController.Delone', '删除角色', '', '100');
INSERT INTO `ph_permission` VALUES ('40', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '20', '1', '1', '菜单管理', 'PermissionIndexController.Get', '菜单管理', '', '100');
INSERT INTO `ph_permission` VALUES ('41', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '40', '1', '0', '添加菜单', 'PermissionAddController.Get', '添加菜单', '', '100');
INSERT INTO `ph_permission` VALUES ('42', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '40', '1', '0', '编辑菜单', 'PermissionEditController.Get', '编辑菜单', '', '100');
INSERT INTO `ph_permission` VALUES ('43', '2017-10-01 20:59:59', '2017-10-01 20:59:59', '1', '1', '1', '40', '1', '0', '删除菜单', 'PermissionIndexController.Delone', '删除菜单', '', '100');
INSERT INTO `ph_permission` VALUES ('50', '2017-10-14 15:11:23', '2017-10-14 15:11:23', '0', '0', '0', '20', '1', '1', '站点配置', 'SiteConfigIndexController.Get', '站点配置', '', '100');
INSERT INTO `ph_permission` VALUES ('51', '2017-10-14 15:12:26', '2017-10-14 15:13:48', '0', '0', '0', '50', '1', '0', '添加站点配置', 'SiteConfigAddController.Get', '添加站点配置', '', '100');
INSERT INTO `ph_permission` VALUES ('52', '2017-10-14 15:12:54', '2017-10-14 15:13:52', '0', '0', '0', '50', '1', '0', '编辑站点配置', 'SiteConfigEditController.Get', '编辑站点配置', '', '100');
INSERT INTO `ph_permission` VALUES ('53', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '50', '1', '0', '删除站点配置', 'SiteConfigIndexController.Delone', '删除站点配置', '', '100');
INSERT INTO `ph_permission` VALUES ('60', '2017-10-14 15:11:23', '2017-10-14 15:11:23', '0', '0', '0', '20', '1', '1', '快捷导航', 'QuickNavIndexController.Get', '快捷导航', '', '100');
INSERT INTO `ph_permission` VALUES ('61', '2017-10-14 15:12:26', '2017-10-14 15:13:48', '0', '0', '0', '60', '1', '0', '添加快捷导航', 'QuickNavAddController.Get', '添加快捷导航', '', '100');
INSERT INTO `ph_permission` VALUES ('62', '2017-10-14 15:12:54', '2017-10-14 15:13:52', '0', '0', '0', '60', '1', '0', '编辑快捷导航', 'QuickNavEditController.Get', '编辑快捷导航', '', '100');
INSERT INTO `ph_permission` VALUES ('63', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '60', '1', '0', '删除快捷导航', 'QuickNavIndexController.Delone', '删除快捷导航', '', '100');
INSERT INTO `ph_permission` VALUES ('100', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '0', '1', '1', 'vip设置', '', 'vip设置', '#xe614;', '100');
INSERT INTO `ph_permission` VALUES ('101', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '100', '1', '1', 'vip等级', 'LevelController.Get', 'vip等级', '', '100');
INSERT INTO `ph_permission` VALUES ('102', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '100', '1', '0', '添加vip等级', 'LevelAddController.Get', '添加vip等级', '', '100');
INSERT INTO `ph_permission` VALUES ('103', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '100', '1', '0', '修改vip等级', 'LevelEditController.Get', '修改vip等级', '', '100');
INSERT INTO `ph_permission` VALUES ('104', '2017-10-14 15:13:22', '2017-10-14 15:13:57', '0', '0', '0', '100', '1', '0', '删除vip等级', 'LevelController.Delone', '删除vip等级', '', '100');
COMMIT;

-- ----------------------------
--  Table structure for `ph_role_permission`
-- ----------------------------
DROP TABLE IF EXISTS `ph_role_permission`;
CREATE TABLE `ph_role_permission` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色',
  `permission_id` bigint(20) NOT NULL COMMENT '权限',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_permission_role` (`role_id`, `permission_id`),
  KEY `ind_permission_role_role` (`role_id`),
  KEY `ind_permission_role_permission` (`permission_id`),
  CONSTRAINT `FK_permission_role_permission` FOREIGN KEY (`permission_id`) REFERENCES `ph_permission` (`id`),
  CONSTRAINT `FK_permission_role_role` FOREIGN KEY (`role_id`) REFERENCES `ph_role` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Records of `ph_role_permission`
-- ----------------------------
BEGIN;

INSERT INTO `ph_role_permission` VALUES ('1', '1', '1');
INSERT INTO `ph_role_permission` VALUES ('2', '1', '2');
INSERT INTO `ph_role_permission` VALUES ('3', '1', '3');
INSERT INTO `ph_role_permission` VALUES ('20', '1', '20');
INSERT INTO `ph_role_permission` VALUES ('21', '1', '21');
INSERT INTO `ph_role_permission` VALUES ('22', '1', '22');
INSERT INTO `ph_role_permission` VALUES ('23', '1', '23');
INSERT INTO `ph_role_permission` VALUES ('24', '1', '24');
INSERT INTO `ph_role_permission` VALUES ('25', '1', '25');
INSERT INTO `ph_role_permission` VALUES ('26', '1', '26');
INSERT INTO `ph_role_permission` VALUES ('30', '1', '30');
INSERT INTO `ph_role_permission` VALUES ('31', '1', '31');
INSERT INTO `ph_role_permission` VALUES ('32', '1', '32');
INSERT INTO `ph_role_permission` VALUES ('33', '1', '33');
INSERT INTO `ph_role_permission` VALUES ('40', '1', '40');
INSERT INTO `ph_role_permission` VALUES ('41', '1', '41');
INSERT INTO `ph_role_permission` VALUES ('42', '1', '42');
INSERT INTO `ph_role_permission` VALUES ('43', '1', '43');
INSERT INTO `ph_role_permission` VALUES ('50', '1', '50');
INSERT INTO `ph_role_permission` VALUES ('51', '1', '51');
INSERT INTO `ph_role_permission` VALUES ('52', '1', '52');
INSERT INTO `ph_role_permission` VALUES ('53', '1', '53');
INSERT INTO `ph_role_permission` VALUES ('60', '1', '60');
INSERT INTO `ph_role_permission` VALUES ('61', '1', '61');
INSERT INTO `ph_role_permission` VALUES ('62', '1', '62');
INSERT INTO `ph_role_permission` VALUES ('63', '1', '63');
INSERT INTO `ph_role_permission` VALUES ('101', '2', '1');
INSERT INTO `ph_role_permission` VALUES ('102', '2', '2');
INSERT INTO `ph_role_permission` VALUES ('103', '2', '3');
INSERT INTO `ph_role_permission` VALUES ('112', '2', '10');
INSERT INTO `ph_role_permission` VALUES ('120', '2', '20');
INSERT INTO `ph_role_permission` VALUES ('121', '2', '21');
INSERT INTO `ph_role_permission` VALUES ('122', '2', '22');
INSERT INTO `ph_role_permission` VALUES ('123', '2', '23');
INSERT INTO `ph_role_permission` VALUES ('124', '2', '24');
INSERT INTO `ph_role_permission` VALUES ('125', '2', '25');
INSERT INTO `ph_role_permission` VALUES ('126', '2', '26');
INSERT INTO `ph_role_permission` VALUES ('150', '2', '50');
INSERT INTO `ph_role_permission` VALUES ('151', '2', '51');
INSERT INTO `ph_role_permission` VALUES ('152', '2', '52');
INSERT INTO `ph_role_permission` VALUES ('153', '2', '53');
INSERT INTO `ph_role_permission` VALUES ('160', '2', '60');
INSERT INTO `ph_role_permission` VALUES ('161', '2', '61');
INSERT INTO `ph_role_permission` VALUES ('162', '2', '62');
INSERT INTO `ph_role_permission` VALUES ('163', '2', '63');
INSERT INTO `ph_role_permission` VALUES ('200', '2', '100');
INSERT INTO `ph_role_permission` VALUES ('201', '2', '101');
INSERT INTO `ph_role_permission` VALUES ('202', '2', '102');
INSERT INTO `ph_role_permission` VALUES ('203', '2', '103');
INSERT INTO `ph_role_permission` VALUES ('204', '2', '104');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
