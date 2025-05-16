/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80033 (8.0.33)
 Source Host           : localhost:3306
 Source Schema         : balala_v4

 Target Server Type    : MySQL
 Target Server Version : 80033 (8.0.33)
 File Encoding         : 65001

 Date: 23/08/2024 17:17:08
*/


SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

USE `balala`;

-- ----------------------------
-- Table structure for dim_date
-- ----------------------------
DROP TABLE IF EXISTS `dim_date`;
CREATE TABLE `dim_date` (
  `date_id` int NOT NULL AUTO_INCREMENT,
  `date` date DEFAULT NULL,
  `year` int DEFAULT NULL,
  `quarter` int DEFAULT NULL,
  `month` int DEFAULT NULL,
  `day` int DEFAULT NULL,
  `day_of_week` int DEFAULT NULL,
  `week_of_year` int DEFAULT NULL,
  `month_name` varchar(20) DEFAULT NULL,
  `weekday_name` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`date_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for dim_object
-- ----------------------------
DROP TABLE IF EXISTS `dim_object`;
CREATE TABLE `dim_object` (
  `space_id` bigint DEFAULT NULL COMMENT '空间id',
  `object_id` bigint DEFAULT NULL COMMENT '空间模块id',
  `object_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '空间模块名称',
  `gmt_create` datetime DEFAULT NULL COMMENT '空间模块-创建时间',
  `gmt_modified` datetime DEFAULT NULL COMMENT '空间模块-更新时间',
  `start_date` datetime DEFAULT NULL COMMENT '纬度-生效日期',
  `end_date` datetime DEFAULT NULL COMMENT '纬度-失效日期',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `object_id` (`object_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='模版维度表';

-- ----------------------------
-- Table structure for dim_space
-- ----------------------------
DROP TABLE IF EXISTS `dim_space`;
CREATE TABLE `dim_space` (
  `space_id` bigint DEFAULT NULL COMMENT '空间id',
  `space_name` varchar(255) DEFAULT NULL COMMENT '空间名称',
  `gmt_create` datetime DEFAULT NULL COMMENT '空间-创建时间',
  `gmt_modified` datetime DEFAULT NULL COMMENT '空间-更新时间',
  `start_date` datetime DEFAULT NULL COMMENT '纬度-生效日期',
  `end_date` datetime DEFAULT NULL COMMENT '纬度-失效日期',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `space_id` (`space_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='工作空间纬度表(全量)';

-- ----------------------------
-- Table structure for dim_user
-- ----------------------------
DROP TABLE IF EXISTS `dim_user`;
CREATE TABLE `dim_user` (
  `user_id` int DEFAULT NULL,
  `user_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户名',
  `user_nick_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户昵称',
  `user_pinyin` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户昵称拼音',
  `gmt_create` datetime DEFAULT NULL COMMENT '用户-创建时间',
  `gmt_modified` datetime DEFAULT NULL COMMENT '用户-更新时间',
  `start_date` datetime DEFAULT NULL COMMENT '纬度-生效日期',
  `end_date` datetime DEFAULT NULL COMMENT '纬度-失效日期',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for dim_version
-- ----------------------------
DROP TABLE IF EXISTS `dim_version`;
CREATE TABLE `dim_version` (
  `space_id` bigint DEFAULT NULL COMMENT '空间id',
  `version_id` bigint DEFAULT NULL COMMENT '版本id',
  `version_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '版本名称',
  `gmt_create` datetime DEFAULT NULL COMMENT '版本-创建时间',
  `gmt_modified` datetime DEFAULT NULL COMMENT '版本-更新时间',
  `start_date` datetime DEFAULT NULL COMMENT '纬度-生效日期',
  `end_date` datetime DEFAULT NULL COMMENT '纬度-失效日期',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `version_id` (`version_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='版本维度表';

-- ----------------------------
-- Table structure for dim_witem_status
-- ----------------------------
DROP TABLE IF EXISTS `dim_witem_status`;
    CREATE TABLE `dim_witem_status` (
    `space_id` bigint DEFAULT NULL COMMENT '空间id',
    `status_id` bigint DEFAULT NULL COMMENT '工作项状态id',
    `status_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '工作项状态名称',
    `status_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '工作项状态Key',
    `status_val` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '工作项状态Val',
    `status_type` int DEFAULT NULL COMMENT '工作项状态类型',
    `flow_scope` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '工作项状态Val',
    `gmt_create` datetime DEFAULT NULL COMMENT '工作项状态-创建时间',
    `gmt_modified` datetime DEFAULT NULL COMMENT '工作项状态-更新时间',
    `start_date` datetime DEFAULT NULL COMMENT '纬度-生效日期',
    `end_date` datetime DEFAULT NULL COMMENT '纬度-失效日期',
    `_id` bigint NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `status_id` (`status_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='任务状态纬度表（全量）';

-- ----------------------------
-- Table structure for dwd_member
-- ----------------------------
DROP TABLE IF EXISTS `dwd_member`;
CREATE TABLE `dwd_member` (
  `member_id` bigint unsigned NOT NULL COMMENT '成员id',
  `space_id` bigint unsigned DEFAULT NULL COMMENT '项目空间id',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '创建用户id',
  `role_id` int unsigned DEFAULT NULL COMMENT '关联角色id 1: 项目空间管理员 2: 项目空间成员',
  `gmt_create` datetime DEFAULT NULL COMMENT '成员-创建时间',
  `gmt_modified` datetime DEFAULT NULL COMMENT '成员-更新时间',
  `start_date` datetime DEFAULT NULL COMMENT '生效日期，格式为yyyy-mm-dd hh:MM:ss',
  `end_date` datetime DEFAULT NULL COMMENT '失效日期，格式为yyyy-mm-dd  hh:MM:ss，9999-12-31 00:00:00 表示有效',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `member_id` (`member_id`,`end_date`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='成员事实明细表（拉链存储，增量）';

-- ----------------------------
-- Table structure for dwd_witem
-- ----------------------------
DROP TABLE IF EXISTS `dwd_witem`;
CREATE TABLE `dwd_witem` (
     `work_item_id` bigint DEFAULT NULL COMMENT '工作项id',
     `space_id` bigint DEFAULT NULL COMMENT '空间id',
     `user_id` bigint DEFAULT NULL COMMENT '创建人id',
     `status_id` bigint DEFAULT NULL COMMENT '工作项状态id',
     `object_id` bigint DEFAULT NULL COMMENT '模块id',
     `version_id` bigint DEFAULT NULL COMMENT '版本id',
     `work_item_type_key` varchar(255) DEFAULT NULL COMMENT '任务类型key',
     `last_status_at` bigint DEFAULT NULL COMMENT '状态变更时间',
     `priority` varchar(255) DEFAULT NULL COMMENT '优先级',
     `plan_start_at` bigint DEFAULT NULL COMMENT '排期-开始时间',
     `plan_complete_at` bigint DEFAULT NULL COMMENT '排期-结束时间',
     `directors` json DEFAULT NULL COMMENT '负责人',
     `node_directors` json DEFAULT NULL COMMENT '节点负责人',
     `participators` json DEFAULT NULL COMMENT '参与人',
     `gmt_create` datetime DEFAULT NULL COMMENT '工作项创建时间',
     `gmt_modified` datetime DEFAULT NULL COMMENT '工作项更新时间',
     `start_date` datetime DEFAULT NULL COMMENT '生效日期，格式为yyyy-mm-dd hh:MM:ss',
     `end_date` datetime DEFAULT NULL COMMENT '失效日期，格式为yyyy-mm-dd  hh:MM:ss，9999-12-31 00:00:00 表示有效',
     `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `idx_witem_id` (`work_item_id`,`end_date`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='任务事实明细表 （全量）拉链存储';

-- ----------------------------
-- Table structure for dwd_witem_flow_node
-- ----------------------------
DROP TABLE IF EXISTS `dwd_witem_flow_node`;
CREATE TABLE `dwd_witem_flow_node` (
  `space_id` bigint DEFAULT NULL COMMENT '空间id',
  `work_item_id` bigint DEFAULT NULL COMMENT '工作项id',
  `node_id` bigint DEFAULT NULL COMMENT '工作流程节点id',
  `node_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '工作流程节点-编码',
  `node_status` int DEFAULT NULL COMMENT '工作流程节点-状态 1: 未开启 2: 进行中 3:已完成',
  `plan_start_at` bigint DEFAULT NULL COMMENT '节点排期-开始时间',
  `plan_complete_at` bigint DEFAULT NULL COMMENT '节点排期-结束时间',
  `directors` json DEFAULT NULL COMMENT '负责人',
  `gmt_create` datetime DEFAULT NULL COMMENT '工作项创建时间',
  `gmt_modified` datetime DEFAULT NULL COMMENT '工作项更新时间',
  `start_date` datetime DEFAULT NULL COMMENT '生效日期，格式为yyyy-mm-dd hh:MM:ss',
  `end_date` datetime DEFAULT NULL COMMENT '失效日期，格式为yyyy-mm-dd  hh:MM:ss，9999-12-31 00:00:00 表示有效',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `idx_witem_id` (`node_id`,`end_date`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='任务事实明细表 （全量）拉链存储';

-- ----------------------------
-- Table structure for dws_mbr_witem_1h
-- ----------------------------
DROP TABLE IF EXISTS `dws_mbr_witem_1h`;
CREATE TABLE `dws_mbr_witem_1h` (
  `space_id` bigint DEFAULT NULL COMMENT '纬度-空间',
  `user_id` bigint DEFAULT NULL COMMENT '纬度-用户id',
  `num` int DEFAULT NULL COMMENT '任务数量',
  `expire_num` int DEFAULT NULL COMMENT '过期数量',
  `todo_num` int DEFAULT NULL COMMENT '待办的任务数量',
  `complete_num` int DEFAULT NULL COMMENT '完成的任务数量',
  `close_num` int DEFAULT NULL COMMENT '关闭的任务数量',
  `abort_num` int DEFAULT NULL COMMENT '终止的任务数量',
  `start_date` datetime DEFAULT NULL COMMENT '汇总的开始时间',
  `end_date` datetime DEFAULT NULL COMMENT '汇总的结束时间',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `space_id` (`space_id`,`user_id`,`start_date`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='版本任务统计表/小时';

-- ----------------------------
-- Table structure for dws_vers_witem_1h
-- ----------------------------
DROP TABLE IF EXISTS `dws_vers_witem_1h`;
CREATE TABLE `dws_vers_witem_1h` (
  `space_id` bigint DEFAULT NULL COMMENT '纬度-空间',
  `version_id` bigint DEFAULT NULL COMMENT '纬度-版本',
  `expire_num` int DEFAULT NULL COMMENT '过期数量',
  `num` int DEFAULT NULL COMMENT '任务数量',
  `todo_num` int DEFAULT NULL COMMENT '待办的任务数量',
  `complete_num` int DEFAULT NULL COMMENT '完成的任务数量',
  `close_num` int DEFAULT NULL COMMENT '关闭的任务数量',
  `abort_num` int DEFAULT NULL COMMENT '终止的任务数量',
  `start_date` datetime DEFAULT NULL COMMENT '汇总的开始时间',
  `end_date` datetime DEFAULT NULL COMMENT '汇总的结束时间',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `space_id` (`space_id`,`version_id`,`start_date`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='版本任务统计表/小时';

-- ----------------------------
-- Table structure for job_variables
-- ----------------------------
DROP TABLE IF EXISTS `job_variables`;
CREATE TABLE `job_variables` (
  `job_name` varchar(225) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `variable_name` varchar(225) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `variable_value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  UNIQUE KEY `job_name` (`job_name`,`variable_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for ods_member_d
-- ----------------------------
DROP TABLE IF EXISTS `ods_member_d`;
CREATE TABLE `ods_member_d` (
  `id` bigint unsigned NOT NULL,
  `space_id` bigint unsigned NOT NULL COMMENT '项目空间id',
  `user_id` bigint unsigned NOT NULL COMMENT '创建用户id',
  `role_id` int unsigned NOT NULL COMMENT '关联角色id 1: 项目空间管理员 2: 项目空间成员',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '备注',
  `ranking` bigint NOT NULL DEFAULT '0' COMMENT '排序值',
  `notify` tinyint NOT NULL DEFAULT '1' COMMENT '是否通知',
  `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  `history_role_id` int unsigned DEFAULT '0' COMMENT '最新一次历史的角色id',
  `_op_ts` bigint DEFAULT NULL COMMENT '写入时间(元字段)',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID(元字段)',
  PRIMARY KEY (`_id`) USING BTREE,
  KEY `_id` (`_id`)
) ENGINE=MyISAM AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for ods_object_d
-- ----------------------------
DROP TABLE IF EXISTS `ods_object_d`;
CREATE TABLE `ods_object_d` (
  `id` bigint unsigned NOT NULL,
  `space_id` bigint unsigned NOT NULL COMMENT '项目空间id',
  `user_id` bigint unsigned NOT NULL COMMENT '创建用户id',
  `work_object_guid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '工作项(工作模块)Guid',
  `work_object_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '工作项(工作模块)名称',
  `work_object_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '空间状态;0:禁用,1:正常,2:未验证',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '备注',
  `describe` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '描述信息',
  `ranking` bigint DEFAULT '0' COMMENT '排序值',
  `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  `_op_ts` bigint DEFAULT NULL COMMENT '写入时间(元字段)',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID(元字段)',
  PRIMARY KEY (`_id`) USING BTREE,
  KEY `_id` (`_id`)
) ENGINE=MyISAM AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for ods_space_d
-- ----------------------------
DROP TABLE IF EXISTS `ods_space_d`;
CREATE TABLE `ods_space_d` (
  `id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL COMMENT '创建用户id',
  `space_guid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '空间Guid',
  `space_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '空间名称',
  `space_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '空间状态;0:禁用,1:正常,2:未验证',
  `remark` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '备注',
  `describe` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '描述信息',
  `notify` tinyint DEFAULT '1' COMMENT '0不通知 1通知',
  `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  `_op_ts` bigint DEFAULT NULL COMMENT '写入时间(元字段)',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID(元字段)',
  PRIMARY KEY (`_id`),
  KEY `id_Idx` (`_id`) USING BTREE,
  KEY `_op_ts` (`_op_ts`)
) ENGINE=MyISAM AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for ods_user_d
-- ----------------------------
DROP TABLE IF EXISTS `ods_user_d`;
CREATE TABLE `ods_user_d` (
  `id` bigint unsigned NOT NULL,
  `user_name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户名',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '中国手机不带国家代码，国际手机号格式为：国家代码-手机号',
  `user_nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户昵称',
  `user_pinyin` varchar(127) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户昵称拼音',
  `user_password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '登录密码;cmf_password加密',
  `user_salt` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '加密盐',
  `user_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '用户状态;0:禁用,1:正常,2:未验证',
  `user_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户登录邮箱',
  `sex` tinyint DEFAULT '0' COMMENT '性别;0:保密,1:男,2:女',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户头像',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '备注',
  `describe` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '描述信息',
  `last_login_ip` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '最后登录ip',
  `last_login_time` bigint unsigned DEFAULT NULL COMMENT '最后登录时间',
  `role` tinyint NOT NULL DEFAULT '0' COMMENT '角色 0普通成员 50系统管理员 100超级管理员',
  `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  `_op_ts` bigint DEFAULT NULL COMMENT '写入时间(元字段)',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID(元字段)',
  PRIMARY KEY (`_id`) USING BTREE,
  KEY `_id` (`_id`),
  KEY `_op_ts` (`_op_ts`)
) ENGINE=MyISAM AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for ods_version_d
-- ----------------------------
DROP TABLE IF EXISTS `ods_version_d`;
CREATE TABLE `ods_version_d` (
  `id` bigint unsigned NOT NULL,
  `space_id` bigint unsigned NOT NULL COMMENT '项目空间id',
  `version_key` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '工作项（版本） KEY',
  `version_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '工作项(工作模块)名称',
  `version_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '空间状态;0:禁用,1:正常,2:未验证',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '备注',
  `ranking` bigint DEFAULT '0' COMMENT '排序值',
  `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  `_op_ts` bigint DEFAULT NULL COMMENT '写入时间(元字段)',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID(元字段)',
  PRIMARY KEY (`_id`) USING BTREE,
  KEY `_id` (`_id`),
  KEY `_op_ts` (`_op_ts`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='项目空间-版本信息表';

-- ----------------------------
-- Table structure for ods_witem_d
-- ----------------------------
DROP TABLE IF EXISTS `ods_witem_d`;
CREATE TABLE `ods_witem_d` (
  `id` bigint unsigned NOT NULL,
  `pid` bigint unsigned DEFAULT NULL COMMENT '上级id 非0表示当前为子任务',
  `space_id` bigint unsigned NOT NULL COMMENT '项目空间id',
  `work_object_id` bigint unsigned NOT NULL COMMENT '项目工作项id',
  `user_id` bigint unsigned NOT NULL COMMENT '创建用户id',
  `work_item_guid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '工作项(工作模块)Guid',
  `work_item_name` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '工作项(工作模块)名称',
  `work_item_status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务状态Val',
  `work_item_status_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务状态Key',
  `work_item_status_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '任务状态Id',
  `work_item_type_id` bigint DEFAULT NULL COMMENT '工作项类型Id',
  `work_item_type_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '工作项类型key',
  `doc` json DEFAULT NULL,
  `flow_template_id` bigint unsigned DEFAULT NULL COMMENT '工作项流程模版id',
  `flow_template_version` bigint unsigned DEFAULT NULL COMMENT '工作项流程模版版本',
  `work_item_type` tinyint DEFAULT NULL COMMENT '[已废弃]工作项类型 ',
  `flow_id` bigint unsigned DEFAULT NULL COMMENT '工作项流程Id',
  `flow_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '工作项流程Key',
  `work_item_flow_id` bigint unsigned DEFAULT NULL COMMENT '工作项流程Id',
  `flow_mode` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '流程模式 stateflow | workflow',
  `work_item_flow_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '工作项流程Key',
  `flow_mode_version` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '[已废弃] 模式版本',
  `flow_mode_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '[已废弃] 模式编码',
  `last_status_at` bigint unsigned DEFAULT '0' COMMENT '历史状态更新时间',
  `last_status_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '历史状态Key',
  `last_status_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '历史状态Id',
  `last_status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '历史状态',
  `is_restart` tinyint unsigned DEFAULT '0' COMMENT '是否为重启任务',
  `restart_at` bigint unsigned DEFAULT '0' COMMENT '重启时间',
  `icon_flags` int unsigned NOT NULL DEFAULT '0',
  `restart_user_id` bigint unsigned DEFAULT '0' COMMENT '重启用户id',
  `comment_num` int unsigned NOT NULL DEFAULT '0' COMMENT '评论数',
  `resume_at` bigint unsigned DEFAULT '0' COMMENT '任务恢复时间',
  `version_id` bigint unsigned DEFAULT '0' COMMENT '版本信息id',
  `reason` json DEFAULT NULL COMMENT '原因',
  `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  `child_num` int DEFAULT '0' COMMENT '子任务数量',
  `count_at` bigint DEFAULT NULL,
  `_op_ts` bigint DEFAULT NULL COMMENT '写入时间(元字段)',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID(元字段)',
  PRIMARY KEY (`_id`),
  KEY `_id` (`_id`),
  KEY `_op_ts` (`_op_ts`)
) ENGINE=MyISAM AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for ods_witem_flow_node_d
-- ----------------------------
DROP TABLE IF EXISTS `ods_witem_flow_node_d`;
CREATE TABLE `ods_witem_flow_node_d` (
  `id` bigint unsigned NOT NULL,
  `space_id` bigint unsigned DEFAULT '0' COMMENT '空间id',
  `work_item_id` bigint unsigned DEFAULT NULL COMMENT '任务id',
  `flow_mode` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '流程模式',
  `flow_id` bigint unsigned DEFAULT NULL COMMENT '工作项流程Id',
  `flow_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '工作项流程Key',
  `flow_template_id` bigint unsigned DEFAULT NULL COMMENT '流程模版id',
  `flow_template_version` bigint unsigned DEFAULT NULL COMMENT '工作项流程模版版本',
  `flow_node_uuid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '节点uuid',
  `flow_node_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '任务状态; 0:未定义 , 1: 未开启 2: 进行中 3:已完成',
  `flow_node_code` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'P0' COMMENT '节点编码',
  `flow_node_passed` tinyint DEFAULT NULL COMMENT '节点是否通过',
  `flow_node_reached` tinyint DEFAULT NULL COMMENT '节点是否到达',
  `flow_mode_version` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '[已废弃]模式版本',
  `flow_mode_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '[已废弃]模式编码',
  `directors` json DEFAULT NULL COMMENT '节点负责人',
  `work_item_role_id` bigint DEFAULT '0' COMMENT '节点的关联角色id',
  `work_item_role_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '节点的关联角色key',
  `start_at` bigint unsigned DEFAULT NULL COMMENT '开始时间',
  `finish_at` bigint unsigned DEFAULT NULL COMMENT '完成时间',
  `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  `plan_start_at` bigint unsigned DEFAULT NULL COMMENT '计划排期时间-开始',
  `plan_complete_at` bigint unsigned DEFAULT NULL COMMENT '计划排期时间-结束',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID(元字段)',
  `_op_ts` bigint DEFAULT NULL COMMENT '写入时间(元字段)',
  PRIMARY KEY (`_id`) USING BTREE,
  KEY `id` (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for ods_witem_status_d
-- ----------------------------
DROP TABLE IF EXISTS `ods_witem_status_d`;
CREATE TABLE `ods_witem_status_d` (
  `id` bigint unsigned NOT NULL,
  `uuid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一识别码',
  `space_id` bigint unsigned NOT NULL COMMENT '项目空间id',
  `user_id` bigint DEFAULT NULL COMMENT '创建人',
  `work_item_type_id` bigint unsigned NOT NULL COMMENT '工作项类型id',
  `key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `val` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `name` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '角色名称',
  `status_type` tinyint DEFAULT NULL COMMENT '状态类型 1:起始 2:过程 3:归档',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态;0:禁用,1:正常,2:未验证',
  `ranking` bigint DEFAULT '0' COMMENT '排序值',
  `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  `is_sys` tinyint DEFAULT NULL COMMENT '是否系统预设 0:否 1:是',
  `flow_scope` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '可作用的流程范围  空(表示全部) |state_flow | work_flow',
  `_op_ts` bigint DEFAULT NULL COMMENT '写入时间(元字段)',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID(元字段)',
  PRIMARY KEY (`_id`) USING BTREE,
  KEY `_id` (`_id`),
  KEY `_op_ts` (`_op_ts`)
) ENGINE=MyISAM AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='工作项状态-全局';

-- ----------------------------
-- Table structure for dws_space_witem_1h
-- ----------------------------
DROP TABLE IF EXISTS `dws_space_witem_1h`;
CREATE TABLE `dws_space_witem_1h` (
  `space_id` bigint DEFAULT NULL COMMENT '纬度-空间',
  `expire_num` int DEFAULT NULL COMMENT '过期数量',
  `num` int DEFAULT NULL COMMENT '任务数量',
  `todo_num` int DEFAULT NULL COMMENT '待办的任务数量',
  `complete_num` int DEFAULT NULL COMMENT '完成的任务数量',
  `close_num` int DEFAULT NULL COMMENT '关闭的任务数量',
  `abort_num` int DEFAULT NULL COMMENT '终止的任务数量',
  `start_date` datetime DEFAULT NULL COMMENT '汇总的开始时间',
  `end_date` datetime DEFAULT NULL COMMENT '汇总的结束时间',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `space_id` (`space_id`,`start_date`)
) ENGINE=InnoDB AUTO_INCREMENT=262141 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='项目任务统计表/小时';

SET FOREIGN_KEY_CHECKS = 1;
