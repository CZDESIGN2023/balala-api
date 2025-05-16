CREATE TABLE `business_id`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `business_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '业务id',
    `max_id`      bigint unsigned DEFAULT NULL COMMENT '最大id',
    `step`        bigint unsigned DEFAULT NULL COMMENT '步长',
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '描述',
    `create_time` datetime                                               DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime                                               DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `table_name`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_business_id` (`business_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='分布式自增主键';

CREATE TABLE `space_work_item_flow_role_v2`
(
    `id`                 bigint unsigned NOT NULL AUTO_INCREMENT,
    `space_id`           bigint unsigned DEFAULT '0' COMMENT '空间id',
    `work_item_id`       bigint unsigned DEFAULT NULL COMMENT '任务id',
    `flow_id`            bigint unsigned DEFAULT NULL COMMENT '工作项流程Id',
    `flow_template_id`   bigint unsigned DEFAULT NULL COMMENT '流程模版id',
    `work_item_role_id`  bigint                                                        DEFAULT '0' COMMENT '节点的关联角色id',
    `work_item_role_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '节点的关联角色key',
    `directors`          json                                                          DEFAULT NULL COMMENT '节点负责人',
    `created_at`         bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`         bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`         bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

CREATE TABLE `work_flow`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一识别码',
    `user_id`           bigint                                                        DEFAULT NULL COMMENT '创建人',
    `space_id`          bigint unsigned NOT NULL DEFAULT '0' COMMENT '空间id',
    `work_item_type_id` bigint unsigned NOT NULL COMMENT '工作项id',
    `name`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '工作流名称',
    `key`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'key',
    `flow_mode`         varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '流程模式',
    `version`           int unsigned NOT NULL DEFAULT '0' COMMENT '当前最大版本',
    `is_sys`            tinyint                                                       DEFAULT '0' COMMENT '是否系统预设',
    `status`            tinyint                                                       DEFAULT NULL COMMENT '状态 0:禁用 1:启用',
    `ranking`           bigint unsigned DEFAULT NULL COMMENT '排序',
    `last_template_id`  bigint                                                        DEFAULT NULL COMMENT '最大版本流程模版id',
    `created_at`        bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`        bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`        bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY                 `work_item_id_idx` (`work_item_type_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='流程模版表-空间';

CREATE TABLE `work_flow_template`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一识别码',
    `user_id`           bigint  DEFAULT NULL COMMENT '创建人',
    `space_id`          bigint unsigned NOT NULL DEFAULT '0' COMMENT '空间id',
    `work_item_type_id` bigint unsigned NOT NULL COMMENT '工作项id',
    `work_flow_id`      bigint unsigned NOT NULL DEFAULT '0' COMMENT '工作流id',
    `flow_mode`         varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '流程模式',
    `version`           int unsigned NOT NULL DEFAULT '1' COMMENT '版本',
    `setting`           json    DEFAULT NULL COMMENT '配置信息',
    `status`            tinyint DEFAULT NULL COMMENT '状态: 0 禁言 1启用',
    `created_at`        bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`        bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`        bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY                 `work_item_id_idx` (`work_item_type_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='流程模版配置表';

CREATE TABLE `work_item_role`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一识别码',
    `space_id`          bigint unsigned NOT NULL COMMENT '项目空间id',
    `user_id`           bigint                                                        DEFAULT NULL COMMENT '创建人',
    `work_item_type_id` bigint unsigned NOT NULL COMMENT '创建用户id',
    `key`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `name`              varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '角色名称',
    `status`            tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态;0:禁用,1:正常,2:未验证',
    `ranking`           bigint                                                        DEFAULT '0' COMMENT '排序值',
    `is_sys`            tinyint                                                       DEFAULT NULL COMMENT '系统预设 0:否 1:是',
    `created_at`        bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`        bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`        bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='工作项角色-空间';

CREATE TABLE `work_item_status`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一识别码',
    `space_id`          bigint unsigned NOT NULL COMMENT '项目空间id',
    `user_id`           bigint                                                        DEFAULT NULL COMMENT '创建人',
    `work_item_type_id` bigint unsigned NOT NULL COMMENT '工作项类型id',
    `key`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `val`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `name`              varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '角色名称',
    `status_type`       tinyint                                                       DEFAULT NULL COMMENT '状态类型 1:起始 2:过程 3:归档',
    `status`            tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态;0:禁用,1:正常,2:未验证',
    `ranking`           bigint                                                        DEFAULT '0' COMMENT '排序值',
    `created_at`        bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`        bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`        bigint unsigned DEFAULT NULL COMMENT '删除时间',
    `is_sys`            tinyint                                                       DEFAULT NULL COMMENT '是否系统预设 0:否 1:是',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='工作项状态-全局';

CREATE TABLE `work_item_type`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一标识',
    `space_id`   bigint unsigned DEFAULT '0' COMMENT '空间id',
    `user_id`    bigint                                                                 DEFAULT NULL COMMENT '创建人',
    `name`       varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '名称',
    `key`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '编码[设置后不允许修改]',
    `flow_mode`  varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '工作流模式',
    `ranking`    bigint unsigned DEFAULT '0' COMMENT '排序',
    `is_sys`     tinyint                                                                DEFAULT '0' COMMENT '是否系统预设',
    `status`     tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态;0:禁用,1:正常,2:未验证',
    `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='任务类型-空间';


ALTER TABLE `space_work_item_v2`
    ADD COLUMN `work_item_status_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务状态Key' AFTER `work_item_status`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `work_item_status_id` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT '任务状态Id' AFTER `work_item_status_key`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `work_item_type_id` bigint NULL DEFAULT NULL COMMENT '工作项类型Id' AFTER `work_item_status_id`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `work_item_type_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '工作项类型key' AFTER `work_item_type_id`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `flow_template_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '工作项流程模版id' AFTER `doc`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `flow_template_version` bigint UNSIGNED NULL DEFAULT NULL COMMENT '工作项流程模版版本' AFTER `flow_template_id`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `flow_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '工作项流程Id' AFTER `work_item_type`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `flow_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '工作项流程Key' AFTER `flow_id`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `work_item_flow_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '工作项流程Id' AFTER `flow_key`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `work_item_flow_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '工作项流程Key' AFTER `flow_mode`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `last_status_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '历史状态Key' AFTER `last_status_at`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `last_status_id` bigint UNSIGNED NOT NULL DEFAULT 0 COMMENT '历史状态Id' AFTER `last_status_key`;
ALTER TABLE `space_work_item_v2` DROP COLUMN `child_num`;
ALTER TABLE `space_work_item_v2`
    ADD COLUMN `child_num` int NULL DEFAULT 0 COMMENT '子任务数量' AFTER `deleted_at`;

ALTER TABLE `space_work_item_flow_v2`
    ADD COLUMN `flow_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '工作项流程Id' AFTER `flow_mode`;
ALTER TABLE `space_work_item_flow_v2`
    ADD COLUMN `flow_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '工作项流程Key' AFTER `flow_id`;
ALTER TABLE `space_work_item_flow_v2`
    ADD COLUMN `flow_template_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT '流程模版id' AFTER `flow_key`;
ALTER TABLE `space_work_item_flow_v2`
    ADD COLUMN `flow_template_version` bigint UNSIGNED NULL DEFAULT NULL COMMENT '工作项流程模版版本' AFTER `flow_template_id`;
ALTER TABLE `space_work_item_flow_v2`
    ADD COLUMN `work_item_role_id` bigint NULL DEFAULT 0 COMMENT '节点的关联角色id' AFTER `directors`;
ALTER TABLE `space_work_item_flow_v2`
    ADD COLUMN `work_item_role_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '节点的关联角色key' AFTER `work_item_role_id`;

ALTER TABLE `space_work_item_v2` MODIFY COLUMN `work_item_status` varchar (255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务状态Val' AFTER `work_item_name`;
ALTER TABLE `space_work_item_v2` MODIFY COLUMN `work_item_type` tinyint NULL DEFAULT NULL COMMENT '[已废弃]工作项类型 ' AFTER `flow_template_version`;
ALTER TABLE `space_work_item_v2` MODIFY COLUMN `flow_mode_version` varchar (100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '[已废弃] 模式版本' AFTER `work_item_flow_key`;
ALTER TABLE `space_work_item_v2` MODIFY COLUMN `flow_mode_code` varchar (100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '[已废弃] 模式编码' AFTER `flow_mode_version`;
ALTER TABLE `space_work_item_v2` MODIFY COLUMN `last_status` varchar (255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '历史状态' AFTER `last_status_id`;
ALTER TABLE `space_work_item_v2` MODIFY COLUMN `created_at` bigint UNSIGNED NULL DEFAULT NULL COMMENT '创建时间' AFTER `version_id`;
ALTER TABLE `space_work_item_v2` MODIFY COLUMN `updated_at` bigint UNSIGNED NULL DEFAULT NULL COMMENT '更新时间' AFTER `created_at`;
ALTER TABLE `space_work_item_v2` MODIFY COLUMN `deleted_at` bigint UNSIGNED NULL DEFAULT NULL COMMENT '删除时间' AFTER `updated_at`;
ALTER TABLE `third_pf_account` MODIFY COLUMN `pf_user_id` bigint NULL DEFAULT NULL COMMENT '平台用户id' AFTER `pf_user_key`;
ALTER TABLE `user` MODIFY COLUMN `user_name` varchar (60) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs NULL DEFAULT NULL COMMENT '用户名' AFTER `id`;
ALTER TABLE `user` MODIFY COLUMN `user_nickname` varchar (50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户昵称' AFTER `mobile`;
ALTER TABLE `space_work_item_flow_v2` MODIFY COLUMN `flow_mode` varchar (30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '流程模式' AFTER `work_item_id`;
ALTER TABLE `space_work_item_flow_v2` MODIFY COLUMN `flow_mode_version` varchar (100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '[已废弃]模式版本' AFTER `flow_node_reached`;
ALTER TABLE `space_work_item_flow_v2` MODIFY COLUMN `flow_mode_code` varchar (100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '[已废弃]模式编码' AFTER `flow_mode_version`;
ALTER TABLE `space_work_item_flow_v2` MODIFY COLUMN `directors` json NULL COMMENT '节点负责人' AFTER `flow_mode_code`;
ALTER TABLE `space_tag` MODIFY COLUMN `tag_name` varchar (50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs NOT NULL DEFAULT '' COMMENT '空间名称';
ALTER TABLE `space_work_object` MODIFY COLUMN `work_object_name` varchar (50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs NOT NULL DEFAULT '' COMMENT '工作项(工作模块)名称';

-- 调整索引
ALTER TABLE `space_work_item_v2`
    ADD INDEX `participators_idx`((cast(json_extract(`doc`,_utf8mb4'$.participators') as char(20) array)) );
ALTER TABLE `config` drop INDEX `idx_config_key`;
ALTER TABLE `config`
    ADD UNIQUE INDEX `idx_config_key`(`config_key` ASC) USING BTREE;

-- 清理废弃的表
DROP TABLE IF EXISTS `member_category`;
DROP TABLE IF EXISTS `space_member_category`;
DROP TABLE IF EXISTS `space_work_item_type`;

-- 初始化业务id
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (1, 'work_flow', 8051, 30, NULL, '2024-05-09 08:07:12', '2024-07-04 11:48:43', 'work_flow');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (2, 'work_flow_template', 9297, 30, NULL, '2024-05-09 08:07:46', '2024-07-04 11:48:43', 'work_flow_template');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (3, 'work_item_status', 8233, 30, NULL, '2024-05-09 09:16:15', '2024-07-04 11:48:43', 'work_item_status');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (4, 'work_item_role', 6313, 30, NULL, '2024-05-09 10:16:10', '2024-07-04 11:48:43', 'work_item_role');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (5, 'space', 1879, 30, NULL, '2024-05-10 12:32:24', '2024-07-04 11:48:43', 'space');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (6, 'user', 2024, 30, NULL, '2024-05-10 12:32:56', '2024-07-03 09:22:41', 'user');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (7, 'space_member', 127, 30, NULL, '2024-05-11 02:01:35', '2024-07-04 11:48:43', 'space_member');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (8, 'work_item_type', 2847, 30, NULL, '2024-05-22 08:28:47', '2024-07-04 11:48:43', 'work_item_type');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (9, 'space_work_object', 2851, 30, NULL, '2024-05-29 08:29:03', '2024-07-04 11:48:43', 'space_work_object');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (10, 'space_work_item', 4863, 30, NULL, '2024-05-29 09:01:53', '2024-07-04 11:48:43', 'space_work_item_v2');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (11, 'space_tag', 2303, 30, NULL, '2024-06-07 10:50:51', '2024-07-04 11:48:43', 'space_tag');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (12, 'file_info', 2458, 30, NULL, '2024-06-12 03:43:22', '2024-07-04 06:59:51', 'file_info');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (13, 'space_file_info', 6, 30, NULL, '2024-06-12 03:43:39', '2024-07-04 06:59:51', 'space_file_info');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (14, 'space_work_version', 2279, 30, NULL, '2024-06-27 02:42:51', '2024-07-04 11:48:43', 'space_work_version');
		